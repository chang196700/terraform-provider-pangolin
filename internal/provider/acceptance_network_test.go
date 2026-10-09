package provider

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// These fixtures model the documented wire contract, including response fields
// deliberately absent from the OpenAPI and the different policy-role ID types.
// They exercise Terraform's plan/state consistency checks, not just Go hooks.
type networkLifecycleFixture struct {
	mu                           sync.Mutex
	t                            *testing.T
	kind                         string
	aiTarget                     bool
	present                      bool
	state                        Object
	policy                       Object
	roles, users, clients, sites []any
	rules                        []any
}

func newNetworkLifecycleFixture(t *testing.T, kind string) *networkLifecycleFixture {
	return &networkLifecycleFixture{t: t, kind: kind}
}
func (f *networkLifecycleFixture) exists() bool { f.mu.Lock(); defer f.mu.Unlock(); return f.present }
func (f *networkLifecycleFixture) remove()      { f.mu.Lock(); defer f.mu.Unlock(); f.present = false }
func (f *networkLifecycleFixture) reply(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(Object{"success": true, "error": false, "status": 200, "data": data})
}
func (f *networkLifecycleFixture) fail(w http.ResponseWriter, r *http.Request, msg string) {
	f.t.Errorf("%s fixture: %s %s: %s", f.kind, r.Method, r.URL.Path, msg)
	http.Error(w, msg, 400)
}

