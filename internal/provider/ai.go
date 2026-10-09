package provider

import (
	"context"
	"fmt"
	"reflect"
	"strconv"

	"github.com/chang196700/terraform-provider-pangolin/internal/client"
)

// The response projections and lifecycle rules in this file were checked
// against Pangolin 1.24.0 (a235635d9e5a6431ea426431e3e862e965928cd3).
func init() { registerDefinitions(aiDefinitions) }

func aiDefinitions() []Definition {
	return []Definition{aiProviderDefinition(), aiModelDefinition(), aiBudgetDefinition(), virtualAPIKeyDefinition()}
}

func aiProviderDefinition() Definition {
	d := Definition{
		Name: "ai_provider", Description: "Manages an AI gateway upstream provider. API keys and custom headers are stored as sensitive Terraform state.",
		IDField: "providerId", ListKey: "providers", Pagination: "page", ImportFields: []string{"org_id", "id"},
		Fields: SpecFields("/org/{orgId}/ai-provider", "put", "/ai-provider/{providerId}", "post"),
		Create: Endpoint{Method: "PUT", Path: "/org/{org_id}/ai-provider", Key: "provider"},
		Read:   Endpoint{Method: "GET", Path: "/ai-provider/{id}", Key: "provider"},
		Update: Endpoint{Method: "POST", Path: "/ai-provider/{id}", Key: "provider"},
		Delete: Endpoint{Method: "DELETE", Path: "/ai-provider/{id}"},
		List:   Endpoint{Method: "GET", Path: "/org/{org_id}/ai-providers"},
	}
	d.Fields["org_id"] = Local(Immutable(Required(String("orgId"))))
	SetField(&d, "type", func(f *Field) { f.Replace = true; f.Update = false })
	SetField(&d, "api_key", func(f *Field) { f.Sensitive = true })
	SetField(&d, "headers", func(f *Field) {
		f.Sensitive = true
		f.Nullable = false
		f.Default = []any{}
	})
	SetField(&d, "capabilities", func(f *Field) { f.Kind = "set" })
	SetField(&d, "routing_mode", func(f *Field) { f.Default = "url" })
	SetField(&d, "enabled", func(f *Field) { f.Default = true })
	SetField(&d, "skip_tls_verification", func(f *Field) { f.Default = false })
	d.Fields["api_key_last_chars"] = Computed(String("apiKeyLastChars"))
	d.Fields["effective_upstream_url"] = Computed(String("effectiveUpstreamUrl"))
	d.Fields["effective_auth_type"] = Computed(String("effectiveAuthType"))
	aiTimestamps(&d)
	d.Validate = validateAIProvider
	d.EnrichRead = func(_ context.Context, _ *client.Client, _ Values, obj Object) error {
		// The server stores an empty header collection as null, and an empty
		// create-time API key as null. Represent these equivalent empty values
		// consistently, including imports and clearing configured credentials.
		if obj["headers"] == nil {
			obj["headers"] = []any{}
		}
		if key, present := obj["apiKey"]; present && key == nil {
			obj["apiKey"] = ""
		}
		return nil
	}
	d.AfterWrite = func(ctx context.Context, c *client.Client, v, prior Values) error {
		// niceId is supported only by Update; apply an explicitly configured
		// value after the server has allocated the provider ID.
		if prior != nil || v["nice_id"] == nil || Text(v["nice_id"]) == "" {
			return nil
		}
		_, err := Call(ctx, c, d.Update, v, Object{"niceId": v["nice_id"]})
		return err
	}
	return d
}

func validateAIProvider(v Values) error {
	providerType, routing := Text(v["type"]), Text(v["routing_mode"])
	if routing == "" {
		routing = "url"
	}
	if routing == "target" && providerType != "custom" {
		return fmt.Errorf("routing_mode target requires type custom")
	}
	if routing == "target" && Text(v["upstream_url"]) != "" {
		return fmt.Errorf("upstream_url must be omitted for target routing")
	}
	if routing != "target" && Text(v["upstream_url"]) == "" {
		return fmt.Errorf("upstream_url is required for URL routing")
	}
	if providerType == "custom" && len(Strings(v["capabilities"])) == 0 {
		return fmt.Errorf("custom providers require at least one capability")
	}
	if v["capabilities"] != nil && len(Strings(v["capabilities"])) == 0 {
		return fmt.Errorf("capabilities must contain at least one capability")
	}
	return nil
}

