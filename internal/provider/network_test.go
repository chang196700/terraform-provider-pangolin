package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/chang196700/terraform-provider-pangolin/internal/client"
)

type networkTestStep struct {
	method, path string
	data         any
	check        func(*testing.T, *http.Request, Object)
}

func networkTestClient(t *testing.T, steps ...networkTestStep) *client.Client {
	t.Helper()
	var mu sync.Mutex
	next := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		if next >= len(steps) {
			t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
			http.Error(w, "unexpected", 500)
			return
		}
		step := steps[next]
		next++
		if r.Method != step.method || r.URL.Path != "/v1"+step.path {
			t.Errorf("request %d = %s %s, wanted %s %s", next, r.Method, r.URL.Path, step.method, step.path)
			http.Error(w, "unexpected", 500)
			return
		}
		var body Object
		if r.Body != nil && r.ContentLength > 0 {
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Error(err)
			}
		}
		if step.check != nil {
			step.check(t, r, body)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(Object{"success": true, "error": false, "status": 200, "data": step.data})
	}))
	t.Cleanup(func() {
		srv.Close()
		mu.Lock()
		defer mu.Unlock()
		if next != len(steps) {
			t.Errorf("performed %d of %d expected requests", next, len(steps))
		}
	})
	c, err := client.New(srv.URL+"/v1", "test-key", time.Second)
	if err != nil {
		t.Fatal(err)
	}
	return c
}

func TestNetworkModeValidation(t *testing.T) {
	for _, tt := range []struct {
		name string
		d    Definition
		v    Values
		fail bool
	}{
		{"HTTP", networkPublicResource(), Values{"mode": "http", "domain_id": "domain"}, false},
		{"TCP", networkPublicResource(), Values{"mode": "tcp", "proxy_port": int64(8888)}, false},
		{"HTTP missing domain", networkPublicResource(), Values{"mode": "http"}, true},
		{"TCP with domain", networkPublicResource(), Values{"mode": "tcp", "proxy_port": int64(8888), "domain_id": "domain"}, true},
		{"shared policy ownership", networkPublicResource(), Values{"mode": "http", "domain_id": "domain", "resource_policy_id": "7", "access": Values{"sso": true}}, true},
		{"Newt generated settings", networkSite(), Values{"type": "newt"}, false},
		{"Newt ignores exit node", networkSite(), Values{"type": "newt", "exit_node_id": "2"}, true},
		{"WireGuard missing keys", networkSite(), Values{"type": "wireguard"}, true},
		{"WireGuard", networkSite(), Values{"type": "wireguard", "exit_node_id": "2", "subnet": "10.0.0.0/24", "pub_key": "public-key"}, false},
		{"target health hostname", networkTarget(), Values{"resource_id": "11", "hc_enabled": true}, true},
		{"target health configured", networkTarget(), Values{"resource_id": "11", "hc_enabled": true, "hc_hostname": "example.internal"}, false},
		{"target AI parent", networkTarget(), Values{"ai_provider_id": "11"}, false},
		{"target both parents", networkTarget(), Values{"resource_id": "11", "ai_provider_id": "12"}, true},
		{"target missing parent", networkTarget(), Values{}, true},
		{"private hostname alias", networkPrivateResource(), Values{"mode": "host", "site_ids": []any{"1"}, "destination": "internal-service"}, true},
		{"private network", networkPrivateResource(), Values{"mode": "cidr", "site_ids": []any{"1"}, "destination": "10.2.0.0/16"}, false},
		{"private SSH derived port", networkPrivateResource(), Values{"mode": "ssh", "site_ids": []any{"1"}, "destination": "192.0.2.1", "destination_port": int64(2222), "tcp_port_range_string": "2222"}, false},
		{"private SSH conflicting port", networkPrivateResource(), Values{"mode": "ssh", "site_ids": []any{"1"}, "destination": "192.0.2.1", "destination_port": int64(2222), "tcp_port_range_string": "22"}, true},
		{"private gateway conflicting destination", networkPrivateResource(), Values{"mode": "gateway", "site_ids": []any{"1"}, "destination": "192.0.2.1"}, true},
		{"private inference conflicting ICMP", networkPrivateResource(), Values{"mode": "inference", "disable_icmp": false}, true},
		{"private HTTP missing port", networkPrivateResource(), Values{"mode": "http", "site_ids": []any{"1"}, "destination": "10.0.0.1", "domain_id": "domain"}, true},
		{"native SSH one site", networkPrivateResource(), Values{"mode": "ssh", "site_ids": []any{"1", "2"}, "auth_daemon_mode": "native"}, true},
		{"overlapping organization subnets", networkOrganization(), Values{"subnet": "10.0.0.0/16", "utility_subnet": "10.0.1.0/24"}, true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.d.Validate(tt.v); (err != nil) != tt.fail {
				t.Fatalf("validation = %v, wanted failure %v", err, tt.fail)
			}
		})
	}
	public := networkPublicResource()
	if err := validateFields(public.Fields, Values{"mode": "http"}); err != nil {
		t.Fatalf("anyOf must retain the HTTP mode: %v", err)
	}
}

