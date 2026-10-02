package http

import (
	"net/http"

	"github.com/ednanf/school-api/internal/domain"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type staffHandler struct {
	service  domain.StaffService
	validate *validator.Validate
}

func NewStaffHandler(service domain.StaffService, validate *validator.Validate) *staffHandler {
	return &staffHandler{service: service, validate: validate}
}

func (h *staffHandler) StaffRoutes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", h.HandleList)
	r.Post("/", h.HandleCreate)
	r.Delete("/{id}", h.HandleDelete)
	r.Get("/{id}", h.HandleGetByID)
	r.Patch("/{id}", h.HandleUpdate)

	return r
}

func (h *staffHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	sendSuccess(w, http.StatusOK, "create hit", nil)
}

func (h *staffHandler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	sendSuccess(w, http.StatusOK, "delete hit", nil)
}

func (h *staffHandler) HandleGetByID(w http.ResponseWriter, r *http.Request) {
	sendSuccess(w, http.StatusOK, "get by id hit", nil)
}

func (h *staffHandler) HandleList(w http.ResponseWriter, r *http.Request) {
	sendSuccess(w, http.StatusOK, "list hit", nil)
}

func (h *staffHandler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	sendSuccess(w, http.StatusOK, "update hit", nil)
}
