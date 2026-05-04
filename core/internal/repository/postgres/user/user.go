package pguser

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/lib/pq"

	"lab-tracker/internal/repository/models"
)

type UserRepo struct {
	db *sql.DB
}

func New(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) CreateUser(
	ctx context.Context,
	username, passwordHash, fullName string,
	groupID *int,
) (*models.User, error) {
	const q = `
		INSERT INTO users (id, username, password_hash, full_name, group_id, role)
		VALUES ($1, $2, $3, $4, $5, 'student')
		RETURNING id, username, password_hash, full_name, group_id, role, created_at
	`
	var u models.User
	err := r.db.QueryRowContext(ctx, q, uuid.New(), username, passwordHash, fullName, groupID).Scan(
		&u.ID, &u.Username, &u.PasswordHash, &u.FullName, &u.GroupID, &u.Role, &u.CreatedAt,
	)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			switch pqErr.Code {
			case "23505":
				return nil, fmt.Errorf("user repo: create: %w", models.ErrUserAlreadyExists)
			case "23503":
				return nil, fmt.Errorf("user repo: create: %w", models.ErrGroupNotFound)
			}
		}
		return nil, fmt.Errorf("user repo: create: %w", err)
	}
	return &u, nil
}

func (r *UserRepo) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	const q = `
		SELECT id, username, password_hash, full_name, group_id, role, created_at
		FROM users
		WHERE username = $1
	`
	var u models.User
	err := r.db.QueryRowContext(ctx, q, username).Scan(
		&u.ID, &u.Username, &u.PasswordHash, &u.FullName, &u.GroupID, &u.Role, &u.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user repo: get: %w", models.ErrUserNotFound)
		}
		return nil, fmt.Errorf("user repo: get: %w", err)
	}

	return &u, nil
}
