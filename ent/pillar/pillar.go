// Code generated by ent, DO NOT EDIT.

package pillar

import (
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

const (
	// Label holds the string label denoting the pillar type in the database.
	Label = "pillar"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// EdgeTasks holds the string denoting the tasks edge name in mutations.
	EdgeTasks = "tasks"
	// Table holds the table name of the pillar in the database.
	Table = "pillars"
	// TasksTable is the table that holds the tasks relation/edge.
	TasksTable = "tasks"
	// TasksInverseTable is the table name for the Task entity.
	// It exists in this package in order to avoid circular dependency with the "task" package.
	TasksInverseTable = "tasks"
	// TasksColumn is the table column denoting the tasks relation/edge.
	TasksColumn = "pillar_tasks"
)

// Columns holds all SQL columns for pillar fields.
var Columns = []string{
	FieldID,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the "pillars"
// table and are not defined as standalone fields in the schema.
var ForeignKeys = []string{
	"board_pillars",
}

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

// OrderOption defines the ordering options for the Pillar queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
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
func newTasksStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(TasksInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, TasksTable, TasksColumn),
	)
}
