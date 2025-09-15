package server

import (
	"approval-system/config"
	"context"
	"log"

	"approval-system/pkg/database"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type (
	server struct {
		app *echo.Echo

		postgres database.DatabasesPostgres
		cfg      *config.Config
	}
)

func (s *server) gracefulShutdown(pctx context.Context, quit <-chan os.Signal) {
	log.Printf("Start service...%s", s.cfg.App.Name) // ✅ แก้ตรงนี้
	<-quit
	log.Printf("Shutdown service...%s", s.cfg.App.Name) // ✅ แก้ตรงนี้

	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	if err := s.app.Shutdown(ctx); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
func (s *server) httpListening() {
	if err := s.app.Start(s.cfg.App.Url); err != http.ErrServerClosed {
		log.Fatalf("Error: %v", err)
	}
}
func Start(pctx context.Context, cfg *config.Config, dbPost database.DatabasesPostgres) {
	s := &server{
		app:      echo.New(),
		postgres: dbPost,
		cfg:      cfg,
	}

	s.app.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Skipper:      middleware.DefaultSkipper,
		ErrorMessage: "Error: Request Timeout",
		Timeout:      30 * time.Second,
	}))

	s.app.Use(middleware.CORSWithConfig(
		middleware.CORSConfig{
			Skipper:      middleware.DefaultSkipper,
			AllowOrigins: []string{"*"},
			AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.PATCH, echo.DELETE},
		},
	))

	s.app.Use(middleware.BodyLimit("10M"))

	switch s.cfg.App.Name {
	case "doc":
		s.docService()
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	s.app.Use(middleware.Logger())

	go s.gracefulShutdown(pctx, quit)

	s.httpListening()
}