func (f *networkLifecycleFixture) handler(w http.ResponseWriter, r *http.Request) {
	f.mu.Lock()
	defer f.mu.Unlock()
	p := strings.TrimPrefix(r.URL.Path, "/v1")
	var b Object
	if r.Body != nil && r.ContentLength > 0 {
		if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
			f.fail(w, r, "bad JSON")
			return
		}
	}
	createPath, idPath, idKey := f.paths()
	if r.Method == "PUT" && p == createPath || f.kind == "public_resource_policy" && r.Method == "POST" && p == createPath {
		if f.present {
			f.fail(w, r, "object already exists")
			return
		}
		if f.kind == "resource_target" && b["hcEnabled"] == true && b["hcHostname"] == nil {
			f.fail(w, r, "health checks require hcHostname at creation")
			return
		}
		f.present = true
		f.state = Object{idKey: 42, "orgId": "test"}
		if f.kind == "organization" {
			f.state[idKey] = "test"
		}
		for k, v := range b {
			f.state[k] = v
		}
		f.roles = []any{}
		f.users = []any{}
		f.clients = []any{}
		f.sites = []any{}
		f.rules = []any{}
		f.policy = Object{"sso": true, "applyRules": false, "emailWhitelistEnabled": false, "idpId": nil, "passwordId": nil, "pincodeId": nil, "headerAuth": nil, "roles": []any{}, "users": []any{}, "emailWhiteList": []any{}, "rules": []any{}}
		switch f.kind {
		case "resource_target":
			if f.aiTarget {
				f.state["providerId"] = 11
				f.state["resourceId"] = nil
			} else {
				f.state["resourceId"] = 11
				f.state["providerId"] = nil
			}
		case "organization":
			f.state["requireTwoFactor"] = false
		case "site":
			if f.state["niceId"] == nil {
				f.state["niceId"] = "generated-site"
			}
			f.state["dockerSocketEnabled"] = false
			f.state["exitNodeSubnet"] = "0.0.0.0/32"
			if f.state["type"] == "newt" {
				f.state["newtId"] = "generated-newt-id"
				f.state["secret"] = "generated-newt-secret"
				f.state["address"] = "100.64.0.7/24"
				f.state["exitNodeSubnet"] = nil
			}
		case "public_resource":
			f.state["niceId"] = "generated-public"
			f.state["resourcePolicyId"] = nil
			f.state["defaultResourcePolicyId"] = 99
			f.state["ssl"] = false
		case "private_resource":
			if f.state["niceId"] == nil {
				f.state["niceId"] = "generated-private"
			}
			f.capturePrivate(b)
			if f.state["mode"] == "ssh" {
				f.state["tcpPortRangeString"] = Text(f.state["destinationPort"])
				f.state["udpPortRangeString"] = ""
				f.state["disableIcmp"] = true
			}
		case "public_resource_policy":
			f.state["niceId"] = "generated-policy"
			f.capturePolicy(b)
		}
		f.reply(w, f.state)
		return
	}
	if !f.present {
		http.Error(w, "missing", 404)
		return
	}
	if r.Method == "DELETE" && (p == idPath || f.kind == "public_resource_policy" && p == "/resource-policy/42") {
		f.present = false
		f.reply(w, nil)
		return
	}
	if r.Method == "GET" && p == idPath {
		o := Object{}
		for k, v := range f.state {
			o[k] = v
		}
		switch f.kind {
		case "organization":
			f.reply(w, Object{"org": o})
			return
		case "site":
			delete(o, "secret")
		case "private_resource":
			for _, key := range []string{"siteIds", "roleIds", "userIds", "clientIds"} {
				delete(o, key)
			}
		case "public_resource_policy":
			for _, key := range []string{"roleIds", "userIds", "emails", "password", "pincode", "headerAuth"} {
				delete(o, key)
			}
			for k, v := range f.policy {
				o[k] = v
			}
		}
		f.reply(w, o)
		return
	}
	if r.Method == "POST" && p == idPath {
		if f.kind == "resource_target" && (b["siteId"] == nil || b["ip"] == nil) {
			f.fail(w, r, "target update requires siteId and ip")
			return
		}
		for k, v := range b {
			f.state[k] = v
		}
		if f.kind == "private_resource" {
			f.capturePrivate(b)
		}
		if f.kind == "public_resource" {
			if f.state["resourcePolicyId"] != nil {
				for _, key := range []string{"sso", "emailWhitelistEnabled", "applyRules", "skipToIdpId"} {
					if _, ok := b[key]; ok {
						f.fail(w, r, "inline access setter used with shared policy")
						return
					}
				}
			}
			for _, key := range []string{"sso", "emailWhitelistEnabled", "applyRules"} {
				if v, ok := b[key]; ok {
					f.policy[key] = v
				}
			}
			if v, ok := b["skipToIdpId"]; ok {
				f.policy["idpId"] = v
			}
		}
		f.reply(w, f.state)
		return
	}
	if f.kind == "private_resource" {
		if r.Method == "GET" && p == "/org/test/private-resources" {
			f.reply(w, Object{"siteResources": []any{Object{"siteResourceId": 42, "siteIds": f.sites}}, "pagination": Object{"total": 1}})
			return
		}
		for _, rel := range []struct {
			name, key string
			rows      []any
		}{{"roles", "roles", f.roles}, {"users", "users", f.users}, {"clients", "clients", f.clients}} {
			if r.Method == "GET" && p == idPath+"/"+rel.name {
				f.reply(w, Object{rel.key: rel.rows})
				return
			}
		}
	}
	if f.kind == "public_resource" {
		if f.state["resourcePolicyId"] != nil && r.Method != "GET" {
			f.fail(w, r, "inline access mutation used with shared policy")
			return
		}
		if r.Method == "GET" && p == idPath+"/policies" {
			f.reply(w, Object{"defaultPolicy": f.policy, "sharedPolicy": nil})
			return
		}
		if r.Method == "GET" && p == idPath+"/rules" {
			f.reply(w, Object{"rules": f.rules, "pagination": Object{"total": len(f.rules)}})
			return
		}
		if r.Method == "POST" {
			switch p {
			case idPath + "/roles":
				f.policy["roles"] = networkFixtureIDRows(b["roleIds"], "roleId")
				f.reply(w, Object{})
				return
			case idPath + "/users":
				f.policy["users"] = networkFixtureIDRows(b["userIds"], "userId")
				f.reply(w, Object{})
				return
			case idPath + "/password":
				if b["password"] == nil {
					f.policy["passwordId"] = nil
				} else {
					f.policy["passwordId"] = 1
				}
				f.reply(w, Object{})
				return
			case idPath + "/pincode":
				if b["pincode"] == nil {
					f.policy["pincodeId"] = nil
				} else {
					f.policy["pincodeId"] = 1
				}
				f.reply(w, Object{})
				return
			case idPath + "/header-auth":
				if b["user"] == nil {
					f.policy["headerAuth"] = nil
				} else {
					f.policy["headerAuth"] = Object{"id": 1}
				}
				f.reply(w, Object{})
				return
			case idPath + "/whitelist":
				if f.policy["emailWhitelistEnabled"] != true {
					f.fail(w, r, "whitelist must be enabled")
					return
				}
				f.policy["emailWhiteList"] = networkFixtureIDRows(b["emails"], "email")
				f.reply(w, Object{})
				return
			}
		}
	}
	if f.kind == "public_resource_policy" && r.Method == "POST" {
		switch p {
		case idPath + "/access-control":
			f.capturePolicy(b)
			f.reply(w, Object{})
			return
		case idPath + "/password", idPath + "/pincode", idPath + "/header-auth", idPath + "/whitelist", idPath + "/rules":
			f.capturePolicy(b)
			f.reply(w, Object{})
			return
		}
	}
	f.fail(w, r, "unimplemented contract route")
}

