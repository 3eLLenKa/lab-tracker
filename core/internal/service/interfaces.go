package service

import (
	"context"
	"github.com/google/uuid"
	"lab-tracker/internal/domain"
	"lab-tracker/internal/repository/models"
)

type UserRepo interface {
	CreateUser(ctx context.Context, username, passwordHash, fullName string, groupID *int) (*models.User, error)
	GetByUsername(ctx context.Context, username string) (*models.User, error)
}

type GroupRepo interface {
	List(ctx context.Context) ([]models.Group, error)
}

type LabWorkRepo interface {
	List(ctx context.Context, search string, limit, offset int) ([]models.LabWork, int, error)
	GetByID(ctx context.Context, id int64) (*models.LabWork, error)
	Create(ctx context.Context, lab models.LabWork) (*models.LabWork, error)
	Update(ctx context.Context, id int64, lab models.LabWork) (*models.LabWork, error)
	Delete(ctx context.Context, id int64) error
}

type AssignmentRepo interface {
	ListStudentAssignments(ctx context.Context, studentID uuid.UUID) ([]domain.StudentAssignment, error)
}

type SubmissionRepo interface {
	GetStudentAssignmentMeta(ctx context.Context, assignmentID int64, studentID uuid.UUID) (*models.Assignment, *models.Submission, error)
	Create(ctx context.Context, input domain.SubmissionInput) (*models.Submission, error)
	UpdateDraft(ctx context.Context, submissionID int64, input domain.SubmissionInput) (*models.Submission, error)
	ListTeacherSubmissions(ctx context.Context, teacherID uuid.UUID) ([]domain.TeacherSubmission, error)
	GetForTeacher(ctx context.Context, submissionID int64, teacherID uuid.UUID) (*models.Submission, error)
	SetStatus(ctx context.Context, submissionID int64, status models.SubmissionStatus) error
}

type GradeRepo interface {
	Save(ctx context.Context, submissionID int64, teacherID uuid.UUID, grade int, comment string) error
}

type ReportRepo interface {
	GetStudentProgress(ctx context.Context, studentID uuid.UUID) (*domain.StudentProgress, error)
	GetAdminStats(ctx context.Context) (*domain.AdminStats, error)
	ListReportRows(ctx context.Context) ([]domain.ReportRow, error)
}
