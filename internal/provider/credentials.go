package provider

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/chang196700/terraform-provider-pangolin/internal/client"
)

func init() { registerDefinitions(credentialDefinitions) }

func credentialDefinitions() []Definition {
	defs := []Definition{organizationAPIKeyDefinition(), siteProvisioningKeyDefinition(), resourceAccessTokenDefinition()}
	for i := range defs {
		defs[i].ListKey = defs[i].List.Key
	}
	return defs
}

func organizationAPIKeyDefinition() Definition {
	d := Definition{
		Name: "organization_api_key", Description: "An organization API key with authoritative action assignments. The API key value is returned only during creation and is stored as sensitive state. Imported keys cannot recover their original secret.", IDField: "apiKeyId",
		Fields: SpecFields("/org/{orgId}/api-key", "put", "", ""),
		Create: Endpoint{Method: "PUT", Path: "/org/{org_id}/api-key"}, Delete: Endpoint{Method: "DELETE", Path: "/org/{org_id}/api-key/{id}"},
		List: Endpoint{Method: "GET", Path: "/org/{org_id}/api-keys", Key: "apiKeys"}, Pagination: "offset-short", ImportFields: []string{"org_id", "id"}, Sources: []string{"organization_api_key", "organization_api_keys"},
	}
	// listOrgApiKeys reports the current page length as total. Reading through
	// the final short page avoids marking an existing later-page key absent.
	d.Fields["org_id"] = identityOrgField()
	d.Fields["action_ids"] = Required(Local(StringSet("actionIds")))
	SetField(&d, "action_ids", func(f *Field) {
		f.Default = nil
		f.Preserve = false
		f.Description = "Nonempty authoritative set of API action identifiers. The server cannot replace this set with an empty set."
	})
	d.Fields["api_key"] = Preserve(Sensitive(Computed(String("apiKey"))))
	SetField(&d, "api_key", func(f *Field) {
		f.Description = "One-time API key secret. The server returns only the secret component; form the Bearer credential as id.api_key (for example, ${pangolin_organization_api_key.example.id}.${pangolin_organization_api_key.example.api_key}). Import cannot recover this secret."
	})
	d.Fields["last_chars"] = Computed(String("lastChars"))
	d.Fields["created_at"] = Computed(String("createdAt"))
	d.Validate = func(v Values) error {
		if v["action_ids"] != nil && len(Strings(v["action_ids"])) == 0 {
			return fmt.Errorf("action_ids must contain at least one action")
		}
		return nil
	}
	// Name is immutable; actions are updated by their dedicated endpoint.
	d.UpdateFunc = func(context.Context, *client.Client, Values, Values) error { return nil }
	d.AfterWrite = func(ctx context.Context, c *client.Client, v, prior Values) error {
		if !identityChanged(prior, v, "action_ids") {
			return nil
		}
		return identityWrite(ctx, c, Endpoint{Method: "POST", Path: "/org/{org_id}/api-key/{id}/actions"}, v, Object{"actionIds": v["action_ids"]})
	}
	d.EnrichRead = func(ctx context.Context, c *client.Client, v Values, obj Object) error {
		path, err := Path("/org/{org_id}/api-key/{id}/actions", v)
		if err != nil {
			return err
		}
		// listApiKeyActions also reports page length as total.
		rows, err := c.List(ctx, path, "actions", "offset-short")
		if err != nil {
			return fmt.Errorf("API key exists but its actions could not be read: %v", err)
		}
		actions := []any{}
		for _, action := range rows {
			if action["actionId"] == nil {
				return fmt.Errorf("API key response contains an invalid action")
			}
			actions = append(actions, Text(action["actionId"]))
		}
		obj["actionIds"] = actions
		return nil
	}
	return d
}

