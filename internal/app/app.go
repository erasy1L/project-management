package app

import (
	"context"
	"os"
	"os/signal"
	"project-management/config"
	"project-management/internal/handler"
	"project-management/internal/repository"
	"project-management/internal/repository/postgres"
	"project-management/internal/service/management"
	"project-management/pkg/log"
	"project-management/pkg/server"
	"syscall"
	"time"
)

func Run() {
	logger := log.LoggerFromContext(context.Background())

	configs, err := config.New()
	if err != nil {
		logger.Err(err).Stack().Msg("failed to load configurations")
		return
	}

	db, err := postgres.New(configs.DB)
	if err != nil {
		logger.Err(err).Stack().Msg("failed to connect to database")
		return
	}
	defer db.Close()

	repositories, err := repository.New(repository.WithPostgresStore(configs.DB))
	if err != nil {
		logger.Err(err).Stack().Msg("failed to create repositories")
		return
	}

	managementService := management.New(
		management.WithProjectRepository(repositories.Project),
		management.WithTaskRepository(repositories.Task),
		management.WithUserRepository(repositories.User),
	)

	handler := handler.New(
		handler.Dependencies{
			ManagementService: managementService,
		},
		handler.WithHTTPHandler())

	server, err := server.New(server.WithHTTPServer(handler.HTTP, configs.APP.Port))
	if err != nil {
		logger.Err(err).Stack().Msg("failed to create server")
		return
	}

	if err := server.Start(); err != nil {
		logger.Err(err).Stack().Msg("failed to start server")
		return
	}

	logger.Info().Msgf("server is running on port %s, swagger is at /swagger/index.html", configs.APP.Port)

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	<-shutdown
	logger.Info().Msg("shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := server.Stop(ctx); err != nil {
		logger.Err(err).Stack().Msg("failed to stop server")
		return
	}

	logger.Info().Msg("server stopped")
}
