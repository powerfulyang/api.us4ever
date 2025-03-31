// Code generated by entimport, DO NOT EDIT.

package schema

import (
	"encoding/json"

	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type File struct {
	ent.Schema
}

func (File) Fields() []ent.Field {
	return []ent.Field{field.String("id"), field.String("bucketId").Optional(), field.String("name"), field.String("type"), field.String("hash"), field.Int("size"), field.String("path"), field.Bool("isPublic"), field.String("description"), field.JSON("tags", json.RawMessage{}), field.JSON("extraData", json.RawMessage{}), field.String("category"), field.String("uploadedBy").Optional(), field.Time("createdAt"), field.Time("updatedAt")}

}
func (File) Edges() []ent.Edge {
	return []ent.Edge{edge.From("bucket", Bucket.Type).Ref("files").Unique().Field("bucketId"), edge.From("user", User.Type).Ref("files").Unique().Field("uploadedBy"), edge.To("Image_compressed", Image.Type), edge.To("Image_original", Image.Type), edge.To("Image_thumbnail320x", Image.Type), edge.To("Image_thumbnail768x", Image.Type), edge.To("Video_file", Video.Type), edge.To("Video_poster", Video.Type)}
}
func (File) Annotations() []schema.Annotation {
	return nil
}
