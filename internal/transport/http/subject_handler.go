package http

import (
	"net/http"

	"github.com/ednanf/school-api/internal/domain"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

// SubjectHandler contains `repo` with a way to communicate with the database and the pointer to the validator instantiated once in `main.go`
type SubjectHandler struct {
	repo     domain.SubjectRepository
	validate *validator.Validate
}

// NewSubjectHandler isa constructor that returns a pointer to a SubjectHandler struct, initializing it with the injected repository and validator dependencies
func NewSubjectHandler(repo domain.SubjectRepository, validate *validator.Validate) *SubjectHandler {
	return &SubjectHandler{repo: repo, validate: validate}
}

// Route paths
func (h *SubjectHandler) SubjectRoutes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", h.HandleList)
	r.Post("/", h.HandleCreate)
	r.Delete("/{id}", h.HandleDelete)
	r.Get("/{id}", h.HandleGetById)
	r.Patch("/{id}", h.HandlePatch)

	return r
}

func (h *SubjectHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	sendSuccess(w, http.StatusOK, "Create hit", nil)
}

func (h *SubjectHandler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	sendSuccess(w, http.StatusOK, "Delete hit", nil)
}

func (h *SubjectHandler) HandleGetById(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	sendSuccess(w, http.StatusOK, "Get by id hit", idStr)
}

func (h *SubjectHandler) HandleList(w http.ResponseWriter, r *http.Request) {
	sendSuccess(w, http.StatusOK, "List hit", nil)
}

func (h *SubjectHandler) HandlePatch(w http.ResponseWriter, r *http.Request) {
	sendSuccess(w, http.StatusOK, "Patch hit", nil)
}
