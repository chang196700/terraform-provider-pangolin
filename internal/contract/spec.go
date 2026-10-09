// Package contract exposes the immutable source contract. Lifecycle decisions
// and corrections are owned by the provider, not inferred from HTTP verbs.
package contract

import (
	"embed"
	"encoding/json"
)

//go:embed spec.json
var files embed.FS

type Schema struct {
	Type             string            `json:"type"`
	Description      string            `json:"description"`
	Nullable         bool              `json:"nullable"`
	Properties       map[string]Schema `json:"properties"`
	Required         []string          `json:"required"`
	Items            *Schema           `json:"items"`
	Enum             []any             `json:"enum"`
	Default          any               `json:"default"`
	AnyOf            []Schema          `json:"anyOf"`
	AllOf            []Schema          `json:"allOf"`
	OneOf            []Schema          `json:"oneOf"`
	Minimum          *float64          `json:"minimum"`
	Maximum          *float64          `json:"maximum"`
	ExclusiveMinimum bool              `json:"exclusiveMinimum"`
	MinLength        int               `json:"minLength"`
	MaxLength        int               `json:"maxLength"`
}
type Operation struct {
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
	Deprecated  bool     `json:"deprecated"`
	RequestBody struct {
		Content map[string]struct {
			Schema Schema `json:"schema"`
		} `json:"content"`
	} `json:"requestBody"`
}
type Document struct {
	Paths map[string]map[string]Operation `json:"paths"`
}

var Spec Document

func init() {
	raw, err := files.ReadFile("spec.json")
	if err != nil {
		panic(err)
	}
	if err = json.Unmarshal(raw, &Spec); err != nil {
		panic(err)
	}
}

// Flatten combines object branches for typed fields; semantic union constraints
// are implemented by resource validation. An anyOf field is required only when
// every branch requires it.
func Flatten(s Schema) Schema {
	if s.Properties == nil {
		s.Properties = map[string]Schema{}
	}
	for _, b := range s.AllOf {
		b = Flatten(b)
		for k, v := range b.Properties {
			s.Properties[k] = v
		}
		s.Required = append(s.Required, b.Required...)
		if s.Type == "" {
			s.Type = b.Type
		}
	}
	branches := append(append([]Schema{}, s.AnyOf...), s.OneOf...)
	counts := map[string]int{}
	for _, b := range branches {
		b = Flatten(b)
		if s.Type == "" {
			s.Type = b.Type
		}
		if b.Nullable {
			s.Nullable = true
		}
		if len(s.Enum) == 0 {
			s.Enum = b.Enum
		}
		for k, v := range b.Properties {
			s.Properties[k] = v
		}
		for _, k := range b.Required {
			counts[k]++
		}
	}
	for k, n := range counts {
		if n == len(branches) {
			s.Required = append(s.Required, k)
		}
	}
	return s
}
