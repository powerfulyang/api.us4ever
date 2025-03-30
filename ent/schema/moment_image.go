// Code generated by entimport, DO NOT EDIT.

package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type MomentImage struct {
	ent.Schema
}

func (MomentImage) Fields() []ent.Field {
	return []ent.Field{field.Uint("id").SchemaType(map[string]string{"postgres": "serial"}), field.String("imageId").Optional(), field.String("momentId").Optional(), field.Int32("sort"), field.Time("createdAt"), field.Time("updatedAt")}
}
func (MomentImage) Edges() []ent.Edge {
	return []ent.Edge{edge.From("image", Image.Type).Ref("moment_images").Unique().Field("imageId"), edge.From("moment", Moment.Type).Ref("moment_images").Unique().Field("momentId")}
}
func (MomentImage) Annotations() []schema.Annotation {
	return nil
}
