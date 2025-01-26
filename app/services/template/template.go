package template

import (
	"base-api/data/models"
	"base-api/infra/context/repository"
	"base-api/utils"
	"context"
	"errors"
)

type template struct {
	*repository.RepositoryContext
}

func (u *template) InsertUser(ctx context.Context, data models.User) (err error) {
	data.Password, err = utils.HashPassword(data.Password)
	if err != nil {
		return err
	}
	return u.TemplateRepository.InsertUser(ctx, data)
}

func (u *template) Login(ctx context.Context, email, password string) (*models.UserResponse, error) {
	data, err := u.TemplateRepository.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, errors.New("user not found")
	}

	if !utils.CheckPasswordHash(password, data.Password) {
		return nil, err
	}

	if data.IsSuspended {
		return nil, errors.New("account is suspsended, please contact customer service.")
	}

	userResponse := models.NewUserResponse()
	userResponseData := userResponse.ToResponse(data)

	return userResponseData, nil
}

func (u *template) GetUserByID(ctx context.Context, id int) (*models.UserResponse, error) {
	data, err := u.TemplateRepository.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, errors.New("user not found")
	}

	userResponse := models.NewUserResponse()
	userResponse.ToResponse(data)
	return userResponse, nil
}
