package payments

import (
	// "database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-refresh-practice/go-refresh-course/types"
	"github.com/go-refresh-practice/go-refresh-course/utils"
	"github.com/gorilla/mux"
)

type Handler struct {
	store        *Store
	pasisClient  *PasisClient
}

func NewHandler(store *Store, pasisClient *PasisClient) *Handler {
	return &Handler{
		store:       store,
		pasisClient: pasisClient,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/payments", h.handleCreatePayment).Methods("POST")
	router.HandleFunc("/payments", h.handleGetPayments).Methods("GET")
	router.HandleFunc("/payments/{id}", h.handleGetPaymentByID).Methods("GET")
	router.HandleFunc("/payments/{id}/status", h.handleCheckPaymentStatus).Methods("GET")
	router.HandleFunc("/payments/wallet/balance", h.handleGetWalletBalance).Methods("GET")
}

/* ------------------ CREATE PAYMENT --------------------- */

type CreatePaymentRequest struct {
	BookingID     int    `json:"booking_id"`
	Amount        string `json:"amount"`
	Currency      string `json:"currency"`
	PaymentMethod string `json:"payment_method"` // e.g., "mobile_money", "card", "bank_transfer"
	PhoneNumber   string `json:"phone_number"`
	Region        string `json:"region"` // e.g., "RW", "US"
}

func (h *Handler) handleCreatePayment(w http.ResponseWriter, r *http.Request) {
	var req CreatePaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid request payload"))
		return
	}

	// Validate required fields
	if req.BookingID == 0 || req.Amount == "" || req.Currency == "" || req.PaymentMethod == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing required fields"))
		return
	}

	// Convert amount string to float64
	amount, err := strconv.ParseFloat(req.Amount, 64)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid amount format"))
		return
	}

	// Initiate payment with Pasis
	ctx := r.Context()
	
	depositParams := DepositParams{
		Amount:      req.Amount,
		Currency:    req.Currency,
		Provider:    req.PaymentMethod,
		Region:      req.Region,
		PhoneNumber: req.PhoneNumber,
		Metadata: map[string]string{
			"booking_id": strconv.Itoa(req.BookingID),
			"source":     "apartment_booking",
		},
	}

	transaction, err := h.pasisClient.InitiateDeposit(ctx, depositParams)
	if err != nil {
		// Check if it's an authentication or validation error based on the error message
		errMsg := strings.ToLower(err.Error())
		if strings.Contains(errMsg, "authentication") || strings.Contains(errMsg, "unauthorized") {
			utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("payment gateway authentication failed"))
			return
		}
		if strings.Contains(errMsg, "validation") || strings.Contains(errMsg, "invalid") {
			utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payment details: %v", err))
			return
		}

		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to process payment: %v", err))
		return
	}

	// Create payment record in database
	payment := types.Payment{
		BookingID:            req.BookingID,
		PaymentType:          "deposit",
		Amount:               amount,                // Now a float64
		Currency:             req.Currency,
		PaymentStatus:        transaction.Status,    // "pending", "completed", "failed"
		PaymentMethod:        req.PaymentMethod,
		TransactionReference: transaction.ID,
		PaidAt:               time.Time{},           // Zero time.Time (will be updated when confirmed)
	}

	createdPayment, err := h.store.CreatePayment(payment)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// Return payment details with transaction info
	response := map[string]interface{}{
		"payment":     createdPayment,
		"transaction": transaction,
		"message":     "Payment initiated successfully",
	}

	utils.WriteJson(w, http.StatusCreated, response)  // Fixed: WriteJson not WriteJSON
}

/* ------------------ GET ALL PAYMENTS --------------------- */

func (h *Handler) handleGetPayments(w http.ResponseWriter, r *http.Request) {
	payments, err := h.store.GetPayments()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJson(w, http.StatusOK, payments)
}

/* ------------------ GET PAYMENT BY ID --------------------- */

func (h *Handler) handleGetPaymentByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payment ID"))
		return
	}

	payment, err := h.store.GetPaymentByID(id)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, err)
		return
	}

	utils.WriteJson(w, http.StatusOK, payment)
}

/* ------------------ CHECK PAYMENT STATUS --------------------- */

func (h *Handler) handleCheckPaymentStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payment ID"))
		return
	}

	// Get payment from database
	payment, err := h.store.GetPaymentByID(id)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, err)
		return
	}

	// Check status with Pasis
	ctx := r.Context()
	transaction, err := h.pasisClient.GetTransactionStatus(ctx, payment.TransactionReference)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to check payment status"))
		return
	}

	// Update payment status if changed
	if transaction.Status != payment.PaymentStatus {
		payment.PaymentStatus = transaction.Status
		
		// If payment is completed, set paid_at timestamp
		if transaction.Status == "completed" && payment.PaidAt.IsZero() {
			payment.PaidAt = time.Now()
			
			// Update in database
			err = h.store.UpdatePaymentStatus(id, transaction.Status, payment.PaidAt)
			if err != nil {
				utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to update payment status"))
				return
			}
		}
	}

	response := map[string]interface{}{
		"payment":     payment,
		"transaction": transaction,
	}

	utils.WriteJson(w, http.StatusOK, response)
}

/* ------------------ GET WALLET BALANCE --------------------- */

func (h *Handler) handleGetWalletBalance(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	wallet, err := h.pasisClient.GetWalletBalance(ctx)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to get wallet balance"))
		return
	}

	utils.WriteJson(w, http.StatusOK, wallet)
}