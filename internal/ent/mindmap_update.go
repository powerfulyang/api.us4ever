// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"api.us4ever/internal/ent/mindmap"
	"api.us4ever/internal/ent/predicate"
	"api.us4ever/internal/ent/user"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/dialect/sql/sqljson"
	"entgo.io/ent/schema/field"
)

// MindmapUpdate is the builder for updating Mindmap entities.
type MindmapUpdate struct {
	config
	hooks    []Hook
	mutation *MindmapMutation
}

// Where appends a list predicates to the MindmapUpdate builder.
func (mu *MindmapUpdate) Where(ps ...predicate.Mindmap) *MindmapUpdate {
	mu.mutation.Where(ps...)
	return mu
}

// SetTitle sets the "title" field.
func (mu *MindmapUpdate) SetTitle(s string) *MindmapUpdate {
	mu.mutation.SetTitle(s)
	return mu
}

// SetNillableTitle sets the "title" field if the given value is not nil.
func (mu *MindmapUpdate) SetNillableTitle(s *string) *MindmapUpdate {
	if s != nil {
		mu.SetTitle(*s)
	}
	return mu
}

// SetContent sets the "content" field.
func (mu *MindmapUpdate) SetContent(jm json.RawMessage) *MindmapUpdate {
	mu.mutation.SetContent(jm)
	return mu
}

// AppendContent appends jm to the "content" field.
func (mu *MindmapUpdate) AppendContent(jm json.RawMessage) *MindmapUpdate {
	mu.mutation.AppendContent(jm)
	return mu
}

// SetSummary sets the "summary" field.
func (mu *MindmapUpdate) SetSummary(s string) *MindmapUpdate {
	mu.mutation.SetSummary(s)
	return mu
}

// SetNillableSummary sets the "summary" field if the given value is not nil.
func (mu *MindmapUpdate) SetNillableSummary(s *string) *MindmapUpdate {
	if s != nil {
		mu.SetSummary(*s)
	}
	return mu
}

// SetIsPublic sets the "isPublic" field.
func (mu *MindmapUpdate) SetIsPublic(b bool) *MindmapUpdate {
	mu.mutation.SetIsPublic(b)
	return mu
}

// SetNillableIsPublic sets the "isPublic" field if the given value is not nil.
func (mu *MindmapUpdate) SetNillableIsPublic(b *bool) *MindmapUpdate {
	if b != nil {
		mu.SetIsPublic(*b)
	}
	return mu
}

// SetTags sets the "tags" field.
func (mu *MindmapUpdate) SetTags(jm json.RawMessage) *MindmapUpdate {
	mu.mutation.SetTags(jm)
	return mu
}

// AppendTags appends jm to the "tags" field.
func (mu *MindmapUpdate) AppendTags(jm json.RawMessage) *MindmapUpdate {
	mu.mutation.AppendTags(jm)
	return mu
}

// SetViews sets the "views" field.
func (mu *MindmapUpdate) SetViews(i int32) *MindmapUpdate {
	mu.mutation.ResetViews()
	mu.mutation.SetViews(i)
	return mu
}

// SetNillableViews sets the "views" field if the given value is not nil.
func (mu *MindmapUpdate) SetNillableViews(i *int32) *MindmapUpdate {
	if i != nil {
		mu.SetViews(*i)
	}
	return mu
}

// AddViews adds i to the "views" field.
func (mu *MindmapUpdate) AddViews(i int32) *MindmapUpdate {
	mu.mutation.AddViews(i)
	return mu
}

// SetLikes sets the "likes" field.
func (mu *MindmapUpdate) SetLikes(i int32) *MindmapUpdate {
	mu.mutation.ResetLikes()
	mu.mutation.SetLikes(i)
	return mu
}

// SetNillableLikes sets the "likes" field if the given value is not nil.
func (mu *MindmapUpdate) SetNillableLikes(i *int32) *MindmapUpdate {
	if i != nil {
		mu.SetLikes(*i)
	}
	return mu
}

// AddLikes adds i to the "likes" field.
func (mu *MindmapUpdate) AddLikes(i int32) *MindmapUpdate {
	mu.mutation.AddLikes(i)
	return mu
}

// SetExtraData sets the "extraData" field.
func (mu *MindmapUpdate) SetExtraData(jm json.RawMessage) *MindmapUpdate {
	mu.mutation.SetExtraData(jm)
	return mu
}

// AppendExtraData appends jm to the "extraData" field.
func (mu *MindmapUpdate) AppendExtraData(jm json.RawMessage) *MindmapUpdate {
	mu.mutation.AppendExtraData(jm)
	return mu
}

