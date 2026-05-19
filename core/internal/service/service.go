package service

import (
	"bytes"
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"lab-tracker/internal/domain"
	"lab-tracker/internal/repository/models"
)

type Service struct {
	log            *zap.Logger
	userRepo       UserRepo
	groupRepo      GroupRepo
	labWorkRepo    LabWorkRepo
	assignmentRepo AssignmentRepo
	submissionRepo SubmissionRepo
	gradeRepo      GradeRepo
	reportRepo     ReportRepo
	jwtSecret      string
}

func New(
	log *zap.Logger,
	userRepo UserRepo,
	groupRepo GroupRepo,
	labWorkRepo LabWorkRepo,
	assignmentRepo AssignmentRepo,
	submissionRepo SubmissionRepo,
	gradeRepo GradeRepo,
	reportRepo ReportRepo,
	jwtSecret string,
) *Service {
	return &Service{
		log:            log,
		userRepo:       userRepo,
		groupRepo:      groupRepo,
		labWorkRepo:    labWorkRepo,
		assignmentRepo: assignmentRepo,
		submissionRepo: submissionRepo,
		gradeRepo:      gradeRepo,
		reportRepo:     reportRepo,
		jwtSecret:      jwtSecret,
	}
}

type tokenClaims struct {
	UserID   uuid.UUID       `json:"user_id"`
	Username string          `json:"username"`
	Role     domain.UserRole `json:"role"`
	jwt.StandardClaims
}

func (s *Service) generateToken(u *domain.User) (string, error) {
	c := &tokenClaims{
		UserID:   u.ID,
		Username: u.Username,
		Role:     u.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, c).SignedString([]byte(s.jwtSecret))
}

func (s *Service) Register(
	ctx context.Context,
	username, password, fullName string,
	groupID int,
) (string, error) {
	log := s.log.With(
		zap.String("op", "service.Register"),
		zap.String("username", username),
	)

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("service: hash password: %w", err)
	}

	var gid *int
	if groupID != 0 {
		gid = &groupID
	}

	u, err := s.userRepo.CreateUser(ctx, username, string(hash), fullName, gid)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrUserAlreadyExists):
			log.Warn("user already exists")
			return "", domain.ErrUserAlreadyExists
		case errors.Is(err, domain.ErrGroupNotFound):
			log.Warn("group not found")
			return "", domain.ErrGroupNotFound
		}
		log.Error("failed to create user", zap.Error(err))
		return "", fmt.Errorf("service: create user: %w", err)
	}

	token, err := s.generateToken(&domain.User{
		ID:       u.ID,
		Username: u.Username,
		Role:     domain.UserRole(u.Role),
	})
	if err != nil {
		log.Error("failed to generate token", zap.Error(err))
		return "", fmt.Errorf("service: generate token: %w", err)
	}

	log.Info("user registered", zap.String("user_id", u.ID.String()))
	return token, nil
}

func (s *Service) Login(ctx context.Context, username, password string) (string, error) {
	log := s.log.With(zap.String("op", "service.Login"), zap.String("username", username))

	u, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			log.Warn("user not found")
			return "", domain.ErrInvalidCredentials
		}
		log.Error("failed to get user", zap.Error(err))
		return "", fmt.Errorf("service: get user: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		log.Warn("invalid password")
		return "", domain.ErrInvalidCredentials
	}

	token, err := s.generateToken(&domain.User{
		ID:       u.ID,
		Username: u.Username,
		Role:     domain.UserRole(u.Role),
	})
	if err != nil {
		log.Error("failed to generate token", zap.Error(err))
		return "", fmt.Errorf("service: generate token: %w", err)
	}

	log.Info("user logged in", zap.String("user_id", u.ID.String()))
	return token, nil
}

func (s *Service) ListGroups(ctx context.Context) ([]domain.Group, error) {
	groups, err := s.groupRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("service: list groups: %w", err)
	}

	result := make([]domain.Group, 0, len(groups))
	for _, g := range groups {
		result = append(result, domain.Group{
			ID:     g.ID,
			Name:   g.Name,
			Course: g.Course,
		})
	}

	return result, nil
}

