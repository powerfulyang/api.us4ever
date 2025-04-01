// Code generated by ent, DO NOT EDIT.

package momentimage

import (
	"time"

	"api.us4ever/internal/ent/predicate"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

// ID filters vertices based on their ID field.
func ID(id uint) predicate.MomentImage {
	return predicate.MomentImage(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id uint) predicate.MomentImage {
	return predicate.MomentImage(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id uint) predicate.MomentImage {
	return predicate.MomentImage(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...uint) predicate.MomentImage {
	return predicate.MomentImage(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...uint) predicate.MomentImage {
	return predicate.MomentImage(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id uint) predicate.MomentImage {
	return predicate.MomentImage(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id uint) predicate.MomentImage {
	return predicate.MomentImage(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id uint) predicate.MomentImage {
	return predicate.MomentImage(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id uint) predicate.MomentImage {
	return predicate.MomentImage(sql.FieldLTE(FieldID, id))
}

// ImageId applies equality check predicate on the "imageId" field. It's identical to ImageIdEQ.
func ImageId(v string) predicate.MomentImage {
	return predicate.MomentImage(sql.FieldEQ(FieldImageId, v))
}

// MomentId applies equality check predicate on the "momentId" field. It's identical to MomentIdEQ.
func MomentId(v string) predicate.MomentImage {
	return predicate.MomentImage(sql.FieldEQ(FieldMomentId, v))
}

// Sort applies equality check predicate on the "sort" field. It's identical to SortEQ.
func Sort(v int32) predicate.MomentImage {
	return predicate.MomentImage(sql.FieldEQ(FieldSort, v))
}

// CreatedAt applies equality check predicate on the "createdAt" field. It's identical to CreatedAtEQ.
func CreatedAt(v time.Time) predicate.MomentImage {
	return predicate.MomentImage(sql.FieldEQ(FieldCreatedAt, v))
}

// UpdatedAt applies equality check predicate on the "updatedAt" field. It's identical to UpdatedAtEQ.
func UpdatedAt(v time.Time) predicate.MomentImage {
	return predicate.MomentImage(sql.FieldEQ(FieldUpdatedAt, v))
}

// ImageIdEQ applies the EQ predicate on the "imageId" field.
func ImageIdEQ(v string) predicate.MomentImage {
	return predicate.MomentImage(sql.FieldEQ(FieldImageId, v))
}

// ImageIdNEQ applies the NEQ predicate on the "imageId" field.
func ImageIdNEQ(v string) predicate.MomentImage {
	return predicate.MomentImage(sql.FieldNEQ(FieldImageId, v))
}

// ImageIdIn applies the In predicate on the "imageId" field.
func ImageIdIn(vs ...string) predicate.MomentImage {
	return predicate.MomentImage(sql.FieldIn(FieldImageId, vs...))
}

// ImageIdNotIn applies the NotIn predicate on the "imageId" field.
func ImageIdNotIn(vs ...string) predicate.MomentImage {
	return predicate.MomentImage(sql.FieldNotIn(FieldImageId, vs...))
}

// ImageIdGT applies the GT predicate on the "imageId" field.
func ImageIdGT(v string) predicate.MomentImage {
	return predicate.MomentImage(sql.FieldGT(FieldImageId, v))
}

// ImageIdGTE applies the GTE predicate on the "imageId" field.
func ImageIdGTE(v string) predicate.MomentImage {
	return predicate.MomentImage(sql.FieldGTE(FieldImageId, v))
}

// ImageIdLT applies the LT predicate on the "imageId" field.
func ImageIdLT(v string) predicate.MomentImage {
	return predicate.MomentImage(sql.FieldLT(FieldImageId, v))
}

// ImageIdLTE applies the LTE predicate on the "imageId" field.
func ImageIdLTE(v string) predicate.MomentImage {
	return predicate.MomentImage(sql.FieldLTE(FieldImageId, v))
}

// ImageIdContains applies the Contains predicate on the "imageId" field.
func ImageIdContains(v string) predicate.MomentImage {
	return predicate.MomentImage(sql.FieldContains(FieldImageId, v))
}

// ImageIdHasPrefix applies the HasPrefix predicate on the "imageId" field.
func ImageIdHasPrefix(v string) predicate.MomentImage {
	return predicate.MomentImage(sql.FieldHasPrefix(FieldImageId, v))
}

// ImageIdHasSuffix applies the HasSuffix predicate on the "imageId" field.
func ImageIdHasSuffix(v string) predicate.MomentImage {
	return predicate.MomentImage(sql.FieldHasSuffix(FieldImageId, v))
}

// ImageIdIsNil applies the IsNil predicate on the "imageId" field.
func ImageIdIsNil() predicate.MomentImage {
	return predicate.MomentImage(sql.FieldIsNull(FieldImageId))
}

// ImageIdNotNil applies the NotNil predicate on the "imageId" field.
func ImageIdNotNil() predicate.MomentImage {
	return predicate.MomentImage(sql.FieldNotNull(FieldImageId))
}

// ImageIdEqualFold applies the EqualFold predicate on the "imageId" field.
func ImageIdEqualFold(v string) predicate.MomentImage {
	return predicate.MomentImage(sql.FieldEqualFold(FieldImageId, v))
}

// ImageIdContainsFold applies the ContainsFold predicate on the "imageId" field.
func ImageIdContainsFold(v string) predicate.MomentImage {
	return predicate.MomentImage(sql.FieldContainsFold(FieldImageId, v))
}

// MomentIdEQ applies the EQ predicate on the "momentId" field.
func MomentIdEQ(v string) predicate.MomentImage {
	return predicate.MomentImage(sql.FieldEQ(FieldMomentId, v))
}

// MomentIdNEQ applies the NEQ predicate on the "momentId" field.
func MomentIdNEQ(v string) predicate.MomentImage {
	return predicate.MomentImage(sql.FieldNEQ(FieldMomentId, v))
}

// MomentIdIn applies the In predicate on the "momentId" field.
func MomentIdIn(vs ...string) predicate.MomentImage {
	return predicate.MomentImage(sql.FieldIn(FieldMomentId, vs...))
}

// MomentIdNotIn applies the NotIn predicate on the "momentId" field.
func MomentIdNotIn(vs ...string) predicate.MomentImage {
	return predicate.MomentImage(sql.FieldNotIn(FieldMomentId, vs...))
}

// MomentIdGT applies the GT predicate on the "momentId" field.
func MomentIdGT(v string) predicate.MomentImage {
	return predicate.MomentImage(sql.FieldGT(FieldMomentId, v))
}

// MomentIdGTE applies the GTE predicate on the "momentId" field.
func MomentIdGTE(v string) predicate.MomentImage {
	return predicate.MomentImage(sql.FieldGTE(FieldMomentId, v))
}

// MomentIdLT applies the LT predicate on the "momentId" field.
func MomentIdLT(v string) predicate.MomentImage {
	return predicate.MomentImage(sql.FieldLT(FieldMomentId, v))
}

// MomentIdLTE applies the LTE predicate on the "momentId" field.
func MomentIdLTE(v string) predicate.MomentImage {
	return predicate.MomentImage(sql.FieldLTE(FieldMomentId, v))
}

// MomentIdContains applies the Contains predicate on the "momentId" field.
func MomentIdContains(v string) predicate.MomentImage {
	return predicate.MomentImage(sql.FieldContains(FieldMomentId, v))
}

// MomentIdHasPrefix applies the HasPrefix predicate on the "momentId" field.
func MomentIdHasPrefix(v string) predicate.MomentImage {
	return predicate.MomentImage(sql.FieldHasPrefix(FieldMomentId, v))
}

// MomentIdHasSuffix applies the HasSuffix predicate on the "momentId" field.
func MomentIdHasSuffix(v string) predicate.MomentImage {
	return predicate.MomentImage(sql.FieldHasSuffix(FieldMomentId, v))
}

// MomentIdIsNil applies the IsNil predicate on the "momentId" field.
func MomentIdIsNil() predicate.MomentImage {
	return predicate.MomentImage(sql.FieldIsNull(FieldMomentId))
}

// MomentIdNotNil applies the NotNil predicate on the "momentId" field.
func MomentIdNotNil() predicate.MomentImage {
	return predicate.MomentImage(sql.FieldNotNull(FieldMomentId))
}

// MomentIdEqualFold applies the EqualFold predicate on the "momentId" field.
func MomentIdEqualFold(v string) predicate.MomentImage {
	return predicate.MomentImage(sql.FieldEqualFold(FieldMomentId, v))
}

// MomentIdContainsFold applies the ContainsFold predicate on the "momentId" field.
func MomentIdContainsFold(v string) predicate.MomentImage {
	return predicate.MomentImage(sql.FieldContainsFold(FieldMomentId, v))
}

// SortEQ applies the EQ predicate on the "sort" field.
func SortEQ(v int32) predicate.MomentImage {
	return predicate.MomentImage(sql.FieldEQ(FieldSort, v))
}

// SortNEQ applies the NEQ predicate on the "sort" field.
func SortNEQ(v int32) predicate.MomentImage {
	return predicate.MomentImage(sql.FieldNEQ(FieldSort, v))
}

// SortIn applies the In predicate on the "sort" field.
func SortIn(vs ...int32) predicate.MomentImage {
	return predicate.MomentImage(sql.FieldIn(FieldSort, vs...))
}

// SortNotIn applies the NotIn predicate on the "sort" field.
func SortNotIn(vs ...int32) predicate.MomentImage {
	return predicate.MomentImage(sql.FieldNotIn(FieldSort, vs...))
}

// SortGT applies the GT predicate on the "sort" field.
func SortGT(v int32) predicate.MomentImage {
	return predicate.MomentImage(sql.FieldGT(FieldSort, v))
}

// SortGTE applies the GTE predicate on the "sort" field.
func SortGTE(v int32) predicate.MomentImage {
	return predicate.MomentImage(sql.FieldGTE(FieldSort, v))
}

// SortLT applies the LT predicate on the "sort" field.
func SortLT(v int32) predicate.MomentImage {
	return predicate.MomentImage(sql.FieldLT(FieldSort, v))
}

// SortLTE applies the LTE predicate on the "sort" field.
func SortLTE(v int32) predicate.MomentImage {
	return predicate.MomentImage(sql.FieldLTE(FieldSort, v))
}

// CreatedAtEQ applies the EQ predicate on the "createdAt" field.
func CreatedAtEQ(v time.Time) predicate.MomentImage {
	return predicate.MomentImage(sql.FieldEQ(FieldCreatedAt, v))
}

// CreatedAtNEQ applies the NEQ predicate on the "createdAt" field.
func CreatedAtNEQ(v time.Time) predicate.MomentImage {
	return predicate.MomentImage(sql.FieldNEQ(FieldCreatedAt, v))
}

// CreatedAtIn applies the In predicate on the "createdAt" field.
func CreatedAtIn(vs ...time.Time) predicate.MomentImage {
	return predicate.MomentImage(sql.FieldIn(FieldCreatedAt, vs...))
}

// CreatedAtNotIn applies the NotIn predicate on the "createdAt" field.
func CreatedAtNotIn(vs ...time.Time) predicate.MomentImage {
	return predicate.MomentImage(sql.FieldNotIn(FieldCreatedAt, vs...))
}

// CreatedAtGT applies the GT predicate on the "createdAt" field.
func CreatedAtGT(v time.Time) predicate.MomentImage {
	return predicate.MomentImage(sql.FieldGT(FieldCreatedAt, v))
}

// CreatedAtGTE applies the GTE predicate on the "createdAt" field.
func CreatedAtGTE(v time.Time) predicate.MomentImage {
	return predicate.MomentImage(sql.FieldGTE(FieldCreatedAt, v))
}

// CreatedAtLT applies the LT predicate on the "createdAt" field.
func CreatedAtLT(v time.Time) predicate.MomentImage {
	return predicate.MomentImage(sql.FieldLT(FieldCreatedAt, v))
}

// CreatedAtLTE applies the LTE predicate on the "createdAt" field.
func CreatedAtLTE(v time.Time) predicate.MomentImage {
	return predicate.MomentImage(sql.FieldLTE(FieldCreatedAt, v))
}

// UpdatedAtEQ applies the EQ predicate on the "updatedAt" field.
func UpdatedAtEQ(v time.Time) predicate.MomentImage {
	return predicate.MomentImage(sql.FieldEQ(FieldUpdatedAt, v))
}

// UpdatedAtNEQ applies the NEQ predicate on the "updatedAt" field.
func UpdatedAtNEQ(v time.Time) predicate.MomentImage {
	return predicate.MomentImage(sql.FieldNEQ(FieldUpdatedAt, v))
}

// UpdatedAtIn applies the In predicate on the "updatedAt" field.
func UpdatedAtIn(vs ...time.Time) predicate.MomentImage {
	return predicate.MomentImage(sql.FieldIn(FieldUpdatedAt, vs...))
}

// UpdatedAtNotIn applies the NotIn predicate on the "updatedAt" field.
func UpdatedAtNotIn(vs ...time.Time) predicate.MomentImage {
	return predicate.MomentImage(sql.FieldNotIn(FieldUpdatedAt, vs...))
}

// UpdatedAtGT applies the GT predicate on the "updatedAt" field.
func UpdatedAtGT(v time.Time) predicate.MomentImage {
	return predicate.MomentImage(sql.FieldGT(FieldUpdatedAt, v))
}

// UpdatedAtGTE applies the GTE predicate on the "updatedAt" field.
func UpdatedAtGTE(v time.Time) predicate.MomentImage {
	return predicate.MomentImage(sql.FieldGTE(FieldUpdatedAt, v))
}

// UpdatedAtLT applies the LT predicate on the "updatedAt" field.
func UpdatedAtLT(v time.Time) predicate.MomentImage {
	return predicate.MomentImage(sql.FieldLT(FieldUpdatedAt, v))
}

// UpdatedAtLTE applies the LTE predicate on the "updatedAt" field.
func UpdatedAtLTE(v time.Time) predicate.MomentImage {
	return predicate.MomentImage(sql.FieldLTE(FieldUpdatedAt, v))
}

// HasImage applies the HasEdge predicate on the "image" edge.
func HasImage() predicate.MomentImage {
	return predicate.MomentImage(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, ImageTable, ImageColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasImageWith applies the HasEdge predicate on the "image" edge with a given conditions (other predicates).
func HasImageWith(preds ...predicate.Image) predicate.MomentImage {
	return predicate.MomentImage(func(s *sql.Selector) {
		step := newImageStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasMoment applies the HasEdge predicate on the "moment" edge.
func HasMoment() predicate.MomentImage {
	return predicate.MomentImage(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, MomentTable, MomentColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasMomentWith applies the HasEdge predicate on the "moment" edge with a given conditions (other predicates).
func HasMomentWith(preds ...predicate.Moment) predicate.MomentImage {
	return predicate.MomentImage(func(s *sql.Selector) {
		step := newMomentStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.MomentImage) predicate.MomentImage {
	return predicate.MomentImage(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.MomentImage) predicate.MomentImage {
	return predicate.MomentImage(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.MomentImage) predicate.MomentImage {
	return predicate.MomentImage(sql.NotPredicates(p))
}
