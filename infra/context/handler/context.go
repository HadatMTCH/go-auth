package handler

import (
	"base-api/app/http/handlers/account"
	"base-api/app/http/handlers/health"
	"base-api/app/http/handlers/template"
	"base-api/infra/context/repository"
	"base-api/infra/context/service"
)

type HandlerContext struct {
	TemplateHandler template.Template
	HealthHandler   health.Health
	AccountHandler  account.Account
}

// initServiceCtx for contextService
func InitHandlerContext(serviceContext *service.ServiceContext, repo *repository.RepositoryContext) *HandlerContext {
	return &HandlerContext{
		TemplateHandler: template.New(serviceContext),
		HealthHandler:   health.New(),
		AccountHandler:  account.New(repo),
	}
}
