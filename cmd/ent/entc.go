//go:build ignore

package main

import (
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"log"
)

func main() {
	err := entc.Generate("../../internal/ent/schema", &gen.Config{
		Hooks: []gen.Hook{
			func(next gen.Generator) gen.Generator {
				return gen.GenerateFunc(func(g *gen.Graph) error {
					for _, node := range g.Nodes {
						for _, field := range node.Fields {
							// ...
							field.StorageKey() // nothing here, just for example
						}
					}
					return next.Generate(g)
				})
			},
		},
	})
	if err != nil {
		log.Fatalf("running ent codegen: %v", err)
	}
}
