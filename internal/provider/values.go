package provider

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type unknownValue struct{}

var unknown = unknownValue{}

func isUnknown(v any) bool { _, ok := v.(unknownValue); return ok }
func hasUnknown(v any) bool {
	if isUnknown(v) {
		return true
	}
	switch x := v.(type) {
	case map[string]any:
		for _, e := range x {
			if hasUnknown(e) {
				return true
			}
		}
	case []any:
		for _, e := range x {
			if hasUnknown(e) {
				return true
			}
		}
	}
	return false
}

func native(v attr.Value) any {
	if v.IsNull() {
		return nil
	}
	if v.IsUnknown() {
		return unknown
	}
	switch x := v.(type) {
	case types.String:
		return x.ValueString()
	case types.Bool:
		return x.ValueBool()
	case types.Int64:
		return x.ValueInt64()
	case types.Float64:
		return x.ValueFloat64()
	case types.Object:
		out := Values{}
		for k, a := range x.Attributes() {
			out[k] = native(a)
		}
		return out
	case types.List:
		out := []any{}
		for _, a := range x.Elements() {
			out = append(out, native(a))
		}
		return out
	case types.Set:
		out := []any{}
		for _, a := range x.Elements() {
			out = append(out, native(a))
		}
		return out
	case types.Map:
		out := Values{}
		for k, a := range x.Elements() {
			out[k] = native(a)
		}
		return out
	}
	panic("unknown Terraform value type")
}
func nativeValue(f Field, v any) (attr.Value, error) {
	t := fieldType(f)
	if v == nil {
		return nullValue(t), nil
	}
	if isUnknown(v) {
		return unknownAttr(t), nil
	}
	switch f.Kind {
	case "string", "id":
		return types.StringValue(Text(v)), nil
	case "bool":
		x, ok := v.(bool)
		if !ok {
			return nil, fmt.Errorf("expected boolean")
		}
		return types.BoolValue(x), nil
	case "int":
		n, e := strconv.ParseInt(Text(v), 10, 64)
		if e != nil {
			return nil, fmt.Errorf("expected integer")
		}
		return types.Int64Value(n), nil
	case "number":
		n, e := strconv.ParseFloat(Text(v), 64)
		if e != nil {
			return nil, fmt.Errorf("expected number")
		}
		return types.Float64Value(n), nil
	case "object":
		m, ok := v.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("expected object")
		}
		a := map[string]attr.Value{}
		for k, c := range f.Fields {
			value, e := nativeValue(c, m[k])
			if e != nil {
				return nil, fmt.Errorf("%s: %w", k, e)
			}
			a[k] = value
		}
		o, d := types.ObjectValue(fieldTypes(f.Fields), a)
		if d.HasError() {
			return nil, fmt.Errorf("invalid object: %v", d)
		}
		return o, nil
	case "list", "set":
		rv := reflect.ValueOf(v)
		if rv.Kind() != reflect.Slice {
			return nil, fmt.Errorf("expected collection")
		}
		a := []attr.Value{}
		for i := 0; i < rv.Len(); i++ {
			value, e := nativeValue(element(f), rv.Index(i).Interface())
			if e != nil {
				return nil, e
			}
			a = append(a, value)
		}
		if f.Kind == "set" {
			o, d := types.SetValue(fieldType(element(f)), a)
			if d.HasError() {
				return nil, fmt.Errorf("invalid set")
			}
			return o, nil
		}
		o, d := types.ListValue(fieldType(element(f)), a)
		if d.HasError() {
			return nil, fmt.Errorf("invalid list")
		}
		return o, nil
	case "map":
		m, ok := v.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("expected map")
		}
		a := map[string]attr.Value{}
		for k, x := range m {
			value, e := nativeValue(element(f), x)
			if e != nil {
				return nil, e
			}
			a[k] = value
		}
		o, d := types.MapValue(fieldType(element(f)), a)
		if d.HasError() {
			return nil, fmt.Errorf("invalid map")
		}
		return o, nil
	}
	return nil, fmt.Errorf("unknown field kind")
}
func nullValue(t attr.Type) attr.Value {
	switch x := t.(type) {
	case basetypes.StringType:
		return types.StringNull()
	case basetypes.BoolType:
		return types.BoolNull()
	case basetypes.Int64Type:
		return types.Int64Null()
	case basetypes.Float64Type:
		return types.Float64Null()
	case types.ObjectType:
		return types.ObjectNull(x.AttrTypes)
	case types.ListType:
		return types.ListNull(x.ElemType)
	case types.SetType:
		return types.SetNull(x.ElemType)
	case types.MapType:
		return types.MapNull(x.ElemType)
	}
	panic("unknown type")
}
func unknownAttr(t attr.Type) attr.Value {
	switch x := t.(type) {
	case basetypes.StringType:
		return types.StringUnknown()
	case basetypes.BoolType:
		return types.BoolUnknown()
	case basetypes.Int64Type:
		return types.Int64Unknown()
	case basetypes.Float64Type:
		return types.Float64Unknown()
	case types.ObjectType:
		return types.ObjectUnknown(x.AttrTypes)
	case types.ListType:
		return types.ListUnknown(x.ElemType)
	case types.SetType:
		return types.SetUnknown(x.ElemType)
	case types.MapType:
		return types.MapUnknown(x.ElemType)
	}
	panic("unknown type")
}

