package template

import (
	"base-api/data/models"
	"base-api/infra/context/repository"
	"context"
)

type Template interface {
	InsertUser(ctx context.Context, data models.User) error
	GetUserByID(ctx context.Context, ID int) (*models.UserResponse, error)
	Login(ctx context.Context, email, password string) (*models.UserResponse, error)
}

func New(ctx *repository.RepositoryContext) Template {
	return &template{
		ctx,
	}
}
