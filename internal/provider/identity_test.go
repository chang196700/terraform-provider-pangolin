package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/chang196700/terraform-provider-pangolin/internal/client"
)

// These fixtures model the public handler responses in the pinned upstream source.
// In particular, role SSH arrays are database JSON strings, IdP GET results
// contain joined objects, and organization-user create returns an empty object.
func identityTestClient(t *testing.T, handler http.HandlerFunc) *client.Client {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	c, err := client.New(server.URL+"/v1", "fixture-key", 5*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	return c
}

func identityRespond(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{"data": data, "success": true, "error": false, "status": 200})
}

func identityRequestBody(t *testing.T, r *http.Request) Object {
	t.Helper()
	var body Object
	dec := json.NewDecoder(r.Body)
	dec.UseNumber()
	if err := dec.Decode(&body); err != nil {
		t.Errorf("decode request: %v", err)
	}
	return body
}

func TestIdentityOrganizationUserResolvesIDThroughOIDCUsername(t *testing.T) {
	d := organizationUserDefinition()
	calls := []string{}
	c := identityTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		calls = append(calls, r.Method+" "+r.URL.Path)
		switch r.URL.Path {
		case "/v1/org/example/user":
			body := identityRequestBody(t, r)
			if _, ok := body["idpId"].(json.Number); !ok {
				t.Errorf("idpId must be numeric: %#v", body)
			}
			if !reflect.DeepEqual(body["roleIds"], []any{json.Number("12"), json.Number("13")}) {
				t.Errorf("roleIds = %#v", body["roleIds"])
			}
			if _, exists := body["roleId"]; exists {
				t.Error("legacy roleId sent alongside roleIds")
			}
			identityRespond(w, Object{})
		case "/v1/org/example/user-by-username":
			if r.URL.Query().Get("username") != "reader+test@example.com" || r.URL.Query().Get("idpId") != "7" {
				t.Errorf("lookup query = %v", r.URL.Query())
			}
			identityRespond(w, Object{"userId": "external-user", "username": "reader+test@example.com", "roleIds": []any{12, 13}})
		default:
			t.Errorf("unexpected path %s", r.URL.Path)
			http.NotFound(w, r)
		}
	})
	obj, err := d.CreateFunc(context.Background(), c, Values{"org_id": "example", "username": "reader+test@example.com", "type": "oidc", "idp_id": "7", "role_ids": []any{"12", "13"}})
	if err != nil {
		t.Fatal(err)
	}
	if Text(obj["userId"]) != "external-user" || len(calls) != 2 {
		t.Fatalf("created = %#v; calls=%v", obj, calls)
	}
}

func TestIdentityRoleChangesAddBeforeRemove(t *testing.T) {
	for _, failAddition := range []bool{false, true} {
		t.Run(fmt.Sprint(failAddition), func(t *testing.T) {
			calls := []string{}
			c := identityTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				calls = append(calls, r.Method+" "+r.URL.Path)
				if failAddition && strings.Contains(r.URL.Path, "/add-role/") {
					w.WriteHeader(http.StatusForbidden)
					return
				}
				identityRespond(w, nil)
			})
			err := reconcileUserRoles(context.Background(), c, Values{"id": "u1"}, []string{"1", "2"}, []string{"2", "3"})
			if failAddition {
				if err == nil || len(calls) != 1 {
					t.Fatalf("failed addition must stop before removal: %v %v", err, calls)
				}
			} else {
				want := []string{"POST /v1/user/u1/add-role/3", "DELETE /v1/user/u1/remove-role/1"}
				if err != nil || !reflect.DeepEqual(calls, want) {
					t.Fatalf("reconciliation: %v %v", err, calls)
				}
			}
		})
	}
}

