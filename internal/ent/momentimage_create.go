// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"api.us4ever/internal/ent/image"
	"api.us4ever/internal/ent/moment"
	"api.us4ever/internal/ent/momentimage"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// MomentImageCreate is the builder for creating a MomentImage entity.
type MomentImageCreate struct {
	config
	mutation *MomentImageMutation
	hooks    []Hook
}

// SetImageId sets the "imageId" field.
func (mic *MomentImageCreate) SetImageId(s string) *MomentImageCreate {
	mic.mutation.SetImageId(s)
	return mic
}

// SetNillableImageId sets the "imageId" field if the given value is not nil.
func (mic *MomentImageCreate) SetNillableImageId(s *string) *MomentImageCreate {
	if s != nil {
		mic.SetImageId(*s)
	}
	return mic
}

// SetMomentId sets the "momentId" field.
func (mic *MomentImageCreate) SetMomentId(s string) *MomentImageCreate {
	mic.mutation.SetMomentId(s)
	return mic
}

// SetNillableMomentId sets the "momentId" field if the given value is not nil.
func (mic *MomentImageCreate) SetNillableMomentId(s *string) *MomentImageCreate {
	if s != nil {
		mic.SetMomentId(*s)
	}
	return mic
}

// SetSort sets the "sort" field.
func (mic *MomentImageCreate) SetSort(i int32) *MomentImageCreate {
	mic.mutation.SetSort(i)
	return mic
}

// SetCreatedAt sets the "createdAt" field.
func (mic *MomentImageCreate) SetCreatedAt(t time.Time) *MomentImageCreate {
	mic.mutation.SetCreatedAt(t)
	return mic
}

// SetUpdatedAt sets the "updatedAt" field.
func (mic *MomentImageCreate) SetUpdatedAt(t time.Time) *MomentImageCreate {
	mic.mutation.SetUpdatedAt(t)
	return mic
}

// SetID sets the "id" field.
func (mic *MomentImageCreate) SetID(u uint) *MomentImageCreate {
	mic.mutation.SetID(u)
	return mic
}

// SetImageID sets the "image" edge to the Image entity by ID.
func (mic *MomentImageCreate) SetImageID(id string) *MomentImageCreate {
	mic.mutation.SetImageID(id)
	return mic
}

// SetNillableImageID sets the "image" edge to the Image entity by ID if the given value is not nil.
func (mic *MomentImageCreate) SetNillableImageID(id *string) *MomentImageCreate {
	if id != nil {
		mic = mic.SetImageID(*id)
	}
	return mic
}

// SetImage sets the "image" edge to the Image entity.
func (mic *MomentImageCreate) SetImage(i *Image) *MomentImageCreate {
	return mic.SetImageID(i.ID)
}

// SetMomentID sets the "moment" edge to the Moment entity by ID.
func (mic *MomentImageCreate) SetMomentID(id string) *MomentImageCreate {
	mic.mutation.SetMomentID(id)
	return mic
}

// SetNillableMomentID sets the "moment" edge to the Moment entity by ID if the given value is not nil.
func (mic *MomentImageCreate) SetNillableMomentID(id *string) *MomentImageCreate {
	if id != nil {
		mic = mic.SetMomentID(*id)
	}
	return mic
}

// SetMoment sets the "moment" edge to the Moment entity.
func (mic *MomentImageCreate) SetMoment(m *Moment) *MomentImageCreate {
	return mic.SetMomentID(m.ID)
}

// Mutation returns the MomentImageMutation object of the builder.
func (mic *MomentImageCreate) Mutation() *MomentImageMutation {
	return mic.mutation
}

