// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/romashorodok/infosec/ent/ent/pillar"
)

// Pillar is the model entity for the Pillar schema.
type Pillar struct {
	config
	// ID of the ent.
	ID           int `json:"id,omitempty"`
	selectValues sql.SelectValues
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Pillar) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case pillar.FieldID:
			values[i] = new(sql.NullInt64)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Pillar fields.
func (pi *Pillar) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case pillar.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			pi.ID = int(value.Int64)
		default:
			pi.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the Pillar.
// This includes values selected through modifiers, order, etc.
func (pi *Pillar) Value(name string) (ent.Value, error) {
	return pi.selectValues.Get(name)
}

// Update returns a builder for updating this Pillar.
// Note that you need to call Pillar.Unwrap() before calling this method if this Pillar
// was returned from a transaction, and the transaction was committed or rolled back.
func (pi *Pillar) Update() *PillarUpdateOne {
	return NewPillarClient(pi.config).UpdateOne(pi)
}

// Unwrap unwraps the Pillar entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (pi *Pillar) Unwrap() *Pillar {
	_tx, ok := pi.config.driver.(*txDriver)
	if !ok {
		panic("ent: Pillar is not a transactional entity")
	}
	pi.config.driver = _tx.drv
	return pi
}

// String implements the fmt.Stringer.
func (pi *Pillar) String() string {
	var builder strings.Builder
	builder.WriteString("Pillar(")
	builder.WriteString(fmt.Sprintf("id=%v", pi.ID))
	builder.WriteByte(')')
	return builder.String()
}

// Pillars is a parsable slice of Pillar.
type Pillars []*Pillar
