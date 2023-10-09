// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"log"
	"reflect"

	"github.com/google/uuid"
	"github.com/romashorodok/infosec/ent/migrate"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/romashorodok/infosec/ent/board"
	"github.com/romashorodok/infosec/ent/grant"
	"github.com/romashorodok/infosec/ent/participant"
	"github.com/romashorodok/infosec/ent/pillar"
	"github.com/romashorodok/infosec/ent/task"
	"github.com/romashorodok/infosec/ent/user"
)

// Client is the client that holds all ent builders.
type Client struct {
	config
	// Schema is the client for creating, migrating and dropping schema.
	Schema *migrate.Schema
	// Board is the client for interacting with the Board builders.
	Board *BoardClient
	// Grant is the client for interacting with the Grant builders.
	Grant *GrantClient
	// Participant is the client for interacting with the Participant builders.
	Participant *ParticipantClient
	// Pillar is the client for interacting with the Pillar builders.
	Pillar *PillarClient
	// Task is the client for interacting with the Task builders.
	Task *TaskClient
	// User is the client for interacting with the User builders.
	User *UserClient
}

// NewClient creates a new client configured with the given options.
func NewClient(opts ...Option) *Client {
	cfg := config{log: log.Println, hooks: &hooks{}, inters: &inters{}}
	cfg.options(opts...)
	client := &Client{config: cfg}
	client.init()
	return client
}

func (c *Client) init() {
	c.Schema = migrate.NewSchema(c.driver)
	c.Board = NewBoardClient(c.config)
	c.Grant = NewGrantClient(c.config)
	c.Participant = NewParticipantClient(c.config)
	c.Pillar = NewPillarClient(c.config)
	c.Task = NewTaskClient(c.config)
	c.User = NewUserClient(c.config)
}

type (
	// config is the configuration for the client and its builder.
	config struct {
		// driver used for executing database requests.
		driver dialect.Driver
		// debug enable a debug logging.
		debug bool
		// log used for logging on debug mode.
		log func(...any)
		// hooks to execute on mutations.
		hooks *hooks
		// interceptors to execute on queries.
		inters *inters
	}
	// Option function to configure the client.
	Option func(*config)
)

// options applies the options on the config object.
func (c *config) options(opts ...Option) {
	for _, opt := range opts {
		opt(c)
	}
	if c.debug {
		c.driver = dialect.Debug(c.driver, c.log)
	}
}

// Debug enables debug logging on the ent.Driver.
func Debug() Option {
	return func(c *config) {
		c.debug = true
	}
}

// Log sets the logging function for debug mode.
func Log(fn func(...any)) Option {
	return func(c *config) {
		c.log = fn
	}
}

// Driver configures the client driver.
func Driver(driver dialect.Driver) Option {
	return func(c *config) {
		c.driver = driver
	}
}

// Open opens a database/sql.DB specified by the driver name and
// the data source name, and returns a new client attached to it.
// Optional parameters can be added for configuring the client.
func Open(driverName, dataSourceName string, options ...Option) (*Client, error) {
	switch driverName {
	case dialect.MySQL, dialect.Postgres, dialect.SQLite:
		drv, err := sql.Open(driverName, dataSourceName)
		if err != nil {
			return nil, err
		}
		return NewClient(append(options, Driver(drv))...), nil
	default:
		return nil, fmt.Errorf("unsupported driver: %q", driverName)
	}
}

// ErrTxStarted is returned when trying to start a new transaction from a transactional client.
var ErrTxStarted = errors.New("ent: cannot start a transaction within a transaction")

// Tx returns a new transactional client. The provided context
// is used until the transaction is committed or rolled back.
func (c *Client) Tx(ctx context.Context) (*Tx, error) {
	if _, ok := c.driver.(*txDriver); ok {
		return nil, ErrTxStarted
	}
	tx, err := newTx(ctx, c.driver)
	if err != nil {
		return nil, fmt.Errorf("ent: starting a transaction: %w", err)
	}
	cfg := c.config
	cfg.driver = tx
	return &Tx{
		ctx:         ctx,
		config:      cfg,
		Board:       NewBoardClient(cfg),
		Grant:       NewGrantClient(cfg),
		Participant: NewParticipantClient(cfg),
		Pillar:      NewPillarClient(cfg),
		Task:        NewTaskClient(cfg),
		User:        NewUserClient(cfg),
	}, nil
}

