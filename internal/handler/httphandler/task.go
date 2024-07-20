package httphandler

import (
	"encoding/json"
	"errors"
	"net/http"
	"project-management/internal/domain/task"
	"project-management/internal/service/management"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type TaskHandler struct {
	managementService *management.Service
}

func NewTaskHandler(managementService *management.Service) *TaskHandler {
	return &TaskHandler{
		managementService: managementService,
	}
}

func (h *TaskHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/", h.create)
	r.Get("/", h.list)

	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", h.get)
		r.Put("/", h.update)
		r.Delete("/", h.delete)
	})

	r.Get("/search", h.search)

	return r
}

// @Summary Create a task
// @Description Create a task
// @Tags tasks
// @Accept json
// @Param body body task.Request true "Task request"
// @Success 201 {string} string "Task ID"
// @Failure 400 {object} []string "Validation errors"
// @Router /tasks [post]
func (h *TaskHandler) create(w http.ResponseWriter, r *http.Request) {
	req := task.Request{}
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

	id, err := h.managementService.CreateTask(r.Context(), req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	render.PlainText(w, r, id)
}

// @Summary Get a task
// @Description Get a task
// @Tags tasks
// @Accept json
// @Param id path string true "Task ID"
// @Success 200 {object} task.Response
// @Failure 400 {string} string "Bad request"
// @Router /tasks/{id} [get]
func (h *TaskHandler) get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	task, err := h.managementService.GetTask(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	render.JSON(w, r, task)
}

// @Summary List tasks
// @Description List tasks
// @Tags tasks
// @Success 200 {object} []task.Response
// @Failure 400 {string} string "Bad request"
// @Router /tasks [get]
func (h *TaskHandler) list(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.managementService.ListTasks(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	render.JSON(w, r, tasks)
}

// @Summary Update a task
// @Description Update a task
// @Tags tasks
// @Accept json
// @Param id path string true "Task ID"
// @Param body body task.UpdateRequest true "Task update request"
// @Success 200 {string} string "Task updated"
// @Failure 400 {string} string "Bad request"
// @Router /tasks/{id} [put]
func (h *TaskHandler) update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	req := task.UpdateRequest{}
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

	err := h.managementService.UpdateTask(r.Context(), id, req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// @Summary Delete a task
// @Description Delete a task
// @Tags tasks
// @Param id path string true "Task ID"
// @Success 200 {string} string "Task deleted"
// @Failure 400 {string} string "Bad request"
// @Router /tasks/{id} [delete]
func (h *TaskHandler) delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	err := h.managementService.DeleteTask(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// @Summary Search tasks
// @Description Search tasks
// @Tags tasks
// @Param query query string true "Query"
// @Param value query string true "Value"
// @Success 200 {object} []task.Response
// @Failure 400 {string} string "Bad request"
// @Failure 404 {string} string "Not found"
// @Router /tasks/search [get]
func (h *TaskHandler) search(w http.ResponseWriter, r *http.Request) {
	filter, val := "", ""

	for k, v := range r.URL.Query() {
		filter, val = k, v[0]
	}

	tasks, err := h.managementService.SearchTasks(r.Context(), filter, val)
	if err != nil {
		if errors.Is(err, task.ErrSearch) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if errors.Is(err, task.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, tasks)
}
