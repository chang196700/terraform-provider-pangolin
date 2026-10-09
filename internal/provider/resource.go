package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/chang196700/terraform-provider-pangolin/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type managedResource struct {
	definition Definition
	client     *client.Client
}

var _ resource.ResourceWithConfigure = (*managedResource)(nil)
var _ resource.ResourceWithImportState = (*managedResource)(nil)
var _ resource.ResourceWithValidateConfig = (*managedResource)(nil)

func (r *managedResource) Metadata(_ context.Context, req resource.MetadataRequest, res *resource.MetadataResponse) {
	res.TypeName = req.ProviderTypeName + "_" + r.definition.Name
}
func (r *managedResource) Schema(_ context.Context, _ resource.SchemaRequest, res *resource.SchemaResponse) {
	res.Schema = resourceSchema(r.definition)
}
func (r *managedResource) Configure(_ context.Context, req resource.ConfigureRequest, res *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		res.Diagnostics.AddError("Invalid provider client", "Provider returned an unexpected client type.")
		return
	}
	r.client = c
}

func (r *managedResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, res *resource.ValidateConfigResponse) {
	var m types.Object
	res.Diagnostics.Append(req.Config.Get(ctx, &m)...)
	if res.Diagnostics.HasError() || m.IsNull() {
		return
	}
	v := native(m).(map[string]any)
	applyDefaults(r.definition.Fields, v)
	if hasUnknown(v) {
		return
	}
	if err := validateFields(r.definition.Fields, v); err != nil {
		res.Diagnostics.AddError("Invalid configuration", err.Error())
		return
	}
	if r.definition.Validate != nil {
		if err := r.definition.Validate(v); err != nil {
			res.Diagnostics.AddError("Invalid configuration", err.Error())
		}
	}
}
func readPlan(ctx context.Context, p tfsdk.Plan, d *diag.Diagnostics) Values {
	var m types.Object
	d.Append(p.Get(ctx, &m)...)
	if d.HasError() {
		return nil
	}
	return native(m).(map[string]any)
}
func readState(ctx context.Context, s tfsdk.State, d *diag.Diagnostics) Values {
	var m types.Object
	d.Append(s.Get(ctx, &m)...)
	if d.HasError() || m.IsNull() {
		return nil
	}
	return native(m).(map[string]any)
}
func saveState(ctx context.Context, s *tfsdk.State, fields map[string]Field, v Values, diags *diag.Diagnostics) {
	for k, x := range v {
		if isUnknown(x) {
			v[k] = nil
		}
	}
	obj, err := objectValue(fields, v)
	if err != nil {
		diags.AddError("Cannot encode resource state", err.Error())
		return
	}
	diags.Append(s.Set(ctx, obj)...)
}
func (r *managedResource) check(diags *diag.Diagnostics) bool {
	if r.client == nil {
		diags.AddError("Provider not configured", "A known endpoint and API key are required.")
		return false
	}
	return true
}

