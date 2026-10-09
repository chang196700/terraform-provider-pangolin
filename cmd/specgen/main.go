// specgen embeds an immutable JSON transcription of the supplied OpenAPI.
package main

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

func main() {
	raw, err := os.ReadFile("openapi.yaml")
	if err != nil {
		panic(err)
	}
	var doc map[string]any
	if err = yaml.Unmarshal(raw, &doc); err != nil {
		panic(err)
	}
	data, err := json.Marshal(doc)
	if err != nil {
		panic(err)
	}
	if err = os.MkdirAll("internal/contract", 0755); err != nil {
		panic(err)
	}
	if err = os.WriteFile("internal/contract/spec.json", append(data, '\n'), 0644); err != nil {
		panic(err)
	}
	fmt.Println("Transcribed openapi.yaml to internal/contract/spec.json")
}
