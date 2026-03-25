package model

import "time"

type Priority string

const (
	PriorityLow    Priority = "low"
	PriorityMedium Priority = "medium"
	PriorityHigh   Priority = "high"
)

type Todo struct {
	ID        string    `db:"id"`
	ListID    string    `db:"list_id"`
	Text      string    `db:"text"`
	Done      bool      `db:"done"`
	Priority  Priority  `db:"priority"`
	DueDate   *string   `db:"due_date"`
	CreatedAt time.Time `db:"created_at"`
}

type TaskList struct {
	ID        string    `db:"id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
	Todos     []Todo
}
