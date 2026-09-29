package http

import (
	"net/http"

	"github.com/ednanf/school-api/internal/domain"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type StaffPositionHandler struct {
	service  domain.StaffPositionService
	validate *validator.Validate
}

func NewStaffPositionHandler(service domain.StaffPositionService, validate *validator.Validate) *StaffPositionHandler {
	return &StaffPositionHandler{service: service, validate: validate}
}

func (h *StaffPositionHandler) StaffPositionRoutes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", h.HandleList)
	r.Post("/", h.HandleCreate)
	r.Delete("/{id}", h.HandleDelete)
	r.Get("/{id}", h.HandleGetByID)
	r.Patch("/{id}", h.HandleUpdate)

	return r
}

func (h *StaffPositionHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	sendSuccess(w, http.StatusOK, "create hit", nil)
}

func (h *StaffPositionHandler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	sendSuccess(w, http.StatusOK, "delete hit", nil)
}

func (h *StaffPositionHandler) HandleGetByID(w http.ResponseWriter, r *http.Request) {
	sendSuccess(w, http.StatusOK, "get by id hit", nil)
}

func (h *StaffPositionHandler) HandleList(w http.ResponseWriter, r *http.Request) {
	sendSuccess(w, http.StatusOK, "list hit", nil)
}

func (h *StaffPositionHandler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	sendSuccess(w, http.StatusOK, "update hit", nil)
}
