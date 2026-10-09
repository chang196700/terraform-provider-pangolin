package provider

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/chang196700/terraform-provider-pangolin/internal/client"
)

func TestAIProviderValidation(t *testing.T) {
	tests := []struct {
		name      string
		values    Values
		wantError bool
	}{
		{"public URL", Values{"type": "openai", "upstream_url": "https://api.openai.com/v1"}, false},
		{"missing URL rejected by upstream refinement", Values{"type": "openai"}, true},
		{"custom target", Values{"type": "custom", "routing_mode": "target", "capabilities": []any{"openai_chat"}}, false},
		{"target only for custom", Values{"type": "openai", "routing_mode": "target"}, true},
		{"target discards URL", Values{"type": "custom", "routing_mode": "target", "upstream_url": "https://example.com", "capabilities": []any{"openai_chat"}}, true},
		{"custom needs capability", Values{"type": "custom", "routing_mode": "target"}, true},
		{"empty capabilities", Values{"type": "openai", "upstream_url": "https://example.com", "capabilities": []any{}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateAIProvider(tt.values); (err != nil) != tt.wantError {
				t.Fatalf("validation error=%v, wantError=%v", err, tt.wantError)
			}
		})
	}
}

func TestAIEmptyCredentialsNormalizeAndClear(t *testing.T) {
	d := aiProviderDefinition()
	obj := Object{"headers": nil, "apiKey": nil}
	if err := d.EnrichRead(context.Background(), nil, Values{}, obj); err != nil {
		t.Fatal(err)
	}
	if obj["headers"] == nil || len(objectArray(obj["headers"])) != 0 || obj["apiKey"] != "" {
		t.Fatalf("empty credentials are not stable: %#v", obj)
	}
	v := Values{}
	applyDefaults(d.Fields, v)
	body, err := Body(d, v, false)
	if err != nil {
		t.Fatal(err)
	}
	if headers, present := body["headers"]; !present || headers == nil || len(objectArray(headers)) != 0 {
		t.Fatalf("removing headers must submit an empty collection: %#v", body)
	}
	for _, amount := range []float64{0, -1} {
		if err := validateAIBudget(Values{"amount": amount}); err == nil {
			t.Errorf("budget accepted nonpositive amount %v", amount)
		}
	}
	if err := validateAIBudget(Values{"amount": 0.01}); err != nil {
		t.Fatalf("fractional positive budget rejected: %v", err)
	}
}

func TestAIBudgetScopeUpdateClearsPreviousScope(t *testing.T) {
	var received Object
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" || r.URL.Path != "/v1/ai-budget/7" {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Error(err)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"error":false,"status":200,"data":{"budget":{"budgetId":7}}}`))
	}))
	defer srv.Close()
	c, err := client.New(srv.URL+"/v1", "test-key", time.Second)
	if err != nil {
		t.Fatal(err)
	}
	d := aiBudgetDefinition()
	v := Values{"org_id": "test-org", "id": "7", "model_id": "13", "amount": float64(20), "unit": "usd"}
	if err := d.UpdateFunc(context.Background(), c, v, Values{"provider_id": "9"}); err != nil {
		t.Fatal(err)
	}
	if received["modelId"] != float64(13) {
		t.Fatalf("model ID was not encoded as numeric: %#v", received)
	}
	for _, key := range []string{"providerId", "resourceId", "siteResourceId", "roleId", "virtualApiKeyId"} {
		if value, ok := received[key]; !ok || value != nil {
			t.Errorf("old scope %s must be explicitly null: %#v", key, received)
		}
	}
	if _, exists := received["orgId"]; exists {
		t.Fatal("org_id leaked into update body")
	}
	if err := validateAIBudget(Values{}); err != nil {
		t.Fatalf("organization scope: %v", err)
	}
	if err := validateAIBudget(Values{"provider_id": "9", "model_id": "13"}); err == nil {
		t.Fatal("multiple scopes accepted")
	}
}

func TestVirtualAPIKeyUpdateDoesNotRenewUnchangedExpiry(t *testing.T) {
	var bodies []Object
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" || r.URL.Path != "/v1/virtual-api-key/key1" {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		var body Object
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Error(err)
		}
		bodies = append(bodies, body)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"error":false,"status":200,"data":{"virtualApiKey":{"virtualApiKeyId":"key1"}}}`))
	}))
	defer srv.Close()
	c, err := client.New(srv.URL+"/v1", "test-key", time.Second)
	if err != nil {
		t.Fatal(err)
	}
	d := virtualAPIKeyDefinition()
	prior := Values{"id": "key1", "valid_for_seconds": int64(3600)}
	for _, lifetime := range []any{int64(3600), int64(7200), nil} {
		v := Values{"id": "key1", "name": "renamed", "all_resources": true, "resource_ids": []any{}, "valid_for_seconds": lifetime}
		if err := d.UpdateFunc(context.Background(), c, v, prior); err != nil {
			t.Fatal(err)
		}
	}
	if _, ok := bodies[0]["validForSeconds"]; ok {
		t.Fatal("renaming renewed the key expiry")
	}
	if bodies[1]["validForSeconds"] != float64(7200) {
		t.Fatalf("changed lifetime not submitted: %#v", bodies[1])
	}
	if value, ok := bodies[2]["validForSeconds"]; !ok || value != nil {
		t.Fatal("removing lifetime must send explicit null")
	}
	for _, body := range bodies {
		if _, ok := body["sendEmail"]; ok {
			t.Fatal("update can trigger an unintended key email")
		}
	}
}

func TestVirtualAPIKeyScopeValidation(t *testing.T) {
	for _, v := range []Values{{"all_resources": true}, {"resource_ids": []any{"12"}}} {
		if err := validateVirtualAPIKey(v); err != nil {
			t.Error(err)
		}
	}
	for _, v := range []Values{{}, {"all_resources": true, "resource_ids": []any{"12"}}} {
		if err := validateVirtualAPIKey(v); err == nil {
			t.Errorf("accepted incompatible key scope: %#v", v)
		}
	}
}

func TestAIResponseProjection(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"error":false,"status":200,"data":{"provider":{"providerId":4,"apiKey":"readable-secret","headers":[{"name":"x-key","value":"secret"}]}}}`))
	}))
	defer srv.Close()
	c, err := client.New(srv.URL+"/v1", "test-key", time.Second)
	if err != nil {
		t.Fatal(err)
	}
	d := aiProviderDefinition()
	obj, err := Call(context.Background(), c, d.Read, Values{"id": "4"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if Text(obj["providerId"]) != "4" || obj["apiKey"] != "readable-secret" {
		t.Fatalf("incorrect provider projection: %#v", obj)
	}
	if !d.Fields["api_key"].Sensitive || !d.Fields["headers"].Sensitive {
		t.Fatal("provider credentials must be sensitive")
	}
	if !virtualAPIKeyDefinition().Fields["secret"].Sensitive {
		t.Fatal("virtual key secret must be sensitive")
	}
}
