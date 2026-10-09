package provider

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// These fixtures model the pinned upstream wire contract independently of the
// provider definitions. In particular IDs, projections and response envelopes
// are literal server fields, not inferred from provider configuration.
type extendedFixture struct {
	mu                                                          sync.Mutex
	exists                                                      bool
	row                                                         Object
	seed                                                        Object
	createPath, itemPath, listPath, listKey, envelopeKey, idKey string
	writeIDOnly                                                 bool
	checkUpdate                                                 func(*testing.T, Object)
	project                                                     func(Object) Object
}

func (f *extendedFixture) remove()       { f.mu.Lock(); defer f.mu.Unlock(); f.exists = false }
func (f *extendedFixture) present() bool { f.mu.Lock(); defer f.mu.Unlock(); return f.exists }

func (f *extendedFixture) handler(t *testing.T) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		f.mu.Lock()
		defer f.mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		respond := func(status int, data any) {
			w.WriteHeader(status)
			_ = json.NewEncoder(w).Encode(Object{"success": true, "error": false, "status": status, "data": data})
		}
		missing := func() {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"success":false,"error":true,"status":404,"data":null}`))
		}
		project := func() Object {
			row := extendedClone(f.row)
			if f.project != nil {
				row = f.project(row)
			}
			return row
		}
		wrap := func(row Object) Object {
			if f.envelopeKey != "" {
				return Object{f.envelopeKey: row}
			}
			return row
		}
		var body Object
		if r.Method == "PUT" || r.Method == "POST" {
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Errorf("invalid request body: %v", err)
				w.WriteHeader(400)
				return
			}
		}
		switch {
		case r.Method == "PUT" && r.URL.Path == f.createPath:
			if f.exists {
				t.Errorf("duplicate create for %s", f.createPath)
				w.WriteHeader(409)
				return
			}
			f.row = extendedClone(f.seed)
			for k, v := range body {
				f.row[k] = v
			}
			f.exists = true
			if f.idKey == "virtualApiKeyId" {
				delete(f.row, "validForSeconds")
			}
			if f.writeIDOnly {
				respond(201, Object{f.idKey: f.row[f.idKey]})
			} else {
				respond(201, wrap(project()))
			}
		case r.Method == "GET" && f.listPath != "" && r.URL.Path == f.listPath:
			rows := []Object{}
			if f.exists {
				rows = append(rows, project())
			}
			respond(200, Object{f.listKey: rows, "pagination": Object{"total": len(rows), "limit": 100, "offset": 0}})
		case r.Method == "GET" && r.URL.Path == f.itemPath:
			if !f.exists {
				missing()
				return
			}
			respond(200, wrap(project()))
		case r.Method == "POST" && r.URL.Path == f.itemPath:
			if !f.exists {
				missing()
				return
			}
			if f.checkUpdate != nil {
				f.checkUpdate(t, body)
			}
			for k, v := range body {
				f.row[k] = v
			}
			if f.writeIDOnly {
				respond(200, Object{f.idKey: f.row[f.idKey]})
			} else {
				row := project()
				if f.idKey == "virtualApiKeyId" {
					delete(row, "secret")
				}
				respond(200, wrap(row))
			}
		case r.Method == "DELETE" && r.URL.Path == f.itemPath:
			if !f.exists {
				missing()
				return
			}
			f.exists = false
			respond(200, nil)
		default:
			t.Errorf("unexpected extended API request %s %s", r.Method, r.URL)
			missing()
		}
	}
}

func extendedClone(in Object) Object {
	out := Object{}
	for k, v := range in {
		out[k] = v
	}
	return out
}

func TestAccExtendedAIProvider(t *testing.T) {
	f := &extendedFixture{
		createPath: "/v1/org/test-org/ai-provider", itemPath: "/v1/ai-provider/41", envelopeKey: "provider", idKey: "providerId",
		seed: Object{"providerId": 41, "orgId": "test-org", "name": "before", "niceId": "generated-name", "type": "openai", "upstreamUrl": "https://api.example.com/v1", "apiKey": "fixture-secret", "apiKeyLastChars": "cret", "routingMode": "url", "authType": "bearer", "capabilities": []any{"openai_chat"}, "headers": nil, "enabled": true, "skipTlsVerification": false, "createdAt": 1000, "updatedAt": 1000},
		project: func(o Object) Object {
			o["effectiveUpstreamUrl"] = o["upstreamUrl"]
			o["effectiveAuthType"] = o["authType"]
			if len(objectArray(o["headers"])) == 0 {
				o["headers"] = nil
			}
			if o["apiKey"] == "" {
				o["apiKey"] = nil
			}
			return o
		},
	}
	attrs := `org_id = "test-org"
