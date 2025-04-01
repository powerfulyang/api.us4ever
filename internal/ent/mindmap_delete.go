// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"

	"api.us4ever/internal/ent/mindmap"
	"api.us4ever/internal/ent/predicate"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// MindmapDelete is the builder for deleting a Mindmap entity.
type MindmapDelete struct {
	config
	hooks    []Hook
	mutation *MindmapMutation
}

// Where appends a list predicates to the MindmapDelete builder.
func (md *MindmapDelete) Where(ps ...predicate.Mindmap) *MindmapDelete {
	md.mutation.Where(ps...)
	return md
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (md *MindmapDelete) Exec(ctx context.Context) (int, error) {
	return withHooks(ctx, md.sqlExec, md.mutation, md.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (md *MindmapDelete) ExecX(ctx context.Context) int {
	n, err := md.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (md *MindmapDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := sqlgraph.NewDeleteSpec(mindmap.Table, sqlgraph.NewFieldSpec(mindmap.FieldID, field.TypeString))
	if ps := md.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, md.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	md.mutation.done = true
	return affected, err
}

// MindmapDeleteOne is the builder for deleting a single Mindmap entity.
type MindmapDeleteOne struct {
	md *MindmapDelete
}

// Where appends a list predicates to the MindmapDelete builder.
func (mdo *MindmapDeleteOne) Where(ps ...predicate.Mindmap) *MindmapDeleteOne {
	mdo.md.mutation.Where(ps...)
	return mdo
}

// Exec executes the deletion query.
func (mdo *MindmapDeleteOne) Exec(ctx context.Context) error {
	n, err := mdo.md.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{mindmap.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (mdo *MindmapDeleteOne) ExecX(ctx context.Context) {
	if err := mdo.Exec(ctx); err != nil {
		panic(err)
	}
}
