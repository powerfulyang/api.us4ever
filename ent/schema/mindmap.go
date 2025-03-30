// Code generated by entimport, DO NOT EDIT.

package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Mindmap struct {
	ent.Schema
}

func (Mindmap) Fields() []ent.Field {
	return []ent.Field{field.String("id"), field.String("title"), field.JSON("content", struct{}{}), field.String("summary"), field.Bool("isPublic"), field.JSON("tags", struct{}{}), field.Int32("views"), field.Int32("likes"), field.JSON("extraData", struct{}{}), field.String("category"), field.String("ownerId").Optional(), field.Time("createdAt"), field.Time("updatedAt")}

}
func (Mindmap) Edges() []ent.Edge {
	return []ent.Edge{edge.From("user", User.Type).Ref("mindmaps").Unique().Field("ownerId")}
}
func (Mindmap) Annotations() []schema.Annotation {
	return nil
}
