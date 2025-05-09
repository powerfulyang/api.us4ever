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
	return []ent.Field{field.Uint("id").StorageKey("id").SchemaType(map[string]string{"postgres": "serial"}), field.String("imageId").Optional().StorageKey("imageId"), field.String("momentId").Optional().StorageKey("momentId"), field.Int32("sort").StorageKey("sort"), field.Time("createdAt").StorageKey("createdAt"), field.Time("updatedAt").StorageKey("updatedAt")}
}
func (MomentImage) Edges() []ent.Edge {
	return []ent.Edge{edge.From("image", Image.Type).Ref("moment_images").Unique().Field("imageId"), edge.From("moment", Moment.Type).Ref("moment_images").Unique().Field("momentId")}
}
func (MomentImage) Annotations() []schema.Annotation {
	return nil
}
