// Code generated by ent, DO NOT EDIT.

package migrate

import (
	"entgo.io/ent/dialect/sql/schema"
	"entgo.io/ent/schema/field"
)

var (
	// PillarsColumns holds the columns for the "pillars" table.
	PillarsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
	}
	// PillarsTable holds the schema information for the "pillars" table.
	PillarsTable = &schema.Table{
		Name:       "pillars",
		Columns:    PillarsColumns,
		PrimaryKey: []*schema.Column{PillarsColumns[0]},
	}
	// Tables holds all the tables in the schema.
	Tables = []*schema.Table{
		PillarsTable,
	}
)

func init() {
}
