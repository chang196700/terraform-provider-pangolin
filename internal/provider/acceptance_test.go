package provider

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// Fixtures are authored from the pinned server handlers, independently of
// Definition, so a wrong HTTP verb, field name or response mapping fails tests.
type lifecycleCase struct {
	Resource, CreateAttributes, UpdateAttributes, ImportID string
	ImportIgnore                                           []string
	Handler                                                http.HandlerFunc
	Remove                                                 func()
	Exists                                                 func() bool
	Checks, UpdateChecks                                   []resource.TestCheckFunc
}

func testProtoV6ProviderFactories() map[string]func() (tfprotov6.ProviderServer, error) {
	return map[string]func() (tfprotov6.ProviderServer, error){"pangolin": providerserver.NewProtocol6WithError(New("test")())}
}

func testHTTPServer(t *testing.T, handler http.HandlerFunc) *httptest.Server {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer lifecycle-fixture-key" {
			t.Errorf("missing API bearer header")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		handler(w, r)
	}))
	t.Cleanup(server.Close)
	return server
}

func runLifecycle(t *testing.T, tc lifecycleCase) {
	t.Helper()
	server := testHTTPServer(t, tc.Handler)
	config := func(attributes string) string {
		return fmt.Sprintf("provider \"pangolin\" {\n endpoint = %q\n api_key = \"lifecycle-fixture-key\"\n}\nresource \"pangolin_%s\" \"test\" {\n%s\n}\n", server.URL+"/v1", tc.Resource, attributes)
	}
	initial, updated := config(tc.CreateAttributes), config(tc.UpdateAttributes)
	if tc.UpdateAttributes == "" {
		updated = initial
	}
	address := "pangolin_" + tc.Resource + ".test"
	checks := append([]resource.TestCheckFunc{resource.TestCheckResourceAttrSet(address, "id")}, tc.Checks...)
	updateChecks := append([]resource.TestCheckFunc{resource.TestCheckResourceAttrSet(address, "id")}, tc.UpdateChecks...)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testProtoV6ProviderFactories(),
		CheckDestroy: func(*terraform.State) error {
			if tc.Exists != nil && tc.Exists() {
				return fmt.Errorf("%s still exists after destroy", tc.Resource)
			}
			return nil
		},
		Steps: []resource.TestStep{
			{Config: initial, Check: resource.ComposeTestCheckFunc(checks...)},
			{Config: initial, PlanOnly: true, ExpectNonEmptyPlan: false},
			{Config: updated, Check: resource.ComposeTestCheckFunc(updateChecks...)},
			{Config: updated, PlanOnly: true, ExpectNonEmptyPlan: false},
			{ResourceName: address, ImportState: true, ImportStateId: tc.ImportID, ImportStateVerify: true, ImportStateVerifyIgnore: tc.ImportIgnore},
			{PreConfig: tc.Remove, Config: updated, Check: resource.ComposeTestCheckFunc(updateChecks...)},
		},
	})
}