func objectValue(fields map[string]Field, v Values) (types.Object, error) {
	a := map[string]attr.Value{}
	for k, f := range fields {
		value, e := nativeValue(f, v[k])
		if e != nil {
			return types.Object{}, fmt.Errorf("%s: %w", k, e)
		}
		a[k] = value
	}
	o, d := types.ObjectValue(fieldTypes(fields), a)
	if d.HasError() {
		return o, fmt.Errorf("invalid state: %v", d)
	}
	return o, nil
}
func copyValues(v Values) Values {
	out := Values{}
	for k, x := range v {
		out[k] = x
	}
	return out
}
func applyDefaults(fields map[string]Field, v Values) {
	for k, f := range fields {
		if (v[k] == nil || isUnknown(v[k])) && f.Default != nil {
			v[k] = f.Default
		}
		if m, ok := v[k].(map[string]any); ok && f.Kind == "object" {
			applyDefaults(f.Fields, m)
		}
	}
}

func Body(d Definition, v Values, create bool) (Object, error) {
	out := Object{}
	for k, f := range d.Fields {
		if (create && !f.Create) || (!create && !f.Update) || f.API == "" {
			continue
		}
		x := v[k]
		if isUnknown(x) {
			return nil, fmt.Errorf("%s is unknown at apply", k)
		}
		if x == nil {
			if !create && f.Nullable {
				out[f.API] = nil
			}
			continue
		}
		value, e := wireValue(f, x)
		if e != nil {
			return nil, fmt.Errorf("%s: %w", k, e)
		}
		out[f.API] = value
	}
	return out, nil
}
func wireValue(f Field, v any) (any, error) {
	if v == nil {
		return nil, nil
	}
	if isUnknown(v) {
		return nil, fmt.Errorf("unknown value")
	}
	switch f.Kind {
	case "id":
		n, e := strconv.ParseInt(Text(v), 10, 64)
		if e != nil || n <= 0 {
			return nil, fmt.Errorf("ID must be a positive decimal integer")
		}
		return n, nil
	case "object":
		m, ok := v.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("expected object")
		}
		o := Object{}
		for k, c := range f.Fields {
			if m[k] == nil && !c.Nullable {
				continue
			}
			x, e := wireValue(c, m[k])
			if e != nil {
				return nil, e
			}
			name := c.API
			if name == "" {
				name = k
			}
			o[name] = x
		}
		return o, nil
	case "list", "set":
		rv := reflect.ValueOf(v)
		if rv.Kind() != reflect.Slice {
			return nil, fmt.Errorf("expected collection")
		}
		out := []any{}
		for i := 0; i < rv.Len(); i++ {
			x, e := wireValue(element(f), rv.Index(i).Interface())
			if e != nil {
				return nil, e
			}
			out = append(out, x)
		}
		return out, nil
	case "map":
		m, ok := v.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("expected map")
		}
		out := Object{}
		for k, x := range m {
			a, e := wireValue(element(f), x)
			if e != nil {
				return nil, e
			}
			out[k] = a
		}
		return out, nil
	default:
		return v, nil
	}
}
func fromWire(f Field, v any) (any, error) {
	if v == nil {
		return nil, nil
	}
	switch f.Kind {
	case "id", "string":
		return Text(v), nil
	case "object":
		m, ok := v.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("expected API object")
		}
		out := Values{}
		for k, c := range f.Fields {
			name := c.ReadAPI
			if name == "" {
				name = c.API
			}
			if name == "" {
				name = k
			}
			value, e := fromWire(c, m[name])
			if e != nil {
				return nil, e
			}
			out[k] = value
		}
		return out, nil
	case "list", "set":
		rv := reflect.ValueOf(v)
		if rv.Kind() != reflect.Slice {
			return nil, fmt.Errorf("expected API collection")
		}
		out := []any{}
		for i := 0; i < rv.Len(); i++ {
			x, e := fromWire(element(f), rv.Index(i).Interface())
			if e != nil {
				return nil, e
			}
			out = append(out, x)
		}
		return out, nil
	default:
		return v, nil
	}
}
func mergeRemote(d Definition, v Values, o Object) error {
	if id, ok := o[d.IDField]; ok && id != nil {
		v["id"] = Text(id)
	}
	for k, f := range d.Fields {
		if k == "id" {
			continue
		}
		name := f.ReadAPI
		if name == "" {
			name = f.API
		}
		raw, exists := o[name]
		if !exists || (raw == nil && f.Preserve) {
			if isUnknown(v[k]) {
				v[k] = nil
			}
			continue
		}
		value, e := fromWire(f, raw)
		if e != nil {
			return fmt.Errorf("%s: %w", k, e)
		}
		v[k] = value
	}
	return nil
}
func validateFields(fields map[string]Field, v Values) error {
	for k, f := range fields {
		x := v[k]
		if x == nil || isUnknown(x) {
			continue
		}
		if len(f.Enum) > 0 {
			ok := false
			for _, s := range f.Enum {
				if Text(x) == s {
					ok = true
				}
			}
			if !ok {
				return fmt.Errorf("%s must be one of %s", k, strings.Join(f.Enum, ", "))
			}
		}
		if f.Kind == "id" {
			if _, e := wireValue(f, x); e != nil {
				return fmt.Errorf("%s: %w", k, e)
			}
		}
		if f.Kind == "string" && (f.MinLength > 0 || f.MaxLength > 0) {
			n := len([]rune(Text(x)))
			if n < f.MinLength || (f.MaxLength > 0 && n > f.MaxLength) {
				return fmt.Errorf("%s has invalid length", k)
			}
		}
		if f.Min != nil || f.Max != nil {
			n, e := strconv.ParseFloat(Text(x), 64)
			if e == nil && ((f.Min != nil && n < *f.Min) || (f.Max != nil && n > *f.Max)) {
				return fmt.Errorf("%s is outside the permitted range", k)
			}
		}
		if f.Kind == "object" {
			if m, ok := x.(map[string]any); ok {
				if e := validateFields(f.Fields, m); e != nil {
					return fmt.Errorf("%s.%w", k, e)
				}
			}
		}
	}
	return nil
}

// decodeJSON is used by contract fixtures and typed arbitrary configuration fields.
func decodeJSON(raw []byte) (Object, error) {
	var o Object
	dec := json.NewDecoder(strings.NewReader(string(raw)))
	dec.UseNumber()
	err := dec.Decode(&o)
	return o, err
}
