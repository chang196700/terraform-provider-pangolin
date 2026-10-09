package provider

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/chang196700/terraform-provider-pangolin/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestCredentialAPIKeyReadsEveryActionPageDespiteWrongTotal(t *testing.T) {
	d := organizationAPIKeyDefinition()
	pages := []string{}
	c := identityTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		offset := r.URL.Query().Get("offset")
		pages = append(pages, offset)
		rows := []any{}
		if offset == "0" {
			for i := 0; i < 100; i++ {
				rows = append(rows, Object{"actionId": fmt.Sprintf("action%d", i)})
			}
		} else {
			rows = append(rows, Object{"actionId": "last-action"})
		}
		identityRespond(w, Object{"actions": rows, "pagination": Object{"total": len(rows), "limit": 100, "offset": offset}})
	})
	obj := Object{"apiKeyId": "key"}
	if err := d.EnrichRead(context.Background(), c, Values{"org_id": "example", "id": "key"}, obj); err != nil {
		t.Fatal(err)
	}
	if len(Strings(obj["actionIds"])) != 101 || !reflect.DeepEqual(pages, []string{"0", "100"}) {
		t.Fatalf("truncated action list: %v pages=%v", obj, pages)
	}
}

func TestCredentialLaterPageIdentitySurvivesIncorrectTotals(t *testing.T) {
	for _, d := range credentialDefinitions() {
		t.Run(d.Name, func(t *testing.T) {
			pages := []string{}
			c := identityTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				offset := r.URL.Query().Get("offset")
				pages = append(pages, offset)
				if r.URL.Query().Get("limit") != "100" {
					t.Errorf("missing bounded page limit: %v", r.URL.Query())
				}
				rows := []any{}
				if offset == "0" {
					for i := 0; i < 100; i++ {
						rows = append(rows, Object{d.IDField: fmt.Sprintf("object-%d", i), "resourceId": 11})
					}
				} else {
					rows = append(rows, Object{d.IDField: "target", "resourceId": 11})
				}
				// Key handlers report page length. Token handler instead reports
				// the number of resources, which may be less than a full token page.
				total := len(rows)
				if d.Name == "public_resource_access_token" {
					total = 1
				}
				identityRespond(w, Object{d.ListKey: rows, "pagination": Object{"total": total, "limit": 100, "offset": offset}})
			})
			obj, err := readRemote(context.Background(), c, d, Values{"org_id": "org", "resource_id": "11", "id": "target"})
			if err != nil {
				t.Fatalf("existing later-page credential must not be marked missing: %v", err)
			}
			if Text(obj[d.IDField]) != "target" || !reflect.DeepEqual(pages, []string{"0", "100"}) {
				t.Fatalf("incorrect total truncated credential read: %#v, pages=%v", obj, pages)
			}
			for _, source := range dataDefinitions() {
				if source.Entity.Name == d.Name && source.IsList && source.Pagination != "offset-short" {
					t.Fatal("list data source lost the endpoint-specific pagination workaround")
				}
			}
		})
	}
}

func TestCredentialSupplemental404DoesNotReportMissingAPIKey(t *testing.T) {
	d := organizationAPIKeyDefinition()
	c := identityTestClient(t, func(w http.ResponseWriter, r *http.Request) { http.NotFound(w, r) })
	err := d.EnrichRead(context.Background(), c, Values{"org_id": "org", "id": "key"}, Object{"apiKeyId": "key"})
	if err == nil || client.IsNotFound(err) {
		t.Fatalf("supplemental actions failure must not delete existing key state: %v", err)
	}
}

