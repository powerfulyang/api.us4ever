// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"

	"api.us4ever/internal/ent/predicate"
	"api.us4ever/internal/ent/video"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// VideoDelete is the builder for deleting a Video entity.
type VideoDelete struct {
	config
	hooks    []Hook
	mutation *VideoMutation
}

// Where appends a list predicates to the VideoDelete builder.
func (vd *VideoDelete) Where(ps ...predicate.Video) *VideoDelete {
	vd.mutation.Where(ps...)
	return vd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (vd *VideoDelete) Exec(ctx context.Context) (int, error) {
	return withHooks(ctx, vd.sqlExec, vd.mutation, vd.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (vd *VideoDelete) ExecX(ctx context.Context) int {
	n, err := vd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (vd *VideoDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := sqlgraph.NewDeleteSpec(video.Table, sqlgraph.NewFieldSpec(video.FieldID, field.TypeString))
	if ps := vd.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, vd.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	vd.mutation.done = true
	return affected, err
}

// VideoDeleteOne is the builder for deleting a single Video entity.
type VideoDeleteOne struct {
	vd *VideoDelete
}

// Where appends a list predicates to the VideoDelete builder.
func (vdo *VideoDeleteOne) Where(ps ...predicate.Video) *VideoDeleteOne {
	vdo.vd.mutation.Where(ps...)
	return vdo
}

// Exec executes the deletion query.
func (vdo *VideoDeleteOne) Exec(ctx context.Context) error {
	n, err := vdo.vd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{video.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (vdo *VideoDeleteOne) ExecX(ctx context.Context) {
	if err := vdo.Exec(ctx); err != nil {
		panic(err)
	}
}