func (f *networkLifecycleFixture) paths() (string, string, string) {
	switch f.kind {
	case "organization":
		return "/org", "/org/test", "orgId"
	case "site":
		return "/org/test/site", "/site/42", "siteId"
	case "public_resource":
		return "/org/test/public-resource", "/public-resource/42", "resourceId"
	case "private_resource":
		return "/org/test/private-resource", "/private-resource/42", "siteResourceId"
	case "resource_target":
		if f.aiTarget {
			return "/ai-provider/11/target", "/target/42", "targetId"
		}
		return "/public-resource/11/target", "/target/42", "targetId"
	default:
		return "/org/test/resource-policy", "/public-resource-policy/42", "resourcePolicyId"
	}
}
func networkFixtureIDRows(v any, key string) []any {
	out := []any{}
	if a, ok := v.([]any); ok {
		for _, id := range a {
			out = append(out, Object{key: id})
		}
	}
	return out
}
func (f *networkLifecycleFixture) capturePrivate(b Object) {
	if v, ok := b["siteIds"].([]any); ok {
		f.sites = v
	}
	if v, ok := b["roleIds"]; ok {
		f.roles = networkFixtureIDRows(v, "roleId")
	}
	if v, ok := b["userIds"]; ok {
		f.users = networkFixtureIDRows(v, "userId")
	}
	if v, ok := b["clientIds"]; ok {
		f.clients = networkFixtureIDRows(v, "clientId")
	}
}
func (f *networkLifecycleFixture) capturePolicy(b Object) {
	for _, key := range []string{"sso", "applyRules", "emailWhitelistEnabled"} {
		if v, ok := b[key]; ok {
			f.policy[key] = v
		}
	}
	if v, ok := b["skipToIdpId"]; ok {
		f.policy["idpId"] = v
	}
	if v, ok := b["roleIds"]; ok {
		f.policy["roles"] = networkFixtureIDRows(v, "roleId")
	}
	if v, ok := b["userIds"]; ok {
		f.policy["users"] = networkFixtureIDRows(v, "userId")
	}
	if v, ok := b["emails"]; ok {
		f.policy["emailWhiteList"] = networkFixtureIDRows(v, "email")
	}
	if v, ok := b["password"]; ok {
		if v == nil {
			f.policy["passwordId"] = nil
		} else {
			f.policy["passwordId"] = 1
		}
	}
	if v, ok := b["pincode"]; ok {
		if v == nil {
			f.policy["pincodeId"] = nil
		} else {
			f.policy["pincodeId"] = 1
		}
	}
	if v, ok := b["headerAuth"]; ok {
		if v == nil {
			f.policy["headerAuth"] = nil
		} else {
			f.policy["headerAuth"] = Object{"id": 1}
		}
	}
	if v, ok := b["rules"]; ok {
		f.policy["rules"] = v
	}
}

func TestAccNetworkOrganization(t *testing.T) {
	f := newNetworkLifecycleFixture(t, "organization")
	runLifecycle(t, lifecycleCase{Resource: f.kind, CreateAttributes: `org_id = "test"
name = "Network"
subnet = "100.64.0.0/24"
utility_subnet = "100.65.0.0/24"`, UpdateAttributes: `org_id = "test"
name = "Network updated"
subnet = "100.64.0.0/24"
utility_subnet = "100.65.0.0/24"
require_two_factor = true`, ImportID: "test", Handler: f.handler, Remove: f.remove, Exists: f.exists,
		UpdateChecks: []resource.TestCheckFunc{resource.TestCheckResourceAttr("pangolin_organization.test", "require_two_factor", "true")}})
}

