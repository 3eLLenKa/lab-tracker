package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	// common
	ErrInvalidRequest = errors.New("invalid request")
	ErrUnauthorized   = errors.New("unauthorized")
	ErrForbidden      = errors.New("forbidden")
	ErrNotFound       = errors.New("not found")
	ErrInternal       = errors.New("internal error")

	// users
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserAlreadyExists  = errors.New("user already exists")

	// groups
	ErrGroupNotFound = errors.New("group not found")
	ErrGroupExists   = errors.New("group already exists")

	// lab works
	ErrLabWorkNotFound = errors.New("lab work not found")

	// assignments
	ErrAssignmentNotFound = errors.New("assignment not found")

	// submissions
	ErrSubmissionNotFound = errors.New("submission not found")
	ErrAlreadySubmitted   = errors.New("already submitted")
	ErrDeadlinePassed     = errors.New("deadline passed")

	// grades
	ErrGradeAlreadyExists = errors.New("grade already exists")
)

type UserRole string

const (
	RoleStudent UserRole = "student"
	RoleTeacher UserRole = "teacher"
	RoleAdmin   UserRole = "admin"
)

type AssignmentStatus string

const (
	AssignmentAssigned   AssignmentStatus = "assigned"
	AssignmentInProgress AssignmentStatus = "in_progress"
	AssignmentCompleted  AssignmentStatus = "completed"
)

type SubmissionStatus string

const (
	SubmissionPending   SubmissionStatus = "pending"
	SubmissionSubmitted SubmissionStatus = "submitted"
	SubmissionChecked   SubmissionStatus = "checked"
)

type User struct {
	ID           uuid.UUID
	GroupID      *int64
	Username     string
	PasswordHash string
	FullName     string
	Role         UserRole
	CreatedAt    time.Time
	UpdatedAt    *time.Time
}

type Group struct {
	ID        int64
	Name      string
	Course    int
	TeacherID *uuid.UUID
}

type LabWork struct {
	ID          int64
	Title       string
	Description string
	Goal        string
	Equipment   string
	Reagents    string
	Procedure   string
	FilePath    *string
	CreatedBy   *uuid.UUID
	CreatedAt   time.Time
}

type Assignment struct {
	ID        int64
	LabWorkID int64
	GroupID   int64
	Deadline  *time.Time
	Status    AssignmentStatus
	CreatedAt time.Time
}

type Submission struct {
	ID           int64
	AssignmentID int64
	StudentID    uuid.UUID
	TextReport   string
	FilePath     *string
	Status       SubmissionStatus
	SubmittedAt  *time.Time
}

type Grade struct {
	ID           int64
	SubmissionID int64
	TeacherID    *uuid.UUID
	Grade        int
	Comment      string
	GradedAt     *time.Time
}