func TestCredentialAccessTokenFiltersServerOrganizationWideResults(t *testing.T) {
	d := resourceAccessTokenDefinition()
	d.ListKey = d.List.Key
	c := identityTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/org/org/access-tokens" {
			t.Errorf("must use explicitly organization-scoped route for root-key support: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		identityRespond(w, Object{"accessTokens": []any{
			Object{"accessTokenId": "other", "resourceId": 22, "title": "another resource"},
			Object{"accessTokenId": "ours", "resourceId": 11, "title": "our resource"},
		}, "pagination": Object{"total": 1, "limit": 100, "offset": 0}})
	})
	if _, err := readRemote(context.Background(), c, d, Values{"org_id": "org", "resource_id": "11", "id": "other"}); !client.IsNotFound(err) {
		t.Fatalf("read accepted token from another resource: %v", err)
	}
	if obj, err := readRemote(context.Background(), c, d, Values{"org_id": "org", "resource_id": "11", "id": "ours"}); err != nil || Text(obj["resourceId"]) != "11" {
		t.Fatalf("read lost correctly scoped token: %#v, %v", obj, err)
	}
	var dd DataDefinition
	for _, candidate := range dataDefinitions() {
		if candidate.Name == "public_resource_access_tokens" {
			dd = candidate
			break
		}
	}
	if dd.Name == "" {
		t.Fatal("token list data source is missing")
	}
	if !dd.Inputs["resource_id"].Required || !dd.Inputs["org_id"].Required {
		t.Fatal("token collection requires both organization and resource scope")
	}
	obj, err := objectValue(dd.fields(), Values{"org_id": "org", "resource_id": "11"})
	if err != nil {
		t.Fatal(err)
	}
	raw, err := obj.ToTerraformValue(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	schema := dataSchema(dd)
	reader := &apiDataSource{definition: dd, client: c}
	res := datasource.ReadResponse{State: tfsdk.State{Schema: schema}}
	reader.Read(context.Background(), datasource.ReadRequest{Config: tfsdk.Config{Schema: schema, Raw: raw}}, &res)
	if res.Diagnostics.HasError() {
		t.Fatal(res.Diagnostics)
	}
	var items types.List
	if diagnostics := res.State.GetAttribute(context.Background(), path.Root("items"), &items); diagnostics.HasError() {
		t.Fatal(diagnostics)
	}
	rows := native(items).([]any)
	if len(rows) != 1 || rows[0].(map[string]any)["id"] != "ours" {
		t.Fatalf("data source leaked tokens outside configured resource: %#v", rows)
	}
}

func TestCredentialAPIKeyPartialCreateRetainsIDAndOneTimeSecret(t *testing.T) {
	d := organizationAPIKeyDefinition()
	d.Fields["id"] = Computed(String(d.IDField))
	c := identityTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "PUT" {
			identityRespond(w, Object{"apiKeyId": "generated", "apiKey": "one-time-secret", "name": "automation", "lastChars": "cret", "createdAt": "2026-01-01T00:00:00Z"})
			return
		}
		w.WriteHeader(http.StatusForbidden)
	})
	values := Values{"org_id": "example", "name": "automation", "action_ids": []any{"getSite"}}
	obj, err := objectValue(d.Fields, values)
	if err != nil {
		t.Fatal(err)
	}
	raw, err := obj.ToTerraformValue(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	schema := resourceSchema(d)
	r := &managedResource{definition: d, client: c}
	res := resource.CreateResponse{State: tfsdk.State{Schema: schema}}
	r.Create(context.Background(), resource.CreateRequest{Plan: tfsdk.Plan{Schema: schema, Raw: raw}}, &res)
	if !res.Diagnostics.HasError() {
		t.Fatal("failed action assignment must return a diagnostic")
	}
	var id, secret string
	if ds := res.State.GetAttribute(context.Background(), path.Root("id"), &id); ds.HasError() {
		t.Fatal(ds)
	}
	if ds := res.State.GetAttribute(context.Background(), path.Root("api_key"), &secret); ds.HasError() {
		t.Fatal(ds)
	}
	if id != "generated" || secret != "one-time-secret" {
		t.Fatalf("partial create lost recovery state: id=%q secretPresent=%t", id, secret != "")
	}
}