func TestAccNetworkSite(t *testing.T) {
	f := newNetworkLifecycleFixture(t, "site")
	runLifecycle(t, lifecycleCase{Resource: f.kind, CreateAttributes: `org_id = "test"
name = "Local site"
type = "local"
nice_id = "managed-site"`, UpdateAttributes: `org_id = "test"
name = "Renamed site"
type = "local"
nice_id = "managed-site"
docker_socket_enabled = true`, ImportID: "test/42", Handler: f.handler, Remove: f.remove, Exists: f.exists,
		Checks: []resource.TestCheckFunc{resource.TestCheckResourceAttr("pangolin_site.test", "nice_id", "managed-site")}})
}

func TestAccNetworkNewtGeneratedCredentials(t *testing.T) {
	f := newNetworkLifecycleFixture(t, "site")
	runLifecycle(t, lifecycleCase{Resource: f.kind, CreateAttributes: `org_id = "test"
name = "Newt site"
type = "newt"`, UpdateAttributes: `org_id = "test"
name = "Newt site updated"
type = "newt"`, ImportID: "test/42", ImportIgnore: []string{"secret"}, Handler: f.handler, Remove: f.remove, Exists: f.exists,
		Checks:       []resource.TestCheckFunc{resource.TestCheckResourceAttr("pangolin_site.test", "secret", "generated-newt-secret"), resource.TestCheckResourceAttr("pangolin_site.test", "address", "100.64.0.7")},
		UpdateChecks: []resource.TestCheckFunc{resource.TestCheckResourceAttr("pangolin_site.test", "secret", "generated-newt-secret")}})
}

func TestAccNetworkPublicResource(t *testing.T) {
	f := newNetworkLifecycleFixture(t, "public_resource")
	runLifecycle(t, lifecycleCase{Resource: f.kind, CreateAttributes: `org_id = "test"
name = "Web"
mode = "http"
domain_id = "example-domain"
nice_id = "managed-web"
access = { sso = true, role_ids = ["2"], password = "initial-password" }`, UpdateAttributes: `org_id = "test"
name = "Web updated"
mode = "http"
domain_id = "example-domain"
nice_id = "managed-web"
ssl = true
access = { sso = false, role_ids = ["3"], password = "updated-password" }`, ImportID: "test/42", ImportIgnore: []string{"access"}, Handler: f.handler, Remove: f.remove, Exists: f.exists,
		Checks:       []resource.TestCheckFunc{resource.TestCheckResourceAttr("pangolin_public_resource.test", "nice_id", "managed-web")},
		UpdateChecks: []resource.TestCheckFunc{resource.TestCheckResourceAttr("pangolin_public_resource.test", "ssl", "true")}})
}

func TestAccNetworkPublicSharedToInline(t *testing.T) {
	f := newNetworkLifecycleFixture(t, "public_resource")
	runLifecycle(t, lifecycleCase{Resource: f.kind, CreateAttributes: `org_id = "test"
name = "Shared web"
mode = "http"
domain_id = "example-domain"
resource_policy_id = "9"`, UpdateAttributes: `org_id = "test"
name = "Inline web"
mode = "http"
domain_id = "example-domain"
access = { sso = false, role_ids = ["2"] }`, ImportID: "test/42", ImportIgnore: []string{"access"}, Handler: f.handler, Remove: f.remove, Exists: f.exists,
		UpdateChecks: []resource.TestCheckFunc{resource.TestCheckNoResourceAttr("pangolin_public_resource.test", "resource_policy_id"), resource.TestCheckResourceAttr("pangolin_public_resource.test", "access.sso", "false")}})
}

func TestAccNetworkPublicInlineToShared(t *testing.T) {
	f := newNetworkLifecycleFixture(t, "public_resource")
	runLifecycle(t, lifecycleCase{Resource: f.kind, CreateAttributes: `org_id = "test"
name = "Inline web"
mode = "http"
domain_id = "example-domain"
access = { sso = false, role_ids = ["2"] }`, UpdateAttributes: `org_id = "test"
name = "Shared web"
mode = "http"
domain_id = "example-domain"
resource_policy_id = "9"`, ImportID: "test/42", Handler: f.handler, Remove: f.remove, Exists: f.exists,
		UpdateChecks: []resource.TestCheckFunc{resource.TestCheckResourceAttr("pangolin_public_resource.test", "resource_policy_id", "9"), resource.TestCheckNoResourceAttr("pangolin_public_resource.test", "access")}})
}