func TestIdentityRoleRetryUsesLiveAssignments(t *testing.T) {
	d := organizationUserDefinition()
	calls := []string{}
	c := identityTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		calls = append(calls, r.Method+" "+r.URL.Path)
		if r.URL.Path == "/v1/org/o1/users" {
			identityRespond(w, Object{"users": []any{Object{"id": "u1", "username": "reader", "idpId": 7}}, "pagination": Object{"total": 1, "page": 1, "pageSize": 100}})
			return
		}
		if r.URL.Path == "/v1/org/o1/user-by-username" {
			identityRespond(w, Object{"userId": "u1", "roleIds": []any{1, 3}})
			return
		}
		identityRespond(w, nil)
	})
	err := d.AfterWrite(context.Background(), c, Values{"org_id": "o1", "id": "u1", "role_ids": []any{"3"}}, Values{"role_ids": []any{"1"}})
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"GET /v1/org/o1/users", "GET /v1/org/o1/user-by-username", "DELETE /v1/user/u1/remove-role/1"}
	if !reflect.DeepEqual(calls, want) {
		t.Fatalf("live assignment retry = %v", calls)
	}
}

func TestIdentityRoleCreateDecodesSerializedSSHLists(t *testing.T) {
	d := roleDefinition()
	c := identityTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		identityRespond(w, Object{"roleId": 11, "orgId": "example", "name": "operator", "sshSudoCommands": `["systemctl status"]`, "sshUnixGroups": `["operators"]`})
	})
	obj, err := d.CreateFunc(context.Background(), c, Values{"org_id": "example", "name": "operator", "reassign_role_id": "1"})
	if err != nil {
		t.Fatal(err)
	}
	v := Values{}
	if err = mergeRemote(d, v, obj); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(Strings(v["ssh_sudo_commands"]), []string{"systemctl status"}) {
		t.Fatalf("decoded commands = %#v", v)
	}
}

func TestIdentityRoleReadsSSHDriftFromPagedRoleList(t *testing.T) {
	d := roleDefinition()
	pages := []string{}
	c := identityTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v1/role/101" {
			identityRespond(w, Object{"roleId": 101, "orgId": "org"})
			return
		}
		page := r.URL.Query().Get("page")
		pages = append(pages, page)
		rows := []any{}
		if page == "1" {
			for i := 1; i <= 100; i++ {
				rows = append(rows, Object{"roleId": i, "allowSsh": false})
			}
		} else {
			rows = append(rows, Object{"roleId": 101, "allowSsh": true})
		}
		identityRespond(w, Object{"roles": rows, "pagination": Object{"total": 101, "page": page, "pageSize": 100}})
	})
	obj, err := d.ReadFunc(context.Background(), c, Values{"org_id": "org", "id": "101", "allow_ssh": false})
	if err != nil {
		t.Fatal(err)
	}
	if obj["allowSsh"] != true || !reflect.DeepEqual(pages, []string{"1", "2"}) {
		t.Fatalf("SSH drift was not read across pages: %#v, %v", obj, pages)
	}
}

func TestIdentityRoleSupplemental404DoesNotReportMissingRole(t *testing.T) {
	d := roleDefinition()
	c := identityTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v1/role/42" {
			identityRespond(w, Object{"roleId": 42, "orgId": "org"})
			return
		}
		http.NotFound(w, r)
	})
	_, err := d.ReadFunc(context.Background(), c, Values{"org_id": "org", "id": "42"})
	if err == nil || client.IsNotFound(err) {
		t.Fatalf("supplemental API 404 must retain the primary role: %v", err)
	}
}

func TestIdentityRoleDeleteChecksReassignmentOrganization(t *testing.T) {
	for _, sameOrganization := range []bool{false, true} {
		t.Run(fmt.Sprint(sameOrganization), func(t *testing.T) {
			d := roleDefinition()
			deleted := false
			c := identityTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				if r.Method == "DELETE" {
					deleted = true
					body := identityRequestBody(t, r)
					if Text(body["roleId"]) != "12" {
						t.Errorf("migration role body = %#v", body)
					}
					identityRespond(w, nil)
					return
				}
				org := "example"
				if !sameOrganization && r.URL.Path == "/v1/role/12" {
					org = "other"
				}
				identityRespond(w, Object{"roleId": 12, "orgId": org})
			})
			err := d.DeleteFunc(context.Background(), c, Values{"id": "11", "org_id": "example", "reassign_role_id": "12"})
			if sameOrganization && (err != nil || !deleted) {
				t.Fatalf("valid migration failed: %v", err)
			}
			if !sameOrganization && (err == nil || deleted) {
				t.Fatalf("cross-organization migration attempted: %v", err)
			}
		})
	}
}