// SetCategory sets the "category" field.
func (mu *MindmapUpdate) SetCategory(s string) *MindmapUpdate {
	mu.mutation.SetCategory(s)
	return mu
}

// SetNillableCategory sets the "category" field if the given value is not nil.
func (mu *MindmapUpdate) SetNillableCategory(s *string) *MindmapUpdate {
	if s != nil {
		mu.SetCategory(*s)
	}
	return mu
}

// SetOwnerId sets the "ownerId" field.
func (mu *MindmapUpdate) SetOwnerId(s string) *MindmapUpdate {
	mu.mutation.SetOwnerId(s)
	return mu
}

// SetNillableOwnerId sets the "ownerId" field if the given value is not nil.
func (mu *MindmapUpdate) SetNillableOwnerId(s *string) *MindmapUpdate {
	if s != nil {
		mu.SetOwnerId(*s)
	}
	return mu
}

// ClearOwnerId clears the value of the "ownerId" field.
func (mu *MindmapUpdate) ClearOwnerId() *MindmapUpdate {
	mu.mutation.ClearOwnerId()
	return mu
}

// SetCreatedAt sets the "createdAt" field.
func (mu *MindmapUpdate) SetCreatedAt(t time.Time) *MindmapUpdate {
	mu.mutation.SetCreatedAt(t)
	return mu
}

// SetNillableCreatedAt sets the "createdAt" field if the given value is not nil.
func (mu *MindmapUpdate) SetNillableCreatedAt(t *time.Time) *MindmapUpdate {
	if t != nil {
		mu.SetCreatedAt(*t)
	}
	return mu
}

// SetUpdatedAt sets the "updatedAt" field.
func (mu *MindmapUpdate) SetUpdatedAt(t time.Time) *MindmapUpdate {
	mu.mutation.SetUpdatedAt(t)
	return mu
}

// SetNillableUpdatedAt sets the "updatedAt" field if the given value is not nil.
func (mu *MindmapUpdate) SetNillableUpdatedAt(t *time.Time) *MindmapUpdate {
	if t != nil {
		mu.SetUpdatedAt(*t)
	}
	return mu
}

// SetUserID sets the "user" edge to the User entity by ID.
func (mu *MindmapUpdate) SetUserID(id string) *MindmapUpdate {
	mu.mutation.SetUserID(id)
	return mu
}

// SetNillableUserID sets the "user" edge to the User entity by ID if the given value is not nil.
func (mu *MindmapUpdate) SetNillableUserID(id *string) *MindmapUpdate {
	if id != nil {
		mu = mu.SetUserID(*id)
	}
	return mu
}

// SetUser sets the "user" edge to the User entity.
func (mu *MindmapUpdate) SetUser(u *User) *MindmapUpdate {
	return mu.SetUserID(u.ID)
}

// Mutation returns the MindmapMutation object of the builder.
func (mu *MindmapUpdate) Mutation() *MindmapMutation {
	return mu.mutation
}

