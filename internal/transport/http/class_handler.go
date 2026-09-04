package http

import (
	"net/http"

	"github.com/ednanf/school-api/internal/domain"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

// Class handler holds routes
type ClassHandler struct {
	repo     domain.ClassRepository
	validate *validator.Validate
}

// Constructor
func NewClassHandler(repo domain.ClassRepository, validate *validator.Validate) *ClassHandler {
	return &ClassHandler{repo: repo, validate: validate}
}

// Route paths
func (h *ClassHandler) ClassRoutes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", h.HandleList)
	r.Post("/", h.HandleCreate)
	r.Delete("/{id}", h.HandleDelete)
	r.Get("/{id}", h.HandleGetById)
	r.Patch("/{id}", h.HandlePatch)

	return r
}

func (h *ClassHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	sendSuccess(w, http.StatusOK, "Class create hit", nil)
}

func (h *ClassHandler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	sendSuccess(w, http.StatusOK, "Class delete hit", id)
}

func (h *ClassHandler) HandleGetById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	sendSuccess(w, http.StatusOK, "GetById hit", id)
}

func (h *ClassHandler) HandleList(w http.ResponseWriter, r *http.Request) {
	sendSuccess(w, http.StatusOK, "List hit", nil)
}

func (h *ClassHandler) HandlePatch(w http.ResponseWriter, r *http.Request) {
	sendSuccess(w, http.StatusOK, "Patch hit", nil)
}
