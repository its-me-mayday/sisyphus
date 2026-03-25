package model

import "time"

type Role string

const (
	RoleViewer Role = "viewer"
	RoleEditor Role = "editor"
	RoleOwner  Role = "owner"
)

type User struct {
	ID           string    `db:"id"`
	Username     string    `db:"username"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash"`
	CreatedAt    time.Time `db:"created_at"`
}

type Permission struct {
	ID        string    `db:"id"`
	UserID    string    `db:"user_id"`
	Resource  string    `db:"resource"`
	Role      Role      `db:"role"`
	CreatedAt time.Time `db:"created_at"`
}
