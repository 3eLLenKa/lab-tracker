package pgassignment

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"

	"lab-tracker/internal/domain"
)

type AssignmentRepo struct {
	db *sql.DB
}

func New(db *sql.DB) *AssignmentRepo {
	return &AssignmentRepo{db: db}
}

func (r *AssignmentRepo) ListStudentAssignments(ctx context.Context, studentID uuid.UUID) ([]domain.StudentAssignment, error) {
	const q = `
		SELECT
			a.id,
			lw.id,
			lw.title,
			COALESCE(lw.description, ''),
			a.deadline,
			a.status,
			s.id,
			s.status,
			s.attempt_number,
			s.text_report,
			s.file_path,
			s.submitted_at,
			g.grade,
			g.comment
		FROM assignments a
		JOIN users u ON u.group_id = a.group_id
		JOIN lab_works lw ON lw.id = a.lab_work_id
		LEFT JOIN submissions s ON s.assignment_id = a.id AND s.student_id = u.id
		LEFT JOIN grades g ON g.submission_id = s.id
		WHERE u.id = $1
		ORDER BY a.deadline NULLS LAST, a.id
	`

	rows, err := r.db.QueryContext(ctx, q, studentID)
	if err != nil {
		return nil, fmt.Errorf("assignment repo: list student assignments: %w", err)
	}
	defer rows.Close()

	result := make([]domain.StudentAssignment, 0)
	for rows.Next() {
		var item domain.StudentAssignment
		var deadline *time.Time
		var submissionStatus *string
		var attemptNumber sql.NullInt64
		var textReport *string
		var submittedAt *time.Time
		var comment *string

		if err := rows.Scan(
			&item.AssignmentID,
			&item.LabWorkID,
			&item.Title,
			&item.Description,
			&deadline,
			&item.AssignmentStatus,
			&item.SubmissionID,
			&submissionStatus,
			&attemptNumber,
			&textReport,
			&item.FilePath,
			&submittedAt,
			&item.Grade,
			&comment,
		); err != nil {
			return nil, fmt.Errorf("assignment repo: scan student assignment: %w", err)
		}

		if deadline != nil {
			formatted := deadline.Format(time.RFC3339)
			item.Deadline = &formatted
		}
		if submissionStatus != nil {
			value := *submissionStatus
			item.SubmissionStatus = &value
		}
		if attemptNumber.Valid {
			value := int(attemptNumber.Int64)
			item.AttemptNumber = &value
		}
		if textReport != nil {
			value := *textReport
			item.TextReport = &value
		}
		if submittedAt != nil {
			formatted := submittedAt.Format(time.RFC3339)
			item.SubmittedAt = &formatted
		}
		if comment != nil {
			value := *comment
			item.TeacherComment = &value
		}

		result = append(result, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("assignment repo: student rows: %w", err)
	}

	return result, nil
}
