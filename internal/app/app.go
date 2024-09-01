package app

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/avran02/kode/config"
	"github.com/avran02/kode/internal/controller"
	"github.com/avran02/kode/internal/repository"
	"github.com/avran02/kode/internal/router"
	"github.com/avran02/kode/internal/service"
	"github.com/avran02/kode/logger"
)

type App struct {
	router       router.Router
	config       *config.Config
	dbConnection *sql.DB
}

func (a *App) Run() error {
	serverEndpoint := fmt.Sprintf("%s:%s", a.config.Server.Host, a.config.Server.Port)
	slog.Info("Starting server at " + serverEndpoint)
	s := http.Server{
		Addr:    serverEndpoint,
		Handler: a.router,
	}

	s.RegisterOnShutdown(func() {
		if err := a.dbConnection.Close(); err != nil {
			slog.Error("can't close db conn: " + err.Error())
		}
	})

	return s.ListenAndServe()
}

func New() *App {
	config := config.New()
	logger.Setup(config.Server)

	dbConnection := repository.MustGetPostgresConnection(config.DB)
	authService := service.NewAuthService(repository.NewUserRepository(dbConnection), config.JWTSecret)
	notesService := service.NewNotesService(repository.NewNotesRepository(dbConnection))

	controller := controller.New(authService, notesService)
	router := router.New(controller)

	return &App{
		router:       router,
		config:       config,
		dbConnection: dbConnection,
	}
}
