package http

import (
	"net/http"

	"github.com/ednanf/school-api/internal/domain"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

// TeacherAssignmentHandler contains `repo` with a way to communicate with the database and the pointer to the validator instantiated once in `main.go`
type TeacherAssignmentHandler struct {
	repo     domain.TeacherAssignmentRepository
	validate *validator.Validate
}

// NewTeacherAssignmentHandler is a constructor that returns a pointer to a TeacherAssignmentHandler struct, initializing it with the injected repository and validator dependencies
func NewTeacherAssignmentHandler(repo domain.TeacherAssignmentRepository, validate *validator.Validate) *TeacherAssignmentHandler {
	return &TeacherAssignmentHandler{repo: repo, validate: validate}
}

// Route paths
func (h *TeacherAssignmentHandler) TeacherAssignmentRoutes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", h.HandleList)
	r.Post("/", h.HandleCreate)
	r.Delete("/{id}", h.HandleDelete)
	r.Get("/{id}", h.HandleGetById)
	r.Patch("/{id}", h.HandleUpdate)

	return r
}

func (h *TeacherAssignmentHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	sendSuccess(w, http.StatusOK, "create hit", nil)
}

func (h *TeacherAssignmentHandler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	sendSuccess(w, http.StatusOK, "delete hit", nil)
}

func (h *TeacherAssignmentHandler) HandleGetById(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	sendSuccess(w, http.StatusOK, "get by id hit", idStr)
}

func (h *TeacherAssignmentHandler) HandleList(w http.ResponseWriter, r *http.Request) {
	sendSuccess(w, http.StatusOK, "list hit", nil)
}

func (h *TeacherAssignmentHandler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	sendSuccess(w, http.StatusOK, "update hit", nil)
}
