package repository

import (
	"database/sql"
	pggroup "lab-tracker/internal/repository/postgres/group"
	pglabwork "lab-tracker/internal/repository/postgres/labwork"
	pguser "lab-tracker/internal/repository/postgres/user"
)

type Repo struct {
	UserRepo    *pguser.UserRepo
	GroupRepo   *pggroup.GroupRepo
	LabWorkRepo *pglabwork.LabWorkRepo
}

func New(db *sql.DB) *Repo {
	return &Repo{
		UserRepo:    pguser.New(db),
		GroupRepo:   pggroup.New(db),
		LabWorkRepo: pglabwork.New(db),
	}
}