func (s *Service) ListLabWorks(ctx context.Context, filter domain.LabWorkFilter) (*domain.LabWorkList, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 10
	}
	if filter.Limit > 20 {
		filter.Limit = 20
	}

	items, total, err := s.labWorkRepo.List(ctx, filter.Search, filter.Limit, (filter.Page-1)*filter.Limit)
	if err != nil {
		return nil, fmt.Errorf("service: list lab works: %w", err)
	}

	totalPages := 0
	if total > 0 {
		totalPages = (total + filter.Limit - 1) / filter.Limit
	}

	return &domain.LabWorkList{
		Items:      mapLabWorks(items),
		Total:      total,
		Page:       filter.Page,
		Limit:      filter.Limit,
		TotalPages: totalPages,
	}, nil
}

func (s *Service) GetLabWork(ctx context.Context, id int64) (*domain.LabWork, error) {
	lab, err := s.labWorkRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, models.ErrLabWorkNotFound) {
			return nil, domain.ErrLabWorkNotFound
		}
		return nil, fmt.Errorf("service: get lab work: %w", err)
	}

	result := mapLabWork(*lab)
	return &result, nil
}

func (s *Service) CreateLabWork(ctx context.Context, input domain.LabWorkInput) (*domain.LabWork, error) {
	if err := validateLabWork(input); err != nil {
		return nil, err
	}

	lab, err := s.labWorkRepo.Create(ctx, models.LabWork{
		Title:       input.Title,
		Description: input.Description,
		Goal:        input.Goal,
		Equipment:   input.Equipment,
		Reagents:    input.Reagents,
		Procedure:   input.Procedure,
		FilePath:    input.FilePath,
		CreatedBy:   input.CreatedBy,
	})
	if err != nil {
		return nil, fmt.Errorf("service: create lab work: %w", err)
	}

	result := mapLabWork(*lab)
	return &result, nil
}

func (s *Service) UpdateLabWork(ctx context.Context, id int64, input domain.LabWorkInput) (*domain.LabWork, error) {
	if err := validateLabWork(input); err != nil {
		return nil, err
	}

	lab, err := s.labWorkRepo.Update(ctx, id, models.LabWork{
		Title:       input.Title,
		Description: input.Description,
		Goal:        input.Goal,
		Equipment:   input.Equipment,
		Reagents:    input.Reagents,
		Procedure:   input.Procedure,
		FilePath:    input.FilePath,
	})
	if err != nil {
		if errors.Is(err, models.ErrLabWorkNotFound) {
			return nil, domain.ErrLabWorkNotFound
		}
		return nil, fmt.Errorf("service: update lab work: %w", err)
	}

	result := mapLabWork(*lab)
	return &result, nil
}

func (s *Service) DeleteLabWork(ctx context.Context, id int64) error {
	if err := s.labWorkRepo.Delete(ctx, id); err != nil {
		if errors.Is(err, models.ErrLabWorkNotFound) {
			return domain.ErrLabWorkNotFound
		}
		return fmt.Errorf("service: delete lab work: %w", err)
	}
	return nil
}

func (s *Service) ListStudentAssignments(ctx context.Context, studentID uuid.UUID) ([]domain.StudentAssignment, error) {
	assignments, err := s.assignmentRepo.ListStudentAssignments(ctx, studentID)
	if err != nil {
		return nil, fmt.Errorf("service: list student assignments: %w", err)
	}
	return assignments, nil
}

func (s *Service) SubmitAssignment(ctx context.Context, input domain.SubmissionInput) error {
	if input.AssignmentID <= 0 || strings.TrimSpace(input.TextReport) == "" {
		return domain.ErrInvalidRequest
	}

	assignment, submission, err := s.submissionRepo.GetStudentAssignmentMeta(ctx, input.AssignmentID, input.StudentID)
	if err != nil {
		if errors.Is(err, models.ErrAssignmentNotFound) {
			return domain.ErrAssignmentNotFound
		}
		return fmt.Errorf("service: get assignment meta: %w", err)
	}

	if assignment.Deadline != nil && assignment.Deadline.Before(time.Now()) {
		return domain.ErrDeadlinePassed
	}

	if submission == nil {
		_, err = s.submissionRepo.Create(ctx, input)
		if err != nil {
			return fmt.Errorf("service: create submission: %w", err)
		}
		return nil
	}

	switch submission.Status {
	case models.SubmissionDraft, models.SubmissionRevision:
	case models.SubmissionSubmitted:
		return domain.ErrAlreadySubmitted
	case models.SubmissionReviewed:
		return domain.ErrSubmissionLocked
	default:
		return domain.ErrInvalidRequest
	}

	if _, err := s.submissionRepo.UpdateDraft(ctx, submission.ID, input); err != nil {
		return fmt.Errorf("service: update submission draft: %w", err)
	}

	return nil
}

