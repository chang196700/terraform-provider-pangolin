package provider

import (
	"context"
	"fmt"
	"reflect"

	"github.com/chang196700/terraform-provider-pangolin/internal/client"
)

// Domain management is registered by Pangolin's Integration API router even
// though the generated OpenAPI document omits the PUT and DELETE operations.
// Keep the lifecycle here so the Terraform resource follows the server's
// actual routes and response fields.
func init() { registerDefinitions(domainDefinitions) }

func domainDefinitions() []Definition { return []Definition{domainDefinition()} }

func domainDefinition() Definition {
	d := Definition{
		Name:         "domain",
		Description:  "Manages an organization domain. Domain type and base domain are immutable; wildcard domain settings can be updated. Import as org_id/domain_id.",
		IDField:      "domainId",
		ListKey:      "domains",
		Pagination:   "offset",
		ImportFields: []string{"org_id", "id"},
		Create:       Endpoint{Method: "PUT", Path: "/org/{org_id}/domain"},
		Read:         Endpoint{Method: "GET", Path: "/org/{org_id}/domain/{id}"},
		Update:       Endpoint{Method: "POST", Path: "/org/{org_id}/domain/{id}"},
		Delete:       Endpoint{Method: "DELETE", Path: "/org/{org_id}/domain/{id}"},
		List:         Endpoint{Method: "GET", Path: "/org/{org_id}/domains"},
		Fields: map[string]Field{
			"org_id":               Local(Immutable(Required(String("orgId")))),
			"type":                 Required(Immutable(String("type"))),
			"base_domain":          Required(Immutable(String("baseDomain"))),
			"cert_resolver":        String("certResolver"),
			"prefer_wildcard_cert": WithDefault(Bool("preferWildcardCert"), false),
			"custom_cert_resolver": Computed(String("customCertResolver")),
			"verified":             Computed(Bool("verified")),
			"failed":               Computed(Bool("failed")),
			"tries":                Computed(Int("tries")),
			"config_managed":       Computed(Bool("configManaged")),
			"error_message":        Computed(String("errorMessage")),
			"last_checked_at":      Computed(Int("lastCheckedAt")),
		},
	}

	SetField(&d, "type", func(f *Field) {
		f.Enum = []string{"ns", "cname", "wildcard"}
		f.Description = "Domain integration type. OSS supports wildcard; SaaS supports ns and cname; Enterprise supports all three. Changes replace the domain."
	})
	SetField(&d, "base_domain", func(f *Field) {
		f.Description = "The base domain to register. Changes replace the domain."
	})
	SetField(&d, "cert_resolver", func(f *Field) {
		f.Nullable = true
		f.Description = "Optional certificate resolver. Only wildcard domains accept updates to this setting."
	})
	SetField(&d, "prefer_wildcard_cert", func(f *Field) {
		f.Description = "Whether wildcard certificate issuance is preferred for wildcard domains."
	})

	d.Validate = func(v Values) error {
		typ := Text(v["type"])
		if typ == "" {
			return nil
		}
		if typ != "ns" && typ != "cname" && typ != "wildcard" {
			return fmt.Errorf("type must be one of ns, cname, wildcard")
		}
		return nil
	}
	d.UpdateFunc = func(ctx context.Context, c *client.Client, plan, prior Values) error {
		// Pangolin's update endpoint rejects settings changes for NS/CNAME
		// domains. Send only attributes that actually changed, preserving the
		// distinction between an omitted value and an explicit null.
		if Text(plan["type"]) != "wildcard" {
			for _, key := range []string{"cert_resolver", "prefer_wildcard_cert"} {
				if !reflect.DeepEqual(plan[key], prior[key]) {
					return fmt.Errorf("%s can only be updated for wildcard domains", key)
				}
			}
			return nil
		}
		fields := map[string]Field{}
		for _, key := range []string{"cert_resolver", "prefer_wildcard_cert"} {
			if reflect.DeepEqual(plan[key], prior[key]) {
				continue
			}
			fields[key] = d.Fields[key]
		}
		if len(fields) == 0 {
			return nil
		}
		body, err := Body(Definition{Fields: fields}, plan, false)
		if err != nil {
			return err
		}
		p, err := Path(d.Update.Path, plan)
		if err != nil {
			return err
		}
		return c.Do(ctx, d.Update.Method, p, body, nil)
	}

	return d
}
