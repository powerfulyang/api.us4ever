// Code generated by entimport, DO NOT EDIT.

package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Keep struct {
	ent.Schema
}

func (Keep) Fields() []ent.Field {
	return []ent.Field{field.String("id"), field.String("title"), field.String("content"), field.String("summary"), field.Bool("isPublic"), field.JSON("tags", struct{}{}), field.Int32("views"), field.Int32("likes"), field.JSON("extraData", struct{}{}), field.String("category"), field.String("ownerId").Optional(), field.Time("createdAt"), field.Time("updatedAt")}

}
func (Keep) Edges() []ent.Edge {
	return []ent.Edge{edge.From("user", User.Type).Ref("keeps").Unique().Field("ownerId")}
}
func (Keep) Annotations() []schema.Annotation {
	return nil
}