func aiModelDefinition() Definition {
	d := Definition{
		Name: "ai_model", Description: "Manages a model allow or block entry on an AI provider. Model keys are unique within a provider.",
		IDField: "modelId", ListKey: "models", Pagination: "page", ImportFields: []string{"provider_id", "id"},
		Fields: SpecFields("/ai-provider/{providerId}/model", "put", "/ai-model/{modelId}", "post"),
		Create: Endpoint{Method: "PUT", Path: "/ai-provider/{provider_id}/model", Key: "model"},
		Read:   Endpoint{Method: "GET", Path: "/ai-model/{id}", Key: "model"},
		Update: Endpoint{Method: "POST", Path: "/ai-model/{id}", Key: "model"},
		Delete: Endpoint{Method: "DELETE", Path: "/ai-model/{id}"},
		List:   Endpoint{Method: "GET", Path: "/ai-provider/{provider_id}/models"},
	}
	d.Fields["provider_id"] = Local(Immutable(Required(ID("providerId"))))
	SetField(&d, "enabled", func(f *Field) { f.Default = true })
	SetField(&d, "list_type", func(f *Field) { f.Default = "allow" })
	aiTimestamps(&d)
	return d
}

var aiBudgetScopes = []string{"provider_id", "model_id", "resource_id", "site_resource_id", "role_id", "virtual_api_key_id"}

func aiBudgetDefinition() Definition {
	d := Definition{
		Name: "ai_budget", Description: "Manages an AI usage budget. Set at most one scope ID; no scope ID selects the organization. A scope, unit and period combination must be unique.",
		IDField: "budgetId", ListKey: "budgets", Pagination: "page", ImportFields: []string{"org_id", "id"},
		Fields: SpecFields("/org/{orgId}/ai-budget", "put", "/ai-budget/{budgetId}", "post"),
		Create: Endpoint{Method: "PUT", Path: "/org/{org_id}/ai-budget", Key: "budget"},
		Read:   Endpoint{Method: "GET", Path: "/ai-budget/{id}", Key: "budget"},
		Update: Endpoint{Method: "POST", Path: "/ai-budget/{id}", Key: "budget"},
		Delete: Endpoint{Method: "DELETE", Path: "/ai-budget/{id}"},
		List:   Endpoint{Method: "GET", Path: "/org/{org_id}/ai-budgets"},
	}
	d.Fields["org_id"] = Local(Immutable(Required(String("orgId"))))
	for _, name := range aiBudgetScopes {
		SetField(&d, name, func(f *Field) {
			if name != "virtual_api_key_id" {
				f.Kind = "id"
			}
			f.Nullable = true
		})
	}
	SetField(&d, "enabled", func(f *Field) { f.Default = true })
	aiTimestamps(&d)
	d.Validate = validateAIBudget
	d.UpdateFunc = func(ctx context.Context, c *client.Client, v, _ Values) error {
		body, err := Body(d, v, false)
		if err != nil {
			return err
		}
		// The API merges a partial update with existing scope fields before
		// validation, so changing scope must clear the old field atomically.
		for _, name := range aiBudgetScopes {
			if v[name] == nil {
				body[d.Fields[name].API] = nil
			}
		}
		_, err = Call(ctx, c, d.Update, v, body)
		return err
	}
	return d
}