name = %q
nice_id = "team-provider"
type = "openai"
upstream_url = "https://api.example.com/v1"
api_key = %q
headers = []`
	runLifecycle(t, lifecycleCase{Resource: "ai_provider", CreateAttributes: fmt.Sprintf(attrs, "before", ""), UpdateAttributes: fmt.Sprintf(attrs, "after", "fixture-secret"), ImportID: "test-org/41", Handler: f.handler(t), Remove: f.remove, Exists: f.present,
		Checks:       []resource.TestCheckFunc{resource.TestCheckResourceAttr("pangolin_ai_provider.test", "nice_id", "team-provider"), resource.TestCheckResourceAttr("pangolin_ai_provider.test", "api_key", ""), resource.TestCheckResourceAttr("pangolin_ai_provider.test", "headers.#", "0")},
		UpdateChecks: []resource.TestCheckFunc{resource.TestCheckResourceAttr("pangolin_ai_provider.test", "name", "after"), resource.TestCheckResourceAttr("pangolin_ai_provider.test", "api_key", "fixture-secret")},
	})
}

func TestAccExtendedAIModel(t *testing.T) {
	f := &extendedFixture{createPath: "/v1/ai-provider/41/model", itemPath: "/v1/ai-model/42", envelopeKey: "model", idKey: "modelId", seed: Object{"modelId": 42, "providerId": 41, "modelKey": "example-model", "name": "before", "enabled": true, "listType": "allow", "createdAt": 1000, "updatedAt": 1000}}
	attrs := "provider_id = \"41\"\nmodel_key = \"example-model\"\nname = %q"
	runLifecycle(t, lifecycleCase{Resource: "ai_model", CreateAttributes: fmt.Sprintf(attrs, "before"), UpdateAttributes: fmt.Sprintf(attrs, "after"), ImportID: "41/42", Handler: f.handler(t), Remove: f.remove, Exists: f.present,
		UpdateChecks: []resource.TestCheckFunc{resource.TestCheckResourceAttr("pangolin_ai_model.test", "name", "after")}})
}

func TestAccExtendedAIBudget(t *testing.T) {
	f := &extendedFixture{createPath: "/v1/org/test-org/ai-budget", itemPath: "/v1/ai-budget/43", envelopeKey: "budget", idKey: "budgetId", seed: Object{"budgetId": 43, "orgId": "test-org", "providerId": nil, "modelId": nil, "resourceId": nil, "siteResourceId": nil, "roleId": nil, "virtualApiKeyId": nil, "amount": 100, "unit": "usd", "period": "monthly", "enforcement": "hard", "enabled": true, "createdAt": 1000, "updatedAt": 1000},
		checkUpdate: func(t *testing.T, body Object) {
			if body["modelId"] == float64(42) {
				if v, ok := body["providerId"]; !ok || v != nil {
					t.Error("budget scope update did not clear providerId")
				}
			}
		},
	}
	create := "org_id = \"test-org\"\nprovider_id = \"41\"\namount = 100\nunit = \"usd\""
	update := "org_id = \"test-org\"\nmodel_id = \"42\"\namount = 200\nunit = \"usd\""
	runLifecycle(t, lifecycleCase{Resource: "ai_budget", CreateAttributes: create, UpdateAttributes: update, ImportID: "test-org/43", Handler: f.handler(t), Remove: f.remove, Exists: f.present,
		UpdateChecks: []resource.TestCheckFunc{resource.TestCheckResourceAttr("pangolin_ai_budget.test", "model_id", "42"), resource.TestCheckResourceAttr("pangolin_ai_budget.test", "amount", "200")}})
}

func TestAccExtendedVirtualAPIKey(t *testing.T) {
	f := &extendedFixture{createPath: "/v1/org/test-org/virtual-api-key", itemPath: "/v1/virtual-api-key/key1", envelopeKey: "virtualApiKey", idKey: "virtualApiKeyId", seed: Object{"virtualApiKeyId": "key1", "orgId": "test-org", "name": "before", "description": nil, "userId": nil, "kind": "manual", "secret": "fixture-key-secret", "lastChars": "cret", "allResources": true, "resourceIds": []any{}, "expiresAt": 3601000, "lastUsedAt": nil, "createdAt": 1000, "createdByUserId": nil},
		checkUpdate: func(t *testing.T, body Object) {
			if _, ok := body["validForSeconds"]; ok {
				t.Error("renaming a virtual API key unexpectedly renewed its expiry")
			}
			if _, ok := body["sendEmail"]; ok {
				t.Error("updating a virtual API key unexpectedly sent email")
			}
		},
	}
	attrs := "org_id = \"test-org\"\nname = %q\nall_resources = true\nvalid_for_seconds = 3600"
	runLifecycle(t, lifecycleCase{Resource: "virtual_api_key", CreateAttributes: fmt.Sprintf(attrs, "before"), UpdateAttributes: fmt.Sprintf(attrs, "after"), ImportID: "test-org/key1", ImportIgnore: []string{"valid_for_seconds"}, Handler: f.handler(t), Remove: f.remove, Exists: f.present,
		Checks: []resource.TestCheckFunc{resource.TestCheckResourceAttr("pangolin_virtual_api_key.test", "secret", "fixture-key-secret")}, UpdateChecks: []resource.TestCheckFunc{resource.TestCheckResourceAttr("pangolin_virtual_api_key.test", "name", "after"), resource.TestCheckResourceAttr("pangolin_virtual_api_key.test", "expires_at", "3601000")}})
}

func TestAccExtendedEventStreamingDestination(t *testing.T) {
	f := &extendedFixture{createPath: "/v1/org/test-org/event-streaming-destination", itemPath: "/v1/org/test-org/event-streaming-destination/44", listPath: "/v1/org/test-org/event-streaming-destinations", listKey: "destinations", idKey: "destinationId", writeIDOnly: true, seed: Object{"destinationId": 44, "orgId": "test-org", "type": "webhook", "config": `{"url":"https://example.com/events"}`, "enabled": true, "sendConnectionLogs": false, "sendRequestLogs": false, "sendActionLogs": false, "sendAccessLogs": false, "sendAISessionLogs": false, "lastError": nil, "lastErrorAt": nil, "createdAt": 1000, "updatedAt": 1000}}
	attrs := "org_id = \"test-org\"\ntype = \"webhook\"\nconfig = jsonencode({url = \"https://example.com/events\"})\nenabled = %t"
	runLifecycle(t, lifecycleCase{Resource: "event_streaming_destination", CreateAttributes: fmt.Sprintf(attrs, true), UpdateAttributes: fmt.Sprintf(attrs, false), ImportID: "test-org/44", Handler: f.handler(t), Remove: f.remove, Exists: f.present,
		UpdateChecks: []resource.TestCheckFunc{resource.TestCheckResourceAttr("pangolin_event_streaming_destination.test", "enabled", "false")}})
}

func TestAccExtendedAlertRule(t *testing.T) {
	f := &extendedFixture{createPath: "/v1/org/test-org/alert-rule", itemPath: "/v1/org/test-org/alert-rule/45", idKey: "alertRuleId", writeIDOnly: true, seed: Object{"alertRuleId": 45, "orgId": "test-org", "name": "before", "eventType": "site_offline", "enabled": true, "cooldownSeconds": 0, "siteIds": []any{}, "healthCheckIds": []any{}, "resourceIds": []any{}, "userIds": []any{}, "roleIds": []any{}, "emails": []any{}, "webhookActions": []any{}, "lastTriggeredAt": nil, "createdAt": 1000, "updatedAt": 1000}, project: extendedAlertProjection}
	attrs := `org_id = "test-org"
