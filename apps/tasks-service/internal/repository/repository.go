package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/its-me-mayday/sisyphus/tasks-service/internal/model"
)

type Repository struct {
	db *sqlx.DB
}

func New(dsn string) (*Repository, error) {
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("connect db: %w", err)
	}
	return &Repository{db: db}, nil
}

func (r *Repository) Migrate(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS task_lists (
			id         TEXT PRIMARY KEY,
			name       TEXT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);

		CREATE TABLE IF NOT EXISTS todos (
			id         TEXT PRIMARY KEY,
			list_id    TEXT NOT NULL REFERENCES task_lists(id) ON DELETE CASCADE,
			text       TEXT NOT NULL,
			done       BOOLEAN NOT NULL DEFAULT FALSE,
			priority   TEXT NOT NULL DEFAULT 'medium',
			due_date   DATE,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);
	`)
	return err
}

// ─── Lists ────────────────────────────────────────────────

func (r *Repository) CreateList(ctx context.Context, name string) (*model.TaskList, error) {
	list := &model.TaskList{
		ID:   uuid.NewString(),
		Name: name,
	}
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO task_lists (id, name) VALUES ($1, $2)`,
		list.ID, list.Name,
	)
	if err != nil {
		return nil, fmt.Errorf("create list: %w", err)
	}
	return r.GetList(ctx, list.ID)
}

func (r *Repository) GetList(ctx context.Context, id string) (*model.TaskList, error) {
	var list model.TaskList
	err := r.db.GetContext(ctx, &list,
		`SELECT id, name, created_at FROM task_lists WHERE id = $1`, id,
	)
	if err != nil {
		return nil, fmt.Errorf("get list: %w", err)
	}
	todos, err := r.getTodosByList(ctx, id)
	if err != nil {
		return nil, err
	}
	list.Todos = todos
	return &list, nil
}

func (r *Repository) ListLists(ctx context.Context) ([]model.TaskList, error) {
	var lists []model.TaskList
	err := r.db.SelectContext(ctx, &lists,
		`SELECT id, name, created_at FROM task_lists ORDER BY created_at ASC`,
	)
	if err != nil {
		return nil, fmt.Errorf("list lists: %w", err)
	}
	for i := range lists {
		todos, err := r.getTodosByList(ctx, lists[i].ID)
		if err != nil {
			return nil, err
		}
		lists[i].Todos = todos
	}
	return lists, nil
}

func (r *Repository) UpdateList(ctx context.Context, id, name string) (*model.TaskList, error) {
	_, err := r.db.ExecContext(ctx,
		`UPDATE task_lists SET name = $1 WHERE id = $2`, name, id,
	)
	if err != nil {
		return nil, fmt.Errorf("update list: %w", err)
	}
	return r.GetList(ctx, id)
}

func (r *Repository) DeleteList(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM task_lists WHERE id = $1`, id,
	)
	return err
}

// ─── Todos ────────────────────────────────────────────────

func (r *Repository) CreateTodo(ctx context.Context, listID, text string, priority model.Priority, dueDate *string) (*model.Todo, error) {
	todo := &model.Todo{
		ID:       uuid.NewString(),
		ListID:   listID,
		Text:     text,
		Priority: priority,
		DueDate:  dueDate,
	}
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO todos (id, list_id, text, priority, due_date) VALUES ($1, $2, $3, $4, $5)`,
		todo.ID, todo.ListID, todo.Text, todo.Priority, todo.DueDate,
	)
	if err != nil {
		return nil, fmt.Errorf("create todo: %w", err)
	}
	return r.getTodo(ctx, todo.ID)
}

func (r *Repository) UpdateTodo(ctx context.Context, id, text string, done bool, priority model.Priority, dueDate *string) (*model.Todo, error) {
	_, err := r.db.ExecContext(ctx,
		`UPDATE todos SET text = $1, done = $2, priority = $3, due_date = $4 WHERE id = $5`,
		text, done, priority, dueDate, id,
	)
	if err != nil {
		return nil, fmt.Errorf("update todo: %w", err)
	}
	return r.getTodo(ctx, id)
}

func (r *Repository) DeleteTodo(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM todos WHERE id = $1`, id,
	)
	return err
}

// ─── Internal ─────────────────────────────────────────────

func (r *Repository) getTodo(ctx context.Context, id string) (*model.Todo, error) {
	var todo model.Todo
	err := r.db.GetContext(ctx, &todo,
		`SELECT id, list_id, text, done, priority, due_date, created_at FROM todos WHERE id = $1`, id,
	)
	if err != nil {
		return nil, fmt.Errorf("get todo: %w", err)
	}
	return &todo, nil
}

func (r *Repository) getTodosByList(ctx context.Context, listID string) ([]model.Todo, error) {
	var todos []model.Todo
	err := r.db.SelectContext(ctx, &todos,
		`SELECT id, list_id, text, done, priority, due_date, created_at FROM todos WHERE list_id = $1 ORDER BY created_at ASC`,
		listID,
	)
	if err != nil {
		return nil, fmt.Errorf("get todos by list: %w", err)
	}
	return todos, nil
}
