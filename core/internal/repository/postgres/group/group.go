package pggroup

import (
	"context"
	"database/sql"
	"fmt"

	"lab-tracker/internal/repository/models"
)

type GroupRepo struct {
	db *sql.DB
}

func New(db *sql.DB) *GroupRepo {
	return &GroupRepo{db: db}
}

func (r *GroupRepo) List(ctx context.Context) ([]models.Group, error) {
	const q = `
		SELECT id, name, course
		FROM groups
		ORDER BY course, name
	`

	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("group repo: list: %w", err)
	}
	defer rows.Close()

	groups := make([]models.Group, 0)
	for rows.Next() {
		var g models.Group
		if err := rows.Scan(&g.ID, &g.Name, &g.Course); err != nil {
			return nil, fmt.Errorf("group repo: list scan: %w", err)
		}
		groups = append(groups, g)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("group repo: list rows: %w", err)
	}

	return groups, nil
}
