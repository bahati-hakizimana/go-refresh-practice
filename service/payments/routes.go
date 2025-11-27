package payments

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/go-refresh-practice/go-refresh-course/middleware"
	"github.com/go-refresh-practice/go-refresh-course/types"
	"github.com/go-refresh-practice/go-refresh-course/utils"
	"github.com/gorilla/mux"

	// "github.com/quarksgroup/paypack-go"
	// "github.com/quarksgroup/paypack-go/api"
	// "github.com/quarksgroup/paypack-go/oauth"
	"github.com/quarksgroup/paypack-go/paypack"
	"github.com/quarksgroup/paypack-go/paypack/api"
	"github.com/quarksgroup/paypack-go/paypack/transport/oauth"
)

type Handler struct {
	store   types.PaymentStore
	paypack *paypack.Client
}

func NewHandler(store types.PaymentStore) *Handler {
	cli := api.NewDefault()
	cli.Http = &http.Client{
		Transport: &oauth.Transport{
			Scheme: oauth.SchemeBearer,
			Source: oauth.ContextTokenSource(),
			Base:   http.DefaultTransport,
		},
	}

	return &Handler{
		store:   store,
		paypack: cli,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {

	router.Handle("/payments",
		middleware.AuthMiddleware(
			middleware.AdminOnly(
				http.HandlerFunc(h.handlerGetPayments),
			)),
	).Methods(http.MethodGet)

	router.Handle("/payments",
		http.HandlerFunc(h.handlerCreatePayment),
	).Methods(http.MethodPost)
}

/* ------------------ Paypack Cashin Helper --------------------- */

type PaypackTxReq struct {
	Amount float64
	Phone  string
	Mode   string
}

func (h *Handler) Cashin(ctx context.Context, tx PaypackTxReq, accessToken string) (*paypack.TransactionResponse, error) {

	ctx = context.WithValue(ctx, paypack.TokenKey{}, &paypack.Token{
		Access: accessToken,
	})

	req := &paypack.TransactionRequest{
		Amount: tx.Amount,
		Phone:  tx.Phone,
		Mode:   tx.Mode,
	}

	res, err := h.paypack.Transaction.Cashin(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("cashin error: %w", err)
	}

	return res, nil
}

/* ------------------ GET ALL PAYMENTS --------------------- */

func (h *Handler) handlerGetPayments(w http.ResponseWriter, r *http.Request) {
	payments, err := h.store.GetPayments()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJson(w, http.StatusOK, payments)
}

/* ------------------ CREATE PAYMENT --------------------- */

func (h *Handler) handlerCreatePayment(w http.ResponseWriter, r *http.Request) {

	var payload types.PaymentPayload

	if err := utils.PulseJson(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		error := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", error))
		return
	}

	// Parse date
	layout := "2006-01-02"
	paidAt, err := time.Parse(layout, payload.PaidAt)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid paidAt date format"))
		return
	}

	// User must send phone in query or you fetch phone from booking
	phone := r.URL.Query().Get("phone")
	if phone == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("phone is required for payment"))
		return
	}

	// For now we allow dev testing —
	accessToken := "YOUR_PAYPACK_ACCESS_TOKEN"

	tx := PaypackTxReq{
		Amount: payload.Amount,
		Phone:  phone,
		Mode:   "development",
	}

	payRes, err := h.Cashin(r.Context(), tx, accessToken)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("payment failed: %v", err))
		return
	}

	payment := types.Payment{
		BookingID:            payload.BookingID,
		PaymentType:          payload.PaymentType,
		Amount:               payload.Amount,
		Currency:             payload.Currency,
		PaymentMethod:        payload.PaymentMethod,
		PaymentStatus:        payRes.Status,
		TransactionReference: payRes.Ref,
		PaidAt:               paidAt,
	}

	createdPayment, err := h.store.CreatePayment(payment)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	response := map[string]interface{}{
		"message": "Payment initiated successfully",
		"payment": createdPayment,
		"paypack": payRes,
	}

	utils.WriteJson(w, http.StatusCreated, response)
}