func TestCredentialSiteKeySendsExplicitUnlimitedAndNormalizesExpiry(t *testing.T) {
	d := siteProvisioningKeyDefinition()
	c := identityTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		body := identityRequestBody(t, r)
		if value, exists := body["maxBatchSize"]; !exists || value != nil {
			t.Errorf("unlimited must be explicit null: %#v", body)
		}
		identityRespond(w, Object{"siteProvisioningKeyId": "key", "name": "bootstrap", "maxBatchSize": nil, "validUntil": "2030-01-01T00:00:00.000Z", "siteProvisioningKey": "one-time-secret", "approveNewSites": true})
	})
	v := Values{"org_id": "example", "name": "bootstrap", "valid_until": "2030-01-01T08:00:00+08:00", "approve_new_sites": true}
	obj, err := d.CreateFunc(context.Background(), c, v)
	if err != nil {
		t.Fatal(err)
	}
	if obj["validUntil"] != v["valid_until"] {
		t.Fatalf("equivalent timestamp changed: %#v", obj["validUntil"])
	}
	body, err := Body(d, Values{"max_batch_size": nil, "valid_until": "", "approve_new_sites": false}, false)
	if err != nil {
		t.Fatal(err)
	}
	if value, exists := body["maxBatchSize"]; !exists || value != nil {
		t.Errorf("update null missing: %#v", body)
	}
	if body["validUntil"] != "" || body["approveNewSites"] != false {
		t.Fatalf("empty and false were omitted: %#v", body)
	}
	if d.Validate(Values{"valid_until": "tomorrow"}) == nil {
		t.Fatal("invalid expiry accepted")
	}
}

func TestCredentialOneTimeSecretsSurviveRefreshAndCannotBeRecoveredByImport(t *testing.T) {
	for _, tc := range []struct {
		definition Definition
		field, api string
	}{
		{organizationAPIKeyDefinition(), "api_key", "apiKey"},
		{siteProvisioningKeyDefinition(), "site_provisioning_key", "siteProvisioningKey"},
		{resourceAccessTokenDefinition(), "access_token", "accessToken"},
	} {
		t.Run(tc.definition.Name, func(t *testing.T) {
			f := tc.definition.Fields[tc.field]
			if !f.Computed || !f.Sensitive || !f.Preserve {
				t.Fatal("one-time secret schema does not preserve sensitive output")
			}
			values := Values{"id": "key", tc.field: "one-time-value"}
			if err := mergeRemote(tc.definition, values, Object{tc.definition.IDField: "key", tc.api: nil}); err != nil {
				t.Fatal(err)
			}
			if values[tc.field] != "one-time-value" {
				t.Fatal("refresh lost one-time secret")
			}
			imported := Values{"id": "key"}
			if err := mergeRemote(tc.definition, imported, Object{tc.definition.IDField: "key"}); err != nil {
				t.Fatal(err)
			}
			if imported[tc.field] != nil {
				t.Fatal("import invented unrecoverable secret")
			}
		})
	}
}

func TestCredentialAccessTokenReconstructsDurationAndPreservesPath(t *testing.T) {
	d := resourceAccessTokenDefinition()
	obj := Object{"accessTokenId": "token", "expiresAt": 2000000000000, "sessionLength": 3600000, "title": nil, "description": nil}
	v := Values{"id": "token", "path": "/private", "title": "", "access_token": "opaque-token"}
	if err := d.EnrichRead(context.Background(), nil, v, obj); err != nil {
		t.Fatal(err)
	}
	if err := mergeRemote(d, v, obj); err != nil {
		t.Fatal(err)
	}
	if Text(v["valid_for_seconds"]) != "3600" || v["path"] != "/private" || v["title"] != "" || v["access_token"] != "opaque-token" {
		t.Fatalf("unexpected token state: %#v", v)
	}
}
