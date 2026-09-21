package http

import (
	"net/http"

	"github.com/ednanf/school-api/internal/domain"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type DepartmentHandler struct {
	repo     domain.DepartmentRepository
	validate *validator.Validate
}

func NewDepartmentHandler(repo domain.DepartmentRepository, validate *validator.Validate) *DepartmentHandler {
	return &DepartmentHandler{repo: repo, validate: validate}
}

func (h *DepartmentHandler) DepartmentRoutes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", h.HandleList)
	r.Post("/", h.HandleCreate)
	r.Delete("/{id}", h.HandleDelete)
	r.Get("/{id}", h.HandleGetById)
	r.Patch("/{id}", h.HandleUpdate)

	return r
}

func (h *DepartmentHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	sendSuccess(w, http.StatusOK, "create hit", nil)
}

func (h *DepartmentHandler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	sendSuccess(w, http.StatusOK, "delete hit", nil)
}

func (h *DepartmentHandler) HandleGetById(w http.ResponseWriter, r *http.Request) {
	sendSuccess(w, http.StatusOK, "get by id hit", nil)
}

func (h *DepartmentHandler) HandleList(w http.ResponseWriter, r *http.Request) {
	sendSuccess(w, http.StatusOK, "list hit", nil)
}

func (h *DepartmentHandler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	sendSuccess(w, http.StatusOK, "update hit", nil)
}
