package template

import (
	"base-api/infra/context/service"
	"net/http"
)

type Template interface {
	RegistrationUser(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	Profile(w http.ResponseWriter, r *http.Request)
}

func New(serviceCtx *service.ServiceContext) Template {
	return &template{
		serviceCtx,
	}
}
