package provider

import (
	"encoding/json"
	"testing"
)

func TestAlertRuleReadReconstructsConfiguration(t *testing.T) {
	configuredJSON := "{\n  \"auth\": {\"token\": \"secret\", \"type\": \"bearer\"}\n}"
	v := Values{"webhook_actions": []any{map[string]any{"webhook_url": "https://example.com/hook", "config": configuredJSON}}}
	obj := Object{
		"eventType": "site_offline", "siteIds": []any{}, "healthCheckIds": []any{}, "resourceIds": []any{},
		"recipients":     []any{map[string]any{"userId": "user1", "roleId": nil, "email": nil}, map[string]any{"roleId": json.Number("7")}, map[string]any{"email": "ops@example.com"}},
		"webhookActions": []any{map[string]any{"webhookUrl": "https://example.com/hook", "enabled": true, "config": map[string]any{"auth": map[string]any{"type": "bearer", "token": "secret"}}}},
	}
	if err := enrichAlertRule(v, obj); err != nil {
		t.Fatal(err)
	}
	if got := Strings(obj["userIds"]); len(got) != 1 || got[0] != "user1" {
		t.Fatalf("users lost: %#v", obj)
	}
	if got := Strings(obj["roleIds"]); len(got) != 1 || got[0] != "7" {
		t.Fatalf("roles lost: %#v", obj)
	}
	if got := Strings(obj["emails"]); len(got) != 1 || got[0] != "ops@example.com" {
		t.Fatalf("emails lost: %#v", obj)
	}
	if obj["allSites"] != true || obj["allHealthChecks"] != false || obj["allResources"] != false {
		t.Fatalf("incorrect source scope: %#v", obj)
	}
	if got := objectArray(obj["webhookActions"])[0]["config"]; got != configuredJSON {
		t.Fatalf("JSON formatting changed: %v", got)
	}
}

func TestAlertRuleReadDetectsWebhookDrift(t *testing.T) {
	v := Values{"webhook_actions": []any{map[string]any{"webhook_url": "https://example.com/hook", "config": `{"token":"old"}`}}}
	obj := Object{"eventType": "health_check_unhealthy", "healthCheckIds": []any{json.Number("5")}, "recipients": []any{}, "webhookActions": []any{map[string]any{"webhookUrl": "https://example.com/hook", "config": map[string]any{"token": "changed"}}}}
	if err := enrichAlertRule(v, obj); err != nil {
		t.Fatal(err)
	}
	if obj["allHealthChecks"] != false {
		t.Fatal("an explicit source was treated as all sources")
	}
	if got := objectArray(obj["webhookActions"])[0]["config"]; got != `{"token":"changed"}` {
		t.Fatalf("remote config drift suppressed: %v", got)
	}
	if got := obj["emails"].([]any); len(got) != 0 {
		t.Fatal("deleted recipients preserved")
	}
}

func TestAlertRuleValidatesEventScope(t *testing.T) {
	for _, v := range []Values{
		{"event_type": "site_online", "all_sites": true},
		{"event_type": "health_check_unhealthy", "health_check_ids": []any{"4"}},
		{"event_type": "resource_degraded", "resource_ids": []any{"9"}},
	} {
		if err := validateAlertRule(v); err != nil {
			t.Error(err)
		}
	}
	for _, v := range []Values{
		{"event_type": "site_online"},
		{"event_type": "site_online", "all_sites": true, "site_ids": []any{"4"}},
		{"event_type": "site_online", "all_sites": true, "resource_ids": []any{"4"}},
		{"event_type": "resource_healthy", "all_resources": true, "webhook_actions": []any{map[string]any{"config": "not JSON"}}},
	} {
		if err := validateAlertRule(v); err == nil {
			t.Errorf("invalid alert accepted: %#v", v)
		}
	}
}

func TestMonitoringListIdentityAndSensitiveFields(t *testing.T) {
	health := healthCheckDefinition()
	if health.IDField != "targetHealthCheckId" || health.ListKey != "healthChecks" || health.Read.Path != "" || health.Pagination != "offset" {
		t.Fatalf("health-check read must use the real API list identity: %#v", health)
	}
	if !health.Fields["hc_headers"].Sensitive {
		t.Fatal("health-check headers may contain credentials")
	}
	event := eventStreamingDefinition()
	if event.ListKey != "destinations" || event.Read.Path != "" || event.List.Path != "/org/{org_id}/event-streaming-destinations" || !event.Fields["config"].Sensitive {
		t.Fatal("event destinations require list reads and sensitive config")
	}
	alert := alertRuleDefinition()
	if !alert.Fields["webhook_actions"].Sensitive {
		t.Fatal("webhook configuration may contain credentials")
	}
}

func TestHealthCheckPositiveBoundsAndHostname(t *testing.T) {
	d := healthCheckDefinition()
	for _, name := range []string{"hc_interval", "hc_unhealthy_interval", "hc_timeout", "hc_healthy_threshold", "hc_unhealthy_threshold"} {
		field := d.Fields[name]
		if field.Min == nil || *field.Min != 1 {
			t.Fatalf("%s must accept the server's positive-integer minimum 1", name)
		}
		if err := validateFields(map[string]Field{name: field}, Values{name: int64(1)}); err != nil {
			t.Errorf("%s rejected 1: %v", name, err)
		}
		if err := validateFields(map[string]Field{name: field}, Values{name: int64(0)}); err == nil {
			t.Errorf("%s accepted 0", name)
		}
	}
	for _, v := range []Values{{"hc_enabled": true}, {"hc_enabled": true, "hc_hostname": "  "}} {
		if err := d.Validate(v); err == nil {
			t.Errorf("enabled health check accepted missing hostname: %#v", v)
		}
	}
	if err := d.Validate(Values{"hc_enabled": false}); err != nil {
		t.Fatal(err)
	}
}
