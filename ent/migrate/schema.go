// Code generated by ent, DO NOT EDIT.

package migrate

import (
	"entgo.io/ent/dialect/sql/schema"
	"entgo.io/ent/schema/field"
)

var (
	// BoardsColumns holds the columns for the "boards" table.
	BoardsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "title", Type: field.TypeString},
	}
	// BoardsTable holds the schema information for the "boards" table.
	BoardsTable = &schema.Table{
		Name:       "boards",
		Columns:    BoardsColumns,
		PrimaryKey: []*schema.Column{BoardsColumns[0]},
	}
	// GrantsColumns holds the columns for the "grants" table.
	GrantsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
	}
	// GrantsTable holds the schema information for the "grants" table.
	GrantsTable = &schema.Table{
		Name:       "grants",
		Columns:    GrantsColumns,
		PrimaryKey: []*schema.Column{GrantsColumns[0]},
	}
	// ParticipantsColumns holds the columns for the "participants" table.
	ParticipantsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "user_participants", Type: field.TypeUUID, Nullable: true},
	}
	// ParticipantsTable holds the schema information for the "participants" table.
	ParticipantsTable = &schema.Table{
		Name:       "participants",
		Columns:    ParticipantsColumns,
		PrimaryKey: []*schema.Column{ParticipantsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "participants_users_participants",
				Columns:    []*schema.Column{ParticipantsColumns[1]},
				RefColumns: []*schema.Column{UsersColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
	}
	// PillarsColumns holds the columns for the "pillars" table.
	PillarsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "board_pillars", Type: field.TypeInt, Nullable: true},
	}
	// PillarsTable holds the schema information for the "pillars" table.
	PillarsTable = &schema.Table{
		Name:       "pillars",
		Columns:    PillarsColumns,
		PrimaryKey: []*schema.Column{PillarsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "pillars_boards_pillars",
				Columns:    []*schema.Column{PillarsColumns[1]},
				RefColumns: []*schema.Column{BoardsColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
	}
	// TasksColumns holds the columns for the "tasks" table.
	TasksColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "title", Type: field.TypeString},
		{Name: "description", Type: field.TypeString, Size: 2147483647},
		{Name: "board_tasks", Type: field.TypeInt, Nullable: true},
		{Name: "pillar_tasks", Type: field.TypeInt, Nullable: true},
	}
	// TasksTable holds the schema information for the "tasks" table.
	TasksTable = &schema.Table{
		Name:       "tasks",
		Columns:    TasksColumns,
		PrimaryKey: []*schema.Column{TasksColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "tasks_boards_tasks",
				Columns:    []*schema.Column{TasksColumns[3]},
				RefColumns: []*schema.Column{BoardsColumns[0]},
				OnDelete:   schema.SetNull,
			},
			{
				Symbol:     "tasks_pillars_tasks",
				Columns:    []*schema.Column{TasksColumns[4]},
				RefColumns: []*schema.Column{PillarsColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
	}
	// UsersColumns holds the columns for the "users" table.
	UsersColumns = []*schema.Column{
		{Name: "id", Type: field.TypeUUID},
		{Name: "username", Type: field.TypeString, Unique: true},
		{Name: "password", Type: field.TypeString},
	}
	// UsersTable holds the schema information for the "users" table.
	UsersTable = &schema.Table{
		Name:       "users",
		Columns:    UsersColumns,
		PrimaryKey: []*schema.Column{UsersColumns[0]},
	}
	// BoardParticipantsColumns holds the columns for the "board_participants" table.
	BoardParticipantsColumns = []*schema.Column{
		{Name: "board_id", Type: field.TypeInt},
		{Name: "participant_id", Type: field.TypeInt},
	}
	// BoardParticipantsTable holds the schema information for the "board_participants" table.
	BoardParticipantsTable = &schema.Table{
		Name:       "board_participants",
		Columns:    BoardParticipantsColumns,
		PrimaryKey: []*schema.Column{BoardParticipantsColumns[0], BoardParticipantsColumns[1]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "board_participants_board_id",
				Columns:    []*schema.Column{BoardParticipantsColumns[0]},
				RefColumns: []*schema.Column{BoardsColumns[0]},
				OnDelete:   schema.Cascade,
			},
			{
				Symbol:     "board_participants_participant_id",
				Columns:    []*schema.Column{BoardParticipantsColumns[1]},
				RefColumns: []*schema.Column{ParticipantsColumns[0]},
				OnDelete:   schema.Cascade,
			},
		},
	}
	// TaskParticipantsColumns holds the columns for the "task_participants" table.
	TaskParticipantsColumns = []*schema.Column{
		{Name: "task_id", Type: field.TypeInt},
		{Name: "participant_id", Type: field.TypeInt},
	}
	// TaskParticipantsTable holds the schema information for the "task_participants" table.
	TaskParticipantsTable = &schema.Table{
		Name:       "task_participants",
		Columns:    TaskParticipantsColumns,
		PrimaryKey: []*schema.Column{TaskParticipantsColumns[0], TaskParticipantsColumns[1]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "task_participants_task_id",
				Columns:    []*schema.Column{TaskParticipantsColumns[0]},
				RefColumns: []*schema.Column{TasksColumns[0]},
				OnDelete:   schema.Cascade,
			},
			{
				Symbol:     "task_participants_participant_id",
				Columns:    []*schema.Column{TaskParticipantsColumns[1]},
				RefColumns: []*schema.Column{ParticipantsColumns[0]},
				OnDelete:   schema.Cascade,
			},
		},
	}
	// Tables holds all the tables in the schema.
	Tables = []*schema.Table{
		BoardsTable,
		GrantsTable,
		ParticipantsTable,
		PillarsTable,
		TasksTable,
		UsersTable,
		BoardParticipantsTable,
		TaskParticipantsTable,
	}
)

func init() {
	ParticipantsTable.ForeignKeys[0].RefTable = UsersTable
	PillarsTable.ForeignKeys[0].RefTable = BoardsTable
	TasksTable.ForeignKeys[0].RefTable = BoardsTable
	TasksTable.ForeignKeys[1].RefTable = PillarsTable
	BoardParticipantsTable.ForeignKeys[0].RefTable = BoardsTable
	BoardParticipantsTable.ForeignKeys[1].RefTable = ParticipantsTable
	TaskParticipantsTable.ForeignKeys[0].RefTable = TasksTable
	TaskParticipantsTable.ForeignKeys[1].RefTable = ParticipantsTable
}
