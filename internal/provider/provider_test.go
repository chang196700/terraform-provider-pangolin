package provider

import (
	"context"
	"github.com/chang196700/terraform-provider-pangolin/internal/client"
	"testing"
)

func TestAllProviderSchemas(t *testing.T) {
	defs := definitions()
	if len(defs) != 24 {
		t.Fatalf("expected 24 resources, got %d", len(defs))
	}
	names := map[string]bool{}
	for _, d := range defs {
		if names[d.Name] {
			t.Fatalf("duplicate resource %s", d.Name)
		}
		names[d.Name] = true
		t.Run(d.Name, func(t *testing.T) {
			if err := verifySchema(context.Background(), d); err != nil {
				t.Fatal(err)
			}
			if d.IDField == "" || len(d.ImportFields) == 0 && d.ImportFunc == nil {
				t.Fatal("resource identity/import missing")
			}
			if d.Delete.Path == "" && d.DeleteFunc == nil {
				t.Fatal("delete lifecycle missing")
			}
		})
	}
	names = map[string]bool{}
	for _, d := range dataDefinitions() {
		if names[d.Name] {
			t.Fatalf("duplicate data source %s", d.Name)
		}
		names[d.Name] = true
		t.Run("data_"+d.Name, func(t *testing.T) {
			if diag := dataSchema(d).ValidateImplementation(context.Background()); diag.HasError() {
				t.Fatal(diag)
			}
		})
	}
}

func TestRelatedRead404DoesNotRemovePrimaryResource(t *testing.T) {
	r := managedResource{definition: Definition{IDField: "entityId", Fields: map[string]Field{}, ReadFunc: func(context.Context, *client.Client, Values) (Object, error) { return Object{"entityId": 1}, nil }, EnrichRead: func(context.Context, *client.Client, Values, Object) error {
		return client.NotFound("/missing-relation-route")
	}}}
	if err := r.refresh(context.Background(), Values{"id": "1"}); err == nil || client.IsNotFound(err) {
		t.Fatalf("related 404 must be an ordinary read diagnostic: %v", err)
	}
}
func TestNullableAndZeroSerialization(t *testing.T) {
	d := Definition{Fields: map[string]Field{"scope": {Kind: "id", API: "scopeId", Update: true, Nullable: true}, "enabled": Bool("enabled"), "limit": Int("limit"), "roles": IDSet("roleIds"), "omitted": String("omitted")}}
	v := Values{"enabled": false, "limit": int64(0), "roles": []any{}}
	b, e := Body(d, v, false)
	if e != nil {
		t.Fatal(e)
	}
	if _, ok := b["scopeId"]; !ok {
		t.Fatal("explicit null omitted")
	}
	if _, ok := b["omitted"]; ok {
		t.Fatal("optional absence was sent")
	}
	if b["enabled"] != false || b["limit"] != int64(0) || len(b["roleIds"].([]any)) != 0 {
		t.Fatal("zero values lost")
	}
}
func TestOptionalOnlyAccess(t *testing.T) {
	f := Field{Kind: "object", OptionalOnly: true, Fields: map[string]Field{"enabled": Bool("enabled")}}
	a := resourceAttribute(f)
	if a.IsComputed() || !a.IsOptional() {
		t.Fatal("access omission must remain a real removal plan")
	}
}
func TestAcronymSnakeCase(t *testing.T) {
	for in, want := range map[string]string{"sendAISessionLogs": "send_ai_session_logs", "siteResourceId": "site_resource_id", "hcHTTPMethod": "hc_http_method"} {
		if got := snake(in); got != want {
			t.Fatalf("%s => %s, want %s", in, got, want)
		}
	}
}
