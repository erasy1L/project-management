package httphandler

import (
	"encoding/json"
	"errors"
	"net/http"
	"project-management/internal/domain/user"
	"project-management/internal/service/management"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type UserHandler struct {
	managementService *management.Service
}

func NewUserHandler(s *management.Service) *UserHandler {
	return &UserHandler{managementService: s}
}

func (h *UserHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", h.list)
	r.Post("/", h.create)

	r.Get("/search", h.search)

	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", h.get)
		r.Put("/", h.update)
		r.Delete("/", h.delete)
		r.Get("/tasks", h.listTasks)
	})

	return r
}

// @Summary List users
// @Description List users
// @Tags users
// @Accept json
// @Success 200 {array} user.Response
// @Failure 400 {object} string
// @Router /users [get]
func (h *UserHandler) list(w http.ResponseWriter, r *http.Request) {
	users, err := h.managementService.ListUsers(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	render.JSON(w, r, users)
}

// @Summary Create a user
// @Description Create a user
// @Tags users
// @Accept json
// @Param body body user.Request true "User request"
// @Success 201 {string} string "User ID"
// @Failure 400 {object} []string "Validation errors"
// @Router /users [post]
func (h *UserHandler) create(w http.ResponseWriter, r *http.Request) {
	req := user.Request{}
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

	id, err := h.managementService.CreateUser(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	render.PlainText(w, r, id)
}

// @Summary Get a user
// @Description Get a user
// @Tags users
// @Accept json
// @Param id path string true "User ID"
// @Success 200 {object} user.Response
// @Failure 400 {string} string "Bad request"
// @Failure 404 {string} string "User not found"
// @Router /users/{id} [get]
func (h *UserHandler) get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	data, err := h.managementService.GetUser(r.Context(), id)
	if err != nil {
		if errors.Is(err, user.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
		}

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	render.JSON(w, r, data)
}

// @Summary Update a user
// @Description Update a user
// @Tags users
// @Accept json
// @Param id path string true "User ID"
// @Param body body user.UpdateRequest true "User request"
// @Success 200 {string} string "User ID"
// @Failure 400 {object} []string "Validation errors"
// @Failure 404 {string} string "User not found"
// @Router /users/{id} [put]
func (h *UserHandler) update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	req := user.UpdateRequest{}
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

	err := h.managementService.UpdateUser(r.Context(), id, req)
	if err != nil {
		if errors.Is(err, user.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
		}

		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

// @Summary Delete a user
// @Description Delete a user
// @Tags users
// @Accept json
// @Param id path string true "User ID"
// @Success 200 {string} string "User deleted"
// @Failure 400 {string} string "Bad request"
// @Failure 404 {string} string "User not found"
// @Router /users/{id} [delete]
func (h *UserHandler) delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	err := h.managementService.DeleteUser(r.Context(), id)
	if err != nil {
		if errors.Is(err, user.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

// @Summary List user tasks
// @Description List user tasks
// @Tags users
// @Accept json
// @Param id path string true "User ID"
// @Success 200 {array} task.Response
// @Failure 400 {string} string "Bad request"
// @Failure 404 {string} string "User not found"
// @Router /users/{id}/tasks [get]
func (h *UserHandler) listTasks(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	tasks, err := h.managementService.SearchTasks(r.Context(), "user_id", id)
	if err != nil {
		if errors.Is(err, user.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
		}

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	render.JSON(w, r, tasks)
}

// @Summary Search users
// @Description Search users
// @Tags users
// @Accept json
// @Param query query string true "Query"
// @Param value query string true "Value"
// @Success 200 {array} user.Response
// @Failure 400 {string} string "Bad request"
// @Failure 404 {string} string "Not found"
// @Router /users/search [get]
func (h *UserHandler) search(w http.ResponseWriter, r *http.Request) {
	filter, val := "", ""

	for k, v := range r.URL.Query() {
		filter, val = k, v[0]
	}

	users, err := h.managementService.SearchUsers(r.Context(), filter, val)
	if err != nil {
		if errors.Is(err, user.ErrSearch) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if errors.Is(err, user.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, users)
}
