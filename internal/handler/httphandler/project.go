package httphandler

import (
	"encoding/json"
	"errors"
	"net/http"
	"project-management/internal/domain/project"
	"project-management/internal/service/management"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type ProjectHandler struct {
	managementService *management.Service
}

func NewProjectHandler(managementService *management.Service) *ProjectHandler {
	return &ProjectHandler{
		managementService: managementService,
	}
}

func (h *ProjectHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/", h.create)
	r.Get("/", h.list)

	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", h.get)
		r.Put("/", h.update)
		r.Delete("/", h.delete)
		r.Get("/tasks", h.listTasks)
	})

	r.Get("/search", h.search)

	return r
}

// @Summary Create a project
// @Description Create a project
// @Tags projects
// @Accept json
// @Param body body project.Request true "Project request"
// @Success 201 {string} string "Project ID"
// @Failure 400 {object} []string "Validation errors"
// @Router /projects [post]
func (h *ProjectHandler) create(w http.ResponseWriter, r *http.Request) {
	req := project.Request{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if errs := req.Validate(); errs != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errs)
		return
	}

	id, err := h.managementService.CreateProject(r.Context(), req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	render.PlainText(w, r, id)
}

// @Summary Get a project
// @Description Get a project
// @Tags projects
// @Param id path string true "Project ID"
// @Success 200 {object} project.Response
// @Failure 400 {string} string "Bad request"
// @Router /projects/{id} [get]
func (h *ProjectHandler) get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	project, err := h.managementService.GetProject(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	render.JSON(w, r, project)
}

// @Summary List projects
// @Description List projects
// @Tags projects
// @Success 200 {array} project.Response
// @Failure 400 {string} string "Bad request"
// @Router /projects [get]
func (h *ProjectHandler) list(w http.ResponseWriter, r *http.Request) {
	projects, err := h.managementService.ListProjects(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	render.JSON(w, r, projects)
}

// @Summary Update a project
// @Description Update a project
// @Tags projects
// @Accept json
// @Param id path string true "Project ID"
// @Param body body project.UpdateRequest true "Project update request"
// @Success 200 {string} string "Project updated"
// @Failure 400 {object} []string "Validation errors"
// @Router /projects/{id} [put]
func (h *ProjectHandler) update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	req := project.UpdateRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if errs := req.Validate(); errs != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errs)
		return
	}

	err := h.managementService.UpdateProject(r.Context(), id, req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// @Summary Delete a project
// @Description Delete a project
// @Tags projects
// @Param id path string true "Project ID"
// @Success 200 {string} string "Project deleted"
// @Failure 400 {string} string "Bad request"
// @Router /projects/{id} [delete]
func (h *ProjectHandler) delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	err := h.managementService.DeleteProject(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// @Summary Search projects
// @Description Search projects
// @Tags projects
// @Param query query string true "Query"
// @Param val query string true "Value"
// @Success 200 {array} project.Response
// @Failure 400 {string} string "Bad request"
// @Failure 404 {string} string "Not found"
// @Router /projects/search [get]
func (h *ProjectHandler) search(w http.ResponseWriter, r *http.Request) {
	var filter, val string

	for k, v := range r.URL.Query() {
		filter, val = k, v[0]
	}

	projects, err := h.managementService.SearchProjects(r.Context(), filter, val)
	if err != nil {
		if errors.Is(err, project.ErrSearch) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if errors.Is(err, project.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, projects)
}

// @Summary List project tasks
// @Description List project tasks
// @Tags projects
// @Param id path string true "Project ID"
// @Success 200 {array} task.Response
// @Failure 400 {string} string "Bad request"
// @Router /projects/{id}/tasks [get]
func (h *ProjectHandler) listTasks(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	tasks, err := h.managementService.SearchTasks(r.Context(), "project_id", id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	render.JSON(w, r, tasks)
}