// BeginTx returns a transactional client with specified options.
func (c *Client) BeginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	if _, ok := c.driver.(*txDriver); ok {
		return nil, errors.New("ent: cannot start a transaction within a transaction")
	}
	tx, err := c.driver.(interface {
		BeginTx(context.Context, *sql.TxOptions) (dialect.Tx, error)
	}).BeginTx(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("ent: starting a transaction: %w", err)
	}
	cfg := c.config
	cfg.driver = &txDriver{tx: tx, drv: c.driver}
	return &Tx{
		ctx:         ctx,
		config:      cfg,
		Board:       NewBoardClient(cfg),
		Grant:       NewGrantClient(cfg),
		Participant: NewParticipantClient(cfg),
		Pillar:      NewPillarClient(cfg),
		Task:        NewTaskClient(cfg),
		User:        NewUserClient(cfg),
	}, nil
}

// Debug returns a new debug-client. It's used to get verbose logging on specific operations.
//
//	client.Debug().
//		Board.
//		Query().
//		Count(ctx)
func (c *Client) Debug() *Client {
	if c.debug {
		return c
	}
	cfg := c.config
	cfg.driver = dialect.Debug(c.driver, c.log)
	client := &Client{config: cfg}
	client.init()
	return client
}

// Close closes the database connection and prevents new queries from starting.
func (c *Client) Close() error {
	return c.driver.Close()
}

// Use adds the mutation hooks to all the entity clients.
// In order to add hooks to a specific client, call: `client.Node.Use(...)`.
func (c *Client) Use(hooks ...Hook) {
	for _, n := range []interface{ Use(...Hook) }{
		c.Board, c.Grant, c.Participant, c.Pillar, c.Task, c.User,
	} {
		n.Use(hooks...)
	}
}

// Intercept adds the query interceptors to all the entity clients.
// In order to add interceptors to a specific client, call: `client.Node.Intercept(...)`.
func (c *Client) Intercept(interceptors ...Interceptor) {
	for _, n := range []interface{ Intercept(...Interceptor) }{
		c.Board, c.Grant, c.Participant, c.Pillar, c.Task, c.User,
	} {
		n.Intercept(interceptors...)
	}
}

// Mutate implements the ent.Mutator interface.
func (c *Client) Mutate(ctx context.Context, m Mutation) (Value, error) {
	switch m := m.(type) {
	case *BoardMutation:
		return c.Board.mutate(ctx, m)
	case *GrantMutation:
		return c.Grant.mutate(ctx, m)
	case *ParticipantMutation:
		return c.Participant.mutate(ctx, m)
	case *PillarMutation:
		return c.Pillar.mutate(ctx, m)
	case *TaskMutation:
		return c.Task.mutate(ctx, m)
	case *UserMutation:
		return c.User.mutate(ctx, m)
	default:
		return nil, fmt.Errorf("ent: unknown mutation type %T", m)
	}
}

// BoardClient is a client for the Board schema.
type BoardClient struct {
	config
}

