package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/chang196700/terraform-provider-pangolin/internal/client"
)

func init() { registerDefinitions(monitoringDefinitions) }

func monitoringDefinitions() []Definition {
	return []Definition{eventStreamingDefinition(), alertRuleDefinition(), healthCheckDefinition()}
}

func eventStreamingDefinition() Definition {
	d := Definition{
		Name: "event_streaming_destination", Description: "Manages an organization event streaming destination. The opaque destination configuration is sensitive because it may contain authentication credentials.",
		IDField: "destinationId", ListKey: "destinations", Pagination: "offset", ImportFields: []string{"org_id", "id"},
		Fields: SpecFields("/org/{orgId}/event-streaming-destination", "put", "/org/{orgId}/event-streaming-destination/{destinationId}", "post"),
		Create: Endpoint{Method: "PUT", Path: "/org/{org_id}/event-streaming-destination"},
		Update: Endpoint{Method: "POST", Path: "/org/{org_id}/event-streaming-destination/{id}"},
		Delete: Endpoint{Method: "DELETE", Path: "/org/{org_id}/event-streaming-destination/{id}"},
		// The 1.24.0 router mounts the plural path. Its OpenAPI registration
		// incorrectly declares the singular path used by the create endpoint.
		List: Endpoint{Method: "GET", Path: "/org/{org_id}/event-streaming-destinations"},
	}
	d.Fields["org_id"] = Local(Immutable(Required(String("orgId"))))
	SetField(&d, "config", func(f *Field) { f.Sensitive = true })
	d.Fields["last_error"] = Computed(String("lastError"))
	d.Fields["last_error_at"] = Computed(Int("lastErrorAt"))
	aiTimestamps(&d)
	return d
}

func alertRuleDefinition() Definition {
	d := Definition{
		Name: "alert_rule", Description: "Manages an alert rule, event sources, email recipients and webhook actions. Source and recipient collections are authoritative. The provider reconciles the API's recipient rows and webhook configuration objects back into configuration fields.",
		IDField: "alertRuleId", ListKey: "alertRules", Pagination: "offset", ImportFields: []string{"org_id", "id"},
		Fields: SpecFields("/org/{orgId}/alert-rule", "put", "/org/{orgId}/alert-rule/{alertRuleId}", "post"),
		Create: Endpoint{Method: "PUT", Path: "/org/{org_id}/alert-rule"},
		Read:   Endpoint{Method: "GET", Path: "/org/{org_id}/alert-rule/{id}"},
		Update: Endpoint{Method: "POST", Path: "/org/{org_id}/alert-rule/{id}"},
		Delete: Endpoint{Method: "DELETE", Path: "/org/{org_id}/alert-rule/{id}"},
		List:   Endpoint{Method: "GET", Path: "/org/{org_id}/alert-rules"},
	}
	d.Fields["org_id"] = Local(Immutable(Required(String("orgId"))))
	for _, name := range []string{"site_ids", "health_check_ids", "resource_ids", "role_ids", "user_ids", "emails"} {
		SetField(&d, name, func(f *Field) {
			f.Kind = "set"
			if name != "user_ids" && name != "emails" && f.Element != nil {
				e := *f.Element
				e.Kind = "id"
				f.Element = &e
			}
		})
	}
	SetField(&d, "webhook_actions", func(f *Field) {
		f.Sensitive = true
		if f.Element != nil {
			// Empty or omitted configuration is stored and returned as null.
			// Use an empty string consistently so removing it clears the action
			// credentials while configured webhook lists retain stable values.
			config := f.Element.Fields["config"]
			config.Default = ""
			config.Sensitive = true
			f.Element.Fields["config"] = config
		}
	})
	d.Fields["last_triggered_at"] = Computed(Int("lastTriggeredAt"))
	aiTimestamps(&d)
	d.Validate = validateAlertRule
	d.EnrichRead = func(_ context.Context, _ *client.Client, v Values, obj Object) error { return enrichAlertRule(v, obj) }
	return d
}

func validateAlertRule(v Values) error {
	event := Text(v["event_type"])
	groups := []struct{ prefix, ids, all string }{
		{"site_", "site_ids", "all_sites"},
		{"health_check_", "health_check_ids", "all_health_checks"},
		{"resource_", "resource_ids", "all_resources"},
	}
	for _, group := range groups {
		count := len(Strings(v[group.ids]))
		all, _ := v[group.all].(bool)
		active := strings.HasPrefix(event, group.prefix)
		if active && !all && count == 0 {
			return fmt.Errorf("%s requires at least one %s or %s", event, group.ids, group.all)
		}
		if all && count > 0 {
			return fmt.Errorf("%s must be empty when %s is true", group.ids, group.all)
		}
		if !active && (count > 0 || all) {
			return fmt.Errorf("%s and %s do not apply to event_type %s", group.ids, group.all, event)
		}
	}
	for _, item := range objectArray(v["webhook_actions"]) {
		if cfg := item["config"]; cfg != nil && Text(cfg) != "" {
			var value map[string]any
			if err := json.Unmarshal([]byte(Text(cfg)), &value); err != nil || value == nil {
				return fmt.Errorf("webhook_actions config must be a JSON object")
			}
		}
	}
	return nil
}

