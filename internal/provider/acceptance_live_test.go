package provider

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/chang196700/terraform-provider-pangolin/internal/client"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// TestAccLiveSite is the only test allowed to use a real Pangolin endpoint.
// Four explicit opt-ins are required before a client or provider is created;
// ordinary tests and the full local mock acceptance suite cannot activate it.
func TestAccLiveSite(t *testing.T) {
	if os.Getenv("PANGOLIN_LIVE_TEST") != "1" {
		t.Skip("set PANGOLIN_LIVE_TEST=1 to enable the live site test")
	}
	org, endpoint, key := os.Getenv("PANGOLIN_LIVE_ORG_ID"), os.Getenv("PANGOLIN_ENDPOINT"), os.Getenv("PANGOLIN_API_KEY")
	if org == "" || endpoint == "" || key == "" {
		t.Skip("live test also requires PANGOLIN_LIVE_ORG_ID, PANGOLIN_ENDPOINT and PANGOLIN_API_KEY")
	}
	if os.Getenv("TF_ACC") != "1" {
		t.Skip("live test requires TF_ACC=1 and an explicitly enabled acceptance run")
	}
	c, err := client.New(endpoint, key, 30*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	random := make([]byte, 8)
	if _, err = rand.Read(random); err != nil {
		t.Fatal(err)
	}
	name := "terraform-acceptance-" + hex.EncodeToString(random)
	config := func(label string) string {
		// Credentials remain in the environment rather than the generated HCL.
		return fmt.Sprintf("provider \"pangolin\" {}\nresource \"pangolin_site\" \"live\" {\n org_id = %q\n name = %q\n type = \"local\"\n}\n", org, label)
	}
	address := "pangolin_site.live"
	var siteID string
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testProtoV6ProviderFactories(),
		CheckDestroy: func(*terraform.State) error {
			if siteID == "" {
				return nil
			}
			path, err := client.ExpandPath("/site/{id}", map[string]any{"id": siteID})
			if err != nil {
				return err
			}
			var obj client.Object
			err = c.Do(context.Background(), "GET", path, nil, &obj)
			if client.IsNotFound(err) {
				return nil
			}
			if err != nil {
				return fmt.Errorf("could not verify live site deletion: %w", err)
			}
			return fmt.Errorf("live site %s still exists after destroy", siteID)
		},
		Steps: []resource.TestStep{
			{Config: config(name), Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr(address, "name", name), resource.TestCheckResourceAttr(address, "type", "local"), func(s *terraform.State) error {
				row := s.RootModule().Resources[address]
				if row == nil || row.Primary.ID == "" {
					return fmt.Errorf("live site has no ID")
				}
				siteID = row.Primary.ID
				return nil
			})},
			{Config: config(name), PlanOnly: true, ExpectNonEmptyPlan: false},
			{Config: config(name + "-updated"), Check: resource.TestCheckResourceAttr(address, "name", name+"-updated")},
			{Config: config(name + "-updated"), PlanOnly: true, ExpectNonEmptyPlan: false},
			{ResourceName: address, ImportState: true, ImportStateVerify: true, ImportStateIdFunc: func(*terraform.State) (string, error) { return org + "/" + siteID, nil }},
		},
	})
}
