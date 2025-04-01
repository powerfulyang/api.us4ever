// Code generated by entimport, DO NOT EDIT.

package schema

import (
	"encoding/json"

	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Todo struct {
	ent.Schema
}

func (Todo) Fields() []ent.Field {
	return []ent.Field{field.String("id").StorageKey("id"), field.String("title").StorageKey("title"), field.String("content").Optional().StorageKey("content"), field.Bool("status").StorageKey("status"), field.Int32("priority").StorageKey("priority"), field.Time("dueDate").Optional().StorageKey("dueDate"), field.Bool("isPublic").StorageKey("isPublic"), field.Bool("pinned").StorageKey("pinned"), field.JSON("extraData", json.RawMessage{}).StorageKey("extraData"), field.String("category").StorageKey("category"), field.String("ownerId").Optional().StorageKey("ownerId"), field.Time("createdAt").StorageKey("createdAt"), field.Time("updatedAt").StorageKey("updatedAt")}

}
func (Todo) Edges() []ent.Edge {
	return []ent.Edge{edge.From("user", User.Type).Ref("todos").Unique().Field("ownerId")}
}
func (Todo) Annotations() []schema.Annotation {
	return nil
}
