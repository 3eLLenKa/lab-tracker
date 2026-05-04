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
