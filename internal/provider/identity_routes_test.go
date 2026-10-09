package provider

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/chang196700/terraform-provider-pangolin/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Optional source verification against the locally downloaded official 1.24.0
// snapshot a235635d9e5a6431ea426431e3e862e965928cd3. This checks route mounts and
// handler/authentication guards independently of OpenAPI and mock fixtures.
// CI without that source snapshot skips this check; unit tests below still run.
func TestPinnedIdentityRouteAvailability(t *testing.T) {
	root := filepath.Join("..", "..", ".cache", "upstream")
	if _, err := os.Stat(filepath.Join(root, "server", "routers", "integration.ts")); os.IsNotExist(err) {
		t.Skip("official source snapshot is not installed")
	}
	read := func(name string) string {
		t.Helper()
		data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(name)))
		if err != nil {
			t.Fatal(err)
		}
		return string(data)
	}
	withoutComments := func(s string) string {
		s = regexp.MustCompile(`(?s)/\*.*?\*/`).ReplaceAllString(s, "")
		return regexp.MustCompile(`(?m)^\s*//.*$`).ReplaceAllString(s, "")
	}
	base := withoutComments(read("server/routers/integration.ts"))
	private := withoutComments(read("server/private/routers/integration.ts"))
	combined := base + private
	mounted := func(code, method, route string) bool {
		return regexp.MustCompile(`authenticated\.` + method + `\(\s*["` + "`" + `]` + regexp.QuoteMeta(route) + `["` + "`" + `]`).MatchString(code)
	}
	for _, tc := range []struct{ method, route string }{
		{"get", "/org/:orgId/users"}, {"get", "/org/:orgId/user-by-username"},
		{"get", "/client/:clientId"}, {"post", "/client/:clientId"}, {"delete", "/client/:clientId"},
		{"get", "/org/:orgId/access-tokens"}, {"get", "/user/:userId"},
	} {
		if !mounted(base, tc.method, tc.route) {
			t.Errorf("working route disappeared: %s %s", tc.method, tc.route)
		}
	}
	for _, tc := range []struct{ method, route string }{
		{"put", "/org/:orgId/user/:userId/client"}, {"post", "/user/:userId/server-admin"},
	} {
		if mounted(combined, tc.method, tc.route) {
			t.Errorf("capability boundary changed; review adapter diagnostics: %s %s", tc.method, tc.route)
		}
	}
	for _, check := range []struct{ name, contains string }{
		{"server/routers/integration.ts", "authenticated.use(verifyApiKey)"},
		{"server/middlewares/integration/verifyApiKey.ts", `authHeader.startsWith("Bearer ")`},
		{"server/private/routers/integration.ts", "export const authenticated = a"},
		{"server/routers/external.ts", "authenticated.use(verifySessionUserMiddleware)"},
		{"server/auth/sessions/verifySession.ts", "req.cookies[SESSION_COOKIE_NAME]"},
		{"server/routers/user/getOrgUser.ts", "if (!req.userOrg)"},
		{"server/routers/user/listUsers.ts", "id: users.userId"},
		{"server/routers/client/deleteClient.ts", "if (client.userId)"},
	} {
		if !strings.Contains(read(check.name), check.contains) {
			t.Errorf("upstream behavior changed: %s no longer contains %q", check.name, check.contains)
		}
	}
	if mounted(base, "delete", "/user/:userId/remove-role/:roleId") || !mounted(private, "delete", "/user/:userId/remove-role/:roleId") {
		t.Error("membership role removal edition boundary changed")
	}
}

func TestIdentityOrganizationUserReadUsesBearerCompatibleRoutes(t *testing.T) {
	for _, mode := range []string{"ok", "lookup404", "wrongID", "absent"} {
		t.Run(mode, func(t *testing.T) {
			c := identityTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				switch r.URL.Path {
				case "/v1/org/org/users":
					rows := []any{}
					if mode != "absent" {
						rows = append(rows, Object{"id": "u1", "username": "renamed", "idpId": 7, "roles": []any{Object{"roleId": 1, "roleName": "reader"}}})
					}
					identityRespond(w, Object{"users": rows, "pagination": Object{"total": len(rows), "page": 1, "pageSize": 100}})
				case "/v1/org/org/user-by-username":
					if r.URL.Query().Get("username") != "renamed" || r.URL.Query().Get("idpId") != "7" {
						t.Error("lookup must use fresh identity from list, including external username drift")
					}
					if mode == "lookup404" {
						http.NotFound(w, r)
						return
					}
					id := "u1"
					if mode == "wrongID" {
						id = "other"
					}
					identityRespond(w, Object{"userId": id, "username": "renamed", "roleIds": []any{1}, "autoProvisioned": true})
				default:
					t.Errorf("unexpected or session-only route: %s", r.URL.Path)
					w.WriteHeader(http.StatusForbidden)
				}
			})
			d := organizationUserDefinition()
			obj, err := d.ReadFunc(context.Background(), c, Values{"org_id": "org", "id": "u1", "username": "old-name"})
			if mode == "ok" {
				if err != nil || obj["autoProvisioned"] != true || Text(obj["userId"]) != "u1" {
					t.Fatalf("read failed: %#v %v", obj, err)
				}
			} else if mode == "absent" {
				if !client.IsNotFound(err) {
					t.Fatalf("missing membership must clear state: %v", err)
				}
			} else if err == nil || client.IsNotFound(err) {
				t.Fatalf("supplemental lookup failure must preserve existing state: %v", err)
			}
		})
	}
}

func TestIdentityOrganizationUsersDataSourceRealListShape(t *testing.T) {
	c := identityTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/org/org/users":
			identityRespond(w, Object{"users": []any{Object{"id": "u1", "username": "reader", "idpId": 7, "roles": []any{Object{"roleId": 3, "roleName": "reader"}}}}, "pagination": Object{"total": 1, "page": 1, "pageSize": 100}})
		case "/v1/org/org/user-by-username":
			identityRespond(w, Object{"userId": "u1", "orgId": "org", "username": "reader", "idpId": 7, "roleIds": []any{3}, "autoProvisioned": true})
		default:
			t.Errorf("unexpected route: %s", r.URL.Path)
			http.NotFound(w, r)
		}
	})
	var dd DataDefinition
	for _, candidate := range dataDefinitions() {
		if candidate.Name == "organization_users" {
			dd = candidate
		}
	}
	obj, err := objectValue(dd.fields(), Values{"org_id": "org"})
	if err != nil {
		t.Fatal(err)
	}
	raw, err := obj.ToTerraformValue(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	schema := dataSchema(dd)
	res := datasource.ReadResponse{State: tfsdk.State{Schema: schema}}
	reader := &apiDataSource{definition: dd, client: c}
	reader.Read(context.Background(), datasource.ReadRequest{Config: tfsdk.Config{Schema: schema, Raw: raw}}, &res)
	if res.Diagnostics.HasError() {
		t.Fatal(res.Diagnostics)
	}
	var items types.List
	if ds := res.State.GetAttribute(context.Background(), path.Root("items"), &items); ds.HasError() {
		t.Fatal(ds)
	}
	rows := native(items).([]any)
	if len(rows) != 1 {
		t.Fatalf("unexpected rows: %#v", rows)
	}
	row := rows[0].(map[string]any)
	if row["id"] != "u1" || row["auto_provisioned"] != true || len(Strings(row["role_ids"])) != 1 || Strings(row["role_ids"])[0] != "3" {
		t.Fatalf("list dropped real identity, roles, or membership metadata: %#v", row)
	}
}