func TestIdentityClientDefaultsPreserveCredentialsAndHostAddress(t *testing.T) {
	d := clientDefinition(false)
	defaults := 0
	c := identityTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v1/org/example/pick-client-defaults" {
			defaults++
			identityRespond(w, Object{"olmId": "olm-generated", "olmSecret": "generated-secret", "subnet": "100.64.1.2"})
			return
		}
		body := identityRequestBody(t, r)
		if body["secret"] != "user-secret" || body["olmId"] != "olm-generated" || body["subnet"] != "100.64.1.2" {
			t.Errorf("create defaults = %#v", body)
		}
		identityRespond(w, Object{"clientId": 8, "name": "agent", "niceId": "auto-name", "subnet": "100.64.1.2/16", "type": "olm"})
	})
	obj, err := d.CreateFunc(context.Background(), c, Values{"org_id": "example", "name": "agent", "type": "olm", "secret": "user-secret", "nice_id": "custom"})
	if err != nil {
		t.Fatal(err)
	}
	if defaults != 1 || obj["secret"] != "user-secret" || obj["subnet"] != "100.64.1.2" || obj["niceId"] != "custom" {
		t.Fatalf("client create result: %#v", obj)
	}
	read := Object{"clientId": 8, "olmId": nil, "subnet": "100.64.1.2/16"}
	if err = d.EnrichRead(context.Background(), c, Values{"olm_id": "olm-generated"}, read); err != nil {
		t.Fatal(err)
	}
	if read["olmId"] != "olm-generated" || defaults != 1 {
		t.Fatalf("read must retain credentials without allocating defaults: %#v", read)
	}
	if err = d.Validate(Values{"subnet": "100.64.1.2/16"}); err == nil {
		t.Fatal("CIDR input unexpectedly accepted")
	}
}

func TestIdentityUserClientNeverGeneratesExistingOLM(t *testing.T) {
	d := clientDefinition(true)
	c := identityTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/pick-client-defaults") {
			identityRespond(w, Object{"olmId": "wrong", "olmSecret": "unused", "subnet": "100.64.1.2"})
			return
		}
		body := identityRequestBody(t, r)
		if body["olmId"] != "existing" {
			t.Errorf("existing olm ID overwritten: %#v", body)
		}
		if _, ok := body["secret"]; ok {
			t.Error("user client must not send generated secret")
		}
		identityRespond(w, Object{"clientId": 4, "userId": "user", "olmId": "existing", "subnet": "100.64.1.2/16"})
	})
	_, err := d.CreateFunc(context.Background(), c, Values{"org_id": "example", "user_id": "user", "olm_id": "existing", "name": "laptop", "type": "olm"})
	if err != nil {
		t.Fatal(err)
	}
}

func TestIdentityClientStateRetryAvoidsDuplicateActions(t *testing.T) {
	calls := []string{}
	c := identityTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		calls = append(calls, r.Method+" "+r.URL.Path)
		if r.Method == "GET" {
			identityRespond(w, Object{"clientId": 42, "blocked": true, "archived": false})
			return
		}
		if r.URL.Path != "/v1/client/42/archive" {
			t.Errorf("already completed block action was repeated: %s", r.URL.Path)
		}
		identityRespond(w, nil)
	})
	err := reconcileClientState(context.Background(), c, Values{"id": "42", "blocked": true, "archived": true}, Values{"blocked": false, "archived": false})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(calls, []string{"GET /v1/client/42", "POST /v1/client/42/archive"}) {
		t.Fatalf("retry did not reconcile current flags: %v", calls)
	}
}

func TestIdentityClientStateFailureStopsSubsequentActions(t *testing.T) {
	calls := []string{}
	c := identityTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		calls = append(calls, r.Method+" "+r.URL.Path)
		if r.Method == "GET" {
			identityRespond(w, Object{"clientId": 42, "blocked": false, "archived": false})
			return
		}
		w.WriteHeader(http.StatusForbidden)
	})
	err := reconcileClientState(context.Background(), c, Values{"id": "42", "blocked": true, "archived": true}, nil)
	if err == nil || !reflect.DeepEqual(calls, []string{"GET /v1/client/42", "POST /v1/client/42/block"}) {
		t.Fatalf("failed block should stop before archive: err=%v, calls=%v", err, calls)
	}
}

