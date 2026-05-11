package repository

import (
	"database/sql"
	pgassignment "lab-tracker/internal/repository/postgres/assignment"
	pggrade "lab-tracker/internal/repository/postgres/grade"
	pggroup "lab-tracker/internal/repository/postgres/group"
	pglabwork "lab-tracker/internal/repository/postgres/labwork"
	pgsubmission "lab-tracker/internal/repository/postgres/submission"
	pguser "lab-tracker/internal/repository/postgres/user"
)

type Repo struct {
	UserRepo       *pguser.UserRepo
	GroupRepo      *pggroup.GroupRepo
	LabWorkRepo    *pglabwork.LabWorkRepo
	AssignmentRepo *pgassignment.AssignmentRepo
	SubmissionRepo *pgsubmission.SubmissionRepo
	GradeRepo      *pggrade.GradeRepo
}

func New(db *sql.DB) *Repo {
	return &Repo{
		UserRepo:       pguser.New(db),
		GroupRepo:      pggroup.New(db),
		LabWorkRepo:    pglabwork.New(db),
		AssignmentRepo: pgassignment.New(db),
		SubmissionRepo: pgsubmission.New(db),
		GradeRepo:      pggrade.New(db),
	}
}
