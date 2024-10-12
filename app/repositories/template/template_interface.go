package template

import (
	"base-api/data/models"
	"base-api/infra/db"
	"context"
)

type Template interface {
	InsertUser(ctx context.Context, data models.User) error
	GetUserByID(ctx context.Context, ID int) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
}

func New(db *db.DB) Template {
	return &template{
		db,
	}
}
