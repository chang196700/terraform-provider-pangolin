package provider

import (
	"github.com/chang196700/terraform-provider-pangolin/internal/contract"
	"regexp"
	"strings"
)

var capitals = regexp.MustCompile(`([a-z0-9])([A-Z])`)
var acronyms = regexp.MustCompile(`([A-Z]+)([A-Z][a-z])`)

func snake(s string) string {
	return strings.ToLower(capitals.ReplaceAllString(acronyms.ReplaceAllString(s, "${1}_${2}"), "${1}_${2}"))
}

// SpecFields supplies typed fields only. The subsystem definitions explicitly
// choose endpoints, identities, state ownership, and operation ordering.
func SpecFields(createPath, createMethod, updatePath, updateMethod string) map[string]Field {
	out := map[string]Field{}
	for i, op := range [][2]string{{createPath, createMethod}, {updatePath, updateMethod}} {
		if op[0] == "" {
			continue
		}
		operation, ok := contract.Spec.Paths[op[0]][strings.ToLower(op[1])]
		if !ok {
			panic("unknown OpenAPI operation: " + op[1] + " " + op[0])
		}
		s := contract.Flatten(operation.RequestBody.Content["application/json"].Schema)
		required := map[string]bool{}
		for _, n := range s.Required {
			required[n] = true
		}
		for name, p := range s.Properties {
			key := snake(name)
			f := fieldFromSpec(name, p)
			f.Required = i == 0 && required[name]
			f.Create = i == 0
			f.Update = i == 1
			if old, ok := out[key]; ok {
				f.Required = old.Required
				f.Create = old.Create
				f.Update = true
				if f.Default == nil {
					f.Default = old.Default
				}
				f.Nullable = f.Nullable || old.Nullable
			}
			out[key] = f
		}
	}
	for k, f := range out {
		if f.Create && !f.Update {
			f.Replace = true
		}
		out[k] = f
	}
	return out
}
func fieldFromSpec(name string, s contract.Schema) Field {
	s = contract.Flatten(s)
	f := Field{Kind: s.Type, API: name, Description: s.Description, Nullable: s.Nullable, Default: s.Default, Min: s.Minimum, Max: s.Maximum, MinLength: s.MinLength, MaxLength: s.MaxLength}
	switch s.Type {
	case "integer":
		f.Kind = "int"
	case "number":
		f.Kind = "number"
	case "boolean":
		f.Kind = "bool"
	case "array":
		f.Kind = "list"
		if s.Items != nil {
			e := fieldFromSpec(strings.TrimSuffix(name, "s"), *s.Items)
			f.Element = &e
		}
	case "object":
		f.Fields = map[string]Field{}
		req := map[string]bool{}
		for _, n := range s.Required {
			req[n] = true
		}
		for k, v := range s.Properties {
			c := fieldFromSpec(k, v)
			c.Required = req[k]
			f.Fields[snake(k)] = c
		}
	}
	if f.Kind == "" {
		f.Kind = "string"
	}
	if (f.Kind == "int" || f.Kind == "number") && (strings.HasSuffix(name, "Id") || name == "id") {
		f.Kind = "id"
	}
	if s.ExclusiveMinimum && f.Min != nil {
		v := *f.Min
		if f.Kind == "int" || f.Kind == "id" {
			v++
		}
		f.Min = &v
	}
	for _, v := range s.Enum {
		f.Enum = append(f.Enum, Text(v))
	}
	lower := strings.ToLower(name)
	if strings.Contains(lower, "secret") || strings.Contains(lower, "password") || strings.Contains(lower, "pincode") || lower == "apikey" || lower == "token" {
		f.Sensitive = true
	}
	return f
}
