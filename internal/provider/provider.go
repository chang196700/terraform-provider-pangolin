package provider

import (
	"context"
	"os"
	"sort"
	"time"

	"github.com/chang196700/terraform-provider-pangolin/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	fw "github.com/hashicorp/terraform-plugin-framework/provider"
	pschema "github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type pangolinProvider struct{ version string }
type providerModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
	APIKey   types.String `tfsdk:"api_key"`
	Timeout  types.Int64  `tfsdk:"request_timeout"`
}

func New(version string) func() fw.Provider {
	return func() fw.Provider { return &pangolinProvider{version: version} }
}
func (p *pangolinProvider) Metadata(_ context.Context, _ fw.MetadataRequest, r *fw.MetadataResponse) {
	r.TypeName = "pangolin"
	r.Version = p.version
}
func (p *pangolinProvider) Schema(_ context.Context, _ fw.SchemaRequest, r *fw.SchemaResponse) {
	r.Schema = pschema.Schema{Description: "Manage Pangolin Integration API infrastructure. Targets Pangolin 1.24.0. Terraform 1.11 or newer is supported.", Attributes: map[string]pschema.Attribute{
		"endpoint":        pschema.StringAttribute{Optional: true, Description: "Absolute API URL including /v1. Defaults to PANGOLIN_ENDPOINT."},
		"api_key":         pschema.StringAttribute{Optional: true, Sensitive: true, Description: "Complete Integration API bearer credential in apiKeyId.apiKeySecret form. Defaults to PANGOLIN_API_KEY. Some operations require a root key; see the permissions documentation."},
		"request_timeout": pschema.Int64Attribute{Optional: true, Description: "HTTP request timeout in seconds. Defaults to 30."},
	}}
}
func (p *pangolinProvider) Configure(ctx context.Context, req fw.ConfigureRequest, res *fw.ConfigureResponse) {
	var m providerModel
	res.Diagnostics.Append(req.Config.Get(ctx, &m)...)
	if res.Diagnostics.HasError() {
		return
	}
	if m.Endpoint.IsUnknown() || m.APIKey.IsUnknown() || m.Timeout.IsUnknown() {
		res.Diagnostics.AddError("Unknown provider configuration", "endpoint, api_key and request_timeout must be known before managing resources.")
		return
	}
	endpoint, key := m.Endpoint.ValueString(), m.APIKey.ValueString()
	if m.Endpoint.IsNull() {
		endpoint = os.Getenv("PANGOLIN_ENDPOINT")
	}
	if m.APIKey.IsNull() {
		key = os.Getenv("PANGOLIN_API_KEY")
	}
	timeout := int64(30)
	if !m.Timeout.IsNull() {
		timeout = m.Timeout.ValueInt64()
	}
	c, err := client.New(endpoint, key, time.Duration(timeout)*time.Second)
	if err != nil {
		res.Diagnostics.AddError("Invalid provider configuration", err.Error())
		return
	}
	c.UserAgent = "terraform-provider-pangolin/" + p.version
	res.ResourceData = c
	res.DataSourceData = c
}
func definitions() []Definition {
	ds := allDefinitions()
	sort.Slice(ds, func(i, j int) bool { return ds[i].Name < ds[j].Name })
	for i := range ds {
		if ds[i].Fields == nil {
			ds[i].Fields = map[string]Field{}
		}
		ds[i].Fields["id"] = Computed(String(ds[i].IDField))
	}
	return ds
}
func (p *pangolinProvider) Resources(_ context.Context) []func() resource.Resource {
	var out []func() resource.Resource
	for _, d := range definitions() {
		out = append(out, func() resource.Resource { return &managedResource{definition: d} })
	}
	return out
}
func (p *pangolinProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	var out []func() datasource.DataSource
	for _, d := range dataDefinitions() {
		out = append(out, func() datasource.DataSource { return &apiDataSource{definition: d} })
	}
	return out
}
