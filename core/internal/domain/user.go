package domain

import "github.com/google/uuid"

type UserRole string

const (
	RoleStudent UserRole = "student"
	RoleTeacher UserRole = "teacher"
	RoleAdmin   UserRole = "admin"
)

type User struct {
	ID        uuid.UUID
	GroupID   *int64
	Username  string
	Password  string
	FullName  string
	Role      UserRole
	CreatedAt string
	UpdatedAt *string
}

type Group struct {
	ID     int64
	Name   string
	Course int
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
	CreatedAt   string
}

type LabWorkList struct {
	Items      []LabWork
	Total      int
	Page       int
	Limit      int
	TotalPages int
}

type LabWorkFilter struct {
	Search string
	Page   int
	Limit  int
}

type LabWorkInput struct {
	Title       string
	Description string
	Goal        string
	Equipment   string
	Reagents    string
	Procedure   string
	FilePath    *string
	CreatedBy   *uuid.UUID
}

type StudentAssignment struct {
	AssignmentID     int64
	LabWorkID        int64
	Title            string
	Description      string
	Deadline         *string
	AssignmentStatus string
	SubmissionID     *int64
	SubmissionStatus *string
	AttemptNumber    *int
	TextReport       *string
	FilePath         *string
	SubmittedAt      *string
	Grade            *int
	TeacherComment   *string
}

type SubmissionInput struct {
	AssignmentID int64
	StudentID    uuid.UUID
	TextReport   string
	FilePath     *string
}

type TeacherSubmission struct {
	SubmissionID   int64
	AssignmentID   int64
	StudentID      uuid.UUID
	StudentName    string
	GroupName      string
	LabWorkTitle   string
	Deadline       *string
	AttemptNumber  int
	TextReport     string
	FilePath       *string
	Status         string
	SubmittedAt    *string
	Grade          *int
	TeacherComment *string
}

type GradeInput struct {
	SubmissionID int64
	TeacherID    uuid.UUID
	Grade        int
	Comment      string
	Status       string
}

type StudentProgress struct {
	TotalAssignments int
	DraftCount       int
	SubmittedCount   int
	RevisionCount    int
	ReviewedCount    int
	CompletionRate   int
	AverageGrade     *float64
}

type AdminStats struct {
	UsersTotal       int
	StudentsTotal    int
	TeachersTotal    int
	GroupsTotal      int
	LabWorksTotal    int
	AssignmentsTotal int
	SubmissionsTotal int
	DraftCount       int
	SubmittedCount   int
	RevisionCount    int
	ReviewedCount    int
	AverageGrade     *float64
}

type ReportRow struct {
	StudentName   string
	GroupName     string
	LabWorkTitle  string
	Status        string
	AttemptNumber int
	SubmittedAt   *string
	Grade         *int
	Comment       *string
}
