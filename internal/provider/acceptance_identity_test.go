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

// The endpoint tables below are deliberately independent of provider definitions.
// They encode upstream handler payloads cited in docs. user_server_admin and
// user_client are specification-adapter cases: their create/destroy fixtures
// intentionally do not claim those operations work on stock Pangolin 1.24.0.
type identityLifecycleFixture struct {
	mu      sync.Mutex
	t       *testing.T
	kind    string
	record  Object
	actions []any
	created int
}

func TestAccIdentityLifecycles(t *testing.T) {
	baseOIDC := `client_id = "oidc-client"
client_secret = "oidc-secret"
auth_url = "https://idp.example/authorize"
token_url = "https://idp.example/token"
identifier_path = "sub"
scopes = "openid email"`
	cases := []struct {
		name, initial, update, importID string
		ignore                          []string
	}{
		{"role", "org_id = \"org\"\nname = \"operators\"\nreassign_role_id = \"1\"\nssh_unix_groups = [\"staff\"]", "org_id = \"org\"\nname = \"maintainers\"\nreassign_role_id = \"1\"\nssh_unix_groups = [\"staff\", \"ops\"]", "org/42", []string{"reassign_role_id"}},
		{"organization_user", "org_id = \"org\"\nusername = \"reader\"\nidp_id = \"7\"\nrole_ids = [\"1\"]\nauto_provisioned = false", "org_id = \"org\"\nusername = \"reader\"\nidp_id = \"7\"\nrole_ids = [\"2\"]\nauto_provisioned = true", "org/u42", nil},
		{"user_server_admin", "user_id = \"u42\"\nserver_admin = false", "user_id = \"u42\"\nserver_admin = true", "u42", nil},
		{"client", "org_id = \"org\"\nname = \"agent\"\nnice_id = \"custom\"", "org_id = \"org\"\nname = \"agent-renamed\"\nnice_id = \"custom\"", "org/42", []string{"secret"}},
		{"user_client", "org_id = \"org\"\nuser_id = \"u42\"\nolm_id = \"existing\"\nname = \"laptop\"\nsubnet = \"100.64.1.2\"", "org_id = \"org\"\nuser_id = \"u42\"\nolm_id = \"existing\"\nname = \"laptop-renamed\"\nsubnet = \"100.64.1.2\"", "org/42", nil},
		{"organization_oidc_identity_provider", "org_id = \"org\"\nname = \"SSO\"\nrole_mapping = \"role-v1\"\n" + baseOIDC, "org_id = \"org\"\nname = \"SSO-renamed\"\nrole_mapping = \"role-v2\"\n" + baseOIDC, "org/42", nil},
		{"global_oidc_identity_provider", "name = \"SSO\"\ndefault_role_mapping = \"default-v1\"\n" + baseOIDC, "name = \"SSO-renamed\"\ndefault_role_mapping = \"default-v2\"\n" + baseOIDC, "42", []string{"redirect_url"}},
		{"identity_provider_organization", "idp_id = \"7\"\norg_id = \"org\"\nrole_mapping = \"role-v1\"", "idp_id = \"7\"\norg_id = \"org\"\nrole_mapping = \"role-v2\"", "7/org", nil},
		{"organization_api_key", "org_id = \"org\"\nname = \"automation\"\naction_ids = [\"getSite\"]", "org_id = \"org\"\nname = \"automation\"\naction_ids = [\"getSite\", \"getClient\"]", "org/key42", []string{"api_key"}},
		{"site_provisioning_key", "org_id = \"org\"\nname = \"bootstrap\"", "org_id = \"org\"\nname = \"bootstrap\"\nmax_batch_size = 10\napprove_new_sites = false\nvalid_until = \"2030-01-01T00:00:00Z\"", "org/key42", []string{"site_provisioning_key"}},
		{"public_resource_access_token", "org_id = \"org\"\nresource_id = \"11\"\ntitle = \"first\"\nvalid_for_seconds = 3600\npath = \"/restricted\"", "org_id = \"org\"\nresource_id = \"11\"\ntitle = \"second\"\nvalid_for_seconds = 3600\npath = \"/restricted\"", "org/11/token42", []string{"access_token", "path"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			fixture := &identityLifecycleFixture{t: t, kind: tc.name}
			if tc.name == "user_server_admin" {
				fixture.record = Object{"userId": "u42", "serverAdmin": false, "type": "internal"}
			}
			runLifecycle(t, lifecycleCase{
				Resource: tc.name, CreateAttributes: tc.initial, UpdateAttributes: tc.update, ImportID: tc.importID, ImportIgnore: tc.ignore,
				Handler: fixture.serve,
				Remove: func() {
					fixture.mu.Lock()
					defer fixture.mu.Unlock()
					if fixture.kind == "user_server_admin" {
						fixture.record["serverAdmin"] = false
					} else {
						fixture.record = nil
					}
				},
				Exists: func() bool {
					fixture.mu.Lock()
					defer fixture.mu.Unlock()
					if fixture.kind == "user_server_admin" {
						return fixture.record["serverAdmin"] == true
					}
					return fixture.record != nil
				},
				Checks: []resource.TestCheckFunc{resource.TestCheckResourceAttrSet("pangolin_"+tc.name+".test", "id")},
			})
		})
	}
}