func TestIdentityOIDCReadFlattensHandlerJoin(t *testing.T) {
	obj := Object{"idp": Object{"idpId": 7, "name": "SSO", "autoProvision": false}, "idpOidcConfig": Object{"idpId": 7, "clientId": "client", "clientSecret": "secret", "variant": "oidc"}, "idpOrg": Object{"orgId": "example", "roleMapping": "mapping"}, "redirectUrl": "https://example/redirect"}
	if err := flattenIdentityProvider(obj); err != nil {
		t.Fatal(err)
	}
	if obj["clientSecret"] != "secret" || obj["roleMapping"] != "mapping" || obj["name"] != "SSO" {
		t.Fatalf("flattened response: %#v", obj)
	}
	for _, org := range []bool{false, true} {
		d := oidcDefinition(org)
		if !d.Fields["client_secret"].Sensitive {
			t.Fatal("OIDC secret must be sensitive")
		}
		values := Values{}
		if err := mergeRemote(d, values, obj); err != nil {
			t.Fatal(err)
		}
		if values["id"] != "7" || values["client_secret"] != "secret" {
			t.Fatalf("OIDC state: %#v", values)
		}
	}
}

func TestIdentityOrganizationIDPListIgnoresUnfilteredGlobalTotal(t *testing.T) {
	d := oidcDefinition(true)
	d.ListKey = d.List.Key
	requests := 0
	c := identityTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		requests++
		if r.URL.Query().Get("limit") != "100" || r.URL.Query().Get("offset") != "0" {
			t.Errorf("unexpected organization IDP pagination: %v", r.URL.Query())
		}
		// listOrgIdps returns organization-filtered rows but an unfiltered count
		// from the global idp table. An empty subsequent page is not an outage.
		identityRespond(w, Object{"idps": []any{Object{"idpId": 7, "orgId": "org", "name": "organization SSO", "variant": "oidc"}}, "pagination": Object{"total": 300, "limit": 100, "offset": 0}})
	})
	rows, err := c.List(context.Background(), "/org/org/idp", d.ListKey, d.Pagination)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 || requests != 1 {
		t.Fatalf("incorrect global total affected organization listing: rows=%v, requests=%d", rows, requests)
	}
	for _, source := range dataDefinitions() {
		if source.Name == "organization_oidc_identity_providers" && source.Pagination != "offset-short" {
			t.Fatal("list data source did not inherit the corrected organization IDP pagination")
		}
	}
}

func TestIdentityServerAdminDestroyDemotes(t *testing.T) {
	c := identityTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" && r.URL.Path == "/v1/user/u1" {
			identityRespond(w, Object{"userId": "u1", "serverAdmin": true})
			return
		}
		if r.Method != "POST" || r.URL.Path != "/v1/user/u1/server-admin" {
			t.Errorf("destroy must demote: %s %s", r.Method, r.URL.Path)
		}
		if body := identityRequestBody(t, r); body["serverAdmin"] != false {
			t.Errorf("demotion body = %#v", body)
		}
		identityRespond(w, Object{"userId": "u1", "serverAdmin": false})
	})
	if err := serverAdminDefinition().DeleteFunc(context.Background(), c, Values{"id": "u1"}); err != nil {
		t.Fatal(err)
	}
}

func TestIdentityServerAdminMissingActionDoesNotForgetExistingAdmin(t *testing.T) {
	c := identityTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" && r.URL.Path == "/v1/user/u1" {
			identityRespond(w, Object{"userId": "u1", "serverAdmin": true})
			return
		}
		http.NotFound(w, r)
	})
	err := serverAdminDefinition().DeleteFunc(context.Background(), c, Values{"id": "u1"})
	if err == nil || client.IsNotFound(err) || !strings.Contains(err.Error(), "1.24.0") {
		t.Fatalf("missing action must preserve admin state with capability diagnostic: %v", err)
	}
}