func (s *Service) ListTeacherSubmissions(ctx context.Context, teacherID uuid.UUID) ([]domain.TeacherSubmission, error) {
	items, err := s.submissionRepo.ListTeacherSubmissions(ctx, teacherID)
	if err != nil {
		return nil, fmt.Errorf("service: list teacher submissions: %w", err)
	}
	return items, nil
}

func (s *Service) SetGrade(ctx context.Context, input domain.GradeInput) error {
	if input.SubmissionID <= 0 || input.Grade < 0 || input.Grade > 100 {
		return domain.ErrInvalidRequest
	}
	nextStatus := models.SubmissionStatus(strings.TrimSpace(input.Status))
	if nextStatus != models.SubmissionRevision && nextStatus != models.SubmissionReviewed {
		return domain.ErrInvalidRequest
	}

	submission, err := s.submissionRepo.GetForTeacher(ctx, input.SubmissionID, input.TeacherID)
	if err != nil {
		if errors.Is(err, models.ErrSubmissionNotFound) {
			return domain.ErrSubmissionNotFound
		}
		return fmt.Errorf("service: validate teacher submission: %w", err)
	}
	switch submission.Status {
	case models.SubmissionSubmitted, models.SubmissionRevision:
	case models.SubmissionReviewed:
		return domain.ErrSubmissionLocked
	default:
		return domain.ErrInvalidRequest
	}

	if err := s.gradeRepo.Save(ctx, input.SubmissionID, input.TeacherID, input.Grade, strings.TrimSpace(input.Comment)); err != nil {
		return fmt.Errorf("service: save grade: %w", err)
	}
	if err := s.submissionRepo.SetStatus(ctx, input.SubmissionID, nextStatus); err != nil {
		return fmt.Errorf("service: update submission status: %w", err)
	}

	return nil
}

func (s *Service) GetStudentProgress(ctx context.Context, studentID uuid.UUID) (*domain.StudentProgress, error) {
	progress, err := s.reportRepo.GetStudentProgress(ctx, studentID)
	if err != nil {
		return nil, fmt.Errorf("service: get student progress: %w", err)
	}
	return progress, nil
}

func (s *Service) GetAdminStats(ctx context.Context) (*domain.AdminStats, error) {
	stats, err := s.reportRepo.GetAdminStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("service: get admin stats: %w", err)
	}
	return stats, nil
}

func (s *Service) ExportAdminReportCSV(ctx context.Context) ([]byte, error) {
	rows, err := s.reportRepo.ListReportRows(ctx)
	if err != nil {
		return nil, fmt.Errorf("service: list report rows: %w", err)
	}

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	if err := writer.Write([]string{
		"student",
		"group",
		"lab_work",
		"status",
		"attempt",
		"submitted_at",
		"grade",
		"comment",
	}); err != nil {
		return nil, fmt.Errorf("service: write csv header: %w", err)
	}

	for _, row := range rows {
		submittedAt := ""
		if row.SubmittedAt != nil {
			submittedAt = *row.SubmittedAt
		}
		grade := ""
		if row.Grade != nil {
			grade = fmt.Sprintf("%d", *row.Grade)
		}
		comment := ""
		if row.Comment != nil {
			comment = *row.Comment
		}

		if err := writer.Write([]string{
			row.StudentName,
			row.GroupName,
			row.LabWorkTitle,
			row.Status,
			fmt.Sprintf("%d", row.AttemptNumber),
			submittedAt,
			grade,
			comment,
		}); err != nil {
			return nil, fmt.Errorf("service: write csv row: %w", err)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, fmt.Errorf("service: flush csv: %w", err)
	}

	return buf.Bytes(), nil
}

func validateLabWork(input domain.LabWorkInput) error {
	if strings.TrimSpace(input.Title) == "" {
		return domain.ErrInvalidRequest
	}
	return nil
}

func mapLabWorks(items []models.LabWork) []domain.LabWork {
	result := make([]domain.LabWork, 0, len(items))
	for _, item := range items {
		result = append(result, mapLabWork(item))
	}
	return result
}

func mapLabWork(item models.LabWork) domain.LabWork {
	return domain.LabWork{
		ID:          item.ID,
		Title:       item.Title,
		Description: item.Description,
		Goal:        item.Goal,
		Equipment:   item.Equipment,
		Reagents:    item.Reagents,
		Procedure:   item.Procedure,
		FilePath:    item.FilePath,
		CreatedBy:   item.CreatedBy,
		CreatedAt:   item.CreatedAt.Format(time.RFC3339),
	}
}
