// Code generated by entimport, DO NOT EDIT.

package schema

import (
	"encoding/json"

	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Video struct {
	ent.Schema
}

func (Video) Fields() []ent.Field {
	return []ent.Field{field.String("id").StorageKey("id"), field.String("hash").StorageKey("hash"), field.Int("size").StorageKey("size"), field.Bool("isPublic").StorageKey("isPublic"), field.String("posterId").Optional().StorageKey("posterId"), field.String("fileId").Optional().StorageKey("fileId"), field.String("uploadedBy").Optional().StorageKey("uploadedBy"), field.Time("createdAt").StorageKey("createdAt"), field.Time("updatedAt").StorageKey("updatedAt"), field.Int32("duration").StorageKey("duration"), field.String("name").StorageKey("name"), field.String("type").StorageKey("type"), field.JSON("extraData", json.RawMessage{}).StorageKey("extraData"), field.String("category").StorageKey("category")}

}
func (Video) Edges() []ent.Edge {
	return []ent.Edge{edge.To("moment_videos", MomentVideo.Type), edge.From("file", File.Type).Ref("Video_file").Unique().Field("fileId"), edge.From("poster", File.Type).Ref("Video_poster").Unique().Field("posterId"), edge.From("user", User.Type).Ref("videos").Unique().Field("uploadedBy")}
}
func (Video) Annotations() []schema.Annotation {
	return nil
}
