// Code generated by entimport, DO NOT EDIT.

package schema

import (
	"encoding/json"

	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Image struct {
	ent.Schema
}

func (Image) Fields() []ent.Field {
	return []ent.Field{field.String("id"), field.String("name"), field.String("type"), field.Int("size"), field.Int32("width"), field.Int32("height"), field.JSON("exif", json.RawMessage{}), field.String("hash"), field.String("address"), field.Bool("isPublic"), field.String("description"), field.JSON("tags", json.RawMessage{}), field.JSON("extraData", json.RawMessage{}), field.String("category"), field.Bytes("thumbnail_10x"), field.String("thumbnail_320x_id").Optional(), field.String("thumbnail_768x_id").Optional(), field.String("compressed_id").Optional(), field.String("original_id").Optional(), field.String("uploadedBy").Optional(), field.Time("createdAt"), field.Time("updatedAt")}

}
func (Image) Edges() []ent.Edge {
	return []ent.Edge{edge.From("compressed", File.Type).Ref("Image_compressed").Unique().Field("compressed_id"), edge.From("original", File.Type).Ref("Image_original").Unique().Field("original_id"), edge.From("thumbnail320x", File.Type).Ref("Image_thumbnail320x").Unique().Field("thumbnail_320x_id"), edge.From("thumbnail768x", File.Type).Ref("Image_thumbnail768x").Unique().Field("thumbnail_768x_id"), edge.From("user", User.Type).Ref("images").Unique().Field("uploadedBy"), edge.To("moment_images", MomentImage.Type)}
}
func (Image) Annotations() []schema.Annotation {
	return nil
}
