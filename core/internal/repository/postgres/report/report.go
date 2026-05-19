package pgreport

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"

	"lab-tracker/internal/domain"
)

type ReportRepo struct {
	db *sql.DB
}

func New(db *sql.DB) *ReportRepo {
	return &ReportRepo{db: db}
}

func (r *ReportRepo) GetStudentProgress(ctx context.Context, studentID uuid.UUID) (*domain.StudentProgress, error) {
	const q = `
		SELECT
			COUNT(a.id),
			COUNT(*) FILTER (WHERE COALESCE(s.status::text, 'draft') = 'draft'),
			COUNT(*) FILTER (WHERE s.status = 'submitted'),
			COUNT(*) FILTER (WHERE s.status = 'revision'),
			COUNT(*) FILTER (WHERE s.status = 'reviewed'),
			AVG(gr.grade)
		FROM assignments a
		JOIN users u ON u.group_id = a.group_id
		LEFT JOIN submissions s ON s.assignment_id = a.id AND s.student_id = u.id
		LEFT JOIN grades gr ON gr.submission_id = s.id
		WHERE u.id = $1
	`

	var progress domain.StudentProgress
	var avg sql.NullFloat64
	if err := r.db.QueryRowContext(ctx, q, studentID).Scan(
		&progress.TotalAssignments,
		&progress.DraftCount,
		&progress.SubmittedCount,
		&progress.RevisionCount,
		&progress.ReviewedCount,
		&avg,
	); err != nil {
		return nil, fmt.Errorf("report repo: student progress: %w", err)
	}

	if progress.TotalAssignments > 0 {
		progress.CompletionRate = progress.ReviewedCount * 100 / progress.TotalAssignments
	}
	if avg.Valid {
		value := avg.Float64
		progress.AverageGrade = &value
	}

	return &progress, nil
}

func (r *ReportRepo) GetAdminStats(ctx context.Context) (*domain.AdminStats, error) {
	const q = `
		SELECT
			(SELECT COUNT(*) FROM users),
			(SELECT COUNT(*) FROM users WHERE role = 'student'),
			(SELECT COUNT(*) FROM users WHERE role = 'teacher'),
			(SELECT COUNT(*) FROM groups),
			(SELECT COUNT(*) FROM lab_works),
			(SELECT COUNT(*) FROM assignments),
			(SELECT COUNT(*) FROM submissions),
			(SELECT COUNT(*) FROM submissions WHERE status = 'draft'),
			(SELECT COUNT(*) FROM submissions WHERE status = 'submitted'),
			(SELECT COUNT(*) FROM submissions WHERE status = 'revision'),
			(SELECT COUNT(*) FROM submissions WHERE status = 'reviewed'),
			(SELECT AVG(grade) FROM grades)
	`

	var stats domain.AdminStats
	var avg sql.NullFloat64
	if err := r.db.QueryRowContext(ctx, q).Scan(
		&stats.UsersTotal,
		&stats.StudentsTotal,
		&stats.TeachersTotal,
		&stats.GroupsTotal,
		&stats.LabWorksTotal,
		&stats.AssignmentsTotal,
		&stats.SubmissionsTotal,
		&stats.DraftCount,
		&stats.SubmittedCount,
		&stats.RevisionCount,
		&stats.ReviewedCount,
		&avg,
	); err != nil {
		return nil, fmt.Errorf("report repo: admin stats: %w", err)
	}

	if avg.Valid {
		value := avg.Float64
		stats.AverageGrade = &value
	}

	return &stats, nil
}

func (r *ReportRepo) ListReportRows(ctx context.Context) ([]domain.ReportRow, error) {
	const q = `
		SELECT
			u.full_name,
			g.name,
			lw.title,
			s.status,
			s.attempt_number,
			s.submitted_at,
			gr.grade,
			gr.comment
		FROM submissions s
		JOIN users u ON u.id = s.student_id
		JOIN assignments a ON a.id = s.assignment_id
		JOIN groups g ON g.id = a.group_id
		JOIN lab_works lw ON lw.id = a.lab_work_id
		LEFT JOIN grades gr ON gr.submission_id = s.id
		ORDER BY g.name, u.full_name, lw.title
	`

	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("report repo: list report rows: %w", err)
	}
	defer rows.Close()

	result := make([]domain.ReportRow, 0)
	for rows.Next() {
		var row domain.ReportRow
		var submittedAt sql.NullTime
		var grade sql.NullInt64
		var comment sql.NullString

		if err := rows.Scan(
			&row.StudentName,
			&row.GroupName,
			&row.LabWorkTitle,
			&row.Status,
			&row.AttemptNumber,
			&submittedAt,
			&grade,
			&comment,
		); err != nil {
			return nil, fmt.Errorf("report repo: scan report row: %w", err)
		}

		if submittedAt.Valid {
			value := submittedAt.Time.Format(time.RFC3339)
			row.SubmittedAt = &value
		}
		if grade.Valid {
			value := int(grade.Int64)
			row.Grade = &value
		}
		if comment.Valid {
			value := comment.String
			row.Comment = &value
		}

		result = append(result, row)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("report repo: report rows: %w", err)
	}

	return result, nil
}
