package provider

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/chang196700/terraform-provider-pangolin/internal/client"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func extraDefinition(t *testing.T, name string) DataDefinition {
	t.Helper()
	for _, d := range extraDataDefinitions() {
		if d.Name == name {
			return d
		}
	}
	t.Fatalf("missing extra data source %s", name)
	return DataDefinition{}
}

func TestExtraDataSourceSchemas(t *testing.T) {
	seen := map[string]bool{}
	for _, d := range extraDataDefinitions() {
		if seen[d.Name] {
			t.Errorf("duplicate data source %s", d.Name)
		}
		seen[d.Name] = true
		if len(d.Outputs) == 0 {
			t.Errorf("data source %s exposes no metadata", d.Name)
		}
		if diagnostics := dataSchema(d).ValidateImplementation(context.Background()); diagnostics.HasError() {
			t.Errorf("%s schema: %v", d.Name, diagnostics)
		}
	}
	if len(seen) != 12 {
		t.Fatalf("expected 12 extra data sources, got %d", len(seen))
	}
	if !extraDefinition(t, "blueprint").Outputs["contents"].Sensitive {
		t.Fatal("blueprint contents can contain secrets")
	}
	for _, name := range []string{"user", "users"} {
		fields := extraDefinition(t, name).Outputs
		for _, secret := range []string{"password_hash", "two_factor_secret", "token"} {
			if _, ok := fields[secret]; ok {
				t.Fatalf("%s exposes secret field %s", name, secret)
			}
		}
	}
}

func TestExtraAndDerivedDataSourceSecretProjections(t *testing.T) {
	protected := map[string][]string{
		"ai_provider":                 {"api_key", "headers"},
		"virtual_api_key":             {"secret"},
		"event_streaming_destination": {"config"},
		"alert_rule":                  {"webhook_actions"},
		"health_check":                {"hc_headers"},
		"blueprint":                   {"contents"},
	}
	for _, d := range dataDefinitions() {
		name := d.Name
		if d.IsList {
			name = strings.TrimSuffix(name, "s")
		}
		for _, key := range protected[name] {
			field, exists := d.Outputs[key]
			if d.IsList && exists {
				t.Errorf("%s list exposes %s instead of omitting credentials", d.Name, key)
			}
			if !d.IsList && (!exists || !field.Sensitive) {
				t.Errorf("%s single read must mark %s sensitive", d.Name, key)
			}
		}
	}
}

func TestDNSRecordsBareArrayResponse(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" || r.URL.Path != "/v1/org/test-org/domain/domain1/dns-records" {
			t.Errorf("unexpected DNS request: %s %s", r.Method, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"error":false,"status":200,"data":[{"id":1,"domainId":"domain1","recordType":"CNAME","baseDomain":"app","value":"example.com","verified":true}]}`))
	}))
	defer srv.Close()
	c, err := client.New(srv.URL+"/v1", "test-key", time.Second)
	if err != nil {
		t.Fatal(err)
	}
	d := extraDefinition(t, "domain_dns_records")
	o, err := d.ReadFunc(context.Background(), c, Values{"org_id": "test-org", "domain_id": "domain1"})
	if err != nil {
		t.Fatal(err)
	}
	rows := objectArray(o[d.ListKey])
	if len(rows) != 1 || Text(rows[0]["id"]) != "1" || rows[0]["verified"] != true {
		t.Fatalf("DNS array incorrectly decoded: %#v", o)
	}
}

func TestResourceNamesBareArrayResponse(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" || r.URL.Path != "/v1/org/test-org/resource-names" {
			t.Errorf("unexpected resource names request: %s %s", r.Method, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"error":false,"status":200,"data":[{"resourceId":7,"name":"Portal"}]}`))
	}))
	defer srv.Close()
	c, err := client.New(srv.URL+"/v1", "test-key", time.Second)
	if err != nil {
		t.Fatal(err)
	}
	d := extraDefinition(t, "resource_names")
	o, err := d.ReadFunc(context.Background(), c, Values{"org_id": "test-org"})
	if err != nil {
		t.Fatal(err)
	}
	rows := objectArray(o[d.ListKey])
	if len(rows) != 1 || Text(rows[0]["resourceId"]) != "7" || Text(rows[0]["name"]) != "Portal" {
		t.Fatalf("resource names incorrectly decoded: %#v", o)
	}
}

