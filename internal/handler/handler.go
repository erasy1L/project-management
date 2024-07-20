package handler

import (
	"project-management/internal/handler/httphandler"
	"project-management/internal/service/management"
	"project-management/pkg/router"

	_ "project-management/docs"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Dependencies struct {
	ManagementService *management.Service
}

type Handler struct {
	deps Dependencies

	HTTP *chi.Mux
}

type Configuration func(h *Handler) error

func New(deps Dependencies, cfgs ...Configuration) Handler {
	h := Handler{
		deps: deps,
	}

	for _, cfg := range cfgs {
		cfg(&h)
	}

	return h
}

// @title Project Management API
// @version 1
// @description This is a simple project management API
// @host localhost:8080
// @BasePath /api/v1
func WithHTTPHandler() Configuration {
	return func(h *Handler) error {
		h.HTTP = router.New()

		userHandler := httphandler.NewUserHandler(h.deps.ManagementService)
		taskHandler := httphandler.NewTaskHandler(h.deps.ManagementService)
		projecthandler := httphandler.NewProjectHandler(h.deps.ManagementService)

		h.HTTP.Get("/swagger/*", httpSwagger.WrapHandler)

		h.HTTP.Route("/api/v1", func(r chi.Router) {
			r.Mount("/users", userHandler.Routes())
			r.Mount("/tasks", taskHandler.Routes())
			r.Mount("/projects", projecthandler.Routes())
		})

		return nil
	}
}