func validateAIBudget(v Values) error {
	if amount := v["amount"]; amount != nil && !isUnknown(amount) {
		value, err := strconv.ParseFloat(Text(amount), 64)
		if err != nil || value <= 0 {
			return fmt.Errorf("amount must be greater than zero")
		}
	}
	count := 0
	for _, name := range aiBudgetScopes {
		if v[name] != nil && Text(v[name]) != "" {
			count++
		}
	}
	if count > 1 {
		return fmt.Errorf("at most one of provider_id, model_id, resource_id, site_resource_id, role_id or virtual_api_key_id may be set")
	}
	return nil
}

func virtualAPIKeyDefinition() Definition {
	d := Definition{
		Name: "virtual_api_key", Description: "Manages a manual virtual API key for public inference resources. The returned secret is sensitive and stored in Terraform state. Updating other attributes does not renew its expiry; changing valid_for_seconds does.",
		IDField: "virtualApiKeyId", ListKey: "virtualApiKeys", Pagination: "page", ImportFields: []string{"org_id", "id"},
		Fields: SpecFields("/org/{orgId}/virtual-api-key", "put", "/virtual-api-key/{virtualApiKeyId}", "post"),
		Create: Endpoint{Method: "PUT", Path: "/org/{org_id}/virtual-api-key", Key: "virtualApiKey"},
		Read:   Endpoint{Method: "GET", Path: "/virtual-api-key/{id}", Key: "virtualApiKey"},
		Update: Endpoint{Method: "POST", Path: "/virtual-api-key/{id}", Key: "virtualApiKey"},
		Delete: Endpoint{Method: "DELETE", Path: "/virtual-api-key/{id}"},
		List:   Endpoint{Method: "GET", Path: "/org/{org_id}/virtual-api-keys"},
	}
	d.Fields["org_id"] = Local(Immutable(Required(String("orgId"))))
	// Email delivery is an action, not persistent key configuration. Never
	// retain delivery flags which would resend a secret on unrelated updates.
	for _, name := range []string{"send_email", "send_to_attributed_user", "emails"} {
		delete(d.Fields, name)
	}
	SetField(&d, "resource_ids", func(f *Field) {
		f.Kind = "set"
		f.Default = []any{}
		if f.Element != nil {
			e := *f.Element
			e.Kind = "id"
			f.Element = &e
		}
	})
	SetField(&d, "valid_for_seconds", func(f *Field) {
		f.Preserve = true
		f.Nullable = true
		f.Description = "Lifetime from creation, or from the most recent change to this value. Removing it clears the expiry. This input is not returned by the API."
	})
	d.Fields["secret"] = Computed(Sensitive(String("secret")))
	d.Fields["kind"] = Computed(String("kind"))
	d.Fields["last_chars"] = Computed(String("lastChars"))
	d.Fields["expires_at"] = Computed(Int("expiresAt"))
	d.Fields["last_used_at"] = Computed(Int("lastUsedAt"))
	d.Fields["created_at"] = Computed(Int("createdAt"))
	d.Fields["created_by_user_id"] = Computed(String("createdByUserId"))
	d.Validate = validateVirtualAPIKey
	d.UpdateFunc = func(ctx context.Context, c *client.Client, v, prior Values) error {
		body, err := Body(d, v, false)
		if err != nil {
			return err
		}
		if reflect.DeepEqual(v["valid_for_seconds"], prior["valid_for_seconds"]) {
			delete(body, "validForSeconds")
		} else if v["valid_for_seconds"] == nil {
			body["validForSeconds"] = nil
		}
		_, err = Call(ctx, c, d.Update, v, body)
		return err
	}
	return d
}

func validateVirtualAPIKey(v Values) error {
	all, _ := v["all_resources"].(bool)
	count := len(Strings(v["resource_ids"]))
	if all && count != 0 {
		return fmt.Errorf("resource_ids must be empty when all_resources is true")
	}
	if !all && count == 0 {
		return fmt.Errorf("set all_resources or select at least one public inference resource in resource_ids")
	}
	return nil
}

func aiTimestamps(d *Definition) {
	d.Fields["created_at"] = Computed(Int("createdAt"))
	d.Fields["updated_at"] = Computed(Int("updatedAt"))
}