func (r *managedResource) Create(ctx context.Context, req resource.CreateRequest, res *resource.CreateResponse) {
	if !r.check(&res.Diagnostics) {
		return
	}
	d := r.definition
	v := readPlan(ctx, req.Plan, &res.Diagnostics)
	if res.Diagnostics.HasError() {
		return
	}
	applyDefaults(d.Fields, v)
	// Computed attributes are unknown in the initial plan; they are not request inputs.
	for k, f := range d.Fields {
		if !f.Required && isUnknown(v[k]) {
			v[k] = nil
		}
	}
	if err := validateFields(d.Fields, v); err != nil {
		res.Diagnostics.AddError("Invalid configuration", err.Error())
		return
	}
	if d.Validate != nil {
		if err := d.Validate(v); err != nil {
			res.Diagnostics.AddError("Invalid configuration", err.Error())
			return
		}
	}
	var o Object
	var err error
	if d.CreateFunc != nil {
		o, err = d.CreateFunc(ctx, r.client, v)
	} else {
		var body Object
		body, err = Body(d, v, true)
		if err == nil {
			o, err = Call(ctx, r.client, d.Create, v, body)
		}
	}
	if err != nil {
		if o != nil && o[d.IDField] != nil {
			v["id"] = Text(o[d.IDField])
			initial := copyValues(v)
			if mergeErr := mergeRemote(d, initial, o); mergeErr == nil {
				saveState(ctx, &res.State, d.Fields, initial, &res.Diagnostics)
			} else {
				saveState(ctx, &res.State, d.Fields, v, &res.Diagnostics)
			}
		}
		res.Diagnostics.AddError("Create failed", err.Error())
		return
	}
	if o[d.IDField] != nil {
		v["id"] = Text(o[d.IDField])
	}
	if v["id"] == nil || Text(v["id"]) == "" {
		res.Diagnostics.AddError("Create returned no identity", "The API created an object without a usable "+d.IDField+". Recover/import it before retrying.")
		return
	}
	initial := copyValues(v)
	if err = mergeRemote(d, initial, o); err != nil {
		saveState(ctx, &res.State, d.Fields, v, &res.Diagnostics)
		res.Diagnostics.AddError("Create response invalid", err.Error())
		return
	}
	saveState(ctx, &res.State, d.Fields, initial, &res.Diagnostics)
	if res.Diagnostics.HasError() {
		return
	}
	if d.AfterWrite != nil {
		if err = d.AfterWrite(ctx, r.client, v, nil); err != nil {
			res.Diagnostics.AddError("Object created; configuration incomplete", err.Error()+" The resource ID has been saved. Apply again to reconcile.")
			return
		}
	}
	// Preserve one-time credentials from Create while using the original desired
	// configuration for post-create setters.
	for k, f := range d.Fields {
		if v[k] == nil && (f.Preserve || f.Computed) {
			v[k] = initial[k]
		}
	}
	if err = r.refresh(ctx, v); err != nil {
		res.Diagnostics.AddError("Object created; refresh failed", err.Error())
		return
	}
	saveState(ctx, &res.State, d.Fields, v, &res.Diagnostics)
}
func (r *managedResource) Read(ctx context.Context, req resource.ReadRequest, res *resource.ReadResponse) {
	if !r.check(&res.Diagnostics) {
		return
	}
	v := readState(ctx, req.State, &res.Diagnostics)
	if res.Diagnostics.HasError() {
		return
	}
	err := r.refresh(ctx, v)
	if client.IsNotFound(err) {
		res.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		res.Diagnostics.AddError("Read failed", err.Error())
		return
	}
	saveState(ctx, &res.State, r.definition.Fields, v, &res.Diagnostics)
}
func (r *managedResource) refresh(ctx context.Context, v Values) error {
	d := r.definition
	o, err := readRemote(ctx, r.client, d, v)
	if err != nil {
		return err
	}
	if d.EnrichRead != nil {
		if err = d.EnrichRead(ctx, r.client, v, o); err != nil {
			// The primary object exists. A missing sub-route or failed relation
			// read must never masquerade as deletion of the managed object.
			return fmt.Errorf("resource exists but related configuration could not be read: %v", err)
		}
	}
	return mergeRemote(d, v, o)
}
func readRemote(ctx context.Context, c *client.Client, d Definition, v Values) (Object, error) {
	if d.ReadFunc != nil {
		return d.ReadFunc(ctx, c, v)
	}
	if d.Read.Path != "" {
		return Call(ctx, c, d.Read, v, nil)
	}
	if d.List.Path == "" {
		return nil, fmt.Errorf("%s has no read operation", d.Name)
	}
	p, err := Path(d.List.Path, v)
	if err != nil {
		return nil, err
	}
	rows, err := c.List(ctx, p, d.ListKey, d.Pagination)
	if err != nil {
		return nil, err
	}
	for _, o := range rows {
		if d.ListFilter != nil && !d.ListFilter(v, o) {
			continue
		}
		if Text(o[d.IDField]) == Text(v["id"]) {
			return o, nil
		}
	}
	return nil, client.NotFound(p)
}
func (r *managedResource) Update(ctx context.Context, req resource.UpdateRequest, res *resource.UpdateResponse) {
	if !r.check(&res.Diagnostics) {
		return
	}
	d := r.definition
	v := readPlan(ctx, req.Plan, &res.Diagnostics)
	prior := readState(ctx, req.State, &res.Diagnostics)
	if res.Diagnostics.HasError() {
		return
	}
	v["id"] = prior["id"]
	applyDefaults(d.Fields, v)
	for k, f := range d.Fields {
		if isUnknown(v[k]) {
			if f.Preserve || !f.Create && !f.Update {
				v[k] = prior[k]
			} else {
				v[k] = prior[k]
			}
		}
	}
	var err error
	err = validateFields(d.Fields, v)
	if err == nil && d.Validate != nil {
		err = d.Validate(v)
	}
	if err == nil {
		if d.UpdateFunc != nil {
			err = d.UpdateFunc(ctx, r.client, v, prior)
		} else if d.Update.Path != "" {
			var body Object
			body, err = Body(d, v, false)
			if err == nil {
				p, e := Path(d.Update.Path, v)
				if e != nil {
					err = e
				} else {
					err = r.client.Do(ctx, d.Update.Method, p, body, nil)
				}
			}
		}
	}
	if err == nil && d.AfterWrite != nil {
		err = d.AfterWrite(ctx, r.client, v, prior)
	}
	if err != nil {
		actual := copyValues(prior)
		if refreshErr := r.refresh(ctx, actual); refreshErr == nil {
			saveState(ctx, &res.State, d.Fields, actual, &res.Diagnostics)
		}
		res.Diagnostics.AddError("Update incomplete", err.Error()+" The existing resource identity is retained.")
		return
	}
	if err = r.refresh(ctx, v); err != nil {
		saveState(ctx, &res.State, d.Fields, v, &res.Diagnostics)
		res.Diagnostics.AddError("Update succeeded; refresh failed", err.Error())
		return
	}
	saveState(ctx, &res.State, d.Fields, v, &res.Diagnostics)
}
func (r *managedResource) Delete(ctx context.Context, req resource.DeleteRequest, res *resource.DeleteResponse) {
	if !r.check(&res.Diagnostics) {
		return
	}
	d := r.definition
	v := readState(ctx, req.State, &res.Diagnostics)
	if res.Diagnostics.HasError() {
		return
	}
	var err error
	if d.DeleteFunc != nil {
		err = d.DeleteFunc(ctx, r.client, v)
	} else {
		var p string
		p, err = Path(d.Delete.Path, v)
		if err == nil {
			err = r.client.Do(ctx, d.Delete.Method, p, nil, nil)
		}
	}
	if err != nil && !client.IsNotFound(err) {
		res.Diagnostics.AddError("Delete failed", err.Error())
	}
}
func (r *managedResource) ImportState(ctx context.Context, req resource.ImportStateRequest, res *resource.ImportStateResponse) {
	d := r.definition
	var v Values
	var err error
	if d.ImportFunc != nil {
		v, err = d.ImportFunc(req.ID)
	} else {
		parts := strings.Split(req.ID, "/")
		if len(parts) != len(d.ImportFields) {
			err = fmt.Errorf("expected import ID %s", strings.Join(d.ImportFields, "/"))
		} else {
			v = Values{}
			for i, k := range d.ImportFields {
				if parts[i] == "" {
					err = fmt.Errorf("empty import ID segment %s", k)
					break
				}
				v[k] = parts[i]
			}
		}
	}
	if err != nil {
		res.Diagnostics.AddError("Invalid import ID", err.Error())
		return
	}
	for k, value := range v {
		res.Diagnostics.Append(res.State.SetAttribute(ctx, path.Root(k), value)...)
	}
}