// Save creates the MomentImage in the database.
func (mic *MomentImageCreate) Save(ctx context.Context) (*MomentImage, error) {
	return withHooks(ctx, mic.sqlSave, mic.mutation, mic.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (mic *MomentImageCreate) SaveX(ctx context.Context) *MomentImage {
	v, err := mic.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (mic *MomentImageCreate) Exec(ctx context.Context) error {
	_, err := mic.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (mic *MomentImageCreate) ExecX(ctx context.Context) {
	if err := mic.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (mic *MomentImageCreate) check() error {
	if _, ok := mic.mutation.Sort(); !ok {
		return &ValidationError{Name: "sort", err: errors.New(`ent: missing required field "MomentImage.sort"`)}
	}
	if _, ok := mic.mutation.CreatedAt(); !ok {
		return &ValidationError{Name: "createdAt", err: errors.New(`ent: missing required field "MomentImage.createdAt"`)}
	}
	if _, ok := mic.mutation.UpdatedAt(); !ok {
		return &ValidationError{Name: "updatedAt", err: errors.New(`ent: missing required field "MomentImage.updatedAt"`)}
	}
	return nil
}

func (mic *MomentImageCreate) sqlSave(ctx context.Context) (*MomentImage, error) {
	if err := mic.check(); err != nil {
		return nil, err
	}
	_node, _spec := mic.createSpec()
	if err := sqlgraph.CreateNode(ctx, mic.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	if _spec.ID.Value != _node.ID {
		id := _spec.ID.Value.(int64)
		_node.ID = uint(id)
	}
	mic.mutation.id = &_node.ID
	mic.mutation.done = true
	return _node, nil
}

func (mic *MomentImageCreate) createSpec() (*MomentImage, *sqlgraph.CreateSpec) {
	var (
		_node = &MomentImage{config: mic.config}
		_spec = sqlgraph.NewCreateSpec(momentimage.Table, sqlgraph.NewFieldSpec(momentimage.FieldID, field.TypeUint))
	)
	if id, ok := mic.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = id
	}
	if value, ok := mic.mutation.Sort(); ok {
		_spec.SetField(momentimage.FieldSort, field.TypeInt32, value)
		_node.Sort = value
	}
	if value, ok := mic.mutation.CreatedAt(); ok {
		_spec.SetField(momentimage.FieldCreatedAt, field.TypeTime, value)
		_node.CreatedAt = value
	}
	if value, ok := mic.mutation.UpdatedAt(); ok {
		_spec.SetField(momentimage.FieldUpdatedAt, field.TypeTime, value)
		_node.UpdatedAt = value
	}
	if nodes := mic.mutation.ImageIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   momentimage.ImageTable,
			Columns: []string{momentimage.ImageColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(image.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.ImageId = nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := mic.mutation.MomentIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   momentimage.MomentTable,
			Columns: []string{momentimage.MomentColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(moment.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.MomentId = nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// MomentImageCreateBulk is the builder for creating many MomentImage entities in bulk.
type MomentImageCreateBulk struct {
	config
	err      error
	builders []*MomentImageCreate
}

// Save creates the MomentImage entities in the database.
func (micb *MomentImageCreateBulk) Save(ctx context.Context) ([]*MomentImage, error) {
	if micb.err != nil {
		return nil, micb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(micb.builders))
	nodes := make([]*MomentImage, len(micb.builders))
	mutators := make([]Mutator, len(micb.builders))
	for i := range micb.builders {
		func(i int, root context.Context) {
			builder := micb.builders[i]
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*MomentImageMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				var err error
				nodes[i], specs[i] = builder.createSpec()
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, micb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, micb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				if specs[i].ID.Value != nil && nodes[i].ID == 0 {
					id := specs[i].ID.Value.(int64)
					nodes[i].ID = uint(id)
				}
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, micb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (micb *MomentImageCreateBulk) SaveX(ctx context.Context) []*MomentImage {
	v, err := micb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (micb *MomentImageCreateBulk) Exec(ctx context.Context) error {
	_, err := micb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (micb *MomentImageCreateBulk) ExecX(ctx context.Context) {
	if err := micb.Exec(ctx); err != nil {
		panic(err)
	}
}