func TestAccIdentityClientStateTransitions(t *testing.T) {
	for _, name := range []string{"client", "user_client"} {
		t.Run(name, func(t *testing.T) {
			fixture := &identityLifecycleFixture{t: t, kind: name}
			base := "org_id = \"org\"\nname = \"state-test\"\n"
			if name == "user_client" {
				base += "user_id = \"u42\"\nolm_id = \"existing\"\nsubnet = \"100.64.1.2\"\n"
			}
			ignore := []string{}
			if name == "client" {
				ignore = append(ignore, "secret")
			}
			address := "pangolin_" + name + ".test"
			runLifecycle(t, lifecycleCase{
				Resource: name, CreateAttributes: base + "blocked = true\narchived = true", UpdateAttributes: base + "blocked = false\narchived = false", ImportID: "org/42", ImportIgnore: ignore,
				Handler:      fixture.serve,
				Remove:       func() { fixture.mu.Lock(); defer fixture.mu.Unlock(); fixture.record = nil },
				Exists:       func() bool { fixture.mu.Lock(); defer fixture.mu.Unlock(); return fixture.record != nil },
				Checks:       []resource.TestCheckFunc{resource.TestCheckResourceAttr(address, "blocked", "true"), resource.TestCheckResourceAttr(address, "archived", "true")},
				UpdateChecks: []resource.TestCheckFunc{resource.TestCheckResourceAttr(address, "blocked", "false"), resource.TestCheckResourceAttr(address, "archived", "false")},
			})
		})
	}
}

