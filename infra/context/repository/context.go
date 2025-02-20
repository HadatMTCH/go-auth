package repository

import (
	"base-api/app/repositories/account"
	"base-api/app/repositories/template"
	"base-api/config"
	"base-api/infra/db"
)

type RepositoryContext struct {
	TemplateRepository template.Template
	AccountRepository  account.Account
	DB                 *db.DB
}

func InitializeRepositoryContext(db *db.DB, config *config.S3Configuration) *RepositoryContext {
	templateRepository := template.New(db)
	accountRepository := account.New(db)

	return &RepositoryContext{
		DB:                 db,
		TemplateRepository: templateRepository,
		AccountRepository:  accountRepository,
	}
}
