package template

import (
	"base-api/data/models"
	"base-api/infra/db"
	"context"
	"database/sql"
	"errors"
	"time"
)

type template struct {
	*db.DB
}

func (u *template) InsertUser(ctx context.Context, data models.User) error {
	query := `
		INSERT INTO users (username, email, password, full_name, is_admin, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`
	err := u.DB.QueryRowContext(
		ctx,
		query,
		data.Username,
		data.Email,
		data.Password,
		data.FullName,
		data.IsAdmin,
		time.Now(),
		time.Now(),
	).Scan(&data.ID)

	return err
}

func (u *template) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	query := `SELECT * FROM users WHERE email = $1 LIMIT 1`

	err := u.DB.Master().GetContext(ctx, &user, query, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

func (u *template) GetUserByID(ctx context.Context, ID int) (*models.User, error) {
	var user models.User
	query := `SELECT * FROM users WHERE id = $1 LIMIT 1`

	err := u.DB.Master().GetContext(ctx, &user, query, ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}
