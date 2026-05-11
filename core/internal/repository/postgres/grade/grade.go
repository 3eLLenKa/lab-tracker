package pggrade

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type GradeRepo struct {
	db *sql.DB
}

func New(db *sql.DB) *GradeRepo {
	return &GradeRepo{db: db}
}

func (r *GradeRepo) Save(ctx context.Context, submissionID int64, teacherID uuid.UUID, grade int, comment string) error {
	const q = `
		INSERT INTO grades (submission_id, teacher_id, grade, comment, graded_at)
		VALUES ($1, $2, $3, $4, NOW())
		ON CONFLICT (submission_id)
		DO UPDATE SET
			teacher_id = EXCLUDED.teacher_id,
			grade = EXCLUDED.grade,
			comment = EXCLUDED.comment,
			graded_at = EXCLUDED.graded_at
	`

	if _, err := r.db.ExecContext(ctx, q, submissionID, teacherID, grade, comment); err != nil {
		return fmt.Errorf("grade repo: save: %w", err)
	}

	return nil
}
