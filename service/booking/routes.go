package booking

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/go-refresh-practice/go-refresh-course/middleware"
	"github.com/go-refresh-practice/go-refresh-course/types"
	"github.com/go-refresh-practice/go-refresh-course/utils"
	"github.com/gorilla/mux"
)

type Handler struct {
	store types.BookingStore
}

func NewHandler(store types.BookingStore) *Handler {
	return &Handler{store: store}
}

func(h *Handler) RegisterRoutes(router *mux.Router){
	// Fetch and get all Bookings admin only
			router.Handle("/bookings",
				middleware.AuthMiddleware(middleware.AdminOnly(http.HandlerFunc(h.handleGetBookings))),
		).Methods(http.MethodGet)

		// create booking

		router.Handle("/bookings",
    http.HandlerFunc(h.handleCreateBookings),
).Methods(http.MethodPost)
}

// Get all bookings for admin only
func(h *Handler)handleGetBookings(w http.ResponseWriter, r *http.Request){
	booking, err := h.store.GetBookings()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJson(w, http.StatusOK, booking)
} 

// make or create bookings
func(h *Handler)handleCreateBookings(w http.ResponseWriter, r *http.Request){
	var payload types.BookingPayload
	if err := utils.PulseJson(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate payload

	if err := utils.Validate.Struct(payload); err != nil {
		error:= err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("Envalid payload %v", error))
		return
	}


	 // Parse dates
    layout := "2006-01-02"
    checkin, err := time.Parse(layout, payload.CheckinDate)
    if err != nil {
        utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid checkin date format"))
        return
    }

    checkout, err := time.Parse(layout, payload.CheckoutDate)
    if err != nil {
        utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid checkout date format"))
        return
    }


	// Create booking struct
    bking := types.Booking{
        ApartmentID:   payload.ApartmentID,
        FirstName:     payload.FirstName,
        LastName:      payload.LastName,
        Email:         payload.Email,
        PhoneNumber:   payload.PhoneNumber,
        CheckinDate:   checkin,
        CheckoutDate:  checkout,
        GuestNumber:   payload.GuestNumber,
        TotalPrice:    payload.TotalPrice,
        BookingAmount: payload.BookingAmount,
        BalanceAmount: payload.BalanceAmount,
        Currency:      payload.Currency,
    }

   // Create booking and get updated booking with ID
createdBooking, err := h.store.CreateBooking(bking)
if err != nil {
	utils.WriteError(w, http.StatusInternalServerError, err)
	return
}

response := map[string]interface{}{
	"message": "Booking initiated successfully",
	"booking": createdBooking, // <-- use the booking with the ID
}

utils.WriteJson(w, http.StatusCreated, response)

}