func siteProvisioningKeyDefinition() Definition {
	d := Definition{
		Name: "site_provisioning_key", Description: "An organization site provisioning key. The secret is returned only during creation, retained in sensitive state, and cannot be recovered by import.", IDField: "siteProvisioningKeyId",
		Fields: SpecFields("/org/{orgId}/site-provisioning-key", "put", "/org/{orgId}/site-provisioning-key/{siteProvisioningKeyId}", "patch"),
		Create: Endpoint{Method: "PUT", Path: "/org/{org_id}/site-provisioning-key"}, Update: Endpoint{Method: "PATCH", Path: "/org/{org_id}/site-provisioning-key/{id}"}, Delete: Endpoint{Method: "DELETE", Path: "/org/{org_id}/site-provisioning-key/{id}"},
		List: Endpoint{Method: "GET", Path: "/org/{org_id}/site-provisioning-keys", Key: "siteProvisioningKeys"}, Pagination: "offset-short", ImportFields: []string{"org_id", "id"}, Sources: []string{"site_provisioning_key", "site_provisioning_keys"},
	}
	// listSiteProvisioningKeys counts only the returned page, not all keys.
	d.Fields["org_id"] = identityOrgField()
	SetField(&d, "max_batch_size", func(f *Field) {
		min, max := float64(1), float64(1000000)
		f.Kind = "int"
		f.Required = false
		f.Nullable = true
		f.Min, f.Max = &min, &max
		f.Description = "Maximum sites per batch, from 1 through 1000000. Null means no limit."
	})
	SetField(&d, "valid_until", func(f *Field) {
		f.Default = ""
		f.Description = "RFC3339 expiry timestamp. Empty string removes expiration."
	})
	d.Fields["site_provisioning_key"] = Preserve(Sensitive(Computed(String("siteProvisioningKey"))))
	for _, field := range []string{"lastChars", "createdAt", "lastUsed"} {
		d.Fields[snake(field)] = Computed(String(field))
	}
	d.Fields["num_used"] = Computed(Int("numUsed"))
	d.Validate = func(v Values) error {
		if until := Text(v["valid_until"]); until != "" {
			if _, err := time.Parse(time.RFC3339Nano, until); err != nil {
				return fmt.Errorf("valid_until must be an RFC3339 timestamp or empty")
			}
		}
		return nil
	}
	d.EnrichRead = func(_ context.Context, _ *client.Client, v Values, obj Object) error {
		if obj["validUntil"] == nil {
			obj["validUntil"] = ""
			return nil
		}
		// The server serializes UTC with millisecond precision. Preserve a
		// semantically identical configured timestamp to prevent a false diff.
		if configured := Text(v["valid_until"]); configured != "" {
			a, ea := time.Parse(time.RFC3339Nano, configured)
			b, eb := time.Parse(time.RFC3339Nano, Text(obj["validUntil"]))
			if ea == nil && eb == nil && a.Equal(b) {
				obj["validUntil"] = configured
			}
		}
		return nil
	}
	d.CreateFunc = func(ctx context.Context, c *client.Client, v Values) (Object, error) {
		body, err := Body(d, v, true)
		if err != nil {
			return nil, err
		}
		// This API requires the key even when its value is null (unlimited).
		if v["max_batch_size"] == nil {
			body["maxBatchSize"] = nil
		}
		obj, err := Call(ctx, c, d.Create, v, body)
		if err != nil {
			return nil, err
		}
		if err = d.EnrichRead(ctx, c, v, obj); err != nil {
			return obj, err
		}
		return obj, nil
	}
	return d
}

func resourceAccessTokenDefinition() Definition {
	d := Definition{
		Name: "public_resource_access_token", Description: "A public resource access token. Settings are immutable, so changing them replaces the token. The one-time token is retained in sensitive state; import cannot recover it. Expired tokens disappear from the API list and will be recreated on the next apply.", IDField: "accessTokenId",
		Fields: SpecFields("/public-resource/{resourceId}/access-token", "post", "", ""),
		Create: Endpoint{Method: "POST", Path: "/public-resource/{resource_id}/access-token"}, Delete: Endpoint{Method: "DELETE", Path: "/access-token/{id}"},
		List: Endpoint{Method: "GET", Path: "/org/{org_id}/access-tokens", Key: "accessTokens"}, Pagination: "offset-short", ImportFields: []string{"org_id", "resource_id", "id"}, Sources: []string{"public_resource_access_token", "public_resource_access_tokens"},
	}
	// listAccessTokens counts accessible resources instead of tokens. That total
	// can be higher or lower than the result and cannot terminate pagination.
	// The same handler prioritizes orgId over resourceId in its query and can
	// return other resources' tokens. Enforce the configured scope locally.
	d.ListFilter = func(v Values, obj Object) bool { return Text(obj["resourceId"]) == Text(v["resource_id"]) }
	d.ListInputs = []string{"resource_id"}
	// The resource-scoped list cannot resolve orgId for root keys in 1.24.0.
	// The explicitly organization-scoped route works with both key types.
	d.Fields["org_id"] = identityOrgField()
	d.Fields["resource_id"] = Required(Immutable(Local(ID("resourceId"))))
	d.Fields["access_token"] = Preserve(Sensitive(Computed(String("accessToken"))))
	d.Fields["expires_at"] = Computed(Int("expiresAt"))
	d.Fields["created_at"] = Computed(Int("createdAt"))
	SetField(&d, "path", func(f *Field) {
		f.Preserve = true
		f.Description = "Optional request path restriction. The token list API omits this value, so it is retained locally and cannot be recovered on import."
	})
	SetField(&d, "valid_for_seconds", func(f *Field) { f.Preserve = true })
	d.EnrichRead = func(_ context.Context, _ *client.Client, v Values, obj Object) error {
		if obj["expiresAt"] != nil {
			ms, err := strconv.ParseInt(Text(obj["sessionLength"]), 10, 64)
			if err != nil || ms <= 0 || ms%1000 != 0 {
				return fmt.Errorf("access token response has invalid sessionLength")
			}
			obj["validForSeconds"] = ms / 1000
		}
		for _, field := range []string{"title", "description", "path"} {
			if obj[field] == nil && v[field] != nil && Text(v[field]) == "" {
				obj[field] = ""
			}
		}
		return nil
	}
	return d
}
