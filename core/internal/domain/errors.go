package domain

import "errors"

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
	ErrSubmissionLocked   = errors.New("submission locked")

	// grades
	ErrGradeAlreadyExists = errors.New("grade already exists")
)
