// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"
	"math"

	"api.us4ever/internal/ent/keep"
	"api.us4ever/internal/ent/predicate"
	"api.us4ever/internal/ent/user"
	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// KeepQuery is the builder for querying Keep entities.
type KeepQuery struct {
	config
	ctx        *QueryContext
	order      []keep.OrderOption
	inters     []Interceptor
	predicates []predicate.Keep
	withUser   *UserQuery
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the KeepQuery builder.
func (kq *KeepQuery) Where(ps ...predicate.Keep) *KeepQuery {
	kq.predicates = append(kq.predicates, ps...)
	return kq
}

// Limit the number of records to be returned by this query.
func (kq *KeepQuery) Limit(limit int) *KeepQuery {
	kq.ctx.Limit = &limit
	return kq
}

// Offset to start from.
func (kq *KeepQuery) Offset(offset int) *KeepQuery {
	kq.ctx.Offset = &offset
	return kq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (kq *KeepQuery) Unique(unique bool) *KeepQuery {
	kq.ctx.Unique = &unique
	return kq
}

// Order specifies how the records should be ordered.
func (kq *KeepQuery) Order(o ...keep.OrderOption) *KeepQuery {
	kq.order = append(kq.order, o...)
	return kq
}

// QueryUser chains the current query on the "user" edge.
func (kq *KeepQuery) QueryUser() *UserQuery {
	query := (&UserClient{config: kq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := kq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := kq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(keep.Table, keep.FieldID, selector),
			sqlgraph.To(user.Table, user.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, keep.UserTable, keep.UserColumn),
		)
		fromU = sqlgraph.SetNeighbors(kq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first Keep entity from the query.
// Returns a *NotFoundError when no Keep was found.
func (kq *KeepQuery) First(ctx context.Context) (*Keep, error) {
	nodes, err := kq.Limit(1).All(setContextOp(ctx, kq.ctx, ent.OpQueryFirst))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{keep.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (kq *KeepQuery) FirstX(ctx context.Context) *Keep {
	node, err := kq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first Keep ID from the query.
// Returns a *NotFoundError when no Keep ID was found.
func (kq *KeepQuery) FirstID(ctx context.Context) (id string, err error) {
	var ids []string
	if ids, err = kq.Limit(1).IDs(setContextOp(ctx, kq.ctx, ent.OpQueryFirstID)); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{keep.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (kq *KeepQuery) FirstIDX(ctx context.Context) string {
	id, err := kq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single Keep entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one Keep entity is found.
// Returns a *NotFoundError when no Keep entities are found.
func (kq *KeepQuery) Only(ctx context.Context) (*Keep, error) {
	nodes, err := kq.Limit(2).All(setContextOp(ctx, kq.ctx, ent.OpQueryOnly))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{keep.Label}
	default:
		return nil, &NotSingularError{keep.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (kq *KeepQuery) OnlyX(ctx context.Context) *Keep {
	node, err := kq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only Keep ID in the query.
// Returns a *NotSingularError when more than one Keep ID is found.
// Returns a *NotFoundError when no entities are found.
func (kq *KeepQuery) OnlyID(ctx context.Context) (id string, err error) {
	var ids []string
	if ids, err = kq.Limit(2).IDs(setContextOp(ctx, kq.ctx, ent.OpQueryOnlyID)); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{keep.Label}
	default:
		err = &NotSingularError{keep.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (kq *KeepQuery) OnlyIDX(ctx context.Context) string {
	id, err := kq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of Keeps.
func (kq *KeepQuery) All(ctx context.Context) ([]*Keep, error) {
	ctx = setContextOp(ctx, kq.ctx, ent.OpQueryAll)
	if err := kq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*Keep, *KeepQuery]()
	return withInterceptors[[]*Keep](ctx, kq, qr, kq.inters)
}

// AllX is like All, but panics if an error occurs.
func (kq *KeepQuery) AllX(ctx context.Context) []*Keep {
	nodes, err := kq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of Keep IDs.
func (kq *KeepQuery) IDs(ctx context.Context) (ids []string, err error) {
	if kq.ctx.Unique == nil && kq.path != nil {
		kq.Unique(true)
	}
	ctx = setContextOp(ctx, kq.ctx, ent.OpQueryIDs)
	if err = kq.Select(keep.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (kq *KeepQuery) IDsX(ctx context.Context) []string {
	ids, err := kq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (kq *KeepQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, kq.ctx, ent.OpQueryCount)
	if err := kq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, kq, querierCount[*KeepQuery](), kq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (kq *KeepQuery) CountX(ctx context.Context) int {
	count, err := kq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (kq *KeepQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, kq.ctx, ent.OpQueryExist)
	switch _, err := kq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("ent: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (kq *KeepQuery) ExistX(ctx context.Context) bool {
	exist, err := kq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the KeepQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (kq *KeepQuery) Clone() *KeepQuery {
	if kq == nil {
		return nil
	}
	return &KeepQuery{
		config:     kq.config,
		ctx:        kq.ctx.Clone(),
		order:      append([]keep.OrderOption{}, kq.order...),
		inters:     append([]Interceptor{}, kq.inters...),
		predicates: append([]predicate.Keep{}, kq.predicates...),
		withUser:   kq.withUser.Clone(),
		// clone intermediate query.
		sql:  kq.sql.Clone(),
		path: kq.path,
	}
}

// WithUser tells the query-builder to eager-load the nodes that are connected to
// the "user" edge. The optional arguments are used to configure the query builder of the edge.
func (kq *KeepQuery) WithUser(opts ...func(*UserQuery)) *KeepQuery {
	query := (&UserClient{config: kq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	kq.withUser = query
	return kq
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		Title string `json:"title,omitempty"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.Keep.Query().
//		GroupBy(keep.FieldTitle).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
func (kq *KeepQuery) GroupBy(field string, fields ...string) *KeepGroupBy {
	kq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &KeepGroupBy{build: kq}
	grbuild.flds = &kq.ctx.Fields
	grbuild.label = keep.Label
	grbuild.scan = grbuild.Scan
	return grbuild
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		Title string `json:"title,omitempty"`
//	}
//
//	client.Keep.Query().
//		Select(keep.FieldTitle).
//		Scan(ctx, &v)
func (kq *KeepQuery) Select(fields ...string) *KeepSelect {
	kq.ctx.Fields = append(kq.ctx.Fields, fields...)
	sbuild := &KeepSelect{KeepQuery: kq}
	sbuild.label = keep.Label
	sbuild.flds, sbuild.scan = &kq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a KeepSelect configured with the given aggregations.
func (kq *KeepQuery) Aggregate(fns ...AggregateFunc) *KeepSelect {
	return kq.Select().Aggregate(fns...)
}

func (kq *KeepQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range kq.inters {
		if inter == nil {
			return fmt.Errorf("ent: uninitialized interceptor (forgotten import ent/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, kq); err != nil {
				return err
			}
		}
	}
	for _, f := range kq.ctx.Fields {
		if !keep.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if kq.path != nil {
		prev, err := kq.path(ctx)
		if err != nil {
			return err
		}
		kq.sql = prev
	}
	return nil
}

func (kq *KeepQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*Keep, error) {
	var (
		nodes       = []*Keep{}
		_spec       = kq.querySpec()
		loadedTypes = [1]bool{
			kq.withUser != nil,
		}
	)
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*Keep).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &Keep{config: kq.config}
		nodes = append(nodes, node)
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(columns, values)
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, kq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	if query := kq.withUser; query != nil {
		if err := kq.loadUser(ctx, query, nodes, nil,
			func(n *Keep, e *User) { n.Edges.User = e }); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func (kq *KeepQuery) loadUser(ctx context.Context, query *UserQuery, nodes []*Keep, init func(*Keep), assign func(*Keep, *User)) error {
	ids := make([]string, 0, len(nodes))
	nodeids := make(map[string][]*Keep)
	for i := range nodes {
		fk := nodes[i].OwnerId
		if _, ok := nodeids[fk]; !ok {
			ids = append(ids, fk)
		}
		nodeids[fk] = append(nodeids[fk], nodes[i])
	}
	if len(ids) == 0 {
		return nil
	}
	query.Where(user.IDIn(ids...))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nodeids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "ownerId" returned %v`, n.ID)
		}
		for i := range nodes {
			assign(nodes[i], n)
		}
	}
	return nil
}

func (kq *KeepQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := kq.querySpec()
	_spec.Node.Columns = kq.ctx.Fields
	if len(kq.ctx.Fields) > 0 {
		_spec.Unique = kq.ctx.Unique != nil && *kq.ctx.Unique
	}
	return sqlgraph.CountNodes(ctx, kq.driver, _spec)
}

func (kq *KeepQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := sqlgraph.NewQuerySpec(keep.Table, keep.Columns, sqlgraph.NewFieldSpec(keep.FieldID, field.TypeString))
	_spec.From = kq.sql
	if unique := kq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	} else if kq.path != nil {
		_spec.Unique = true
	}
	if fields := kq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, keep.FieldID)
		for i := range fields {
			if fields[i] != keep.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
		if kq.withUser != nil {
			_spec.Node.AddColumnOnce(keep.FieldOwnerId)
		}
	}
	if ps := kq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := kq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := kq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := kq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (kq *KeepQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(kq.driver.Dialect())
	t1 := builder.Table(keep.Table)
	columns := kq.ctx.Fields
	if len(columns) == 0 {
		columns = keep.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if kq.sql != nil {
		selector = kq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if kq.ctx.Unique != nil && *kq.ctx.Unique {
		selector.Distinct()
	}
	for _, p := range kq.predicates {
		p(selector)
	}
	for _, p := range kq.order {
		p(selector)
	}
	if offset := kq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := kq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// KeepGroupBy is the group-by builder for Keep entities.
type KeepGroupBy struct {
	selector
	build *KeepQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (kgb *KeepGroupBy) Aggregate(fns ...AggregateFunc) *KeepGroupBy {
	kgb.fns = append(kgb.fns, fns...)
	return kgb
}

// Scan applies the selector query and scans the result into the given value.
func (kgb *KeepGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, kgb.build.ctx, ent.OpQueryGroupBy)
	if err := kgb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*KeepQuery, *KeepGroupBy](ctx, kgb.build, kgb, kgb.build.inters, v)
}

func (kgb *KeepGroupBy) sqlScan(ctx context.Context, root *KeepQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(kgb.fns))
	for _, fn := range kgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*kgb.flds)+len(kgb.fns))
		for _, f := range *kgb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*kgb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := kgb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// KeepSelect is the builder for selecting fields of Keep entities.
type KeepSelect struct {
	*KeepQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (ks *KeepSelect) Aggregate(fns ...AggregateFunc) *KeepSelect {
	ks.fns = append(ks.fns, fns...)
	return ks
}

// Scan applies the selector query and scans the result into the given value.
func (ks *KeepSelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, ks.ctx, ent.OpQuerySelect)
	if err := ks.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*KeepQuery, *KeepSelect](ctx, ks.KeepQuery, ks, ks.inters, v)
}

func (ks *KeepSelect) sqlScan(ctx context.Context, root *KeepQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(ks.fns))
	for _, fn := range ks.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*ks.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := ks.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}
