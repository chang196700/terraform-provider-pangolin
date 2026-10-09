package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/chang196700/terraform-provider-pangolin/internal/client"
)

// Field describes a typed Terraform attribute and its explicit API mapping.
// Kind is string, id, int, number, bool, object, list, set, or map.
type Field struct {
	Kind, API, ReadAPI, Description                                                          string
	Required, Computed, Sensitive, Replace, Create, Update, Preserve, Nullable, OptionalOnly bool
	Default                                                                                  any
	Fields                                                                                   map[string]Field
	Element                                                                                  *Field
	Enum                                                                                     []string
	Min, Max                                                                                 *float64
	MinLength, MaxLength                                                                     int
}

type Values = map[string]any
type Object = client.Object

type Endpoint struct{ Method, Path, Key string }

// Hooks implement the non-CRUD semantics of each Pangolin entity. Values use
// Terraform snake_case names; Objects use the API's original names.
type Definition struct {
	Name, Description, IDField, ListKey, Pagination string
	Fields                                          map[string]Field
	Create, Read, Update, Delete, List              Endpoint
	ImportFields                                    []string
	Validate                                        func(Values) error
	CreateFunc                                      func(context.Context, *client.Client, Values) (Object, error)
	ReadFunc                                        func(context.Context, *client.Client, Values) (Object, error)
	UpdateFunc                                      func(context.Context, *client.Client, Values, Values) error
	DeleteFunc                                      func(context.Context, *client.Client, Values) error
	AfterWrite                                      func(context.Context, *client.Client, Values, Values) error
	EnrichRead                                      func(context.Context, *client.Client, Values, Object) error
	ImportFunc                                      func(string) (Values, error)
	ListFilter                                      func(Values, Object) bool
	ListInputs                                      []string
	// Additional read-only attributes and list filters for data sources.
	DataFields map[string]Field
	Sources    []string
}

func String(api string) Field { return Field{Kind: "string", API: api, Create: true, Update: true} }
func ID(api string) Field     { f := String(api); f.Kind = "id"; return f }
func Bool(api string) Field   { f := String(api); f.Kind = "bool"; return f }
func Int(api string) Field    { f := String(api); f.Kind = "int"; return f }
func Number(api string) Field { f := String(api); f.Kind = "number"; return f }
func StringSet(api string) Field {
	e := String("")
	return Field{Kind: "set", API: api, Element: &e, Default: []any{}, Create: true, Update: true}
}
func IDSet(api string) Field           { f := StringSet(api); e := ID(""); f.Element = &e; return f }
func Required(f Field) Field           { f.Required = true; return f }
func Immutable(f Field) Field          { f.Replace = true; f.Update = false; return f }
func Sensitive(f Field) Field          { f.Sensitive = true; return f }
func Computed(f Field) Field           { f.Computed = true; f.Create = false; f.Update = false; return f }
func Local(f Field) Field              { f.Create = false; f.Update = false; f.Preserve = true; return f }
func WithDefault(f Field, v any) Field { f.Default = v; return f }
func Preserve(f Field) Field           { f.Preserve = true; return f }
func SetField(d *Definition, name string, fn func(*Field)) {
	f := d.Fields[name]
	fn(&f)
	d.Fields[name] = f
}
func MergeFields(a, b map[string]Field) map[string]Field {
	for k, v := range b {
		a[k] = v
	}
	return a
}

func Text(v any) string {
	if v == nil {
		return ""
	}
	switch x := v.(type) {
	case string:
		return x
	case json.Number:
		return x.String()
	case float64:
		return strconv.FormatFloat(x, 'f', -1, 64)
	default:
		return fmt.Sprint(v)
	}
}
func Strings(v any) []string {
	var out []string
	switch a := v.(type) {
	case []any:
		for _, x := range a {
			out = append(out, Text(x))
		}
	case []string:
		return a
	}
	return out
}
func Path(template string, v Values) (string, error) { return client.ExpandPath(template, v) }
func Call(ctx context.Context, c *client.Client, e Endpoint, v Values, body any) (Object, error) {
	p, err := Path(e.Path, v)
	if err != nil {
		return nil, err
	}
	var data Object
	if err = c.Do(ctx, e.Method, p, body, &data); err != nil {
		return nil, err
	}
	if e.Key != "" {
		o, ok := data[e.Key].(map[string]any)
		if !ok {
			return nil, fmt.Errorf("response is missing %s object", e.Key)
		}
		return o, nil
	}
	return data, nil
}
func Require(v Values, names ...string) error {
	for _, n := range names {
		if v[n] == nil || Text(v[n]) == "" {
			return fmt.Errorf("%s is required", n)
		}
	}
	return nil
}
func ExactlyOne(v Values, names ...string) error {
	n := 0
	for _, k := range names {
		if v[k] != nil && Text(v[k]) != "" {
			n++
		}
	}
	if n != 1 {
		return fmt.Errorf("exactly one of %s must be set", strings.Join(names, ", "))
	}
	return nil
}

// Definitions are registered by subsystem files, keeping lifecycle ownership explicit.
var definitionFactories []func() []Definition

func registerDefinitions(f func() []Definition) { definitionFactories = append(definitionFactories, f) }
func allDefinitions() []Definition {
	var out []Definition
	for _, f := range definitionFactories {
		out = append(out, f()...)
	}
	return out
}