func TestAccNetworkPublicRemoveAccess(t *testing.T) {
	f := newNetworkLifecycleFixture(t, "public_resource")
	runLifecycle(t, lifecycleCase{Resource: f.kind, CreateAttributes: `org_id = "test"
name = "Web"
mode = "http"
domain_id = "example-domain"
access = { sso = false, role_ids = ["2"], password = "initial-password", email_whitelist_enabled = true, emails = ["user@example.com"] }`, UpdateAttributes: `org_id = "test"
name = "Web"
mode = "http"
domain_id = "example-domain"`, ImportID: "test/42", Handler: f.handler, Remove: f.remove, Exists: f.exists,
		UpdateChecks: []resource.TestCheckFunc{resource.TestCheckNoResourceAttr("pangolin_public_resource.test", "access")}})
}

func TestAccNetworkPolicyRemoveSecrets(t *testing.T) {
	f := newNetworkLifecycleFixture(t, "public_resource_policy")
	runLifecycle(t, lifecycleCase{Resource: f.kind, CreateAttributes: `org_id = "test"
name = "Restricted"
password = "initial-password"
pincode = "123456"
header_auth = { user = "test-user", password = "header-password", extended_compatibility = false }`, UpdateAttributes: `org_id = "test"
name = "Restricted"`, ImportID: "test/42", Handler: f.handler, Remove: f.remove, Exists: f.exists,
		UpdateChecks: []resource.TestCheckFunc{resource.TestCheckNoResourceAttr("pangolin_public_resource_policy.test", "password"), resource.TestCheckNoResourceAttr("pangolin_public_resource_policy.test", "pincode"), resource.TestCheckNoResourceAttr("pangolin_public_resource_policy.test", "header_auth")}})
}

func TestAccNetworkPrivateResource(t *testing.T) {
	f := newNetworkLifecycleFixture(t, "private_resource")
	runLifecycle(t, lifecycleCase{Resource: f.kind, CreateAttributes: `org_id = "test"
name = "Private service"
mode = "host"
destination = "192.0.2.10"
site_ids = ["1"]
role_ids = ["2"]`, UpdateAttributes: `org_id = "test"
name = "Private service updated"
mode = "host"
destination = "192.0.2.11"
site_ids = ["2"]
role_ids = ["3"]`, ImportID: "test/42", Handler: f.handler, Remove: f.remove, Exists: f.exists,
		UpdateChecks: []resource.TestCheckFunc{resource.TestCheckResourceAttr("pangolin_private_resource.test", "destination", "192.0.2.11")}})
}

func TestAccNetworkResourceTarget(t *testing.T) {
	f := newNetworkLifecycleFixture(t, "resource_target")
	runLifecycle(t, lifecycleCase{Resource: f.kind, CreateAttributes: `resource_id = "11"
site_id = "1"
ip = "192.0.2.10"
port = 8080`, UpdateAttributes: `resource_id = "11"
site_id = "1"
ip = "192.0.2.10"
port = 8081
hc_enabled = true
hc_hostname = "example.internal"
hc_path = "/health"`, ImportID: "public/11/42", Handler: f.handler, Remove: f.remove, Exists: f.exists,
		UpdateChecks: []resource.TestCheckFunc{resource.TestCheckResourceAttr("pangolin_resource_target.test", "port", "8081")}})
}

func TestAccNetworkAIProviderTarget(t *testing.T) {
	f := newNetworkLifecycleFixture(t, "resource_target")
	f.aiTarget = true
	runLifecycle(t, lifecycleCase{Resource: f.kind, CreateAttributes: `ai_provider_id = "11"
site_id = "1"
ip = "192.0.2.20"
port = 11434
hc_enabled = true
hc_hostname = "inference.internal"
hc_interval = 1
hc_timeout = 1`, UpdateAttributes: `ai_provider_id = "11"
site_id = "1"
ip = "192.0.2.21"
port = 11434
hc_enabled = true
hc_hostname = "inference.internal"
hc_interval = 1
hc_timeout = 1
hc_path = "/health"`, ImportID: "ai/11/42", Handler: f.handler, Remove: f.remove, Exists: f.exists,
		Checks:       []resource.TestCheckFunc{resource.TestCheckResourceAttr("pangolin_resource_target.test", "ai_provider_id", "11"), resource.TestCheckNoResourceAttr("pangolin_resource_target.test", "resource_id")},
		UpdateChecks: []resource.TestCheckFunc{resource.TestCheckResourceAttr("pangolin_resource_target.test", "ip", "192.0.2.21")}})
}