// NewBoardClient returns a client for the Board from the given config.
func NewBoardClient(c config) *BoardClient {
	return &BoardClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `board.Hooks(f(g(h())))`.
func (c *BoardClient) Use(hooks ...Hook) {
	c.hooks.Board = append(c.hooks.Board, hooks...)
}

// Intercept adds a list of query interceptors to the interceptors stack.
// A call to `Intercept(f, g, h)` equals to `board.Intercept(f(g(h())))`.
func (c *BoardClient) Intercept(interceptors ...Interceptor) {
	c.inters.Board = append(c.inters.Board, interceptors...)
}

// Create returns a builder for creating a Board entity.
func (c *BoardClient) Create() *BoardCreate {
	mutation := newBoardMutation(c.config, OpCreate)
	return &BoardCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of Board entities.
func (c *BoardClient) CreateBulk(builders ...*BoardCreate) *BoardCreateBulk {
	return &BoardCreateBulk{config: c.config, builders: builders}
}

// MapCreateBulk creates a bulk creation builder from the given slice. For each item in the slice, the function creates
// a builder and applies setFunc on it.
func (c *BoardClient) MapCreateBulk(slice any, setFunc func(*BoardCreate, int)) *BoardCreateBulk {
	rv := reflect.ValueOf(slice)
	if rv.Kind() != reflect.Slice {
		return &BoardCreateBulk{err: fmt.Errorf("calling to BoardClient.MapCreateBulk with wrong type %T, need slice", slice)}
	}
	builders := make([]*BoardCreate, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		builders[i] = c.Create()
		setFunc(builders[i], i)
	}
	return &BoardCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for Board.
func (c *BoardClient) Update() *BoardUpdate {
	mutation := newBoardMutation(c.config, OpUpdate)
	return &BoardUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *BoardClient) UpdateOne(b *Board) *BoardUpdateOne {
	mutation := newBoardMutation(c.config, OpUpdateOne, withBoard(b))
	return &BoardUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *BoardClient) UpdateOneID(id int) *BoardUpdateOne {
	mutation := newBoardMutation(c.config, OpUpdateOne, withBoardID(id))
	return &BoardUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for Board.
func (c *BoardClient) Delete() *BoardDelete {
	mutation := newBoardMutation(c.config, OpDelete)
	return &BoardDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a builder for deleting the given entity.
func (c *BoardClient) DeleteOne(b *Board) *BoardDeleteOne {
	return c.DeleteOneID(b.ID)
}

// DeleteOneID returns a builder for deleting the given entity by its id.
func (c *BoardClient) DeleteOneID(id int) *BoardDeleteOne {
	builder := c.Delete().Where(board.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &BoardDeleteOne{builder}
}

// Query returns a query builder for Board.
func (c *BoardClient) Query() *BoardQuery {
	return &BoardQuery{
		config: c.config,
		ctx:    &QueryContext{Type: TypeBoard},
		inters: c.Interceptors(),
	}
}

// Get returns a Board entity by its id.
func (c *BoardClient) Get(ctx context.Context, id int) (*Board, error) {
	return c.Query().Where(board.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *BoardClient) GetX(ctx context.Context, id int) *Board {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// QueryTasks queries the tasks edge of a Board.
func (c *BoardClient) QueryTasks(b *Board) *TaskQuery {
	query := (&TaskClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := b.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(board.Table, board.FieldID, id),
			sqlgraph.To(task.Table, task.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, board.TasksTable, board.TasksColumn),
		)
		fromV = sqlgraph.Neighbors(b.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryParticipants queries the participants edge of a Board.
func (c *BoardClient) QueryParticipants(b *Board) *ParticipantQuery {
	query := (&ParticipantClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := b.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(board.Table, board.FieldID, id),
			sqlgraph.To(participant.Table, participant.FieldID),
			sqlgraph.Edge(sqlgraph.M2M, false, board.ParticipantsTable, board.ParticipantsPrimaryKey...),
		)
		fromV = sqlgraph.Neighbors(b.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *BoardClient) Hooks() []Hook {
	return c.hooks.Board
}

// Interceptors returns the client interceptors.
func (c *BoardClient) Interceptors() []Interceptor {
	return c.inters.Board
}

func (c *BoardClient) mutate(ctx context.Context, m *BoardMutation) (Value, error) {
	switch m.Op() {
	case OpCreate:
		return (&BoardCreate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdate:
		return (&BoardUpdate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdateOne:
		return (&BoardUpdateOne{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpDelete, OpDeleteOne:
		return (&BoardDelete{config: c.config, hooks: c.Hooks(), mutation: m}).Exec(ctx)
	default:
		return nil, fmt.Errorf("ent: unknown Board mutation op: %q", m.Op())
	}
}

// GrantClient is a client for the Grant schema.
type GrantClient struct {
	config
}

// NewGrantClient returns a client for the Grant from the given config.
func NewGrantClient(c config) *GrantClient {
	return &GrantClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `grant.Hooks(f(g(h())))`.
func (c *GrantClient) Use(hooks ...Hook) {
	c.hooks.Grant = append(c.hooks.Grant, hooks...)
}

// Intercept adds a list of query interceptors to the interceptors stack.
// A call to `Intercept(f, g, h)` equals to `grant.Intercept(f(g(h())))`.
func (c *GrantClient) Intercept(interceptors ...Interceptor) {
	c.inters.Grant = append(c.inters.Grant, interceptors...)
}

// Create returns a builder for creating a Grant entity.
func (c *GrantClient) Create() *GrantCreate {
	mutation := newGrantMutation(c.config, OpCreate)
	return &GrantCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of Grant entities.
func (c *GrantClient) CreateBulk(builders ...*GrantCreate) *GrantCreateBulk {
	return &GrantCreateBulk{config: c.config, builders: builders}
}

// MapCreateBulk creates a bulk creation builder from the given slice. For each item in the slice, the function creates
// a builder and applies setFunc on it.
func (c *GrantClient) MapCreateBulk(slice any, setFunc func(*GrantCreate, int)) *GrantCreateBulk {
	rv := reflect.ValueOf(slice)
	if rv.Kind() != reflect.Slice {
		return &GrantCreateBulk{err: fmt.Errorf("calling to GrantClient.MapCreateBulk with wrong type %T, need slice", slice)}
	}
	builders := make([]*GrantCreate, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		builders[i] = c.Create()
		setFunc(builders[i], i)
	}
	return &GrantCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for Grant.
func (c *GrantClient) Update() *GrantUpdate {
	mutation := newGrantMutation(c.config, OpUpdate)
	return &GrantUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *GrantClient) UpdateOne(gr *Grant) *GrantUpdateOne {
	mutation := newGrantMutation(c.config, OpUpdateOne, withGrant(gr))
	return &GrantUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *GrantClient) UpdateOneID(id int) *GrantUpdateOne {
	mutation := newGrantMutation(c.config, OpUpdateOne, withGrantID(id))
	return &GrantUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for Grant.
func (c *GrantClient) Delete() *GrantDelete {
	mutation := newGrantMutation(c.config, OpDelete)
	return &GrantDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a builder for deleting the given entity.
func (c *GrantClient) DeleteOne(gr *Grant) *GrantDeleteOne {
	return c.DeleteOneID(gr.ID)
}

// DeleteOneID returns a builder for deleting the given entity by its id.
func (c *GrantClient) DeleteOneID(id int) *GrantDeleteOne {
	builder := c.Delete().Where(grant.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &GrantDeleteOne{builder}
}

// Query returns a query builder for Grant.
func (c *GrantClient) Query() *GrantQuery {
	return &GrantQuery{
		config: c.config,
		ctx:    &QueryContext{Type: TypeGrant},
		inters: c.Interceptors(),
	}
}

// Get returns a Grant entity by its id.
func (c *GrantClient) Get(ctx context.Context, id int) (*Grant, error) {
	return c.Query().Where(grant.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *GrantClient) GetX(ctx context.Context, id int) *Grant {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// Hooks returns the client hooks.
func (c *GrantClient) Hooks() []Hook {
	return c.hooks.Grant
}

// Interceptors returns the client interceptors.
func (c *GrantClient) Interceptors() []Interceptor {
	return c.inters.Grant
}

func (c *GrantClient) mutate(ctx context.Context, m *GrantMutation) (Value, error) {
	switch m.Op() {
	case OpCreate:
		return (&GrantCreate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdate:
		return (&GrantUpdate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdateOne:
		return (&GrantUpdateOne{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpDelete, OpDeleteOne:
		return (&GrantDelete{config: c.config, hooks: c.Hooks(), mutation: m}).Exec(ctx)
	default:
		return nil, fmt.Errorf("ent: unknown Grant mutation op: %q", m.Op())
	}
}

// ParticipantClient is a client for the Participant schema.
type ParticipantClient struct {
	config
}

// NewParticipantClient returns a client for the Participant from the given config.
func NewParticipantClient(c config) *ParticipantClient {
	return &ParticipantClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `participant.Hooks(f(g(h())))`.
func (c *ParticipantClient) Use(hooks ...Hook) {
	c.hooks.Participant = append(c.hooks.Participant, hooks...)
}

// Intercept adds a list of query interceptors to the interceptors stack.
// A call to `Intercept(f, g, h)` equals to `participant.Intercept(f(g(h())))`.
func (c *ParticipantClient) Intercept(interceptors ...Interceptor) {
	c.inters.Participant = append(c.inters.Participant, interceptors...)
}

// Create returns a builder for creating a Participant entity.
func (c *ParticipantClient) Create() *ParticipantCreate {
	mutation := newParticipantMutation(c.config, OpCreate)
	return &ParticipantCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of Participant entities.
func (c *ParticipantClient) CreateBulk(builders ...*ParticipantCreate) *ParticipantCreateBulk {
	return &ParticipantCreateBulk{config: c.config, builders: builders}
}

// MapCreateBulk creates a bulk creation builder from the given slice. For each item in the slice, the function creates
// a builder and applies setFunc on it.
func (c *ParticipantClient) MapCreateBulk(slice any, setFunc func(*ParticipantCreate, int)) *ParticipantCreateBulk {
	rv := reflect.ValueOf(slice)
	if rv.Kind() != reflect.Slice {
		return &ParticipantCreateBulk{err: fmt.Errorf("calling to ParticipantClient.MapCreateBulk with wrong type %T, need slice", slice)}
	}
	builders := make([]*ParticipantCreate, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		builders[i] = c.Create()
		setFunc(builders[i], i)
	}
	return &ParticipantCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for Participant.
func (c *ParticipantClient) Update() *ParticipantUpdate {
	mutation := newParticipantMutation(c.config, OpUpdate)
	return &ParticipantUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *ParticipantClient) UpdateOne(pa *Participant) *ParticipantUpdateOne {
	mutation := newParticipantMutation(c.config, OpUpdateOne, withParticipant(pa))
	return &ParticipantUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *ParticipantClient) UpdateOneID(id int) *ParticipantUpdateOne {
	mutation := newParticipantMutation(c.config, OpUpdateOne, withParticipantID(id))
	return &ParticipantUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for Participant.
func (c *ParticipantClient) Delete() *ParticipantDelete {
	mutation := newParticipantMutation(c.config, OpDelete)
	return &ParticipantDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a builder for deleting the given entity.
func (c *ParticipantClient) DeleteOne(pa *Participant) *ParticipantDeleteOne {
	return c.DeleteOneID(pa.ID)
}

// DeleteOneID returns a builder for deleting the given entity by its id.
func (c *ParticipantClient) DeleteOneID(id int) *ParticipantDeleteOne {
	builder := c.Delete().Where(participant.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &ParticipantDeleteOne{builder}
}

// Query returns a query builder for Participant.
func (c *ParticipantClient) Query() *ParticipantQuery {
	return &ParticipantQuery{
		config: c.config,
		ctx:    &QueryContext{Type: TypeParticipant},
		inters: c.Interceptors(),
	}
}

// Get returns a Participant entity by its id.
func (c *ParticipantClient) Get(ctx context.Context, id int) (*Participant, error) {
	return c.Query().Where(participant.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *ParticipantClient) GetX(ctx context.Context, id int) *Participant {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// QueryBoards queries the boards edge of a Participant.
func (c *ParticipantClient) QueryBoards(pa *Participant) *BoardQuery {
	query := (&BoardClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := pa.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(participant.Table, participant.FieldID, id),
			sqlgraph.To(board.Table, board.FieldID),
			sqlgraph.Edge(sqlgraph.M2M, true, participant.BoardsTable, participant.BoardsPrimaryKey...),
		)
		fromV = sqlgraph.Neighbors(pa.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryTasks queries the tasks edge of a Participant.
func (c *ParticipantClient) QueryTasks(pa *Participant) *TaskQuery {
	query := (&TaskClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := pa.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(participant.Table, participant.FieldID, id),
			sqlgraph.To(task.Table, task.FieldID),
			sqlgraph.Edge(sqlgraph.M2M, true, participant.TasksTable, participant.TasksPrimaryKey...),
		)
		fromV = sqlgraph.Neighbors(pa.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *ParticipantClient) Hooks() []Hook {
	return c.hooks.Participant
}

// Interceptors returns the client interceptors.
func (c *ParticipantClient) Interceptors() []Interceptor {
	return c.inters.Participant
}

func (c *ParticipantClient) mutate(ctx context.Context, m *ParticipantMutation) (Value, error) {
	switch m.Op() {
	case OpCreate:
		return (&ParticipantCreate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdate:
		return (&ParticipantUpdate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdateOne:
		return (&ParticipantUpdateOne{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpDelete, OpDeleteOne:
		return (&ParticipantDelete{config: c.config, hooks: c.Hooks(), mutation: m}).Exec(ctx)
	default:
		return nil, fmt.Errorf("ent: unknown Participant mutation op: %q", m.Op())
	}
}

// PillarClient is a client for the Pillar schema.
type PillarClient struct {
	config
}

// NewPillarClient returns a client for the Pillar from the given config.
func NewPillarClient(c config) *PillarClient {
	return &PillarClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `pillar.Hooks(f(g(h())))`.
func (c *PillarClient) Use(hooks ...Hook) {
	c.hooks.Pillar = append(c.hooks.Pillar, hooks...)
}

// Intercept adds a list of query interceptors to the interceptors stack.
// A call to `Intercept(f, g, h)` equals to `pillar.Intercept(f(g(h())))`.
func (c *PillarClient) Intercept(interceptors ...Interceptor) {
	c.inters.Pillar = append(c.inters.Pillar, interceptors...)
}

// Create returns a builder for creating a Pillar entity.
func (c *PillarClient) Create() *PillarCreate {
	mutation := newPillarMutation(c.config, OpCreate)
	return &PillarCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of Pillar entities.
func (c *PillarClient) CreateBulk(builders ...*PillarCreate) *PillarCreateBulk {
	return &PillarCreateBulk{config: c.config, builders: builders}
}

// MapCreateBulk creates a bulk creation builder from the given slice. For each item in the slice, the function creates
// a builder and applies setFunc on it.
func (c *PillarClient) MapCreateBulk(slice any, setFunc func(*PillarCreate, int)) *PillarCreateBulk {
	rv := reflect.ValueOf(slice)
	if rv.Kind() != reflect.Slice {
		return &PillarCreateBulk{err: fmt.Errorf("calling to PillarClient.MapCreateBulk with wrong type %T, need slice", slice)}
	}
	builders := make([]*PillarCreate, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		builders[i] = c.Create()
		setFunc(builders[i], i)
	}
	return &PillarCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for Pillar.
func (c *PillarClient) Update() *PillarUpdate {
	mutation := newPillarMutation(c.config, OpUpdate)
	return &PillarUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *PillarClient) UpdateOne(pi *Pillar) *PillarUpdateOne {
	mutation := newPillarMutation(c.config, OpUpdateOne, withPillar(pi))
	return &PillarUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *PillarClient) UpdateOneID(id int) *PillarUpdateOne {
	mutation := newPillarMutation(c.config, OpUpdateOne, withPillarID(id))
	return &PillarUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for Pillar.
func (c *PillarClient) Delete() *PillarDelete {
	mutation := newPillarMutation(c.config, OpDelete)
	return &PillarDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a builder for deleting the given entity.
func (c *PillarClient) DeleteOne(pi *Pillar) *PillarDeleteOne {
	return c.DeleteOneID(pi.ID)
}

// DeleteOneID returns a builder for deleting the given entity by its id.
func (c *PillarClient) DeleteOneID(id int) *PillarDeleteOne {
	builder := c.Delete().Where(pillar.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &PillarDeleteOne{builder}
}

// Query returns a query builder for Pillar.
func (c *PillarClient) Query() *PillarQuery {
	return &PillarQuery{
		config: c.config,
		ctx:    &QueryContext{Type: TypePillar},
		inters: c.Interceptors(),
	}
}

// Get returns a Pillar entity by its id.
func (c *PillarClient) Get(ctx context.Context, id int) (*Pillar, error) {
	return c.Query().Where(pillar.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *PillarClient) GetX(ctx context.Context, id int) *Pillar {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// QueryTasks queries the tasks edge of a Pillar.
func (c *PillarClient) QueryTasks(pi *Pillar) *TaskQuery {
	query := (&TaskClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := pi.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(pillar.Table, pillar.FieldID, id),
			sqlgraph.To(task.Table, task.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, pillar.TasksTable, pillar.TasksColumn),
		)
		fromV = sqlgraph.Neighbors(pi.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *PillarClient) Hooks() []Hook {
	return c.hooks.Pillar
}

// Interceptors returns the client interceptors.
func (c *PillarClient) Interceptors() []Interceptor {
	return c.inters.Pillar
}

func (c *PillarClient) mutate(ctx context.Context, m *PillarMutation) (Value, error) {
	switch m.Op() {
	case OpCreate:
		return (&PillarCreate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdate:
		return (&PillarUpdate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdateOne:
		return (&PillarUpdateOne{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpDelete, OpDeleteOne:
		return (&PillarDelete{config: c.config, hooks: c.Hooks(), mutation: m}).Exec(ctx)
	default:
		return nil, fmt.Errorf("ent: unknown Pillar mutation op: %q", m.Op())
	}
}

// TaskClient is a client for the Task schema.
type TaskClient struct {
	config
}

// NewTaskClient returns a client for the Task from the given config.
func NewTaskClient(c config) *TaskClient {
	return &TaskClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `task.Hooks(f(g(h())))`.
func (c *TaskClient) Use(hooks ...Hook) {
	c.hooks.Task = append(c.hooks.Task, hooks...)
}

// Intercept adds a list of query interceptors to the interceptors stack.
// A call to `Intercept(f, g, h)` equals to `task.Intercept(f(g(h())))`.
func (c *TaskClient) Intercept(interceptors ...Interceptor) {
	c.inters.Task = append(c.inters.Task, interceptors...)
}

// Create returns a builder for creating a Task entity.
func (c *TaskClient) Create() *TaskCreate {
	mutation := newTaskMutation(c.config, OpCreate)
	return &TaskCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of Task entities.
func (c *TaskClient) CreateBulk(builders ...*TaskCreate) *TaskCreateBulk {
	return &TaskCreateBulk{config: c.config, builders: builders}
}

// MapCreateBulk creates a bulk creation builder from the given slice. For each item in the slice, the function creates
// a builder and applies setFunc on it.
func (c *TaskClient) MapCreateBulk(slice any, setFunc func(*TaskCreate, int)) *TaskCreateBulk {
	rv := reflect.ValueOf(slice)
	if rv.Kind() != reflect.Slice {
		return &TaskCreateBulk{err: fmt.Errorf("calling to TaskClient.MapCreateBulk with wrong type %T, need slice", slice)}
	}
	builders := make([]*TaskCreate, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		builders[i] = c.Create()
		setFunc(builders[i], i)
	}
	return &TaskCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for Task.
func (c *TaskClient) Update() *TaskUpdate {
	mutation := newTaskMutation(c.config, OpUpdate)
	return &TaskUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *TaskClient) UpdateOne(t *Task) *TaskUpdateOne {
	mutation := newTaskMutation(c.config, OpUpdateOne, withTask(t))
	return &TaskUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *TaskClient) UpdateOneID(id int) *TaskUpdateOne {
	mutation := newTaskMutation(c.config, OpUpdateOne, withTaskID(id))
	return &TaskUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for Task.
func (c *TaskClient) Delete() *TaskDelete {
	mutation := newTaskMutation(c.config, OpDelete)
	return &TaskDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a builder for deleting the given entity.
func (c *TaskClient) DeleteOne(t *Task) *TaskDeleteOne {
	return c.DeleteOneID(t.ID)
}

// DeleteOneID returns a builder for deleting the given entity by its id.
func (c *TaskClient) DeleteOneID(id int) *TaskDeleteOne {
	builder := c.Delete().Where(task.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &TaskDeleteOne{builder}
}

// Query returns a query builder for Task.
func (c *TaskClient) Query() *TaskQuery {
	return &TaskQuery{
		config: c.config,
		ctx:    &QueryContext{Type: TypeTask},
		inters: c.Interceptors(),
	}
}

// Get returns a Task entity by its id.
func (c *TaskClient) Get(ctx context.Context, id int) (*Task, error) {
	return c.Query().Where(task.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *TaskClient) GetX(ctx context.Context, id int) *Task {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// QueryParticipants queries the participants edge of a Task.
func (c *TaskClient) QueryParticipants(t *Task) *ParticipantQuery {
	query := (&ParticipantClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := t.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(task.Table, task.FieldID, id),
			sqlgraph.To(participant.Table, participant.FieldID),
			sqlgraph.Edge(sqlgraph.M2M, false, task.ParticipantsTable, task.ParticipantsPrimaryKey...),
		)
		fromV = sqlgraph.Neighbors(t.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *TaskClient) Hooks() []Hook {
	return c.hooks.Task
}

// Interceptors returns the client interceptors.
func (c *TaskClient) Interceptors() []Interceptor {
	return c.inters.Task
}

func (c *TaskClient) mutate(ctx context.Context, m *TaskMutation) (Value, error) {
	switch m.Op() {
	case OpCreate:
		return (&TaskCreate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdate:
		return (&TaskUpdate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdateOne:
		return (&TaskUpdateOne{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpDelete, OpDeleteOne:
		return (&TaskDelete{config: c.config, hooks: c.Hooks(), mutation: m}).Exec(ctx)
	default:
		return nil, fmt.Errorf("ent: unknown Task mutation op: %q", m.Op())
	}
}

// UserClient is a client for the User schema.
type UserClient struct {
	config
}

// NewUserClient returns a client for the User from the given config.
func NewUserClient(c config) *UserClient {
	return &UserClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `user.Hooks(f(g(h())))`.
func (c *UserClient) Use(hooks ...Hook) {
	c.hooks.User = append(c.hooks.User, hooks...)
}

// Intercept adds a list of query interceptors to the interceptors stack.
// A call to `Intercept(f, g, h)` equals to `user.Intercept(f(g(h())))`.
func (c *UserClient) Intercept(interceptors ...Interceptor) {
	c.inters.User = append(c.inters.User, interceptors...)
}

// Create returns a builder for creating a User entity.
func (c *UserClient) Create() *UserCreate {
	mutation := newUserMutation(c.config, OpCreate)
	return &UserCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of User entities.
func (c *UserClient) CreateBulk(builders ...*UserCreate) *UserCreateBulk {
	return &UserCreateBulk{config: c.config, builders: builders}
}

// MapCreateBulk creates a bulk creation builder from the given slice. For each item in the slice, the function creates
// a builder and applies setFunc on it.
func (c *UserClient) MapCreateBulk(slice any, setFunc func(*UserCreate, int)) *UserCreateBulk {
	rv := reflect.ValueOf(slice)
	if rv.Kind() != reflect.Slice {
		return &UserCreateBulk{err: fmt.Errorf("calling to UserClient.MapCreateBulk with wrong type %T, need slice", slice)}
	}
	builders := make([]*UserCreate, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		builders[i] = c.Create()
		setFunc(builders[i], i)
	}
	return &UserCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for User.
func (c *UserClient) Update() *UserUpdate {
	mutation := newUserMutation(c.config, OpUpdate)
	return &UserUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *UserClient) UpdateOne(u *User) *UserUpdateOne {
	mutation := newUserMutation(c.config, OpUpdateOne, withUser(u))
	return &UserUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *UserClient) UpdateOneID(id uuid.UUID) *UserUpdateOne {
	mutation := newUserMutation(c.config, OpUpdateOne, withUserID(id))
	return &UserUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for User.
func (c *UserClient) Delete() *UserDelete {
	mutation := newUserMutation(c.config, OpDelete)
	return &UserDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a builder for deleting the given entity.
func (c *UserClient) DeleteOne(u *User) *UserDeleteOne {
	return c.DeleteOneID(u.ID)
}

// DeleteOneID returns a builder for deleting the given entity by its id.
func (c *UserClient) DeleteOneID(id uuid.UUID) *UserDeleteOne {
	builder := c.Delete().Where(user.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &UserDeleteOne{builder}
}

// Query returns a query builder for User.
func (c *UserClient) Query() *UserQuery {
	return &UserQuery{
		config: c.config,
		ctx:    &QueryContext{Type: TypeUser},
		inters: c.Interceptors(),
	}
}

// Get returns a User entity by its id.
func (c *UserClient) Get(ctx context.Context, id uuid.UUID) (*User, error) {
	return c.Query().Where(user.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *UserClient) GetX(ctx context.Context, id uuid.UUID) *User {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// QueryParticipants queries the Participants edge of a User.
func (c *UserClient) QueryParticipants(u *User) *ParticipantQuery {
	query := (&ParticipantClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := u.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(user.Table, user.FieldID, id),
			sqlgraph.To(participant.Table, participant.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, user.ParticipantsTable, user.ParticipantsColumn),
		)
		fromV = sqlgraph.Neighbors(u.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *UserClient) Hooks() []Hook {
	return c.hooks.User
}

// Interceptors returns the client interceptors.
func (c *UserClient) Interceptors() []Interceptor {
	return c.inters.User
}

func (c *UserClient) mutate(ctx context.Context, m *UserMutation) (Value, error) {
	switch m.Op() {
	case OpCreate:
		return (&UserCreate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdate:
		return (&UserUpdate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdateOne:
		return (&UserUpdateOne{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpDelete, OpDeleteOne:
		return (&UserDelete{config: c.config, hooks: c.Hooks(), mutation: m}).Exec(ctx)
	default:
		return nil, fmt.Errorf("ent: unknown User mutation op: %q", m.Op())
	}
}

// hooks and interceptors per client, for fast access.
type (
	hooks struct {
		Board, Grant, Participant, Pillar, Task, User []ent.Hook
	}
	inters struct {
		Board, Grant, Participant, Pillar, Task, User []ent.Interceptor
	}
)