func (f *identityLifecycleFixture) serve(w http.ResponseWriter, r *http.Request) {
	f.mu.Lock()
	defer f.mu.Unlock()
	path := strings.TrimPrefix(r.URL.Path, "/v1")
	var body Object
	if r.Body != nil && r.ContentLength > 0 {
		body = identityRequestBody(f.t, r)
	}
	write := func(data any) { identityRespond(w, data) }
	bad := func() {
		f.t.Errorf("unexpected %s request: %s %s, body %#v", f.kind, r.Method, path, body)
		http.Error(w, "unexpected fixture request", 400)
	}
	set := func(allowed string) bool {
		for name := range body {
			if !strings.Contains(" "+allowed+" ", " "+name+" ") {
				bad()
				return false
			}
		}
		for name, value := range body {
			f.record[name] = value
		}
		return true
	}
	get := func() {
		if f.record == nil {
			http.NotFound(w, r)
			return
		}
		write(f.record)
	}
	list := func(key string) {
		values := []any{}
		if f.record != nil {
			values = append(values, f.record)
		}
		write(Object{key: values, "pagination": Object{"total": len(values), "limit": 100, "offset": 0}})
	}
	del := func() { f.record = nil; write(nil) }
	switch f.kind {
	case "role":
		if r.Method == "GET" && path == "/org/org/roles" {
			if r.URL.Query().Get("page") != "1" || r.URL.Query().Get("pageSize") != "100" {
				bad()
				return
			}
			rows := []any{}
			if f.record != nil {
				row := Object{}
				for k, v := range f.record {
					row[k] = v
				}
				row["allowSsh"] = false
				rows = append(rows, row)
			}
			write(Object{"roles": rows, "pagination": Object{"total": len(rows), "page": 1, "pageSize": 100}})
			return
		}
		if r.Method == "GET" && path == "/role/1" {
			write(Object{"roleId": 1, "orgId": "org", "name": "default"})
			return
		}
		if r.Method == "GET" && path == "/role/42" {
			get()
			return
		}
		if r.Method == "DELETE" && path == "/role/42" {
			if Text(body["roleId"]) != "1" {
				bad()
				return
			}
			del()
			return
		}
		if (r.Method == "PUT" && path == "/org/org/role") || (r.Method == "POST" && path == "/role/42") {
			if r.Method == "PUT" {
				f.record = Object{"roleId": 42, "orgId": "org", "isAdmin": false, "requireDeviceApproval": false, "sshSudoMode": "full", "sshCreateHomeDir": true, "sshSudoCommands": "[]", "sshUnixGroups": "[]"}
			}
			if !set("name description requireDeviceApproval allowSsh sshSudoMode sshSudoCommands sshCreateHomeDir sshUnixGroups") {
				return
			}
			for _, key := range []string{"sshSudoCommands", "sshUnixGroups"} {
				if value, ok := body[key]; ok {
					data, _ := json.Marshal(value)
					f.record[key] = string(data)
				}
			}
			delete(f.record, "allowSsh")
			write(f.record)
			return
		}
	case "organization_user":
		if path == "/org/org/users" && r.Method == "GET" {
			rows := []any{}
			if f.record != nil {
				rows = append(rows, Object{"id": "u42", "orgId": "org", "username": f.record["username"], "idpId": f.record["idpId"], "roles": []any{Object{"roleId": 1, "roleName": "reader"}}})
			}
			write(Object{"users": rows, "pagination": Object{"total": len(rows), "page": 1, "pageSize": 100}})
			return
		}
		if r.Method == "PUT" && path == "/org/org/user" {
			if _, ok := body["idpId"].(json.Number); !ok || body["type"] != "oidc" {
				bad()
				return
			}
			f.record = Object{"userId": "u42", "orgId": "org", "autoProvisioned": false}
			if !set("email username name type idpId roleIds") {
				return
			}
			write(Object{})
			return
		}
		if r.Method == "GET" && path == "/org/org/user-by-username" {
			if r.URL.Query().Get("username") != "reader" || r.URL.Query().Get("idpId") != "7" {
				bad()
				return
			}
			get()
			return
		}
		if path == "/org/org/user/u42" {
			switch r.Method {
			case "GET":
				// The mounted 1.24.0 handler incorrectly requires req.userOrg.
				w.WriteHeader(http.StatusForbidden)
				return
			case "DELETE":
				del()
				return
			case "POST":
				if !set("autoProvisioned") {
					return
				}
				write(f.record)
				return
			}
		}
		if strings.HasPrefix(path, "/user/u42/add-role/") && r.Method == "POST" {
			role := strings.TrimPrefix(path, "/user/u42/add-role/")
			roles := Strings(f.record["roleIds"])
			for _, existing := range roles {
				if role == existing {
					bad()
					return
				}
			}
			f.record["roleIds"] = append(roles, role)
			write(nil)
			return
		}
		if strings.HasPrefix(path, "/user/u42/remove-role/") && r.Method == "DELETE" {
			role := strings.TrimPrefix(path, "/user/u42/remove-role/")
			remaining := []string{}
			for _, existing := range Strings(f.record["roleIds"]) {
				if existing != role {
					remaining = append(remaining, existing)
				}
			}
			if len(remaining) == 0 {
				bad()
				return
			}
			f.record["roleIds"] = remaining
			write(nil)
			return
		}
	case "user_server_admin":
		if path == "/user/u42" && r.Method == "GET" {
			get()
			return
		}
		if path == "/user/u42/server-admin" && r.Method == "POST" {
			if !set("serverAdmin") {
				return
			}
			write(f.record)
			return
		}
	case "client", "user_client":
		if path == "/org/org/pick-client-defaults" && r.Method == "GET" {
			write(Object{"olmId": "olm42", "olmSecret": "client-secret", "subnet": "100.64.1.2"})
			return
		}
		createPath := "/org/org/client"
		if f.kind == "user_client" {
			createPath = "/org/org/user/u42/client"
		}
		if path == createPath && r.Method == "PUT" {
			if body["subnet"] != "100.64.1.2" || body["type"] != "olm" {
				bad()
				return
			}
			f.record = Object{"clientId": 42, "orgId": "org", "niceId": "generated-client", "archived": false, "blocked": false}
			allowed := "name olmId secret subnet type"
			if f.kind == "user_client" {
				allowed = "name olmId subnet type"
				f.record["userId"] = "u42"
			}
			if !set(allowed) {
				return
			}
			delete(f.record, "secret")
			f.record["subnet"] = "100.64.1.2/16"
			write(f.record)
			return
		}
		if path == "/client/42" {
			switch r.Method {
			case "GET":
				get()
				return
			case "DELETE":
				del()
				return
			case "POST":
				if !set("name niceId") {
					return
				}
				write(f.record)
				return
			}
		}
		if r.Method == "POST" {
			field, enabled := "", false
			switch path {
			case "/client/42/archive":
				field, enabled = "archived", true
			case "/client/42/unarchive":
				field = "archived"
			case "/client/42/block":
				field, enabled = "blocked", true
			case "/client/42/unblock":
				field = "blocked"
			}
			if field != "" {
				if f.record == nil || f.record[field] == enabled {
					bad()
					return
				}
				f.record[field] = enabled
				write(nil)
				return
			}
		}
	case "organization_oidc_identity_provider", "global_oidc_identity_provider":
		prefix := ""
		if f.kind == "organization_oidc_identity_provider" {
			prefix = "/org/org"
		}
		if path == prefix+"/idp/oidc" && r.Method == "PUT" {
			f.record = Object{"idpId": 42, "variant": "oidc", "autoProvision": false}
			allowed := "name clientId clientSecret authUrl tokenUrl identifierPath emailPath namePath scopes autoProvision variant tags"
			if prefix != "" {
				allowed += " roleMapping orgMapping"
			}
			if !set(allowed) {
				return
			}
			write(Object{"idpId": 42, "redirectUrl": "https://pangolin.example/idp/42/callback"})
			return
		}
		if path == prefix+"/idp/42/oidc" && r.Method == "POST" {
			allowed := "name clientId clientSecret authUrl tokenUrl identifierPath emailPath namePath scopes autoProvision tags"
			if prefix != "" {
				allowed += " roleMapping orgMapping"
			} else {
				allowed += " defaultRoleMapping defaultOrgMapping variant"
			}
			if !set(allowed) {
				return
			}
			write(nil)
			return
		}
		if path == prefix+"/idp/42" {
			if r.Method == "DELETE" {
				del()
				return
			}
			if r.Method == "GET" {
				if f.record == nil {
					http.NotFound(w, r)
					return
				}
				idp, oidc, org := Object{}, Object{}, Object{}
				for key, value := range f.record {
					switch key {
					case "roleMapping", "orgMapping":
						org[key] = value
					case "clientId", "clientSecret", "authUrl", "tokenUrl", "identifierPath", "emailPath", "namePath", "scopes", "variant":
						oidc[key] = value
					default:
						idp[key] = value
					}
				}
				out := Object{"idp": idp, "idpOidcConfig": oidc}
				if prefix != "" {
					org["orgId"] = "org"
					out["idpOrg"] = org
					out["redirectUrl"] = "https://pangolin.example/idp/42/callback"
				}
				write(out)
				return
			}
		}
	case "identity_provider_organization":
		if path == "/idp/7/org" && r.Method == "GET" {
			list("policies")
			return
		}
		if path == "/idp/7/org/org" {
			if r.Method == "DELETE" {
				del()
				return
			}
			if r.Method == "PUT" || r.Method == "POST" {
				if r.Method == "PUT" {
					f.record = Object{"orgId": "org", "idpId": 7}
				}
				if !set("roleMapping orgMapping") {
					return
				}
				write(Object{})
				return
			}
		}
	case "organization_api_key":
		if path == "/org/org/api-keys" && r.Method == "GET" {
			list("apiKeys")
			return
		}
		if path == "/org/org/api-key" && r.Method == "PUT" {
			f.record = Object{"apiKeyId": "key42", "orgId": "org", "createdAt": "2026-01-01T00:00:00Z", "lastChars": "cret"}
			if !set("name") {
				return
			}
			f.actions = nil
			out := Object{}
			for k, v := range f.record {
				out[k] = v
			}
			out["apiKey"] = "one-time-api-secret"
			write(out)
			return
		}
		if path == "/org/org/api-key/key42" && r.Method == "DELETE" {
			del()
			return
		}
		if path == "/org/org/api-key/key42/actions" {
			if r.Method == "POST" {
				ids, ok := body["actionIds"].([]any)
				if !ok || len(ids) == 0 {
					bad()
					return
				}
				f.actions = ids
				write(nil)
				return
			}
			if r.Method == "GET" {
				rows := []any{}
				for _, id := range f.actions {
					rows = append(rows, Object{"actionId": id})
				}
				write(Object{"actions": rows, "pagination": Object{"total": len(rows), "limit": 100, "offset": 0}})
				return
			}
		}
	case "site_provisioning_key":
		if path == "/org/org/site-provisioning-keys" && r.Method == "GET" {
			list("siteProvisioningKeys")
			return
		}
		if path == "/org/org/site-provisioning-key" && r.Method == "PUT" {
			if _, present := body["maxBatchSize"]; !present {
				bad()
				return
			}
			f.record = Object{"siteProvisioningKeyId": "key42", "orgId": "org", "createdAt": "2026-01-01T00:00:00Z", "lastChars": "cret", "lastUsed": nil, "numUsed": 0, "approveNewSites": true, "validUntil": nil}
			if !set("name maxBatchSize validUntil approveNewSites") {
				return
			}
			if f.record["validUntil"] == "" {
				f.record["validUntil"] = nil
			}
			out := Object{}
			for k, v := range f.record {
				out[k] = v
			}
			out["siteProvisioningKey"] = "one-time-site-secret"
			write(out)
			return
		}
		if path == "/org/org/site-provisioning-key/key42" {
			if r.Method == "DELETE" {
				del()
				return
			}
			if r.Method == "PATCH" {
				if !set("maxBatchSize validUntil approveNewSites") {
					return
				}
				if f.record["validUntil"] == "" {
					f.record["validUntil"] = nil
				}
				write(f.record)
				return
			}
		}
	case "public_resource_access_token":
		if path == "/org/org/access-tokens" && r.Method == "GET" {
			list("accessTokens")
			return
		}
		if path == "/public-resource/11/access-token" && r.Method == "POST" {
			f.record = Object{"accessTokenId": "token42", "orgId": "org", "resourceId": 11, "sessionLength": 3600000, "expiresAt": 2000000000000, "createdAt": 1900000000000, "persistSession": false}
			if !set("validForSeconds title path description persistSession userId") {
				return
			}
			if Text(body["validForSeconds"]) != "3600" {
				bad()
				return
			}
			out := Object{}
			for k, v := range f.record {
				out[k] = v
			}
			delete(f.record, "validForSeconds")
			delete(f.record, "path")
			out["accessToken"] = "one-time-access-secret"
			write(out)
			return
		}
		if path == "/access-token/token42" && r.Method == "DELETE" {
			del()
			return
		}
	}
	bad()
}

func (f *identityLifecycleFixture) String() string { return fmt.Sprintf("%s fixture", f.kind) }