func TestAccNetworkPrivateSSHPortReplacement(t *testing.T) {
	f := newNetworkLifecycleFixture(t, "private_resource")
	runLifecycle(t, lifecycleCase{Resource: f.kind, CreateAttributes: `org_id = "test"
name = "SSH"
mode = "ssh"
destination = "192.0.2.10"
destination_port = 22
site_ids = ["1"]`, UpdateAttributes: `org_id = "test"
name = "SSH"
mode = "ssh"
destination = "192.0.2.10"
destination_port = 2222
site_ids = ["1"]`, ImportID: "test/42", Handler: f.handler, Remove: f.remove, Exists: f.exists,
		Checks:       []resource.TestCheckFunc{resource.TestCheckResourceAttr("pangolin_private_resource.test", "tcp_port_range_string", "22")},
		UpdateChecks: []resource.TestCheckFunc{resource.TestCheckResourceAttr("pangolin_private_resource.test", "tcp_port_range_string", "2222")}})
}

func TestAccNetworkTargetDataSources(t *testing.T) {
	f := newNetworkLifecycleFixture(t, "resource_target")
	server := testHTTPServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			f.fail(w, r, "data source must be read-only")
			return
		}
		switch r.URL.Path {
		case "/v1/target/42":
			f.reply(w, Object{"targetId": 42, "providerId": 11, "resourceId": nil, "siteId": 1, "ip": "192.0.2.10", "port": 11434})
		case "/v1/target/43":
			f.reply(w, Object{"targetId": 43, "providerId": nil, "resourceId": 12, "siteId": 1, "ip": "192.0.2.11", "port": 8080})
		case "/v1/ai-provider/11/targets":
			id := 42
			if r.URL.Query().Get("offset") == "1" {
				id = 44
			}
			f.reply(w, Object{"targets": []any{Object{"targetId": id, "providerId": 11, "resourceId": nil, "ip": "192.0.2.10", "port": 11434}}, "pagination": Object{"total": 2}})
		default:
			f.fail(w, r, "unexpected read route")
		}
	})
	config := fmt.Sprintf(`provider "pangolin" {
endpoint = %q
api_key = "lifecycle-fixture-key"
}

data "pangolin_resource_target" "ai" { id = "42" }
data "pangolin_resource_target" "public" { id = "43" }
data "pangolin_ai_provider_targets" "all" { ai_provider_id = "11" }`, server.URL+"/v1")
	resource.Test(t, resource.TestCase{ProtoV6ProviderFactories: testProtoV6ProviderFactories(), Steps: []resource.TestStep{{Config: config, Check: resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("data.pangolin_resource_target.ai", "ai_provider_id", "11"),
		resource.TestCheckResourceAttr("data.pangolin_resource_target.public", "resource_id", "12"),
		resource.TestCheckResourceAttr("data.pangolin_ai_provider_targets.all", "items.#", "2"),
		resource.TestCheckResourceAttr("data.pangolin_ai_provider_targets.all", "items.1.id", "44"),
	)}}})
}

func TestAccNetworkPublicResourcePolicy(t *testing.T) {
	f := newNetworkLifecycleFixture(t, "public_resource_policy")
	runLifecycle(t, lifecycleCase{Resource: f.kind, CreateAttributes: `org_id = "test"
name = "Restricted"
nice_id = "managed-policy"
role_ids = ["2"]
password = "initial-password"`, UpdateAttributes: `org_id = "test"
name = "Restricted updated"
nice_id = "managed-policy"
sso = false
role_ids = ["3"]
password = "updated-password"`, ImportID: "test/42", ImportIgnore: []string{"password"}, Handler: f.handler, Remove: f.remove, Exists: f.exists,
		UpdateChecks: []resource.TestCheckFunc{resource.TestCheckResourceAttr("pangolin_public_resource_policy.test", "sso", "false")}})
}