// Pangolin's GET returns recipient rows instead of the input arrays and omits
// the three all* flags. An empty source collection means all sources for the
// active event family; that follows the server's create-time validation.
func enrichAlertRule(v Values, obj Object) error {
	if recipients, ok := obj["recipients"]; ok {
		obj["userIds"], obj["roleIds"], obj["emails"] = []any{}, []any{}, []any{}
		for _, recipient := range objectArray(recipients) {
			for _, pair := range [][2]string{{"userId", "userIds"}, {"roleId", "roleIds"}, {"email", "emails"}} {
				if value := recipient[pair[0]]; value != nil {
					obj[pair[1]] = append(obj[pair[1]].([]any), value)
				}
			}
		}
	}
	event := Text(obj["eventType"])
	for _, group := range [][3]string{{"site_", "siteIds", "allSites"}, {"health_check_", "healthCheckIds", "allHealthChecks"}, {"resource_", "resourceIds", "allResources"}} {
		if _, exists := obj[group[2]]; !exists {
			obj[group[2]] = strings.HasPrefix(event, group[0]) && len(Strings(obj[group[1]])) == 0
		}
	}
	priorWebhooks := objectArray(v["webhook_actions"])
	for _, webhook := range objectArray(obj["webhookActions"]) {
		cfg := webhook["config"]
		if cfg == nil {
			webhook["config"] = ""
			continue
		}
		encoded, err := json.Marshal(cfg)
		if err != nil {
			return fmt.Errorf("decode webhook configuration: %w", err)
		}
		if text, ok := cfg.(string); ok {
			encoded = []byte(text)
		}
		// Preserve JSON whitespace/key ordering when the remote value is
		// semantically identical, as Terraform requires configured strings to
		// remain byte-for-byte equal following Apply.
		for _, previous := range priorWebhooks {
			if Text(previous["webhook_url"]) == Text(webhook["webhookUrl"]) {
				var configured, actual any
				if json.Unmarshal([]byte(Text(previous["config"])), &configured) == nil && json.Unmarshal(encoded, &actual) == nil && reflect.DeepEqual(configured, actual) {
					encoded = []byte(Text(previous["config"]))
					break
				}
			}
		}
		webhook["config"] = string(encoded)
	}
	return nil
}

func objectArray(v any) []Object {
	switch values := v.(type) {
	case []Object:
		return values
	case []any:
		out := make([]Object, 0, len(values))
		for _, value := range values {
			if obj, ok := value.(map[string]any); ok {
				out = append(out, obj)
			}
		}
		return out
	default:
		return nil
	}
}

func healthCheckDefinition() Definition {
	d := Definition{
		Name: "health_check", Description: "Manages a standalone site health check. Read refreshes the paginated organization list using targetHealthCheckId, the identifier returned by the API.",
		IDField: "targetHealthCheckId", ListKey: "healthChecks", Pagination: "offset", ImportFields: []string{"org_id", "id"},
		Fields: SpecFields("/org/{orgId}/health-check", "put", "/org/{orgId}/health-check/{healthCheckId}", "post"),
		Create: Endpoint{Method: "PUT", Path: "/org/{org_id}/health-check"},
		Update: Endpoint{Method: "POST", Path: "/org/{org_id}/health-check/{id}"},
		Delete: Endpoint{Method: "DELETE", Path: "/org/{org_id}/health-check/{id}"},
		List:   Endpoint{Method: "GET", Path: "/org/{org_id}/health-checks"},
	}
	d.Fields["org_id"] = Local(Immutable(Required(String("orgId"))))
	SetField(&d, "site_id", func(f *Field) { f.Kind = "id" })
	SetField(&d, "hc_headers", func(f *Field) { f.Sensitive = true })
	d.Fields["hc_health"] = Computed(String("hcHealth"))
	d.Fields["site_name"] = Computed(String("siteName"))
	d.Fields["site_nice_id"] = Computed(String("siteNiceId"))
	d.Validate = func(v Values) error {
		if enabled, _ := v["hc_enabled"].(bool); enabled && strings.TrimSpace(Text(v["hc_hostname"])) == "" {
			return fmt.Errorf("hc_hostname is required when hc_enabled is true")
		}
		return nil
	}
	return d
}
