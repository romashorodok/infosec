//go:build ignore
// +build ignore

// https://entgo.io/blog/2021/09/10/openapi-generator/
// https://entgo.io/blog/2021/10/11/generating-ent-schemas-from-existing-sql-databases/

package main

import (
	"log"

	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"github.com/masseelch/elk"
	"github.com/masseelch/elk/spec"
	entlocal "github.com/romashorodok/infosec/ent"
)

func main() {
	ex, err := elk.NewExtension(
		elk.GenerateSpec("openapi.json",
			elk.SpecSecuritySchemes(
				map[string]spec.SecurityScheme{
					entlocal.SecurityBearerAuth: {
						Type:   "http",
						Scheme: "bearer",
					},
				},
			),
		),
	)
	if err != nil {
		log.Fatalf("creating elk extension: %v", err)
	}
	err = entc.Generate("./schema", &gen.Config{}, entc.Extensions(ex))
	if err != nil {
		log.Fatalf("running ent codegen: %v", err)
	}
}
