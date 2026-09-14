package http

import (
	"net/http"

	"github.com/ednanf/school-api/internal/domain"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

// TeacherHandler contains `repo` with a way to communicate with the database and the pointer to the validator instantiated in `main.go`

type TeacherHandler struct {
	repo     domain.TeacherRepository
	validate *validator.Validate
}

// NewTeacherHandler is a constructor that returns a pointer to a TeacherHandler struct, initializing it with the injected repository and validator dependencies
func NewTeacherHandler(repo domain.TeacherRepository, validate *validator.Validate) *TeacherHandler {
	return &TeacherHandler{repo: repo, validate: validate}
}

// Route paths
func (h *TeacherHandler) TeacherRoutes() chi.Router {
	r := chi.NewRouter()

	r.Post("/", h.HandleCreate)
	r.Get("/", h.HandleList)
	r.Delete("/{id}", h.HandleDelete)
	r.Get("/{id}", h.HandleGetById)
	r.Patch("/{id}", h.HandleUpdate)

	return r
}

func (h *TeacherHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	sendSuccess(w, http.StatusOK, "Create hit", nil)
}

func (h *TeacherHandler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	sendSuccess(w, http.StatusOK, "Delete hit", nil)
}

func (h *TeacherHandler) HandleGetById(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	sendSuccess(w, http.StatusOK, "Get by id hit", idStr)
}

func (h *TeacherHandler) HandleList(w http.ResponseWriter, r *http.Request) {
	sendSuccess(w, http.StatusOK, "List hit", nil)
}

func (h *TeacherHandler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	sendSuccess(w, http.StatusOK, "Update hit", nil)
}
