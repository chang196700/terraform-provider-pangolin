package provider

import (
	"context"
	"fmt"

	"github.com/chang196700/terraform-provider-pangolin/internal/client"
)

func networkAddAIFields(d *Definition, private bool) {
	providerFields := map[string]Field{
		"provider_id": Required(ID("providerId")),
		"access_mode": WithDefault(String("accessMode"), "inherit"),
		"enabled":     WithDefault(Bool("enabled"), true),
	}
	f := providerFields["access_mode"]
	f.Enum = []string{"inherit", "select"}
	providerFields["access_mode"] = f
	provider := Field{Kind: "object", Fields: providerFields}
	modelFields := map[string]Field{"model_id": Required(ID("modelId")), "list_type": Required(String("listType"))}
	f = modelFields["list_type"]
	f.Enum = []string{"allow", "block"}
	modelFields["list_type"] = f
	model := Field{Kind: "object", Fields: modelFields}
	// Associations have their own setters. Supplying them during resource
	// creation would force inherit/enabled regardless of the desired values.
	d.Fields["ai_providers"] = Field{Kind: "set", API: "aiProviders", Element: &provider, Default: []any{}, Description: "Authoritative AI provider attachments for inference mode. Providers are configured before model selections."}
	d.Fields["ai_models"] = Field{Kind: "set", API: "aiModels", Element: &model, Default: []any{}, Description: "Authoritative allow/block selections from attached providers in select mode."}
}

func networkValidateAI(v Values) error {
	providers := networkRows(v["ai_providers"])
	models := networkRows(v["ai_models"])
	if Text(v["mode"]) != "inference" && (len(providers) > 0 || len(models) > 0) {
		return fmt.Errorf("ai_providers and ai_models require inference mode")
	}
	seen := map[string]bool{}
	selectMode := false
	for _, provider := range providers {
		id := Text(provider["provider_id"])
		if seen[id] {
			return fmt.Errorf("ai_providers must contain each provider_id only once")
		}
		seen[id] = true
		if Text(provider["access_mode"]) == "select" {
			selectMode = true
		}
	}
	if len(models) > 0 && !selectMode {
		return fmt.Errorf("ai_models requires at least one ai_providers attachment with access_mode = select")
	}
	return nil
}

func networkWriteAI(ctx context.Context, c *client.Client, d Definition, plan, prior Values, private bool) error {
	if Text(plan["mode"]) != "inference" {
		return nil
	}
	base := "/resource/{id}"
	if private {
		base = "/site-resource/{id}"
	}
	changedProviders := prior == nil || networkChanged(plan, prior, "ai_providers")
	if changedProviders {
		body, err := networkPayload(d, plan, "ai_providers")
		if err != nil {
			return err
		}
		providers := body["aiProviders"]
		if providers == nil {
			providers = []any{}
		}
		if err = networkSend(ctx, c, "POST", base+"/ai-providers", plan, Object{"providers": providers}); err != nil {
			return err
		}
	}
	selectMode := false
	for _, provider := range networkRows(plan["ai_providers"]) {
		if Text(provider["access_mode"]) == "select" {
			selectMode = true
		}
	}
	// Upstream rejects even an empty model setter without a select-mode
	// provider. Provider replacement itself prunes orphaned selections.
	if selectMode && (changedProviders || networkChanged(plan, prior, "ai_models")) {
		body, err := networkPayload(d, plan, "ai_models")
		if err != nil {
			return err
		}
		models := body["aiModels"]
		if models == nil {
			models = []any{}
		}
		if err = networkSend(ctx, c, "POST", base+"/ai-models", plan, Object{"models": models}); err != nil {
			return err
		}
	}
	return nil
}

func networkReadAI(ctx context.Context, c *client.Client, v Values, o Object, private bool) error {
	if Text(o["mode"]) != "inference" {
		o["aiProviders"] = []any{}
		o["aiModels"] = []any{}
		return nil
	}
	base := "/resource/{id}"
	if private {
		base = "/site-resource/{id}"
	}
	providers, err := networkGetRows(ctx, c, base+"/ai-providers", v, "providers")
	if err != nil {
		return err
	}
	models, err := networkGetRows(ctx, c, base+"/ai-models", v, "models")
	if err != nil {
		return err
	}
	o["aiProviders"] = providers
	o["aiModels"] = models
	return nil
}

func networkSend(ctx context.Context, c *client.Client, method, path string, v Values, body any) error {
	p, err := Path(path, v)
	if err != nil {
		return err
	}
	return c.Do(ctx, method, p, body, nil)
}
