package routers

import (
	infra "base-api/infra/context"
	"base-api/infra/middleware"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func InitialRouter(infra infra.InfraContextInterface, e *echo.Echo) *echo.Echo {
	e.Use(middleware.RequestLogger())
	e.GET("/health", infra.Handler().HealthHandler.Check)

	// Swagger
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// API group
	api := e.Group("/api")

	// Auth routes
	auth := api.Group("/auth")
	auth.POST("/register", infra.Handler().TemplateHandler.RegistrationUser)
	auth.POST("/login", infra.Handler().TemplateHandler.Login)

	// Profile routes (with middleware)
	profile := api.Group("/profile")
	profile.Use(infra.Middleware().TokenMiddleware.TokenAuthorize())
	profile.GET("", infra.Handler().TemplateHandler.Profile)

	return e
}