func TestNetworkPolicyCreateAndUpdateIDWireTypes(t *testing.T) {
	d := networkPolicy()
	v := Values{"org_id": "org", "name": "Test", "role_ids": []any{"4"}, "user_ids": []any{}, "sso": true, "rules": []any{}, "emails": []any{}}
	c := networkTestClient(t,
		networkTestStep{"POST", "/org/org/resource-policy", Object{"resourcePolicyId": 9}, func(t *testing.T, _ *http.Request, b Object) {
			if !reflect.DeepEqual(b["roleIds"], []any{"4"}) {
				t.Errorf("create role IDs must be strings: %#v", b["roleIds"])
			}
		}},
		networkTestStep{"POST", "/public-resource-policy/9/access-control", Object{}, func(t *testing.T, _ *http.Request, b Object) {
			if !reflect.DeepEqual(b["roleIds"], []any{float64(4)}) {
				t.Errorf("update role IDs must be numeric: %#v", b["roleIds"])
			}
			if x, ok := b["skipToIdpId"]; !ok || x != nil {
				t.Error("cleared identity-provider selection must be explicit null")
			}
		}},
	)
	created, err := d.CreateFunc(context.Background(), c, v)
	if err != nil {
		t.Fatal(err)
	}
	v["id"] = Text(created["resourcePolicyId"])
	prior := copyValues(v)
	prior["role_ids"] = []any{"3"}
	if err = d.AfterWrite(context.Background(), c, v, prior); err != nil {
		t.Fatal(err)
	}
}

func TestNetworkSharedPolicyChangeIsIsolatedAndDoesNotWriteSharedAccess(t *testing.T) {
	d := networkPublicResource()
	prior := Values{"id": "8", "mode": "http", "domain_id": "domain", "resource_policy_id": nil, "name": "Old", "access": Values{"sso": true}}
	plan := Values{"id": "8", "mode": "http", "domain_id": "domain", "resource_policy_id": "9", "name": "New", "access": nil}
	c := networkTestClient(t,
		networkTestStep{"POST", "/public-resource/8", Object{}, func(t *testing.T, _ *http.Request, b Object) {
			if len(b) != 1 || b["resourcePolicyId"] != float64(9) {
				t.Errorf("policy switch mixed with settings: %#v", b)
			}
		}},
		networkTestStep{"POST", "/public-resource/8", []any{Object{"resourceId": 8}}, func(t *testing.T, _ *http.Request, b Object) {
			if len(b) != 1 || b["name"] != "New" {
				t.Errorf("unexpected settings: %#v", b)
			}
		}},
	)
	if err := d.UpdateFunc(context.Background(), c, plan, prior); err != nil {
		t.Fatal(err)
	}
	if err := d.AfterWrite(context.Background(), c, plan, prior); err != nil {
		t.Fatal(err)
	}
}

func TestNetworkPolicyReadDoesNotExposeHashesAndDetectsRemovedSecret(t *testing.T) {
	d := networkPolicy()
	prior := Values{"password": "old-value", "header_auth": Values{"user": "test", "password": "secret", "extended_compatibility": false}}
	o := Object{"passwordId": nil, "pincodeId": 12, "headerAuth": Object{"id": 4, "extendedCompability": false}, "roles": []any{Object{"roleId": 3, "isAdmin": false}}, "users": []any{}, "emailWhiteList": []any{}, "rules": []any{}}
	if err := d.EnrichRead(context.Background(), nil, prior, o); err != nil {
		t.Fatal(err)
	}
	if err := mergeRemote(d, prior, o); err != nil {
		t.Fatal(err)
	}
	if prior["password"] != nil {
		t.Fatal("deleted password was not detected")
	}
	if prior["header_auth"].(map[string]any)["password"] != "secret" {
		t.Fatal("write-only header credentials were overwritten by hash metadata")
	}
	if !reflect.DeepEqual(prior["role_ids"], []any{"3"}) {
		t.Fatalf("roles not flattened: %#v", prior["role_ids"])
	}
}