name = %q
event_type = "site_offline"
all_sites = true
emails = ["ops@example.com"]
webhook_actions = [
  { webhook_url = "https://example.com/hook", config = jsonencode({token = "fixture-token"}) },
  { webhook_url = "https://example.com/no-auth" }
]`
	runLifecycle(t, lifecycleCase{Resource: "alert_rule", CreateAttributes: fmt.Sprintf(attrs, "before"), UpdateAttributes: fmt.Sprintf(attrs, "after"), ImportID: "test-org/45", Handler: f.handler(t), Remove: f.remove, Exists: f.present,
		Checks: []resource.TestCheckFunc{resource.TestCheckResourceAttr("pangolin_alert_rule.test", "emails.#", "1"), resource.TestCheckResourceAttr("pangolin_alert_rule.test", "all_sites", "true")}, UpdateChecks: []resource.TestCheckFunc{resource.TestCheckResourceAttr("pangolin_alert_rule.test", "name", "after")}})
}

func extendedAlertProjection(row Object) Object {
	for _, key := range []string{"allSites", "allHealthChecks", "allResources"} {
		delete(row, key)
	}
	recipients := []any{}
	for _, pair := range [][2]string{{"userIds", "userId"}, {"roleIds", "roleId"}, {"emails", "email"}} {
		for _, value := range Strings(row[pair[0]]) {
			recipients = append(recipients, Object{pair[1]: value})
		}
		delete(row, pair[0])
	}
	row["recipients"] = recipients
	webhooks := []any{}
	for _, value := range objectArray(row["webhookActions"]) {
		webhook := extendedClone(value)
		if webhook["enabled"] == nil {
			webhook["enabled"] = true
		}
		if text, ok := webhook["config"].(string); ok {
			var decoded any
			_ = json.Unmarshal([]byte(text), &decoded)
			webhook["config"] = decoded
		}
		webhooks = append(webhooks, webhook)
	}
	row["webhookActions"] = webhooks
	return row
}

func TestAccExtendedHealthCheck(t *testing.T) {
	f := &extendedFixture{createPath: "/v1/org/test-org/health-check", itemPath: "/v1/org/test-org/health-check/46", listPath: "/v1/org/test-org/health-checks", listKey: "healthChecks", idKey: "targetHealthCheckId", writeIDOnly: true, seed: Object{"targetHealthCheckId": 46, "name": "before", "siteId": 3, "siteName": "site", "siteNiceId": "site-3", "hcEnabled": true, "hcHealth": "healthy", "hcMode": "http", "hcHostname": "example.com", "hcPort": 443, "hcPath": "/health", "hcScheme": "https", "hcMethod": "GET", "hcInterval": 30, "hcUnhealthyInterval": 30, "hcTimeout": 1, "hcHeaders": nil, "hcFollowRedirects": true, "hcStatus": nil, "hcTlsServerName": nil, "hcHealthyThreshold": 1, "hcUnhealthyThreshold": 1}}
	attrs := `org_id = "test-org"
name = %q
site_id = "3"
hc_enabled = true
hc_hostname = "example.com"
hc_port = 443
hc_path = "/health"
hc_scheme = "https"
hc_timeout = %d`
	runLifecycle(t, lifecycleCase{Resource: "health_check", CreateAttributes: fmt.Sprintf(attrs, "before", 1), UpdateAttributes: fmt.Sprintf(attrs, "after", 3), ImportID: "test-org/46", Handler: f.handler(t), Remove: f.remove, Exists: f.present,
		Checks: []resource.TestCheckFunc{resource.TestCheckResourceAttr("pangolin_health_check.test", "id", "46")}, UpdateChecks: []resource.TestCheckFunc{resource.TestCheckResourceAttr("pangolin_health_check.test", "name", "after"), resource.TestCheckResourceAttr("pangolin_health_check.test", "hc_timeout", "3")}})
}
