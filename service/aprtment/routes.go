package aprtment

import (
	// "encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
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
	var payload types.CreateApartmentPayload
	if err := utils.PulseJson(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// check if apartment with the specific code not exist already
	_, err := h.store.GetApartmentByCode(payload.Code)
	if err == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("Apartment with this code %s already exist", payload.Code))
		return
	}

	// Validate payload
	if err := utils.Validate.Struct(payload); err != nil {
		error := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("Envalid payload %v", error))
		return
	}


	// Create apartment

	aptment := types.Apartment{
		Name:        payload.Name,
		Code:       payload.Code,
		Rooms:      payload.Rooms,
		Description: payload.Description,
		Price:      payload.Price,
	}

	_,err = h.store.CreateApartment(aptment)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	response := map[string]interface{}{
		"message" : "Apartment created successfully",
		"name": aptment.Name,
		"code": aptment.Code,
		"rooms": aptment.Rooms,
		"description": aptment.Description,
		"price": aptment.Price,
		"status": aptment.Status,
		"createdAt": aptment.CreatedAt,
	}

	utils.WriteJson(w, http.StatusCreated, response)

}
