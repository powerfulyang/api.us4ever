// Code generated by ent, DO NOT EDIT.

package keep

import (
	"time"

	"api.us4ever/internal/ent/predicate"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

// ID filters vertices based on their ID field.
func ID(id string) predicate.Keep {
	return predicate.Keep(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id string) predicate.Keep {
	return predicate.Keep(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id string) predicate.Keep {
	return predicate.Keep(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...string) predicate.Keep {
	return predicate.Keep(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...string) predicate.Keep {
	return predicate.Keep(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id string) predicate.Keep {
	return predicate.Keep(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id string) predicate.Keep {
	return predicate.Keep(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id string) predicate.Keep {
	return predicate.Keep(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id string) predicate.Keep {
	return predicate.Keep(sql.FieldLTE(FieldID, id))
}

// IDEqualFold applies the EqualFold predicate on the ID field.
func IDEqualFold(id string) predicate.Keep {
	return predicate.Keep(sql.FieldEqualFold(FieldID, id))
}

// IDContainsFold applies the ContainsFold predicate on the ID field.
func IDContainsFold(id string) predicate.Keep {
	return predicate.Keep(sql.FieldContainsFold(FieldID, id))
}

// Title applies equality check predicate on the "title" field. It's identical to TitleEQ.
func Title(v string) predicate.Keep {
	return predicate.Keep(sql.FieldEQ(FieldTitle, v))
}

// Content applies equality check predicate on the "content" field. It's identical to ContentEQ.
func Content(v string) predicate.Keep {
	return predicate.Keep(sql.FieldEQ(FieldContent, v))
}

// Summary applies equality check predicate on the "summary" field. It's identical to SummaryEQ.
func Summary(v string) predicate.Keep {
	return predicate.Keep(sql.FieldEQ(FieldSummary, v))
}

// IsPublic applies equality check predicate on the "isPublic" field. It's identical to IsPublicEQ.
func IsPublic(v bool) predicate.Keep {
	return predicate.Keep(sql.FieldEQ(FieldIsPublic, v))
}

// Views applies equality check predicate on the "views" field. It's identical to ViewsEQ.
func Views(v int32) predicate.Keep {
	return predicate.Keep(sql.FieldEQ(FieldViews, v))
}

// Likes applies equality check predicate on the "likes" field. It's identical to LikesEQ.
func Likes(v int32) predicate.Keep {
	return predicate.Keep(sql.FieldEQ(FieldLikes, v))
}

// Category applies equality check predicate on the "category" field. It's identical to CategoryEQ.
func Category(v string) predicate.Keep {
	return predicate.Keep(sql.FieldEQ(FieldCategory, v))
}

// OwnerId applies equality check predicate on the "ownerId" field. It's identical to OwnerIdEQ.
func OwnerId(v string) predicate.Keep {
	return predicate.Keep(sql.FieldEQ(FieldOwnerId, v))
}

// CreatedAt applies equality check predicate on the "createdAt" field. It's identical to CreatedAtEQ.
func CreatedAt(v time.Time) predicate.Keep {
	return predicate.Keep(sql.FieldEQ(FieldCreatedAt, v))
}

// UpdatedAt applies equality check predicate on the "updatedAt" field. It's identical to UpdatedAtEQ.
func UpdatedAt(v time.Time) predicate.Keep {
	return predicate.Keep(sql.FieldEQ(FieldUpdatedAt, v))
}

// TitleEQ applies the EQ predicate on the "title" field.
func TitleEQ(v string) predicate.Keep {
	return predicate.Keep(sql.FieldEQ(FieldTitle, v))
}

// TitleNEQ applies the NEQ predicate on the "title" field.
func TitleNEQ(v string) predicate.Keep {
	return predicate.Keep(sql.FieldNEQ(FieldTitle, v))
}

// TitleIn applies the In predicate on the "title" field.
func TitleIn(vs ...string) predicate.Keep {
	return predicate.Keep(sql.FieldIn(FieldTitle, vs...))
}

// TitleNotIn applies the NotIn predicate on the "title" field.
func TitleNotIn(vs ...string) predicate.Keep {
	return predicate.Keep(sql.FieldNotIn(FieldTitle, vs...))
}

// TitleGT applies the GT predicate on the "title" field.
func TitleGT(v string) predicate.Keep {
	return predicate.Keep(sql.FieldGT(FieldTitle, v))
}

// TitleGTE applies the GTE predicate on the "title" field.
func TitleGTE(v string) predicate.Keep {
	return predicate.Keep(sql.FieldGTE(FieldTitle, v))
}

// TitleLT applies the LT predicate on the "title" field.
func TitleLT(v string) predicate.Keep {
	return predicate.Keep(sql.FieldLT(FieldTitle, v))
}

// TitleLTE applies the LTE predicate on the "title" field.
func TitleLTE(v string) predicate.Keep {
	return predicate.Keep(sql.FieldLTE(FieldTitle, v))
}

// TitleContains applies the Contains predicate on the "title" field.
func TitleContains(v string) predicate.Keep {
	return predicate.Keep(sql.FieldContains(FieldTitle, v))
}

// TitleHasPrefix applies the HasPrefix predicate on the "title" field.
func TitleHasPrefix(v string) predicate.Keep {
	return predicate.Keep(sql.FieldHasPrefix(FieldTitle, v))
}

// TitleHasSuffix applies the HasSuffix predicate on the "title" field.
func TitleHasSuffix(v string) predicate.Keep {
	return predicate.Keep(sql.FieldHasSuffix(FieldTitle, v))
}

// TitleEqualFold applies the EqualFold predicate on the "title" field.
func TitleEqualFold(v string) predicate.Keep {
	return predicate.Keep(sql.FieldEqualFold(FieldTitle, v))
}

// TitleContainsFold applies the ContainsFold predicate on the "title" field.
func TitleContainsFold(v string) predicate.Keep {
	return predicate.Keep(sql.FieldContainsFold(FieldTitle, v))
}

// ContentEQ applies the EQ predicate on the "content" field.
func ContentEQ(v string) predicate.Keep {
	return predicate.Keep(sql.FieldEQ(FieldContent, v))
}

// ContentNEQ applies the NEQ predicate on the "content" field.
func ContentNEQ(v string) predicate.Keep {
	return predicate.Keep(sql.FieldNEQ(FieldContent, v))
}

// ContentIn applies the In predicate on the "content" field.
func ContentIn(vs ...string) predicate.Keep {
	return predicate.Keep(sql.FieldIn(FieldContent, vs...))
}

// ContentNotIn applies the NotIn predicate on the "content" field.
func ContentNotIn(vs ...string) predicate.Keep {
	return predicate.Keep(sql.FieldNotIn(FieldContent, vs...))
}

// ContentGT applies the GT predicate on the "content" field.
func ContentGT(v string) predicate.Keep {
	return predicate.Keep(sql.FieldGT(FieldContent, v))
}

// ContentGTE applies the GTE predicate on the "content" field.
func ContentGTE(v string) predicate.Keep {
	return predicate.Keep(sql.FieldGTE(FieldContent, v))
}

// ContentLT applies the LT predicate on the "content" field.
func ContentLT(v string) predicate.Keep {
	return predicate.Keep(sql.FieldLT(FieldContent, v))
}

// ContentLTE applies the LTE predicate on the "content" field.
func ContentLTE(v string) predicate.Keep {
	return predicate.Keep(sql.FieldLTE(FieldContent, v))
}

// ContentContains applies the Contains predicate on the "content" field.
func ContentContains(v string) predicate.Keep {
	return predicate.Keep(sql.FieldContains(FieldContent, v))
}

// ContentHasPrefix applies the HasPrefix predicate on the "content" field.
func ContentHasPrefix(v string) predicate.Keep {
	return predicate.Keep(sql.FieldHasPrefix(FieldContent, v))
}

// ContentHasSuffix applies the HasSuffix predicate on the "content" field.
func ContentHasSuffix(v string) predicate.Keep {
	return predicate.Keep(sql.FieldHasSuffix(FieldContent, v))
}

// ContentEqualFold applies the EqualFold predicate on the "content" field.
func ContentEqualFold(v string) predicate.Keep {
	return predicate.Keep(sql.FieldEqualFold(FieldContent, v))
}

// ContentContainsFold applies the ContainsFold predicate on the "content" field.
func ContentContainsFold(v string) predicate.Keep {
	return predicate.Keep(sql.FieldContainsFold(FieldContent, v))
}

// SummaryEQ applies the EQ predicate on the "summary" field.
func SummaryEQ(v string) predicate.Keep {
	return predicate.Keep(sql.FieldEQ(FieldSummary, v))
}

// SummaryNEQ applies the NEQ predicate on the "summary" field.
func SummaryNEQ(v string) predicate.Keep {
	return predicate.Keep(sql.FieldNEQ(FieldSummary, v))
}

// SummaryIn applies the In predicate on the "summary" field.
func SummaryIn(vs ...string) predicate.Keep {
	return predicate.Keep(sql.FieldIn(FieldSummary, vs...))
}

// SummaryNotIn applies the NotIn predicate on the "summary" field.
func SummaryNotIn(vs ...string) predicate.Keep {
	return predicate.Keep(sql.FieldNotIn(FieldSummary, vs...))
}

// SummaryGT applies the GT predicate on the "summary" field.
func SummaryGT(v string) predicate.Keep {
	return predicate.Keep(sql.FieldGT(FieldSummary, v))
}

// SummaryGTE applies the GTE predicate on the "summary" field.
func SummaryGTE(v string) predicate.Keep {
	return predicate.Keep(sql.FieldGTE(FieldSummary, v))
}

// SummaryLT applies the LT predicate on the "summary" field.
func SummaryLT(v string) predicate.Keep {
	return predicate.Keep(sql.FieldLT(FieldSummary, v))
}

// SummaryLTE applies the LTE predicate on the "summary" field.
func SummaryLTE(v string) predicate.Keep {
	return predicate.Keep(sql.FieldLTE(FieldSummary, v))
}

// SummaryContains applies the Contains predicate on the "summary" field.
func SummaryContains(v string) predicate.Keep {
	return predicate.Keep(sql.FieldContains(FieldSummary, v))
}

// SummaryHasPrefix applies the HasPrefix predicate on the "summary" field.
func SummaryHasPrefix(v string) predicate.Keep {
	return predicate.Keep(sql.FieldHasPrefix(FieldSummary, v))
}

// SummaryHasSuffix applies the HasSuffix predicate on the "summary" field.
func SummaryHasSuffix(v string) predicate.Keep {
	return predicate.Keep(sql.FieldHasSuffix(FieldSummary, v))
}

// SummaryEqualFold applies the EqualFold predicate on the "summary" field.
func SummaryEqualFold(v string) predicate.Keep {
	return predicate.Keep(sql.FieldEqualFold(FieldSummary, v))
}

// SummaryContainsFold applies the ContainsFold predicate on the "summary" field.
func SummaryContainsFold(v string) predicate.Keep {
	return predicate.Keep(sql.FieldContainsFold(FieldSummary, v))
}

// IsPublicEQ applies the EQ predicate on the "isPublic" field.
func IsPublicEQ(v bool) predicate.Keep {
	return predicate.Keep(sql.FieldEQ(FieldIsPublic, v))
}

// IsPublicNEQ applies the NEQ predicate on the "isPublic" field.
func IsPublicNEQ(v bool) predicate.Keep {
	return predicate.Keep(sql.FieldNEQ(FieldIsPublic, v))
}

// ViewsEQ applies the EQ predicate on the "views" field.
func ViewsEQ(v int32) predicate.Keep {
	return predicate.Keep(sql.FieldEQ(FieldViews, v))
}

// ViewsNEQ applies the NEQ predicate on the "views" field.
func ViewsNEQ(v int32) predicate.Keep {
	return predicate.Keep(sql.FieldNEQ(FieldViews, v))
}

// ViewsIn applies the In predicate on the "views" field.
func ViewsIn(vs ...int32) predicate.Keep {
	return predicate.Keep(sql.FieldIn(FieldViews, vs...))
}

// ViewsNotIn applies the NotIn predicate on the "views" field.
func ViewsNotIn(vs ...int32) predicate.Keep {
	return predicate.Keep(sql.FieldNotIn(FieldViews, vs...))
}

// ViewsGT applies the GT predicate on the "views" field.
func ViewsGT(v int32) predicate.Keep {
	return predicate.Keep(sql.FieldGT(FieldViews, v))
}

// ViewsGTE applies the GTE predicate on the "views" field.
func ViewsGTE(v int32) predicate.Keep {
	return predicate.Keep(sql.FieldGTE(FieldViews, v))
}

// ViewsLT applies the LT predicate on the "views" field.
func ViewsLT(v int32) predicate.Keep {
	return predicate.Keep(sql.FieldLT(FieldViews, v))
}

// ViewsLTE applies the LTE predicate on the "views" field.
func ViewsLTE(v int32) predicate.Keep {
	return predicate.Keep(sql.FieldLTE(FieldViews, v))
}

// LikesEQ applies the EQ predicate on the "likes" field.
func LikesEQ(v int32) predicate.Keep {
	return predicate.Keep(sql.FieldEQ(FieldLikes, v))
}

// LikesNEQ applies the NEQ predicate on the "likes" field.
func LikesNEQ(v int32) predicate.Keep {
	return predicate.Keep(sql.FieldNEQ(FieldLikes, v))
}

// LikesIn applies the In predicate on the "likes" field.
func LikesIn(vs ...int32) predicate.Keep {
	return predicate.Keep(sql.FieldIn(FieldLikes, vs...))
}

// LikesNotIn applies the NotIn predicate on the "likes" field.
func LikesNotIn(vs ...int32) predicate.Keep {
	return predicate.Keep(sql.FieldNotIn(FieldLikes, vs...))
}

// LikesGT applies the GT predicate on the "likes" field.
func LikesGT(v int32) predicate.Keep {
	return predicate.Keep(sql.FieldGT(FieldLikes, v))
}

// LikesGTE applies the GTE predicate on the "likes" field.
func LikesGTE(v int32) predicate.Keep {
	return predicate.Keep(sql.FieldGTE(FieldLikes, v))
}

// LikesLT applies the LT predicate on the "likes" field.
func LikesLT(v int32) predicate.Keep {
	return predicate.Keep(sql.FieldLT(FieldLikes, v))
}

// LikesLTE applies the LTE predicate on the "likes" field.
func LikesLTE(v int32) predicate.Keep {
	return predicate.Keep(sql.FieldLTE(FieldLikes, v))
}

// CategoryEQ applies the EQ predicate on the "category" field.
func CategoryEQ(v string) predicate.Keep {
	return predicate.Keep(sql.FieldEQ(FieldCategory, v))
}

// CategoryNEQ applies the NEQ predicate on the "category" field.
func CategoryNEQ(v string) predicate.Keep {
	return predicate.Keep(sql.FieldNEQ(FieldCategory, v))
}

// CategoryIn applies the In predicate on the "category" field.
func CategoryIn(vs ...string) predicate.Keep {
	return predicate.Keep(sql.FieldIn(FieldCategory, vs...))
}

// CategoryNotIn applies the NotIn predicate on the "category" field.
func CategoryNotIn(vs ...string) predicate.Keep {
	return predicate.Keep(sql.FieldNotIn(FieldCategory, vs...))
}

// CategoryGT applies the GT predicate on the "category" field.
func CategoryGT(v string) predicate.Keep {
	return predicate.Keep(sql.FieldGT(FieldCategory, v))
}

// CategoryGTE applies the GTE predicate on the "category" field.
func CategoryGTE(v string) predicate.Keep {
	return predicate.Keep(sql.FieldGTE(FieldCategory, v))
}

// CategoryLT applies the LT predicate on the "category" field.
func CategoryLT(v string) predicate.Keep {
	return predicate.Keep(sql.FieldLT(FieldCategory, v))
}

// CategoryLTE applies the LTE predicate on the "category" field.
func CategoryLTE(v string) predicate.Keep {
	return predicate.Keep(sql.FieldLTE(FieldCategory, v))
}

// CategoryContains applies the Contains predicate on the "category" field.
func CategoryContains(v string) predicate.Keep {
	return predicate.Keep(sql.FieldContains(FieldCategory, v))
}

// CategoryHasPrefix applies the HasPrefix predicate on the "category" field.
func CategoryHasPrefix(v string) predicate.Keep {
	return predicate.Keep(sql.FieldHasPrefix(FieldCategory, v))
}

// CategoryHasSuffix applies the HasSuffix predicate on the "category" field.
func CategoryHasSuffix(v string) predicate.Keep {
	return predicate.Keep(sql.FieldHasSuffix(FieldCategory, v))
}

// CategoryEqualFold applies the EqualFold predicate on the "category" field.
func CategoryEqualFold(v string) predicate.Keep {
	return predicate.Keep(sql.FieldEqualFold(FieldCategory, v))
}

// CategoryContainsFold applies the ContainsFold predicate on the "category" field.
func CategoryContainsFold(v string) predicate.Keep {
	return predicate.Keep(sql.FieldContainsFold(FieldCategory, v))
}

// OwnerIdEQ applies the EQ predicate on the "ownerId" field.
func OwnerIdEQ(v string) predicate.Keep {
	return predicate.Keep(sql.FieldEQ(FieldOwnerId, v))
}

// OwnerIdNEQ applies the NEQ predicate on the "ownerId" field.
func OwnerIdNEQ(v string) predicate.Keep {
	return predicate.Keep(sql.FieldNEQ(FieldOwnerId, v))
}

// OwnerIdIn applies the In predicate on the "ownerId" field.
func OwnerIdIn(vs ...string) predicate.Keep {
	return predicate.Keep(sql.FieldIn(FieldOwnerId, vs...))
}

// OwnerIdNotIn applies the NotIn predicate on the "ownerId" field.
func OwnerIdNotIn(vs ...string) predicate.Keep {
	return predicate.Keep(sql.FieldNotIn(FieldOwnerId, vs...))
}

// OwnerIdGT applies the GT predicate on the "ownerId" field.
func OwnerIdGT(v string) predicate.Keep {
	return predicate.Keep(sql.FieldGT(FieldOwnerId, v))
}

// OwnerIdGTE applies the GTE predicate on the "ownerId" field.
func OwnerIdGTE(v string) predicate.Keep {
	return predicate.Keep(sql.FieldGTE(FieldOwnerId, v))
}

// OwnerIdLT applies the LT predicate on the "ownerId" field.
func OwnerIdLT(v string) predicate.Keep {
	return predicate.Keep(sql.FieldLT(FieldOwnerId, v))
}

// OwnerIdLTE applies the LTE predicate on the "ownerId" field.
func OwnerIdLTE(v string) predicate.Keep {
	return predicate.Keep(sql.FieldLTE(FieldOwnerId, v))
}

// OwnerIdContains applies the Contains predicate on the "ownerId" field.
func OwnerIdContains(v string) predicate.Keep {
	return predicate.Keep(sql.FieldContains(FieldOwnerId, v))
}

// OwnerIdHasPrefix applies the HasPrefix predicate on the "ownerId" field.
func OwnerIdHasPrefix(v string) predicate.Keep {
	return predicate.Keep(sql.FieldHasPrefix(FieldOwnerId, v))
}

// OwnerIdHasSuffix applies the HasSuffix predicate on the "ownerId" field.
func OwnerIdHasSuffix(v string) predicate.Keep {
	return predicate.Keep(sql.FieldHasSuffix(FieldOwnerId, v))
}

// OwnerIdIsNil applies the IsNil predicate on the "ownerId" field.
func OwnerIdIsNil() predicate.Keep {
	return predicate.Keep(sql.FieldIsNull(FieldOwnerId))
}

// OwnerIdNotNil applies the NotNil predicate on the "ownerId" field.
func OwnerIdNotNil() predicate.Keep {
	return predicate.Keep(sql.FieldNotNull(FieldOwnerId))
}

// OwnerIdEqualFold applies the EqualFold predicate on the "ownerId" field.
func OwnerIdEqualFold(v string) predicate.Keep {
	return predicate.Keep(sql.FieldEqualFold(FieldOwnerId, v))
}

// OwnerIdContainsFold applies the ContainsFold predicate on the "ownerId" field.
func OwnerIdContainsFold(v string) predicate.Keep {
	return predicate.Keep(sql.FieldContainsFold(FieldOwnerId, v))
}

// CreatedAtEQ applies the EQ predicate on the "createdAt" field.
func CreatedAtEQ(v time.Time) predicate.Keep {
	return predicate.Keep(sql.FieldEQ(FieldCreatedAt, v))
}

// CreatedAtNEQ applies the NEQ predicate on the "createdAt" field.
func CreatedAtNEQ(v time.Time) predicate.Keep {
	return predicate.Keep(sql.FieldNEQ(FieldCreatedAt, v))
}

// CreatedAtIn applies the In predicate on the "createdAt" field.
func CreatedAtIn(vs ...time.Time) predicate.Keep {
	return predicate.Keep(sql.FieldIn(FieldCreatedAt, vs...))
}

// CreatedAtNotIn applies the NotIn predicate on the "createdAt" field.
func CreatedAtNotIn(vs ...time.Time) predicate.Keep {
	return predicate.Keep(sql.FieldNotIn(FieldCreatedAt, vs...))
}

// CreatedAtGT applies the GT predicate on the "createdAt" field.
func CreatedAtGT(v time.Time) predicate.Keep {
	return predicate.Keep(sql.FieldGT(FieldCreatedAt, v))
}

// CreatedAtGTE applies the GTE predicate on the "createdAt" field.
func CreatedAtGTE(v time.Time) predicate.Keep {
	return predicate.Keep(sql.FieldGTE(FieldCreatedAt, v))
}

// CreatedAtLT applies the LT predicate on the "createdAt" field.
func CreatedAtLT(v time.Time) predicate.Keep {
	return predicate.Keep(sql.FieldLT(FieldCreatedAt, v))
}

// CreatedAtLTE applies the LTE predicate on the "createdAt" field.
func CreatedAtLTE(v time.Time) predicate.Keep {
	return predicate.Keep(sql.FieldLTE(FieldCreatedAt, v))
}

// UpdatedAtEQ applies the EQ predicate on the "updatedAt" field.
func UpdatedAtEQ(v time.Time) predicate.Keep {
	return predicate.Keep(sql.FieldEQ(FieldUpdatedAt, v))
}

// UpdatedAtNEQ applies the NEQ predicate on the "updatedAt" field.
func UpdatedAtNEQ(v time.Time) predicate.Keep {
	return predicate.Keep(sql.FieldNEQ(FieldUpdatedAt, v))
}

// UpdatedAtIn applies the In predicate on the "updatedAt" field.
func UpdatedAtIn(vs ...time.Time) predicate.Keep {
	return predicate.Keep(sql.FieldIn(FieldUpdatedAt, vs...))
}

// UpdatedAtNotIn applies the NotIn predicate on the "updatedAt" field.
func UpdatedAtNotIn(vs ...time.Time) predicate.Keep {
	return predicate.Keep(sql.FieldNotIn(FieldUpdatedAt, vs...))
}

// UpdatedAtGT applies the GT predicate on the "updatedAt" field.
func UpdatedAtGT(v time.Time) predicate.Keep {
	return predicate.Keep(sql.FieldGT(FieldUpdatedAt, v))
}

// UpdatedAtGTE applies the GTE predicate on the "updatedAt" field.
func UpdatedAtGTE(v time.Time) predicate.Keep {
	return predicate.Keep(sql.FieldGTE(FieldUpdatedAt, v))
}

// UpdatedAtLT applies the LT predicate on the "updatedAt" field.
func UpdatedAtLT(v time.Time) predicate.Keep {
	return predicate.Keep(sql.FieldLT(FieldUpdatedAt, v))
}

// UpdatedAtLTE applies the LTE predicate on the "updatedAt" field.
func UpdatedAtLTE(v time.Time) predicate.Keep {
	return predicate.Keep(sql.FieldLTE(FieldUpdatedAt, v))
}

// HasUser applies the HasEdge predicate on the "user" edge.
func HasUser() predicate.Keep {
	return predicate.Keep(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, UserTable, UserColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasUserWith applies the HasEdge predicate on the "user" edge with a given conditions (other predicates).
func HasUserWith(preds ...predicate.User) predicate.Keep {
	return predicate.Keep(func(s *sql.Selector) {
		step := newUserStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.Keep) predicate.Keep {
	return predicate.Keep(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.Keep) predicate.Keep {
	return predicate.Keep(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.Keep) predicate.Keep {
	return predicate.Keep(sql.NotPredicates(p))
}
