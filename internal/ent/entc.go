//go:build ignore

package ent

import (
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"log"
)

func main() {
	err := entc.Generate("./schema", &gen.Config{
		Hooks: []gen.Hook{
			func(next gen.Generator) gen.Generator {
				return gen.GenerateFunc(func(g *gen.Graph) error {
					for _, node := range g.Nodes {
						for _, field := range node.Fields {
							// ...
							field.StorageKey() // nothing here
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
