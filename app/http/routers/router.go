package routers

import (
	infra "base-api/infra/context"

	"github.com/gorilla/mux"
)

const (
	GET   = "GET"
	POST  = "POST"
	PUT   = "PUT"
	PATCH = "PATCH"
)

// InitialRouter for object routers
func InitialRouter(infra infra.InfraContextInterface, r *mux.Router) *mux.Router {
	s := r.PathPrefix("/api").Subrouter()
	auth := s.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/register", infra.Handler().TemplateHandler.RegistrationUser).Methods(POST)
	auth.HandleFunc("/login", infra.Handler().TemplateHandler.Login).Methods(POST)

	profile := s.PathPrefix("/profile").Subrouter()
	profile.Use(infra.Middleware().TokenMiddleware.TokenAuthorize)
	profile.HandleFunc("/", infra.Handler().TemplateHandler.Profile).Methods(GET)

	// router
	// s.HandleFunc("/sample", sampleHandler.SampleMethod).Methods(GET)

	// SubRouter
	// c := r.PathPrefix("/sample-prefix").Subrouter()
	// c.HandleFunc("/sample-route", sampleHandler.SampleMethod).Methods(POST)
	return r
}