func TestNetworkPublicAccessReadUsesInlinePolicyAndPreservesKnownSecrets(t *testing.T) {
	v := Values{"id": "8", "access": Values{"password": "known-password", "sso": true}}
	c := networkTestClient(t, networkTestStep{"GET", "/public-resource/8/policies", Object{"defaultPolicy": Object{"sso": false, "idpId": 4, "passwordId": 11, "pincodeId": nil, "roles": []any{Object{"roleId": 1, "isAdmin": true}, Object{"roleId": 2}}, "users": []any{Object{"userId": "u"}}, "emailWhiteList": []any{}, "rules": []any{}}}, nil})
	o := Object{"resourceId": 8, "resourcePolicyId": nil, "mode": "http"}
	d := networkPublicResource()
	if err := d.EnrichRead(context.Background(), c, v, o); err != nil {
		t.Fatal(err)
	}
	if err := mergeRemote(d, v, o); err != nil {
		t.Fatal(err)
	}
	a := v["access"].(map[string]any)
	if a["password"] != "known-password" || a["sso"] != false || a["skip_to_idp_id"] != "4" {
		t.Fatalf("incorrect inline policy refresh: %#v", a)
	}
	if !reflect.DeepEqual(a["role_ids"], []any{"2"}) {
		t.Fatalf("implicit admin leaked into managed roles: %#v", a["role_ids"])
	}
}

func TestNetworkPrivateReadWalksPagesToRecoverSites(t *testing.T) {
	d := networkPrivateResource()
	v := Values{"id": "8", "org_id": "org"}
	o := Object{"siteResourceId": 8, "mode": "host"}
	c := networkTestClient(t,
		networkTestStep{"GET", "/private-resource/8/roles", Object{"roles": []any{Object{"roleId": 1, "isAdmin": true}, Object{"roleId": 3}}}, nil},
		networkTestStep{"GET", "/private-resource/8/users", Object{"users": []any{Object{"userId": "user"}}}, nil},
		networkTestStep{"GET", "/private-resource/8/clients", Object{"clients": []any{}}, nil},
		networkTestStep{"GET", "/org/org/private-resources", Object{"siteResources": []any{Object{"siteResourceId": 7}}, "pagination": Object{"total": 2}}, func(t *testing.T, r *http.Request, _ Object) {
			if r.URL.Query().Get("page") != "1" {
				t.Error("expected first page")
			}
		}},
		networkTestStep{"GET", "/org/org/private-resources", Object{"siteResources": []any{Object{"siteResourceId": 8, "siteIds": []any{2, 3}}}, "pagination": Object{"total": 2}}, func(t *testing.T, r *http.Request, _ Object) {
			if r.URL.Query().Get("page") != "2" {
				t.Error("expected second page")
			}
		}},
	)
	if err := d.EnrichRead(context.Background(), c, v, o); err != nil {
		t.Fatal(err)
	}
	if err := mergeRemote(d, v, o); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(v["site_ids"], []any{"2", "3"}) || !reflect.DeepEqual(v["role_ids"], []any{"3"}) {
		t.Fatalf("relationship drift not reflected: %#v", v)
	}
}

func TestNetworkInferenceAssociationsOrderAndTypes(t *testing.T) {
	d := networkPublicResource()
	v := Values{"id": "8", "mode": "inference", "ai_providers": []any{Values{"provider_id": "5", "access_mode": "select", "enabled": true}}, "ai_models": []any{Values{"model_id": "6", "list_type": "allow"}}}
	c := networkTestClient(t,
		networkTestStep{"POST", "/resource/8/ai-providers", Object{}, func(t *testing.T, _ *http.Request, b Object) {
			p := networkRows(b["providers"])
			if len(p) != 1 || p[0]["providerId"] != float64(5) || p[0]["accessMode"] != "select" {
				t.Errorf("invalid provider payload: %#v", b)
			}
		}},
		networkTestStep{"POST", "/resource/8/ai-models", Object{}, func(t *testing.T, _ *http.Request, b Object) {
			m := networkRows(b["models"])
			if len(m) != 1 || m[0]["modelId"] != float64(6) {
				t.Errorf("invalid model payload: %#v", b)
			}
		}},
	)
	if err := networkWriteAI(context.Background(), c, d, v, nil, false); err != nil {
		t.Fatal(err)
	}
	v["ai_providers"] = []any{Values{"provider_id": "5", "access_mode": "inherit", "enabled": true}}
	if err := networkValidateAI(v); err == nil {
		t.Fatal("model list with no select provider accepted")
	}
}

