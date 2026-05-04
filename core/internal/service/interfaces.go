package service

import (
	"context"
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