func TestCatalogModelLookupAndQueryEscaping(t *testing.T) {
	requests := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		if r.URL.Query().Get("query") != "a&b / model?" {
			t.Errorf("catalog query was not preserved: %s", r.URL.RawQuery)
		}
		if requests == 1 && r.URL.Path != "/v1/ai-provider/8/catalog-models" {
			t.Errorf("wrong provider route: %s", r.URL.Path)
		}
		if requests == 2 && (r.URL.Path != "/v1/org/test-org/ai-catalog-models" || r.URL.Query().Get("type") != "openai") {
			t.Errorf("wrong provider-type route: %s", r.URL)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"error":false,"status":200,"data":{"models":[{"model":"example-model"}]}}`))
	}))
	defer srv.Close()
	c, err := client.New(srv.URL+"/v1", "test-key", time.Second)
	if err != nil {
		t.Fatal(err)
	}
	d := extraDefinition(t, "ai_catalog_models")
	for _, input := range []Values{{"provider_id": "8", "query": "a&b / model?"}, {"org_id": "test-org", "type": "openai", "query": "a&b / model?"}} {
		o, err := d.ReadFunc(context.Background(), c, input)
		if err != nil {
			t.Fatal(err)
		}
		if rows := objectArray(o["models"]); len(rows) != 1 || rows[0]["model"] != "example-model" {
			t.Fatalf("unexpected catalog: %#v", o)
		}
	}
	for _, input := range []Values{{}, {"provider_id": "8", "org_id": "test-org"}, {"org_id": "test-org"}} {
		if _, err := d.ReadFunc(context.Background(), c, input); err == nil {
			t.Fatalf("ambiguous lookup accepted: %#v", input)
		}
	}
	if requests != 2 {
		t.Fatalf("invalid lookups sent HTTP requests: %d", requests)
	}
}

func TestUserDevicesFilterAndPagination(t *testing.T) {
	requests := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		q := r.URL.Query()
		if q.Get("page") != fmt.Sprint(requests) || q.Get("status") != "active,blocked" || q.Get("online") != "false" || q.Get("query") != "desk & laptop" {
			t.Errorf("filters lost across pages: %s", r.URL)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprintf(w, `{"success":true,"error":false,"status":200,"data":{"devices":[{"clientId":%d,"name":"device"}],"pagination":{"total":2,"page":%d,"pageSize":1}}}`, requests, requests)
	}))
	defer srv.Close()
	c, err := client.New(srv.URL+"/v1", "test-key", time.Second)
	if err != nil {
		t.Fatal(err)
	}
	d := extraDefinition(t, "user_devices")
	o, err := d.ReadFunc(context.Background(), c, Values{"org_id": "test-org", "status": []any{"active", "blocked"}, "online": false, "query": "desk & laptop"})
	if err != nil {
		t.Fatal(err)
	}
	if len(objectArray(o["devices"])) != 2 || requests != 2 {
		t.Fatalf("not all device pages read: %#v", o)
	}
}

// The responses below are independent fixtures of the server's distinct shapes:
// direct objects, nested lists, bare DNS arrays, and the certificate wildcard.
func TestAccExtraDataSources(t *testing.T) {
	responses := map[string]string{
		"/v1/org/test-org/domain/domain1":                      `{"domainId":"domain1","baseDomain":"example.com","type":"ns","verified":true,"configManaged":false,"lastCheckedAt":1000}`,
		"/v1/org/test-org/domains":                             `{"domains":[{"domainId":"domain1","baseDomain":"example.com","type":"ns","verified":true}],"pagination":{"total":1,"limit":100,"offset":0}}`,
		"/v1/org/test-org/resource-names":                      `[{"resourceId":7,"name":"Portal"}]`,
		"/v1/org/test-org/domain/domain1/dns-records":          `[{"id":2,"domainId":"domain1","recordType":"CNAME","baseDomain":"app","value":"example.com","verified":true}]`,
		"/v1/org/test-org/certificate/domain1/app.example.com": `{"certId":3,"domain":"*.example.com","domainType":"wildcard","wildcard":true,"status":"valid","expiresAt":2000,"renewalCount":1}`,
		"/v1/org/test-org/blueprint/4":                         `{"blueprintId":4,"orgId":"test-org","name":"example","source":"api","succeeded":true,"createdAt":1000,"message":"applied","contents":"name: example"}`,
		"/v1/org/test-org/blueprints":                          `{"blueprints":[{"blueprintId":4,"orgId":"test-org","name":"example","source":"api","succeeded":true,"createdAt":1000}],"pagination":{"total":1,"limit":100,"offset":0}}`,
		"/v1/org/test-org/invitations":                         `{"invitations":[{"inviteId":"invite1","email":"invited@example.com","expiresAt":2000,"roles":[{"roleId":5,"roleName":"Member"}]}],"pagination":{"total":1,"limit":100,"offset":0}}`,
		"/v1/org/test-org/user-devices":                        `{"devices":[{"clientId":6,"orgId":"test-org","name":"Laptop","online":false,"megabytesIn":1.5,"fingerprintArch":"amd64"}],"pagination":{"total":1,"page":1,"pageSize":100}}`,
		"/v1/ai-provider/7/catalog-models":                     `{"models":[{"model":"example-model"}]}`,
		"/v1/user/user1":                                       `{"userId":"user1","email":"user@example.com","username":"example","dateCreated":"2026-01-01","type":"internal","serverAdmin":false,"emailVerified":true}`,
		"/v1/users":                                            `{"users":[{"id":"user1","email":"user@example.com","username":"example","dateCreated":"2026-01-01","type":"internal","serverAdmin":false}],"pagination":{"total":1,"page":1,"pageSize":100}}`,
	}
	server := testHTTPServer(t, func(w http.ResponseWriter, r *http.Request) {
		body, ok := responses[r.URL.Path]
		if r.Method != "GET" || !ok {
			t.Errorf("unexpected extra data source request: %s %s", r.Method, r.URL)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if r.URL.Path == "/v1/org/test-org/user-devices" && r.URL.Query().Get("online") != "false" {
			t.Error("device filter lost by data source runtime")
		}
		if r.URL.Path == "/v1/users" && r.URL.Query().Get("idp_id") != "internal" {
			t.Error("user filter lost by data source runtime")
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprintf(w, `{"success":true,"error":false,"status":200,"data":%s}`, body)
	})
	config := fmt.Sprintf(`
provider "pangolin" {
  endpoint = %q
  api_key = "lifecycle-fixture-key"
}
data "pangolin_domain" "test" {
  org_id = "test-org"
  domain_id = "domain1"
}
data "pangolin_domains" "test" { org_id = "test-org" }
data "pangolin_resource_names" "test" { org_id = "test-org" }
data "pangolin_domain_dns_records" "test" {
  org_id = "test-org"
  domain_id = "domain1"
}
data "pangolin_certificate" "test" {
  org_id = "test-org"
  domain_id = "domain1"
  domain = "app.example.com"
}
data "pangolin_blueprint" "test" {
  org_id = "test-org"
  blueprint_id = "4"
}
data "pangolin_blueprints" "test" { org_id = "test-org" }
data "pangolin_invitations" "test" { org_id = "test-org" }
data "pangolin_user_devices" "test" {
  org_id = "test-org"
  online = false
}
data "pangolin_ai_catalog_models" "test" { provider_id = "7" }
data "pangolin_user" "test" { user_id = "user1" }
data "pangolin_users" "test" { idp_id = "internal" }
`, server.URL+"/v1")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testProtoV6ProviderFactories(),
		Steps: []resource.TestStep{{Config: config, Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr("data.pangolin_domain.test", "base_domain", "example.com"),
			resource.TestCheckResourceAttr("data.pangolin_domains.test", "items.0.domain_id", "domain1"),
			resource.TestCheckResourceAttr("data.pangolin_resource_names.test", "items.0.resource_id", "7"),
			resource.TestCheckResourceAttr("data.pangolin_domain_dns_records.test", "items.0.id", "2"),
			resource.TestCheckResourceAttr("data.pangolin_certificate.test", "domain", "app.example.com"),
			resource.TestCheckResourceAttr("data.pangolin_certificate.test", "certificate_domain", "*.example.com"),
			resource.TestCheckResourceAttr("data.pangolin_blueprint.test", "contents", "name: example"),
			resource.TestCheckResourceAttr("data.pangolin_blueprints.test", "items.0.blueprint_id", "4"),
			resource.TestCheckResourceAttr("data.pangolin_invitations.test", "items.0.roles.0.role_id", "5"),
			resource.TestCheckResourceAttr("data.pangolin_user_devices.test", "items.0.client_id", "6"),
			resource.TestCheckResourceAttr("data.pangolin_user_devices.test", "items.0.megabytes_in", "1.5"),
			resource.TestCheckResourceAttr("data.pangolin_ai_catalog_models.test", "items.0.model", "example-model"),
			resource.TestCheckResourceAttr("data.pangolin_user.test", "user_id", "user1"),
			resource.TestCheckResourceAttr("data.pangolin_users.test", "items.0.user_id", "user1"),
		)}},
	})
}