func TestNetworkRuleReconciliationRetainsIDsAndAcceptsNullDeleteResponse(t *testing.T) {
	current := []any{
		Object{"ruleId": 1, "action": "ACCEPT", "match": "IP", "value": "192.0.2.1", "priority": 1, "enabled": true},
		Object{"ruleId": 2, "action": "DROP", "match": "IP", "value": "192.0.2.2", "priority": 2, "enabled": true},
		Object{"ruleId": 3, "action": "DROP", "match": "IP", "value": "192.0.2.3", "priority": 3, "enabled": true},
	}
	desired := []any{
		Values{"action": "ACCEPT", "match": "IP", "value": "192.0.2.1", "priority": int64(1), "enabled": true},
		Values{"action": "DROP", "match": "IP", "value": "192.0.2.22", "priority": int64(2), "enabled": true},
		Values{"action": "PASS", "match": "IP", "value": "192.0.2.4", "priority": int64(4), "enabled": true},
	}
	c := networkTestClient(t,
		networkTestStep{"GET", "/public-resource/8/rules", Object{"rules": current, "pagination": Object{"total": 3}}, nil},
		networkTestStep{"POST", "/public-resource/8/rule/2", Object{}, nil},
		networkTestStep{"PUT", "/public-resource/8/rule", Object{}, nil},
		networkTestStep{"DELETE", "/public-resource/8/rule/3", nil, nil},
	)
	if err := networkReconcilePublicRules(context.Background(), c, Values{"id": "8"}, desired); err != nil {
		t.Fatal(err)
	}
}

func TestNetworkNewtSecretIsPreservedAndAddressNormalized(t *testing.T) {
	d := networkSite()
	v := Values{"id": "4", "type": "newt", "secret": "create-only-secret"}
	o := Object{"siteId": 4, "type": "newt", "address": "100.64.0.7/24", "newtId": "agent", "exitNodeId": 2, "exitNodeSubnet": "100.10.0.0/24"}
	if err := d.EnrichRead(context.Background(), nil, v, o); err != nil {
		t.Fatal(err)
	}
	if err := mergeRemote(d, v, o); err != nil {
		t.Fatal(err)
	}
	if v["secret"] != "create-only-secret" || v["address"] != "100.64.0.7" || v["exit_node_id"] != nil {
		t.Fatalf("wrong Newt state: %#v", v)
	}
	if !d.Fields["secret"].Sensitive {
		t.Fatal("Newt secret must be sensitive")
	}
}

func TestNetworkPolicyMissingSecretCanBeExplicitlyCleared(t *testing.T) {
	d := networkPolicy()
	body, err := networkPayload(d, Values{"password": nil}, "password")
	if err != nil {
		t.Fatal(err)
	}
	if value, ok := body["password"]; !ok || value != nil {
		t.Fatalf("must send null, received %s", fmt.Sprint(body))
	}
}

func TestNetworkTargetImportsAndHealthCheckBoundaries(t *testing.T) {
	d := networkTarget()
	for _, tc := range []struct {
		id, parent string
		valid      bool
	}{{"public/11/42", "resource_id", true}, {"ai/11/42", "ai_provider_id", true}, {"11/42", "", false}, {"public/no/42", "", false}, {"ai/11/0", "", false}, {"unknown/11/42", "", false}} {
		v, err := d.ImportFunc(tc.id)
		if (err == nil) != tc.valid {
			t.Errorf("import %q: %v", tc.id, err)
		}
		if tc.valid && (v[tc.parent] != "11" || v["id"] != "42") {
			t.Errorf("import %q produced %#v", tc.id, v)
		}
	}
	for _, key := range []string{"hc_interval", "hc_unhealthy_interval", "hc_timeout", "hc_healthy_threshold", "hc_unhealthy_threshold"} {
		if err := validateFields(d.Fields, Values{key: int64(1)}); err != nil {
			t.Errorf("%s must accept 1: %v", key, err)
		}
		if err := validateFields(d.Fields, Values{key: int64(0)}); err == nil {
			t.Errorf("%s must reject 0", key)
		}
	}
}

func TestNetworkReadRelationshipsRejectMalformedResponses(t *testing.T) {
	c := networkTestClient(t, networkTestStep{"GET", "/private-resource/8/roles", Object{"roles": []any{"unexpected"}}, nil})
	if _, err := networkGetRows(context.Background(), c, "/private-resource/{id}/roles", Values{"id": "8"}, "roles"); err == nil {
		t.Fatal("malformed authorization response must fail instead of clearing state")
	}
}

func TestNetworkManagedImportDoesNotReadOrAdoptAccess(t *testing.T) {
	d := networkPublicResource()
	c := networkTestClient(t)
	v := Values{"org_id": "test", "id": "42"}
	o := Object{"resourceId": 42, "mode": "http", "resourcePolicyId": nil, "defaultResourcePolicyId": 99}
	if err := d.EnrichRead(context.Background(), c, v, o); err != nil {
		t.Fatal(err)
	}
	if err := mergeRemote(d, v, o); err != nil {
		t.Fatal(err)
	}
	if v["access"] != nil {
		t.Fatal("managed import must not implicitly take ownership of inline authorization")
	}
}