func TestAccNetworkPublicAccessDataSources(t *testing.T) {
	f := newNetworkLifecycleFixture(t, "public_resource")
	server := testHTTPServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			f.fail(w, r, "data source must never mutate authorization")
			return
		}
		switch r.URL.Path {
		case "/v1/public-resource/42":
			f.reply(w, Object{"resourceId": 42, "orgId": "test", "name": "Inline", "mode": "http", "resourcePolicyId": nil, "defaultResourcePolicyId": 99})
		case "/v1/public-resource/43":
			f.reply(w, Object{"resourceId": 43, "orgId": "test", "name": "Shared", "mode": "http", "resourcePolicyId": 7, "defaultResourcePolicyId": 100})
		case "/v1/public-resource/44":
			f.reply(w, Object{"resourceId": 44, "orgId": "test", "name": "Raw TCP", "mode": "tcp", "resourcePolicyId": nil, "defaultResourcePolicyId": nil})
		case "/v1/public-resource/42/policies":
			f.reply(w, Object{"defaultPolicy": Object{"sso": false, "applyRules": true, "emailWhitelistEnabled": true, "idpId": 5, "roles": []any{Object{"roleId": 1, "isAdmin": true}, Object{"roleId": 2}}, "users": []any{Object{"userId": "user-1"}}, "emailWhiteList": []any{Object{"email": "allowed@example.com"}}, "rules": []any{Object{"ruleId": 8, "action": "ACCEPT", "match": "IP", "value": "192.0.2.0/24", "priority": 1, "enabled": true}}, "passwordId": 11, "pincodeId": 12, "headerAuth": Object{"id": 13, "extendedCompatibility": false}}, "sharedPolicy": nil})
		case "/v1/org/test/public-resources":
			// The real list endpoint intentionally omits policy ownership IDs.
			f.reply(w, Object{"resources": []any{Object{"resourceId": 42, "mode": "http"}, Object{"resourceId": 43, "mode": "http"}, Object{"resourceId": 44, "mode": "tcp"}}, "pagination": Object{"total": 3}})
		default:
			f.fail(w, r, "unexpected route; shared/TCP access must not be read as inline")
		}
	})
	config := fmt.Sprintf(`provider "pangolin" {
endpoint = %q
api_key = "lifecycle-fixture-key"
}
data "pangolin_public_resource" "inline" {
org_id = "test"
id = "42"
}
data "pangolin_public_resource" "shared" {
org_id = "test"
id = "43"
}
data "pangolin_public_resource" "tcp" {
org_id = "test"
id = "44"
}
data "pangolin_public_resources" "all" { org_id = "test" }`, server.URL+"/v1")
	checks := []resource.TestCheckFunc{
		resource.TestCheckNoResourceAttr("data.pangolin_public_resource.shared", "access"),
		resource.TestCheckNoResourceAttr("data.pangolin_public_resource.tcp", "access"),
		resource.TestCheckResourceAttr("data.pangolin_public_resources.all", "items.#", "3"),
		resource.TestCheckNoResourceAttr("data.pangolin_public_resources.all", "items.1.access"),
		resource.TestCheckResourceAttr("data.pangolin_public_resources.all", "items.1.resource_policy_id", "7"),
	}
	for address, prefix := range map[string]string{"data.pangolin_public_resource.inline": "access.", "data.pangolin_public_resources.all": "items.0.access."} {
		checks = append(checks,
			resource.TestCheckResourceAttr(address, prefix+"sso", "false"),
			resource.TestCheckResourceAttr(address, prefix+"skip_to_idp_id", "5"),
			resource.TestCheckResourceAttr(address, prefix+"role_ids.#", "1"),
			resource.TestCheckTypeSetElemAttr(address, prefix+"role_ids.*", "2"),
			resource.TestCheckTypeSetElemAttr(address, prefix+"user_ids.*", "user-1"),
			resource.TestCheckTypeSetElemAttr(address, prefix+"emails.*", "allowed@example.com"),
			resource.TestCheckResourceAttr(address, prefix+"rules.#", "1"),
			resource.TestCheckNoResourceAttr(address, prefix+"password"),
			resource.TestCheckNoResourceAttr(address, prefix+"pincode"),
			resource.TestCheckNoResourceAttr(address, prefix+"header_auth"),
		)
	}
	resource.Test(t, resource.TestCase{ProtoV6ProviderFactories: testProtoV6ProviderFactories(), Steps: []resource.TestStep{{Config: config, Check: resource.ComposeTestCheckFunc(checks...)}}})
}
