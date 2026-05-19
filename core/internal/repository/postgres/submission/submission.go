package pgsubmission

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"lab-tracker/internal/domain"
	"lab-tracker/internal/repository/models"
)

type SubmissionRepo struct {
	db *sql.DB
}

func New(db *sql.DB) *SubmissionRepo {
	return &SubmissionRepo{db: db}
}

func (r *SubmissionRepo) GetStudentAssignmentMeta(ctx context.Context, assignmentID int64, studentID uuid.UUID) (*models.Assignment, *models.Submission, error) {
	const q = `
		SELECT
			a.id,
			a.lab_work_id,
			a.group_id,
			a.deadline,
			a.status,
			a.created_at,
			s.id,
			s.assignment_id,
			s.student_id,
			s.text_report,
			s.file_path,
			s.status,
			s.attempt_number,
			s.submitted_at
		FROM assignments a
		JOIN users u ON u.group_id = a.group_id
		LEFT JOIN submissions s ON s.assignment_id = a.id AND s.student_id = u.id
		WHERE a.id = $1 AND u.id = $2
	`

	assignment := &models.Assignment{}
	var submissionID sql.NullInt64
	var submissionAssignmentID sql.NullInt64
	var submissionStudentID uuid.NullUUID
	var textReport sql.NullString
	var filePath sql.NullString
	var status sql.NullString
	var attemptNumber sql.NullInt64
	var submittedAt sql.NullTime

	err := r.db.QueryRowContext(ctx, q, assignmentID, studentID).Scan(
		&assignment.ID,
		&assignment.LabWorkID,
		&assignment.GroupID,
		&assignment.Deadline,
		&assignment.Status,
		&assignment.CreatedAt,
		&submissionID,
		&submissionAssignmentID,
		&submissionStudentID,
		&textReport,
		&filePath,
		&status,
		&attemptNumber,
		&submittedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, fmt.Errorf("submission repo: assignment meta: %w", models.ErrAssignmentNotFound)
		}
		err = r.db.QueryRowContext(ctx, `
			SELECT a.id, a.lab_work_id, a.group_id, a.deadline, a.status, a.created_at
			FROM assignments a
			JOIN users u ON u.group_id = a.group_id
			WHERE a.id = $1 AND u.id = $2
		`, assignmentID, studentID).Scan(
			&assignment.ID,
			&assignment.LabWorkID,
			&assignment.GroupID,
			&assignment.Deadline,
			&assignment.Status,
			&assignment.CreatedAt,
		)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, nil, fmt.Errorf("submission repo: assignment meta: %w", models.ErrAssignmentNotFound)
			}
			return nil, nil, fmt.Errorf("submission repo: assignment meta fallback: %w", err)
		}
		return assignment, nil, nil
	}

	if !submissionID.Valid {
		return assignment, nil, nil
	}

	submission := &models.Submission{
		ID:           submissionID.Int64,
		AssignmentID: submissionAssignmentID.Int64,
		StudentID:    submissionStudentID.UUID,
		TextReport:   textReport.String,
		Status:       models.SubmissionStatus(status.String),
	}
	if attemptNumber.Valid {
		submission.AttemptNumber = int(attemptNumber.Int64)
	}
	if filePath.Valid {
		value := filePath.String
		submission.FilePath = &value
	}
	if submittedAt.Valid {
		value := submittedAt.Time
		submission.SubmittedAt = &value
	}

	return assignment, submission, nil
}

func (r *SubmissionRepo) Create(ctx context.Context, input domain.SubmissionInput) (*models.Submission, error) {
	const q = `
		INSERT INTO submissions (assignment_id, student_id, text_report, file_path, status, attempt_number, submitted_at)
		VALUES ($1, $2, $3, $4, 'submitted', 1, NOW())
		RETURNING id, assignment_id, student_id, text_report, file_path, status, attempt_number, submitted_at
	`

	var submission models.Submission
	if err := r.db.QueryRowContext(
		ctx,
		q,
		input.AssignmentID,
		input.StudentID,
		input.TextReport,
		input.FilePath,
	).Scan(
		&submission.ID,
		&submission.AssignmentID,
		&submission.StudentID,
		&submission.TextReport,
		&submission.FilePath,
		&submission.Status,
		&submission.AttemptNumber,
		&submission.SubmittedAt,
	); err != nil {
		return nil, fmt.Errorf("submission repo: create: %w", err)
	}

	return &submission, nil
}

