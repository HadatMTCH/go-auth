package http_server

import (
	"base-api/app/http/routers"
	infra "base-api/infra/context"
	"context"
	"errors"
	"flag"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/spf13/cobra"
)

type httpServer struct {
	infraContext infra.InfraContextInterface
	e            *echo.Echo
}

type HTTPServer interface {
	RunHTTP(cmd *cobra.Command, args []string) error
}

func New() HTTPServer {
	return &httpServer{
		e: echo.New(),
	}
}

func (h httpServer) initializeRoutes() {
	h.e.Use(middleware.Logger())
	h.e.Use(middleware.Recover())
	routers.InitialRouter(h.infraContext, h.e)
}

func (h httpServer) SetGracefulTimeout() time.Duration {
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*time.Duration(h.infraContext.Config().Server.GraceFulTimeout), "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()
	return wait
}

func (h httpServer) RunHTTP(cmd *cobra.Command, args []string) error {
	h.infraContext = infra.New()

	// Echo configuration
	h.e.HideBanner = true
	h.e.HidePort = true

	// Initialize routes
	h.initializeRoutes()

	// Start server
	go func() {
		if err := h.e.Start(h.infraContext.Config().Server.Addr); err != nil && !errors.Is(err, http.ErrServerClosed) {
			h.e.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), h.SetGracefulTimeout())
	defer cancel()

	if err := h.e.Shutdown(ctx); err != nil {
		h.e.Logger.Fatal(err)
	}

	return nil
}