// ClearUser clears the "user" edge to the User entity.
func (mu *MindmapUpdate) ClearUser() *MindmapUpdate {
	mu.mutation.ClearUser()
	return mu
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (mu *MindmapUpdate) Save(ctx context.Context) (int, error) {
	return withHooks(ctx, mu.sqlSave, mu.mutation, mu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (mu *MindmapUpdate) SaveX(ctx context.Context) int {
	affected, err := mu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (mu *MindmapUpdate) Exec(ctx context.Context) error {
	_, err := mu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (mu *MindmapUpdate) ExecX(ctx context.Context) {
	if err := mu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (mu *MindmapUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := sqlgraph.NewUpdateSpec(mindmap.Table, mindmap.Columns, sqlgraph.NewFieldSpec(mindmap.FieldID, field.TypeString))
	if ps := mu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := mu.mutation.Title(); ok {
		_spec.SetField(mindmap.FieldTitle, field.TypeString, value)
	}
	if value, ok := mu.mutation.Content(); ok {
		_spec.SetField(mindmap.FieldContent, field.TypeJSON, value)
	}
	if value, ok := mu.mutation.AppendedContent(); ok {
		_spec.AddModifier(func(u *sql.UpdateBuilder) {
			sqljson.Append(u, mindmap.FieldContent, value)
		})
	}
	if value, ok := mu.mutation.Summary(); ok {
		_spec.SetField(mindmap.FieldSummary, field.TypeString, value)
	}
	if value, ok := mu.mutation.IsPublic(); ok {
		_spec.SetField(mindmap.FieldIsPublic, field.TypeBool, value)
	}
	if value, ok := mu.mutation.Tags(); ok {
		_spec.SetField(mindmap.FieldTags, field.TypeJSON, value)
	}
	if value, ok := mu.mutation.AppendedTags(); ok {
		_spec.AddModifier(func(u *sql.UpdateBuilder) {
			sqljson.Append(u, mindmap.FieldTags, value)
		})
	}
	if value, ok := mu.mutation.Views(); ok {
		_spec.SetField(mindmap.FieldViews, field.TypeInt32, value)
	}
	if value, ok := mu.mutation.AddedViews(); ok {
		_spec.AddField(mindmap.FieldViews, field.TypeInt32, value)
	}
	if value, ok := mu.mutation.Likes(); ok {
		_spec.SetField(mindmap.FieldLikes, field.TypeInt32, value)
	}
	if value, ok := mu.mutation.AddedLikes(); ok {
		_spec.AddField(mindmap.FieldLikes, field.TypeInt32, value)
	}
	if value, ok := mu.mutation.ExtraData(); ok {
		_spec.SetField(mindmap.FieldExtraData, field.TypeJSON, value)
	}
	if value, ok := mu.mutation.AppendedExtraData(); ok {
		_spec.AddModifier(func(u *sql.UpdateBuilder) {
			sqljson.Append(u, mindmap.FieldExtraData, value)
		})
	}
	if value, ok := mu.mutation.Category(); ok {
		_spec.SetField(mindmap.FieldCategory, field.TypeString, value)
	}
	if value, ok := mu.mutation.CreatedAt(); ok {
		_spec.SetField(mindmap.FieldCreatedAt, field.TypeTime, value)
	}
	if value, ok := mu.mutation.UpdatedAt(); ok {
		_spec.SetField(mindmap.FieldUpdatedAt, field.TypeTime, value)
	}
	if mu.mutation.UserCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   mindmap.UserTable,
			Columns: []string{mindmap.UserColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(user.FieldID, field.TypeString),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := mu.mutation.UserIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   mindmap.UserTable,
			Columns: []string{mindmap.UserColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(user.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, mu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{mindmap.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	mu.mutation.done = true
	return n, nil
}

// MindmapUpdateOne is the builder for updating a single Mindmap entity.
type MindmapUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *MindmapMutation
}

// SetTitle sets the "title" field.
func (muo *MindmapUpdateOne) SetTitle(s string) *MindmapUpdateOne {
	muo.mutation.SetTitle(s)
	return muo
}

// SetNillableTitle sets the "title" field if the given value is not nil.
func (muo *MindmapUpdateOne) SetNillableTitle(s *string) *MindmapUpdateOne {
	if s != nil {
		muo.SetTitle(*s)
	}
	return muo
}

// SetContent sets the "content" field.
func (muo *MindmapUpdateOne) SetContent(jm json.RawMessage) *MindmapUpdateOne {
	muo.mutation.SetContent(jm)
	return muo
}

// AppendContent appends jm to the "content" field.
func (muo *MindmapUpdateOne) AppendContent(jm json.RawMessage) *MindmapUpdateOne {
	muo.mutation.AppendContent(jm)
	return muo
}

// SetSummary sets the "summary" field.
func (muo *MindmapUpdateOne) SetSummary(s string) *MindmapUpdateOne {
	muo.mutation.SetSummary(s)
	return muo
}

// SetNillableSummary sets the "summary" field if the given value is not nil.
func (muo *MindmapUpdateOne) SetNillableSummary(s *string) *MindmapUpdateOne {
	if s != nil {
		muo.SetSummary(*s)
	}
	return muo
}

// SetIsPublic sets the "isPublic" field.
func (muo *MindmapUpdateOne) SetIsPublic(b bool) *MindmapUpdateOne {
	muo.mutation.SetIsPublic(b)
	return muo
}

// SetNillableIsPublic sets the "isPublic" field if the given value is not nil.
func (muo *MindmapUpdateOne) SetNillableIsPublic(b *bool) *MindmapUpdateOne {
	if b != nil {
		muo.SetIsPublic(*b)
	}
	return muo
}

// SetTags sets the "tags" field.
func (muo *MindmapUpdateOne) SetTags(jm json.RawMessage) *MindmapUpdateOne {
	muo.mutation.SetTags(jm)
	return muo
}

// AppendTags appends jm to the "tags" field.
func (muo *MindmapUpdateOne) AppendTags(jm json.RawMessage) *MindmapUpdateOne {
	muo.mutation.AppendTags(jm)
	return muo
}

// SetViews sets the "views" field.
func (muo *MindmapUpdateOne) SetViews(i int32) *MindmapUpdateOne {
	muo.mutation.ResetViews()
	muo.mutation.SetViews(i)
	return muo
}

// SetNillableViews sets the "views" field if the given value is not nil.
func (muo *MindmapUpdateOne) SetNillableViews(i *int32) *MindmapUpdateOne {
	if i != nil {
		muo.SetViews(*i)
	}
	return muo
}

// AddViews adds i to the "views" field.
func (muo *MindmapUpdateOne) AddViews(i int32) *MindmapUpdateOne {
	muo.mutation.AddViews(i)
	return muo
}

// SetLikes sets the "likes" field.
func (muo *MindmapUpdateOne) SetLikes(i int32) *MindmapUpdateOne {
	muo.mutation.ResetLikes()
	muo.mutation.SetLikes(i)
	return muo
}

// SetNillableLikes sets the "likes" field if the given value is not nil.
func (muo *MindmapUpdateOne) SetNillableLikes(i *int32) *MindmapUpdateOne {
	if i != nil {
		muo.SetLikes(*i)
	}
	return muo
}

// AddLikes adds i to the "likes" field.
func (muo *MindmapUpdateOne) AddLikes(i int32) *MindmapUpdateOne {
	muo.mutation.AddLikes(i)
	return muo
}

// SetExtraData sets the "extraData" field.
func (muo *MindmapUpdateOne) SetExtraData(jm json.RawMessage) *MindmapUpdateOne {
	muo.mutation.SetExtraData(jm)
	return muo
}

// AppendExtraData appends jm to the "extraData" field.
func (muo *MindmapUpdateOne) AppendExtraData(jm json.RawMessage) *MindmapUpdateOne {
	muo.mutation.AppendExtraData(jm)
	return muo
}

// SetCategory sets the "category" field.
func (muo *MindmapUpdateOne) SetCategory(s string) *MindmapUpdateOne {
	muo.mutation.SetCategory(s)
	return muo
}

// SetNillableCategory sets the "category" field if the given value is not nil.
func (muo *MindmapUpdateOne) SetNillableCategory(s *string) *MindmapUpdateOne {
	if s != nil {
		muo.SetCategory(*s)
	}
	return muo
}

// SetOwnerId sets the "ownerId" field.
func (muo *MindmapUpdateOne) SetOwnerId(s string) *MindmapUpdateOne {
	muo.mutation.SetOwnerId(s)
	return muo
}

// SetNillableOwnerId sets the "ownerId" field if the given value is not nil.
func (muo *MindmapUpdateOne) SetNillableOwnerId(s *string) *MindmapUpdateOne {
	if s != nil {
		muo.SetOwnerId(*s)
	}
	return muo
}

// ClearOwnerId clears the value of the "ownerId" field.
func (muo *MindmapUpdateOne) ClearOwnerId() *MindmapUpdateOne {
	muo.mutation.ClearOwnerId()
	return muo
}

// SetCreatedAt sets the "createdAt" field.
func (muo *MindmapUpdateOne) SetCreatedAt(t time.Time) *MindmapUpdateOne {
	muo.mutation.SetCreatedAt(t)
	return muo
}

// SetNillableCreatedAt sets the "createdAt" field if the given value is not nil.
func (muo *MindmapUpdateOne) SetNillableCreatedAt(t *time.Time) *MindmapUpdateOne {
	if t != nil {
		muo.SetCreatedAt(*t)
	}
	return muo
}

// SetUpdatedAt sets the "updatedAt" field.
func (muo *MindmapUpdateOne) SetUpdatedAt(t time.Time) *MindmapUpdateOne {
	muo.mutation.SetUpdatedAt(t)
	return muo
}

// SetNillableUpdatedAt sets the "updatedAt" field if the given value is not nil.
func (muo *MindmapUpdateOne) SetNillableUpdatedAt(t *time.Time) *MindmapUpdateOne {
	if t != nil {
		muo.SetUpdatedAt(*t)
	}
	return muo
}

// SetUserID sets the "user" edge to the User entity by ID.
func (muo *MindmapUpdateOne) SetUserID(id string) *MindmapUpdateOne {
	muo.mutation.SetUserID(id)
	return muo
}

// SetNillableUserID sets the "user" edge to the User entity by ID if the given value is not nil.
func (muo *MindmapUpdateOne) SetNillableUserID(id *string) *MindmapUpdateOne {
	if id != nil {
		muo = muo.SetUserID(*id)
	}
	return muo
}

// SetUser sets the "user" edge to the User entity.
func (muo *MindmapUpdateOne) SetUser(u *User) *MindmapUpdateOne {
	return muo.SetUserID(u.ID)
}

// Mutation returns the MindmapMutation object of the builder.
func (muo *MindmapUpdateOne) Mutation() *MindmapMutation {
	return muo.mutation
}

// ClearUser clears the "user" edge to the User entity.
func (muo *MindmapUpdateOne) ClearUser() *MindmapUpdateOne {
	muo.mutation.ClearUser()
	return muo
}

// Where appends a list predicates to the MindmapUpdate builder.
func (muo *MindmapUpdateOne) Where(ps ...predicate.Mindmap) *MindmapUpdateOne {
	muo.mutation.Where(ps...)
	return muo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (muo *MindmapUpdateOne) Select(field string, fields ...string) *MindmapUpdateOne {
	muo.fields = append([]string{field}, fields...)
	return muo
}

// Save executes the query and returns the updated Mindmap entity.
func (muo *MindmapUpdateOne) Save(ctx context.Context) (*Mindmap, error) {
	return withHooks(ctx, muo.sqlSave, muo.mutation, muo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (muo *MindmapUpdateOne) SaveX(ctx context.Context) *Mindmap {
	node, err := muo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (muo *MindmapUpdateOne) Exec(ctx context.Context) error {
	_, err := muo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (muo *MindmapUpdateOne) ExecX(ctx context.Context) {
	if err := muo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (muo *MindmapUpdateOne) sqlSave(ctx context.Context) (_node *Mindmap, err error) {
	_spec := sqlgraph.NewUpdateSpec(mindmap.Table, mindmap.Columns, sqlgraph.NewFieldSpec(mindmap.FieldID, field.TypeString))
	id, ok := muo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "Mindmap.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := muo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, mindmap.FieldID)
		for _, f := range fields {
			if !mindmap.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != mindmap.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := muo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := muo.mutation.Title(); ok {
		_spec.SetField(mindmap.FieldTitle, field.TypeString, value)
	}
	if value, ok := muo.mutation.Content(); ok {
		_spec.SetField(mindmap.FieldContent, field.TypeJSON, value)
	}
	if value, ok := muo.mutation.AppendedContent(); ok {
		_spec.AddModifier(func(u *sql.UpdateBuilder) {
			sqljson.Append(u, mindmap.FieldContent, value)
		})
	}
	if value, ok := muo.mutation.Summary(); ok {
		_spec.SetField(mindmap.FieldSummary, field.TypeString, value)
	}
	if value, ok := muo.mutation.IsPublic(); ok {
		_spec.SetField(mindmap.FieldIsPublic, field.TypeBool, value)
	}
	if value, ok := muo.mutation.Tags(); ok {
		_spec.SetField(mindmap.FieldTags, field.TypeJSON, value)
	}
	if value, ok := muo.mutation.AppendedTags(); ok {
		_spec.AddModifier(func(u *sql.UpdateBuilder) {
			sqljson.Append(u, mindmap.FieldTags, value)
		})
	}
	if value, ok := muo.mutation.Views(); ok {
		_spec.SetField(mindmap.FieldViews, field.TypeInt32, value)
	}
	if value, ok := muo.mutation.AddedViews(); ok {
		_spec.AddField(mindmap.FieldViews, field.TypeInt32, value)
	}
	if value, ok := muo.mutation.Likes(); ok {
		_spec.SetField(mindmap.FieldLikes, field.TypeInt32, value)
	}
	if value, ok := muo.mutation.AddedLikes(); ok {
		_spec.AddField(mindmap.FieldLikes, field.TypeInt32, value)
	}
	if value, ok := muo.mutation.ExtraData(); ok {
		_spec.SetField(mindmap.FieldExtraData, field.TypeJSON, value)
	}
	if value, ok := muo.mutation.AppendedExtraData(); ok {
		_spec.AddModifier(func(u *sql.UpdateBuilder) {
			sqljson.Append(u, mindmap.FieldExtraData, value)
		})
	}
	if value, ok := muo.mutation.Category(); ok {
		_spec.SetField(mindmap.FieldCategory, field.TypeString, value)
	}
	if value, ok := muo.mutation.CreatedAt(); ok {
		_spec.SetField(mindmap.FieldCreatedAt, field.TypeTime, value)
	}
	if value, ok := muo.mutation.UpdatedAt(); ok {
		_spec.SetField(mindmap.FieldUpdatedAt, field.TypeTime, value)
	}
	if muo.mutation.UserCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   mindmap.UserTable,
			Columns: []string{mindmap.UserColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(user.FieldID, field.TypeString),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := muo.mutation.UserIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   mindmap.UserTable,
			Columns: []string{mindmap.UserColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(user.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &Mindmap{config: muo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, muo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{mindmap.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	muo.mutation.done = true
	return _node, nil
}