func (r *SubmissionRepo) UpdateDraft(ctx context.Context, submissionID int64, input domain.SubmissionInput) (*models.Submission, error) {
	const q = `
		UPDATE submissions
		SET text_report = $1,
		    file_path = $2,
		    status = 'submitted',
		    attempt_number = CASE
		        WHEN status = 'revision' THEN attempt_number + 1
		        ELSE GREATEST(attempt_number, 1)
		    END,
		    submitted_at = NOW()
		WHERE id = $3
		RETURNING id, assignment_id, student_id, text_report, file_path, status, attempt_number, submitted_at
	`

	var submission models.Submission
	if err := r.db.QueryRowContext(ctx, q, input.TextReport, input.FilePath, submissionID).Scan(
		&submission.ID,
		&submission.AssignmentID,
		&submission.StudentID,
		&submission.TextReport,
		&submission.FilePath,
		&submission.Status,
		&submission.AttemptNumber,
		&submission.SubmittedAt,
	); err != nil {
		return nil, fmt.Errorf("submission repo: update draft: %w", err)
	}

	return &submission, nil
}

func (r *SubmissionRepo) ListTeacherSubmissions(ctx context.Context, teacherID uuid.UUID) ([]domain.TeacherSubmission, error) {
	const q = `
		SELECT
			s.id,
			s.assignment_id,
			u.id,
			u.full_name,
			g.name,
			lw.title,
			a.deadline,
			s.attempt_number,
			COALESCE(s.text_report, ''),
			s.file_path,
			s.status,
			s.submitted_at,
			gr.grade,
			gr.comment
		FROM submissions s
		JOIN users u ON u.id = s.student_id
		JOIN assignments a ON a.id = s.assignment_id
		JOIN groups g ON g.id = a.group_id
		JOIN lab_works lw ON lw.id = a.lab_work_id
		LEFT JOIN grades gr ON gr.submission_id = s.id
		WHERE g.teacher_id = $1
		ORDER BY s.submitted_at DESC NULLS LAST, s.id DESC
	`

	rows, err := r.db.QueryContext(ctx, q, teacherID)
	if err != nil {
		return nil, fmt.Errorf("submission repo: list teacher submissions: %w", err)
	}
	defer rows.Close()

	result := make([]domain.TeacherSubmission, 0)
	for rows.Next() {
		var item domain.TeacherSubmission
		var deadline *time.Time
		var submittedAt *time.Time
		var comment *string

		if err := rows.Scan(
			&item.SubmissionID,
			&item.AssignmentID,
			&item.StudentID,
			&item.StudentName,
			&item.GroupName,
			&item.LabWorkTitle,
			&deadline,
			&item.AttemptNumber,
			&item.TextReport,
			&item.FilePath,
			&item.Status,
			&submittedAt,
			&item.Grade,
			&comment,
		); err != nil {
			return nil, fmt.Errorf("submission repo: scan teacher submission: %w", err)
		}

		if deadline != nil {
			formatted := deadline.Format(time.RFC3339)
			item.Deadline = &formatted
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
		return nil, fmt.Errorf("submission repo: teacher rows: %w", err)
	}

	return result, nil
}

func (r *SubmissionRepo) GetForTeacher(ctx context.Context, submissionID int64, teacherID uuid.UUID) (*models.Submission, error) {
	const q = `
		SELECT s.id, s.assignment_id, s.student_id, COALESCE(s.text_report, ''), s.file_path, s.status, s.attempt_number, s.submitted_at
		FROM submissions s
		JOIN assignments a ON a.id = s.assignment_id
		JOIN groups g ON g.id = a.group_id
		WHERE s.id = $1 AND g.teacher_id = $2
	`

	var submission models.Submission
	if err := r.db.QueryRowContext(ctx, q, submissionID, teacherID).Scan(
		&submission.ID,
		&submission.AssignmentID,
		&submission.StudentID,
		&submission.TextReport,
		&submission.FilePath,
		&submission.Status,
		&submission.AttemptNumber,
		&submission.SubmittedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("submission repo: get for teacher: %w", models.ErrSubmissionNotFound)
		}
		return nil, fmt.Errorf("submission repo: get for teacher: %w", err)
	}

	return &submission, nil
}

func (r *SubmissionRepo) SetStatus(ctx context.Context, submissionID int64, status models.SubmissionStatus) error {
	const q = `UPDATE submissions SET status = $1 WHERE id = $2`
	if _, err := r.db.ExecContext(ctx, q, status, submissionID); err != nil {
		return fmt.Errorf("submission repo: set status: %w", err)
	}
	return nil
}
