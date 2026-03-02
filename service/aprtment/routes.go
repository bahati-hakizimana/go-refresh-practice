

package aprtment

import (
	"fmt"
	"net/http"
	"strconv"

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
	// Public
	router.HandleFunc("/apartments/public", h.handleGetApartmentsPublic).Methods(http.MethodGet)

	// Protected - any authenticated user
	router.Handle("/apartments",
		middleware.AuthMiddleware(http.HandlerFunc(h.handleGetApartments)),
	).Methods(http.MethodGet)

	router.Handle("/apartments/{id}",
		middleware.AuthMiddleware(http.HandlerFunc(h.handleGetApartmentByID)),
	).Methods(http.MethodGet)

	// Admin only
	router.Handle("/apartments",
		middleware.AuthMiddleware(middleware.AdminOnly(http.HandlerFunc(h.handleCreateApartment))),
	).Methods(http.MethodPost)

	router.Handle("/apartments/{id}",
		middleware.AuthMiddleware(middleware.AdminOnly(http.HandlerFunc(h.handleUpdateApartment))),
	).Methods(http.MethodPatch)

	router.Handle("/apartments/{id}",
		middleware.AuthMiddleware(middleware.AdminOnly(http.HandlerFunc(h.handleDeleteApartment))),
	).Methods(http.MethodDelete)
}

// ─── Handlers ─────────────────────────────────────────────────────────────────

func (h *Handler) handleGetApartmentsPublic(w http.ResponseWriter, r *http.Request) {
	apartments, err := h.store.GetPublicApartments()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJson(w, http.StatusOK, apartments)
}

func (h *Handler) handleGetApartments(w http.ResponseWriter, r *http.Request) {
	apartments, err := h.store.GetApartments()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJson(w, http.StatusOK, apartments)
}

func (h *Handler) handleGetApartmentByID(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromVars(r)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	apartment, err := h.store.GetApartmentByID(id)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("apartment not found"))
		return
	}

	utils.WriteJson(w, http.StatusOK, types.ApartmentResponse{
		Message:   "Apartment fetched successfully",
		Apartment: *apartment,
	})
}

func (h *Handler) handleCreateApartment(w http.ResponseWriter, r *http.Request) {
	var payload types.CreateApartmentPayload
	if err := utils.PulseJson(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest,
			fmt.Errorf("invalid payload: %v", err.(validator.ValidationErrors)))
		return
	}

	// Duplicate code check
	if _, err := h.store.GetApartmentByCode(payload.Code); err == nil {
		utils.WriteError(w, http.StatusConflict,
			fmt.Errorf("apartment with code '%s' already exists", payload.Code))
		return
	}

	created, err := h.store.CreateApartment(types.Apartment{
		Name:        payload.Name,
		Code:        payload.Code,
		Rooms:       payload.Rooms,
		Description: payload.Description,
		Price:       payload.Price,
	})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJson(w, http.StatusCreated, types.ApartmentResponse{
		Message:   "Apartment created successfully",
		Apartment: created,
	})
}

func (h *Handler) handleUpdateApartment(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromVars(r)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	var payload types.UpdateApartmentPayload
	if err := utils.PulseJson(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest,
			fmt.Errorf("invalid payload: %v", err.(validator.ValidationErrors)))
		return
	}

	updated, err := h.store.UpdateApartment(id, payload)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "apartment not found" {
			status = http.StatusNotFound
		}
		utils.WriteError(w, status, err)
		return
	}

	utils.WriteJson(w, http.StatusOK, types.ApartmentResponse{
		Message:   "Apartment updated successfully",
		Apartment: updated,
	})
}

func (h *Handler) handleDeleteApartment(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromVars(r)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	apartment, err := h.store.DeleteApartment(id)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "apartment not found" {
			status = http.StatusNotFound
		}
		utils.WriteError(w, status, err)
		return
	}

	utils.WriteJson(w, http.StatusOK, types.ApartmentResponse{
		Message:   "Apartment deleted successfully",
		Apartment: apartment,
	})
}

// ─── Helpers ──────────────────────────────────────────────────────────────────

func parseIDFromVars(r *http.Request) (int, error) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil || id <= 0 {
		return 0, fmt.Errorf("invalid apartment ID")
	}
	return id, nil
}