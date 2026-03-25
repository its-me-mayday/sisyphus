package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/its-me-mayday/sisyphus/acl-service/internal/model"
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
		CREATE TABLE IF NOT EXISTS users (
			id            TEXT PRIMARY KEY,
			username      TEXT NOT NULL UNIQUE,
			email         TEXT NOT NULL UNIQUE,
			password_hash TEXT NOT NULL,
			created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);

		CREATE TABLE IF NOT EXISTS permissions (
			id         TEXT PRIMARY KEY,
			user_id    TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			resource   TEXT NOT NULL,
			role       TEXT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			UNIQUE(user_id, resource)
		);
	`)
	return err
}

// ─── Users ────────────────────────────────────────────────

func (r *Repository) CreateUser(ctx context.Context, username, email, passwordHash string) (*model.User, error) {
	user := &model.User{
		ID:           uuid.NewString(),
		Username:     username,
		Email:        email,
		PasswordHash: passwordHash,
	}
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO users (id, username, email, password_hash) VALUES ($1, $2, $3, $4)`,
		user.ID, user.Username, user.Email, user.PasswordHash,
	)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	return r.GetUser(ctx, user.ID)
}

func (r *Repository) GetUser(ctx context.Context, id string) (*model.User, error) {
	var user model.User
	err := r.db.GetContext(ctx, &user,
		`SELECT id, username, email, password_hash, created_at FROM users WHERE id = $1`, id,
	)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	return &user, nil
}

func (r *Repository) DeleteUser(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM users WHERE id = $1`, id)
	return err
}

// ─── Permissions ──────────────────────────────────────────

func (r *Repository) GrantPermission(ctx context.Context, userID, resource string, role model.Role) (*model.Permission, error) {
	p := &model.Permission{
		ID:       uuid.NewString(),
		UserID:   userID,
		Resource: resource,
		Role:     role,
	}
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO permissions (id, user_id, resource, role)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id, resource) DO UPDATE SET role = EXCLUDED.role
	`, p.ID, p.UserID, p.Resource, p.Role)
	if err != nil {
		return nil, fmt.Errorf("grant permission: %w", err)
	}
	return r.getPermission(ctx, p.UserID, p.Resource)
}

func (r *Repository) RevokePermission(ctx context.Context, userID, resource string) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM permissions WHERE user_id = $1 AND resource = $2`, userID, resource,
	)
	return err
}

func (r *Repository) CheckPermission(ctx context.Context, userID, resource string, role model.Role) (bool, error) {
	roleRank := map[model.Role]int{
		model.RoleViewer: 1,
		model.RoleEditor: 2,
		model.RoleOwner:  3,
	}
	var current model.Permission
	err := r.db.GetContext(ctx, &current,
		`SELECT role FROM permissions WHERE user_id = $1 AND resource = $2`, userID, resource,
	)
	if err != nil {
		return false, nil
	}
	return roleRank[current.Role] >= roleRank[role], nil
}

func (r *Repository) ListPermissions(ctx context.Context, userID string) ([]model.Permission, error) {
	var perms []model.Permission
	err := r.db.SelectContext(ctx, &perms,
		`SELECT id, user_id, resource, role, created_at FROM permissions WHERE user_id = $1 ORDER BY created_at ASC`,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("list permissions: %w", err)
	}
	return perms, nil
}

// ─── Internal ─────────────────────────────────────────────

func (r *Repository) getPermission(ctx context.Context, userID, resource string) (*model.Permission, error) {
	var p model.Permission
	err := r.db.GetContext(ctx, &p,
		`SELECT id, user_id, resource, role, created_at FROM permissions WHERE user_id = $1 AND resource = $2`,
		userID, resource,
	)
	if err != nil {
		return nil, fmt.Errorf("get permission: %w", err)
	}
	return &p, nil
}
