package main

import (
	"flag"
	"github.com/chang196700/terraform-provider-pangolin/internal/provider"
	"log"
)

func main() {
	check := flag.Bool("check", false, "Fail if generated documentation is stale")
	flag.Parse()
	if err := provider.WriteDocumentation(".", *check); err != nil {
		log.Fatal(err)
	}
}
