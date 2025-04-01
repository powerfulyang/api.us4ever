// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"
	"math"

	"api.us4ever/internal/ent/mindmap"
	"api.us4ever/internal/ent/predicate"
	"api.us4ever/internal/ent/user"
	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// MindmapQuery is the builder for querying Mindmap entities.
type MindmapQuery struct {
	config
	ctx        *QueryContext
	order      []mindmap.OrderOption
	inters     []Interceptor
	predicates []predicate.Mindmap
	withUser   *UserQuery
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the MindmapQuery builder.
func (mq *MindmapQuery) Where(ps ...predicate.Mindmap) *MindmapQuery {
	mq.predicates = append(mq.predicates, ps...)
	return mq
}

// Limit the number of records to be returned by this query.
func (mq *MindmapQuery) Limit(limit int) *MindmapQuery {
	mq.ctx.Limit = &limit
	return mq
}

// Offset to start from.
func (mq *MindmapQuery) Offset(offset int) *MindmapQuery {
	mq.ctx.Offset = &offset
	return mq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (mq *MindmapQuery) Unique(unique bool) *MindmapQuery {
	mq.ctx.Unique = &unique
	return mq
}

// Order specifies how the records should be ordered.
func (mq *MindmapQuery) Order(o ...mindmap.OrderOption) *MindmapQuery {
	mq.order = append(mq.order, o...)
	return mq
}

// QueryUser chains the current query on the "user" edge.
func (mq *MindmapQuery) QueryUser() *UserQuery {
	query := (&UserClient{config: mq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := mq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := mq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(mindmap.Table, mindmap.FieldID, selector),
			sqlgraph.To(user.Table, user.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, mindmap.UserTable, mindmap.UserColumn),
		)
		fromU = sqlgraph.SetNeighbors(mq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first Mindmap entity from the query.
// Returns a *NotFoundError when no Mindmap was found.
func (mq *MindmapQuery) First(ctx context.Context) (*Mindmap, error) {
	nodes, err := mq.Limit(1).All(setContextOp(ctx, mq.ctx, ent.OpQueryFirst))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{mindmap.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (mq *MindmapQuery) FirstX(ctx context.Context) *Mindmap {
	node, err := mq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first Mindmap ID from the query.
// Returns a *NotFoundError when no Mindmap ID was found.
func (mq *MindmapQuery) FirstID(ctx context.Context) (id string, err error) {
	var ids []string
	if ids, err = mq.Limit(1).IDs(setContextOp(ctx, mq.ctx, ent.OpQueryFirstID)); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{mindmap.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (mq *MindmapQuery) FirstIDX(ctx context.Context) string {
	id, err := mq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single Mindmap entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one Mindmap entity is found.
// Returns a *NotFoundError when no Mindmap entities are found.
func (mq *MindmapQuery) Only(ctx context.Context) (*Mindmap, error) {
	nodes, err := mq.Limit(2).All(setContextOp(ctx, mq.ctx, ent.OpQueryOnly))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{mindmap.Label}
	default:
		return nil, &NotSingularError{mindmap.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (mq *MindmapQuery) OnlyX(ctx context.Context) *Mindmap {
	node, err := mq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only Mindmap ID in the query.
// Returns a *NotSingularError when more than one Mindmap ID is found.
// Returns a *NotFoundError when no entities are found.
func (mq *MindmapQuery) OnlyID(ctx context.Context) (id string, err error) {
	var ids []string
	if ids, err = mq.Limit(2).IDs(setContextOp(ctx, mq.ctx, ent.OpQueryOnlyID)); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{mindmap.Label}
	default:
		err = &NotSingularError{mindmap.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (mq *MindmapQuery) OnlyIDX(ctx context.Context) string {
	id, err := mq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of Mindmaps.
func (mq *MindmapQuery) All(ctx context.Context) ([]*Mindmap, error) {
	ctx = setContextOp(ctx, mq.ctx, ent.OpQueryAll)
	if err := mq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*Mindmap, *MindmapQuery]()
	return withInterceptors[[]*Mindmap](ctx, mq, qr, mq.inters)
}

// AllX is like All, but panics if an error occurs.
func (mq *MindmapQuery) AllX(ctx context.Context) []*Mindmap {
	nodes, err := mq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of Mindmap IDs.
func (mq *MindmapQuery) IDs(ctx context.Context) (ids []string, err error) {
	if mq.ctx.Unique == nil && mq.path != nil {
		mq.Unique(true)
	}
	ctx = setContextOp(ctx, mq.ctx, ent.OpQueryIDs)
	if err = mq.Select(mindmap.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (mq *MindmapQuery) IDsX(ctx context.Context) []string {
	ids, err := mq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (mq *MindmapQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, mq.ctx, ent.OpQueryCount)
	if err := mq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, mq, querierCount[*MindmapQuery](), mq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (mq *MindmapQuery) CountX(ctx context.Context) int {
	count, err := mq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (mq *MindmapQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, mq.ctx, ent.OpQueryExist)
	switch _, err := mq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("ent: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (mq *MindmapQuery) ExistX(ctx context.Context) bool {
	exist, err := mq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the MindmapQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (mq *MindmapQuery) Clone() *MindmapQuery {
	if mq == nil {
		return nil
	}
	return &MindmapQuery{
		config:     mq.config,
		ctx:        mq.ctx.Clone(),
		order:      append([]mindmap.OrderOption{}, mq.order...),
		inters:     append([]Interceptor{}, mq.inters...),
		predicates: append([]predicate.Mindmap{}, mq.predicates...),
		withUser:   mq.withUser.Clone(),
		// clone intermediate query.
		sql:  mq.sql.Clone(),
		path: mq.path,
	}
}

// WithUser tells the query-builder to eager-load the nodes that are connected to
// the "user" edge. The optional arguments are used to configure the query builder of the edge.
func (mq *MindmapQuery) WithUser(opts ...func(*UserQuery)) *MindmapQuery {
	query := (&UserClient{config: mq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	mq.withUser = query
	return mq
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
//	client.Mindmap.Query().
//		GroupBy(mindmap.FieldTitle).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
func (mq *MindmapQuery) GroupBy(field string, fields ...string) *MindmapGroupBy {
	mq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &MindmapGroupBy{build: mq}
	grbuild.flds = &mq.ctx.Fields
	grbuild.label = mindmap.Label
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
//	client.Mindmap.Query().
//		Select(mindmap.FieldTitle).
//		Scan(ctx, &v)
func (mq *MindmapQuery) Select(fields ...string) *MindmapSelect {
	mq.ctx.Fields = append(mq.ctx.Fields, fields...)
	sbuild := &MindmapSelect{MindmapQuery: mq}
	sbuild.label = mindmap.Label
	sbuild.flds, sbuild.scan = &mq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a MindmapSelect configured with the given aggregations.
func (mq *MindmapQuery) Aggregate(fns ...AggregateFunc) *MindmapSelect {
	return mq.Select().Aggregate(fns...)
}

func (mq *MindmapQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range mq.inters {
		if inter == nil {
			return fmt.Errorf("ent: uninitialized interceptor (forgotten import ent/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, mq); err != nil {
				return err
			}
		}
	}
	for _, f := range mq.ctx.Fields {
		if !mindmap.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if mq.path != nil {
		prev, err := mq.path(ctx)
		if err != nil {
			return err
		}
		mq.sql = prev
	}
	return nil
}

func (mq *MindmapQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*Mindmap, error) {
	var (
		nodes       = []*Mindmap{}
		_spec       = mq.querySpec()
		loadedTypes = [1]bool{
			mq.withUser != nil,
		}
	)
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*Mindmap).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &Mindmap{config: mq.config}
		nodes = append(nodes, node)
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(columns, values)
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, mq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	if query := mq.withUser; query != nil {
		if err := mq.loadUser(ctx, query, nodes, nil,
			func(n *Mindmap, e *User) { n.Edges.User = e }); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func (mq *MindmapQuery) loadUser(ctx context.Context, query *UserQuery, nodes []*Mindmap, init func(*Mindmap), assign func(*Mindmap, *User)) error {
	ids := make([]string, 0, len(nodes))
	nodeids := make(map[string][]*Mindmap)
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

func (mq *MindmapQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := mq.querySpec()
	_spec.Node.Columns = mq.ctx.Fields
	if len(mq.ctx.Fields) > 0 {
		_spec.Unique = mq.ctx.Unique != nil && *mq.ctx.Unique
	}
	return sqlgraph.CountNodes(ctx, mq.driver, _spec)
}

func (mq *MindmapQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := sqlgraph.NewQuerySpec(mindmap.Table, mindmap.Columns, sqlgraph.NewFieldSpec(mindmap.FieldID, field.TypeString))
	_spec.From = mq.sql
	if unique := mq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	} else if mq.path != nil {
		_spec.Unique = true
	}
	if fields := mq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, mindmap.FieldID)
		for i := range fields {
			if fields[i] != mindmap.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
		if mq.withUser != nil {
			_spec.Node.AddColumnOnce(mindmap.FieldOwnerId)
		}
	}
	if ps := mq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := mq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := mq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := mq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (mq *MindmapQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(mq.driver.Dialect())
	t1 := builder.Table(mindmap.Table)
	columns := mq.ctx.Fields
	if len(columns) == 0 {
		columns = mindmap.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if mq.sql != nil {
		selector = mq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if mq.ctx.Unique != nil && *mq.ctx.Unique {
		selector.Distinct()
	}
	for _, p := range mq.predicates {
		p(selector)
	}
	for _, p := range mq.order {
		p(selector)
	}
	if offset := mq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := mq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// MindmapGroupBy is the group-by builder for Mindmap entities.
type MindmapGroupBy struct {
	selector
	build *MindmapQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (mgb *MindmapGroupBy) Aggregate(fns ...AggregateFunc) *MindmapGroupBy {
	mgb.fns = append(mgb.fns, fns...)
	return mgb
}

// Scan applies the selector query and scans the result into the given value.
func (mgb *MindmapGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, mgb.build.ctx, ent.OpQueryGroupBy)
	if err := mgb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*MindmapQuery, *MindmapGroupBy](ctx, mgb.build, mgb, mgb.build.inters, v)
}

func (mgb *MindmapGroupBy) sqlScan(ctx context.Context, root *MindmapQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(mgb.fns))
	for _, fn := range mgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*mgb.flds)+len(mgb.fns))
		for _, f := range *mgb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*mgb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := mgb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// MindmapSelect is the builder for selecting fields of Mindmap entities.
type MindmapSelect struct {
	*MindmapQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (ms *MindmapSelect) Aggregate(fns ...AggregateFunc) *MindmapSelect {
	ms.fns = append(ms.fns, fns...)
	return ms
}

// Scan applies the selector query and scans the result into the given value.
func (ms *MindmapSelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, ms.ctx, ent.OpQuerySelect)
	if err := ms.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*MindmapQuery, *MindmapSelect](ctx, ms.MindmapQuery, ms, ms.inters, v)
}

func (ms *MindmapSelect) sqlScan(ctx context.Context, root *MindmapQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(ms.fns))
	for _, fn := range ms.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*ms.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := ms.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}
