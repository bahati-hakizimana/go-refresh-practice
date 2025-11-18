package aprtment

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-refresh-practice/go-refresh-course/middleware"
	"github.com/go-refresh-practice/go-refresh-course/types"
	"github.com/go-refresh-practice/go-refresh-course/utils"
	"github.com/gorilla/mux"
)


type Handler struct {

	store types.ApartmentStore
}

func NewHandler(store types.ApartmentStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	// GET /apartments 
	router.Handle("/apartments",
		middleware.AuthMiddleware(http.HandlerFunc(h.handleGetApartments)),
	).Methods(http.MethodGet)

	// POST /apartments (admin only)
	router.Handle("/apartments",
		middleware.AuthMiddleware(middleware.AdminOnly(http.HandlerFunc(h.handlerCreateApartment))),
	).Methods(http.MethodPost)
}

func(h *Handler) handleGetApartments(w http.ResponseWriter, r *http.Request){

	ap, err := h.store.GetApartments()
  if err != nil {
	utils.WriteError(w, http.StatusInternalServerError, err)
	return
  }

  utils.WriteJson(w, http.StatusOK, ap)
}

func (h *Handler) handlerCreateApartment(w http.ResponseWriter, r *http.Request) {
	var apt types.Apartment

	// Parse JSON body into apt
	if err := json.NewDecoder(r.Body).Decode(&apt); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validation
	if apt.Name == "" || apt.Code == "" {
		utils.WriteError(w, http.StatusBadRequest, errInvalidApartment())
		return
	}

	created, err := h.store.CreateApartment(apt)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJson(w, http.StatusCreated, created)
}

func errInvalidApartment() error {
	return fmt.Errorf("name and code are required")
}