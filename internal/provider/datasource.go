package provider

import (
	"context"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/chang196700/terraform-provider-pangolin/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type DataDefinition struct {
	Name, Description   string
	Entity              Definition
	IsList              bool
	Inputs, Outputs     map[string]Field
	Endpoint            Endpoint
	ListKey, Pagination string
	ReadFunc            func(context.Context, *client.Client, Values) (Object, error)
}
type apiDataSource struct {
	definition DataDefinition
	client     *client.Client
}

var pathParameter = regexp.MustCompile(`\{([^}]+)\}`)

func dataDefinitions() []DataDefinition {
	var out []DataDefinition
	for _, d := range definitions() {
		if d.Name == "user_server_admin" {
			continue
		}
		// Domain has explicit data sources in datasource_extra.go with the
		// complete server response projection (including DNS verification
		// state). Do not generate a second `domain`/`domains` data source from
		// the managed resource definition.
		if d.Name == "domain" {
			continue
		}
		if d.Name == "public_resource" {
			managedEnrich := d.EnrichRead
			d.EnrichRead = func(ctx context.Context, c *client.Client, v Values, o Object) error {
				// List rows omit both policy identities. Read details before deciding
				// whether the access object describes an inline policy.
				if _, known := o["resourcePolicyId"]; !known {
					details, err := Call(ctx, c, Endpoint{Method: "GET", Path: "/public-resource/{id}"}, v, nil)
					if err != nil {
						return err
					}
					for key := range o {
						delete(o, key)
					}
					for key, value := range details {
						o[key] = value
					}
				}
				lookup := copyValues(v)
				lookup["access"] = nil
				o["access"] = nil
				mode := Text(o["mode"])
				if o["resourcePolicyId"] == nil && o["defaultResourcePolicyId"] != nil && mode != "tcp" && mode != "udp" {
					// A data source always observes readable inline authorization.
					// Managed imports retain their separate, opt-in ownership rule.
					lookup["access"] = Values{}
				}
				return managedEnrich(ctx, c, lookup, o)
			}
		}
		fields := map[string]Field{}
		for k, f := range d.Fields {
			fields[k] = f
		}
		for k, f := range d.DataFields {
			fields[k] = f
		}
		inputs := map[string]Field{}
		for _, k := range d.ImportFields {
			if f, ok := fields[k]; ok {
				f.Required = true
				f.Computed = false
				inputs[k] = f
			}
		}
		if len(inputs) == 0 {
			inputs["id"] = Required(String(""))
		}
		name := d.Name
		if len(d.Sources) > 0 {
			name = d.Sources[0]
		}
		out = append(out, DataDefinition{Name: name, Description: "Reads an existing " + strings.ReplaceAll(d.Name, "_", " ") + ".", Entity: d, Inputs: inputs, Outputs: fields})
		if d.List.Path != "" {
			filters := map[string]Field{}
			for _, m := range pathParameter.FindAllStringSubmatch(d.List.Path, -1) {
				f, ok := fields[m[1]]
				if !ok {
					f = String("")
				}
				f.Required = true
				f.Computed = false
				filters[m[1]] = f
			}
			for _, key := range d.ListInputs {
				f, ok := fields[key]
				if !ok {
					panic("unknown list input " + d.Name + "." + key)
				}
				f.Required = true
				f.Computed = false
				filters[key] = f
			}
			listFields := map[string]Field{}
			for k, f := range fields {
				if !f.Sensitive {
					listFields[k] = f
				}
			}
			plural := name + "s"
			if strings.HasSuffix(name, "policy") {
				plural = strings.TrimSuffix(name, "y") + "ies"
			}
			if len(d.Sources) > 1 {
				plural = d.Sources[1]
			}
			out = append(out, DataDefinition{Name: plural, Description: "Lists " + strings.ReplaceAll(d.Name, "_", " ") + " objects across all API pages. Secret fields are excluded.", Entity: d, IsList: true, Inputs: filters, Outputs: listFields, Endpoint: d.List, ListKey: d.ListKey, Pagination: d.Pagination})
		}
	}
	out = append(out, extraDataDefinitions()...)
	out = append(out, targetDataDefinitions()...)
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}
func (d DataDefinition) fields() map[string]Field {
	out := map[string]Field{}
	for k, f := range d.Inputs {
		out[k] = f
	}
	if d.IsList {
		e := Field{Kind: "object", Fields: d.Outputs}
		out["items"] = Computed(Field{Kind: "list", Element: &e})
	} else {
		for k, f := range d.Outputs {
			if _, ok := out[k]; !ok {
				out[k] = Computed(f)
			}
		}
	}
	return out
}
func dataSchema(d DataDefinition) schema.Schema {
	attrs := map[string]schema.Attribute{}
	for k, f := range d.fields() {
		_, input := d.Inputs[k]
		attrs[k] = dataAttribute(f, input)
	}
	return schema.Schema{Description: d.Description, Attributes: attrs}
}
func dataAttribute(f Field, input bool) schema.Attribute {
	required, optional, computed := input && f.Required, input && !f.Required, !input
	desc := f.Description
	if desc == "" {
		desc = "Pangolin " + f.API + "."
	}
	switch f.Kind {
	case "id", "string":
		return schema.StringAttribute{Required: required, Optional: optional, Computed: computed, Sensitive: f.Sensitive, Description: desc}
	case "bool":
		return schema.BoolAttribute{Required: required, Optional: optional, Computed: computed, Sensitive: f.Sensitive, Description: desc}
	case "int":
		return schema.Int64Attribute{Required: required, Optional: optional, Computed: computed, Sensitive: f.Sensitive, Description: desc}
	case "number":
		return schema.Float64Attribute{Required: required, Optional: optional, Computed: computed, Sensitive: f.Sensitive, Description: desc}
	case "object":
		a := map[string]schema.Attribute{}
		for k, c := range f.Fields {
			a[k] = dataAttribute(c, false)
		}
		return schema.SingleNestedAttribute{Required: required, Optional: optional, Computed: computed, Sensitive: f.Sensitive, Description: desc, Attributes: a}
	case "list":
		return schema.ListAttribute{Required: required, Optional: optional, Computed: computed, Sensitive: f.Sensitive, Description: desc, ElementType: fieldType(element(f))}
	case "set":
		return schema.SetAttribute{Required: required, Optional: optional, Computed: computed, Sensitive: f.Sensitive, Description: desc, ElementType: fieldType(element(f))}
	case "map":
		return schema.MapAttribute{Required: required, Optional: optional, Computed: computed, Sensitive: f.Sensitive, Description: desc, ElementType: fieldType(element(f))}
	}
	panic("unsupported data source type")
}
func (d *apiDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, res *datasource.MetadataResponse) {
	res.TypeName = req.ProviderTypeName + "_" + d.definition.Name
}
func (d *apiDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, res *datasource.SchemaResponse) {
	res.Schema = dataSchema(d.definition)
}
func (d *apiDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, res *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		res.Diagnostics.AddError("Unexpected client type", "Provider client is unavailable.")
		return
	}
	d.client = c
}
func (d *apiDataSource) Read(ctx context.Context, req datasource.ReadRequest, res *datasource.ReadResponse) {
	if d.client == nil {
		res.Diagnostics.AddError("Provider not configured", "Configure endpoint and api_key first.")
		return
	}
	var model types.Object
	res.Diagnostics.Append(req.Config.Get(ctx, &model)...)
	if res.Diagnostics.HasError() {
		return
	}
	v := native(model).(map[string]any)
	dd := d.definition
	var err error
	if dd.IsList {
		var rows []Object
		if dd.ReadFunc != nil {
			var data Object
			data, err = dd.ReadFunc(ctx, d.client, v)
			if err == nil {
				rows = objectArray(data[dd.ListKey])
				if data[dd.ListKey] == nil {
					err = fmt.Errorf("list response is missing %s", dd.ListKey)
				}
			}
		} else {
			p, e := Path(dd.Endpoint.Path, v)
			err = e
			if err == nil {
				rows, err = d.client.List(ctx, p, dd.ListKey, dd.Pagination)
			}
		}
		if err == nil {
			items := []any{}
			for _, o := range rows {
				if dd.Entity.ListFilter != nil && !dd.Entity.ListFilter(v, o) {
					continue
				}
				item := Values{}
				for k, x := range v {
					item[k] = x
				}
				entity := dd.Entity
				if entity.Name != "" {
					entity.Fields = dd.Outputs
					if o[entity.IDField] != nil {
						item["id"] = Text(o[entity.IDField])
					}
					if entity.EnrichRead != nil {
						if err = entity.EnrichRead(ctx, d.client, item, o); err != nil {
							break
						}
					}
					err = mergeRemote(entity, item, o)
				} else {
					err = projectFields(dd.Outputs, item, o)
				}
				if err != nil {
					break
				}
				items = append(items, item)
			}
			v["items"] = items
		}
	} else {
		var o Object
		if dd.ReadFunc != nil {
			o, err = dd.ReadFunc(ctx, d.client, v)
		} else if dd.Entity.Name != "" {
			entity := dd.Entity
			if v["id"] == nil && entity.IDField == "orgId" {
				v["id"] = v["org_id"]
			}
			o, err = readRemote(ctx, d.client, entity, v)
			if err == nil && entity.EnrichRead != nil {
				err = entity.EnrichRead(ctx, d.client, v, o)
			}
		} else {
			o, err = Call(ctx, d.client, dd.Endpoint, v, nil)
		}
		if err == nil {
			if dd.Entity.Name != "" {
				entity := dd.Entity
				entity.Fields = dd.Outputs
				err = mergeRemote(entity, v, o)
			} else {
				err = projectFields(dd.Outputs, v, o)
			}
		}
	}
	if err != nil {
		res.Diagnostics.AddError("Data source read failed", err.Error())
		return
	}
	for k, x := range v {
		if isUnknown(x) {
			v[k] = nil
		}
	}
	obj, e := objectValue(dd.fields(), v)
	if e != nil {
		res.Diagnostics.AddError("Invalid API response", e.Error())
		return
	}
	res.Diagnostics.Append(res.State.Set(ctx, obj)...)
}
func projectFields(fields map[string]Field, v Values, o Object) error {
	for k, f := range fields {
		name := f.ReadAPI
		if name == "" {
			name = f.API
		}
		value, e := fromWire(f, o[name])
		if e != nil {
			return fmt.Errorf("%s: %w", k, e)
		}
		v[k] = value
	}
	return nil
}

var _ attr.Value = types.Object{}
