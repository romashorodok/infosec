// Code generated by ent, DO NOT EDIT.

package participant

import (
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

const (
	// Label holds the string label denoting the participant type in the database.
	Label = "participant"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// EdgeBoards holds the string denoting the boards edge name in mutations.
	EdgeBoards = "boards"
	// EdgeTasks holds the string denoting the tasks edge name in mutations.
	EdgeTasks = "tasks"
	// EdgeUser holds the string denoting the user edge name in mutations.
	EdgeUser = "user"
	// Table holds the table name of the participant in the database.
	Table = "participants"
	// BoardsTable is the table that holds the boards relation/edge. The primary key declared below.
	BoardsTable = "board_participants"
	// BoardsInverseTable is the table name for the Board entity.
	// It exists in this package in order to avoid circular dependency with the "board" package.
	BoardsInverseTable = "boards"
	// TasksTable is the table that holds the tasks relation/edge. The primary key declared below.
	TasksTable = "task_participants"
	// TasksInverseTable is the table name for the Task entity.
	// It exists in this package in order to avoid circular dependency with the "task" package.
	TasksInverseTable = "tasks"
	// UserTable is the table that holds the user relation/edge.
	UserTable = "participants"
	// UserInverseTable is the table name for the User entity.
	// It exists in this package in order to avoid circular dependency with the "user" package.
	UserInverseTable = "users"
	// UserColumn is the table column denoting the user relation/edge.
	UserColumn = "user_participants"
)

// Columns holds all SQL columns for participant fields.
var Columns = []string{
	FieldID,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the "participants"
// table and are not defined as standalone fields in the schema.
var ForeignKeys = []string{
	"user_participants",
}

var (
	// BoardsPrimaryKey and BoardsColumn2 are the table columns denoting the
	// primary key for the boards relation (M2M).
	BoardsPrimaryKey = []string{"board_id", "participant_id"}
	// TasksPrimaryKey and TasksColumn2 are the table columns denoting the
	// primary key for the tasks relation (M2M).
	TasksPrimaryKey = []string{"task_id", "participant_id"}
)

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	for i := range ForeignKeys {
		if column == ForeignKeys[i] {
			return true
		}
	}
	return false
}

// OrderOption defines the ordering options for the Participant queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByBoardsCount orders the results by boards count.
func ByBoardsCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newBoardsStep(), opts...)
	}
}

// ByBoards orders the results by boards terms.
func ByBoards(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newBoardsStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}

// ByTasksCount orders the results by tasks count.
func ByTasksCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newTasksStep(), opts...)
	}
}

// ByTasks orders the results by tasks terms.
func ByTasks(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newTasksStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}

// ByUserField orders the results by user field.
func ByUserField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newUserStep(), sql.OrderByField(field, opts...))
	}
}
func newBoardsStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(BoardsInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2M, true, BoardsTable, BoardsPrimaryKey...),
	)
}
func newTasksStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(TasksInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2M, true, TasksTable, TasksPrimaryKey...),
	)
}
func newUserStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(UserInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, UserTable, UserColumn),
	)
}
