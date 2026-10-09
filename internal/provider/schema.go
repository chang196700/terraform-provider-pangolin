package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func fieldType(f Field) attr.Type {
	switch f.Kind {
	case "string", "id":
		return types.StringType
	case "bool":
		return types.BoolType
	case "int":
		return types.Int64Type
	case "number":
		return types.Float64Type
	case "object":
		return types.ObjectType{AttrTypes: fieldTypes(f.Fields)}
	case "list":
		return types.ListType{ElemType: fieldType(element(f))}
	case "set":
		return types.SetType{ElemType: fieldType(element(f))}
	case "map":
		return types.MapType{ElemType: fieldType(element(f))}
	default:
		panic("unsupported field kind: " + f.Kind)
	}
}
func element(f Field) Field {
	if f.Element == nil {
		return String("")
	}
	return *f.Element
}
func fieldTypes(fields map[string]Field) map[string]attr.Type {
	out := map[string]attr.Type{}
	for k, f := range fields {
		out[k] = fieldType(f)
	}
	return out
}
func resourceSchema(d Definition) schema.Schema {
	return schema.Schema{Description: d.Description, Attributes: resourceAttributes(d.Fields)}
}
func resourceAttributes(fields map[string]Field) map[string]schema.Attribute {
	out := map[string]schema.Attribute{}
	for k, f := range fields {
		out[k] = resourceAttribute(f)
	}
	return out
}
func resourceAttribute(f Field) schema.Attribute {
	required := f.Required
	computed := f.Computed || (!required && !f.Nullable && !f.OptionalOnly)
	optional := !required && !f.Computed
	desc := f.Description
	if desc == "" {
		desc = "Pangolin " + f.API + "."
	}
	if f.Sensitive {
		desc += " Sensitive; stored in Terraform state."
	}
	if f.Replace {
		desc += " Changes replace the resource."
	}
	var def attr.Value
	if f.Default != nil && !required {
		def, _ = nativeValue(f, f.Default)
		computed = true
	}
	switch f.Kind {
	case "string", "id":
		a := schema.StringAttribute{Required: required, Optional: optional, Computed: computed, Sensitive: f.Sensitive, Description: desc}
		if def != nil {
			a.Default = stringdefault.StaticString(def.(types.String).ValueString())
		}
		if f.Replace {
			a.PlanModifiers = []planmodifier.String{stringplanmodifier.RequiresReplace()}
		}
		return a
	case "bool":
		a := schema.BoolAttribute{Required: required, Optional: optional, Computed: computed, Sensitive: f.Sensitive, Description: desc}
		if def != nil {
			a.Default = booldefault.StaticBool(def.(types.Bool).ValueBool())
		}
		if f.Replace {
			a.PlanModifiers = []planmodifier.Bool{boolplanmodifier.RequiresReplace()}
		}
		return a
	case "int":
		a := schema.Int64Attribute{Required: required, Optional: optional, Computed: computed, Sensitive: f.Sensitive, Description: desc}
		if def != nil {
			a.Default = int64default.StaticInt64(def.(types.Int64).ValueInt64())
		}
		if f.Replace {
			a.PlanModifiers = []planmodifier.Int64{int64planmodifier.RequiresReplace()}
		}
		return a
	case "number":
		a := schema.Float64Attribute{Required: required, Optional: optional, Computed: computed, Sensitive: f.Sensitive, Description: desc}
		if def != nil {
			a.Default = float64default.StaticFloat64(def.(types.Float64).ValueFloat64())
		}
		if f.Replace {
			a.PlanModifiers = []planmodifier.Float64{float64planmodifier.RequiresReplace()}
		}
		return a
	case "object":
		a := schema.SingleNestedAttribute{Required: required, Optional: optional, Computed: computed, Sensitive: f.Sensitive, Description: desc, Attributes: resourceAttributes(f.Fields)}
		if def != nil {
			a.Default = objectdefault.StaticValue(def.(types.Object))
		}
		if f.Replace {
			a.PlanModifiers = []planmodifier.Object{objectplanmodifier.RequiresReplace()}
		}
		return a
	case "list":
		if element(f).Kind == "object" {
			a := schema.ListNestedAttribute{Required: required, Optional: optional, Computed: computed, Sensitive: f.Sensitive, Description: desc, NestedObject: schema.NestedAttributeObject{Attributes: resourceAttributes(element(f).Fields)}}
			if def != nil {
				a.Default = listdefault.StaticValue(def.(types.List))
			}
			if f.Replace {
				a.PlanModifiers = []planmodifier.List{listplanmodifier.RequiresReplace()}
			}
			return a
		}
		a := schema.ListAttribute{Required: required, Optional: optional, Computed: computed, Sensitive: f.Sensitive, Description: desc, ElementType: fieldType(element(f))}
		if def != nil {
			a.Default = listdefault.StaticValue(def.(types.List))
		}
		if f.Replace {
			a.PlanModifiers = []planmodifier.List{listplanmodifier.RequiresReplace()}
		}
		return a
	case "set":
		if element(f).Kind == "object" {
			a := schema.SetNestedAttribute{Required: required, Optional: optional, Computed: computed, Sensitive: f.Sensitive, Description: desc, NestedObject: schema.NestedAttributeObject{Attributes: resourceAttributes(element(f).Fields)}}
			if def != nil {
				a.Default = setdefault.StaticValue(def.(types.Set))
			}
			if f.Replace {
				a.PlanModifiers = []planmodifier.Set{setplanmodifier.RequiresReplace()}
			}
			return a
		}
		a := schema.SetAttribute{Required: required, Optional: optional, Computed: computed, Sensitive: f.Sensitive, Description: desc, ElementType: fieldType(element(f))}
		if def != nil {
			a.Default = setdefault.StaticValue(def.(types.Set))
		}
		if f.Replace {
			a.PlanModifiers = []planmodifier.Set{setplanmodifier.RequiresReplace()}
		}
		return a
	case "map":
		a := schema.MapAttribute{Required: required, Optional: optional, Computed: computed, Sensitive: f.Sensitive, Description: desc, ElementType: fieldType(element(f))}
		if def != nil {
			a.Default = mapdefault.StaticValue(def.(types.Map))
		}
		if f.Replace {
			a.PlanModifiers = []planmodifier.Map{mapplanmodifier.RequiresReplace()}
		}
		return a
	}
	panic(fmt.Sprintf("unsupported kind %s", f.Kind))
}

// verifySchema is also used by the documentation/contract checks.
func verifySchema(ctx context.Context, d Definition) error {
	diagnostics := resourceSchema(d).ValidateImplementation(ctx)
	if diagnostics.HasError() {
		return fmt.Errorf("%s: %v", d.Name, diagnostics)
	}
	return nil
}
