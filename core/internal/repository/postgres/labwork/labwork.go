package pglabwork

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"lab-tracker/internal/repository/models"
)

type LabWorkRepo struct {
	db *sql.DB
}

func New(db *sql.DB) *LabWorkRepo {
	return &LabWorkRepo{db: db}
}

func (r *LabWorkRepo) List(ctx context.Context, search string, limit, offset int) ([]models.LabWork, int, error) {
	where, args := buildWhere(search)
	countQuery := "SELECT COUNT(*) FROM lab_works" + where

	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("labwork repo: count: %w", err)
	}

	args = append(args, limit, offset)
	q := `
		SELECT id, title, description, goal, equipment, reagents, procedure, file_path, created_by, created_at
		FROM lab_works
	` + where + `
		ORDER BY id DESC
		LIMIT $` + fmt.Sprint(len(args)-1) + ` OFFSET $` + fmt.Sprint(len(args))

	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("labwork repo: list: %w", err)
	}
	defer rows.Close()

	items := make([]models.LabWork, 0)
	for rows.Next() {
		lab, err := scanLabWork(rows)
		if err != nil {
			return nil, 0, err
		}
		items = append(items, lab)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("labwork repo: list rows: %w", err)
	}

	return items, total, nil
}

func (r *LabWorkRepo) GetByID(ctx context.Context, id int64) (*models.LabWork, error) {
	const q = `
		SELECT id, title, description, goal, equipment, reagents, procedure, file_path, created_by, created_at
		FROM lab_works
		WHERE id = $1
	`

	row := r.db.QueryRowContext(ctx, q, id)
	lab, err := scanLabWork(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("labwork repo: get: %w", models.ErrLabWorkNotFound)
		}
		return nil, fmt.Errorf("labwork repo: get: %w", err)
	}

	return &lab, nil
}

func (r *LabWorkRepo) Create(ctx context.Context, lab models.LabWork) (*models.LabWork, error) {
	const q = `
		INSERT INTO lab_works (title, description, goal, equipment, reagents, procedure, file_path, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, title, description, goal, equipment, reagents, procedure, file_path, created_by, created_at
	`

	row := r.db.QueryRowContext(
		ctx,
		q,
		lab.Title,
		lab.Description,
		lab.Goal,
		lab.Equipment,
		lab.Reagents,
		lab.Procedure,
		lab.FilePath,
		lab.CreatedBy,
	)
	created, err := scanLabWork(row)
	if err != nil {
		return nil, fmt.Errorf("labwork repo: create: %w", err)
	}

	return &created, nil
}

func (r *LabWorkRepo) Update(ctx context.Context, id int64, lab models.LabWork) (*models.LabWork, error) {
	const q = `
		UPDATE lab_works
		SET title = $1,
		    description = $2,
		    goal = $3,
		    equipment = $4,
		    reagents = $5,
		    procedure = $6,
		    file_path = $7
		WHERE id = $8
		RETURNING id, title, description, goal, equipment, reagents, procedure, file_path, created_by, created_at
	`

	row := r.db.QueryRowContext(
		ctx,
		q,
		lab.Title,
		lab.Description,
		lab.Goal,
		lab.Equipment,
		lab.Reagents,
		lab.Procedure,
		lab.FilePath,
		id,
	)
	updated, err := scanLabWork(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("labwork repo: update: %w", models.ErrLabWorkNotFound)
		}
		return nil, fmt.Errorf("labwork repo: update: %w", err)
	}

	return &updated, nil
}

func (r *LabWorkRepo) Delete(ctx context.Context, id int64) error {
	const q = `DELETE FROM lab_works WHERE id = $1`

	res, err := r.db.ExecContext(ctx, q, id)
	if err != nil {
		return fmt.Errorf("labwork repo: delete: %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("labwork repo: delete affected: %w", err)
	}
	if affected == 0 {
		return fmt.Errorf("labwork repo: delete: %w", models.ErrLabWorkNotFound)
	}

	return nil
}

type scanner interface {
	Scan(dest ...any) error
}

func scanLabWork(row scanner) (models.LabWork, error) {
	var lab models.LabWork
	err := row.Scan(
		&lab.ID,
		&lab.Title,
		&lab.Description,
		&lab.Goal,
		&lab.Equipment,
		&lab.Reagents,
		&lab.Procedure,
		&lab.FilePath,
		&lab.CreatedBy,
		&lab.CreatedAt,
	)
	return lab, err
}

func buildWhere(search string) (string, []any) {
	search = strings.TrimSpace(search)
	if search == "" {
		return "", nil
	}
	return " WHERE title ILIKE $1", []any{"%" + search + "%"}
}
