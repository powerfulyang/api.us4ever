// Code generated by entimport, DO NOT EDIT.

package schema

import (
	"encoding/json"

	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Moment struct {
	ent.Schema
}

func (Moment) Fields() []ent.Field {
	return []ent.Field{field.String("id"), field.String("content"), field.Bool("isPublic"), field.JSON("tags", json.RawMessage{}), field.Int32("views"), field.Int32("likes"), field.JSON("extraData", json.RawMessage{}), field.String("category"), field.String("ownerId").Optional(), field.Time("createdAt"), field.Time("updatedAt")}

}
func (Moment) Edges() []ent.Edge {
	return []ent.Edge{edge.To("moment_images", MomentImage.Type), edge.To("moment_videos", MomentVideo.Type), edge.From("user", User.Type).Ref("moments").Unique().Field("ownerId")}
}
func (Moment) Annotations() []schema.Annotation {
	return nil
}
