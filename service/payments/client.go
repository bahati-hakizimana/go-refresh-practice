package payments

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/pasisltd/go-sdk"
)

// PasisClient wraps the Pasis SDK client with our business logic
type PasisClient struct {
	client *pasis.Client
}

// Transaction represents a Pasis transaction response
type Transaction struct {
	ID       string                 `json:"id"`
	Status   string                 `json:"status"`
	Amount   string                 `json:"amount"`
	Currency string                 `json:"currency"`
	Type     string                 `json:"type"`
	Provider string                 `json:"provider"`
	Metadata map[string]interface{} `json:"metadata"`
}

// Wallet represents wallet information
type Wallet struct {
	ID               string `json:"id"`
	Balance          string `json:"balance"`
	Currency         string `json:"currency"`
	AvailableBalance string `json:"available_balance"`
}

// TransactionList represents paginated transaction list
type TransactionList struct {
	Data       []Transaction `json:"data"`
	Pagination struct {
		Total       int `json:"total"`
		CurrentPage int `json:"current_page"`
		PerPage     int `json:"per_page"`
	} `json:"pagination"`
}

// MerchantProfile represents merchant information
type MerchantProfile struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
}

// NewPasisClient creates a new Pasis client instance
func NewPasisClient(appKey, secretKey string) *PasisClient {
	// Create Pasis client with app credentials
	client := pasis.NewClient(appKey, secretKey)

	return &PasisClient{
		client: client,
	}
}

// DepositParams contains parameters for initiating a deposit
type DepositParams struct {
	Amount      string
	Currency    string
	Provider    string // e.g., "mobile_money", "card"
	Region      string
	PhoneNumber string
	Metadata    map[string]string
}

// WithdrawParams contains parameters for initiating a withdrawal
type WithdrawParams struct {
	Amount      string
	Currency    string
	Provider    string // e.g., "bank_transfer", "mobile_money"
	Region      string
	PhoneNumber string
	Metadata    map[string]string
}

// InitiateDeposit processes a payment deposit through Pasis
func (pc *PasisClient) InitiateDeposit(ctx context.Context, params DepositParams) (*Transaction, error) {
	depositReq := &pasis.DepositRequest{
		Amount:      params.Amount,
		Currency:    params.Currency,
		Provider:    params.Provider,
		Region:      params.Region,
		PhoneNumber: params.PhoneNumber,
		Metadata:    params.Metadata,
	}

	pasisTxn, err := pc.client.Deposit(ctx, depositReq)
	if err != nil {
		// Check for authentication errors
		var authErr *pasis.AuthError
		if errors.As(err, &authErr) {
			log.Printf("Pasis Authentication Error: %v", authErr)
			return nil, fmt.Errorf("authentication failed: %w", authErr)
		}
		
		// Check for validation errors
		var valErr *pasis.ValidationError
		if errors.As(err, &valErr) {
			log.Printf("Pasis Validation Error: %v", valErr)
			return nil, fmt.Errorf("validation failed: %w", valErr)
		}
		
		// Check for API errors
		var apiErr *pasis.APIError
		if errors.As(err, &apiErr) {
			log.Printf("Pasis API Error (status %d): %s", apiErr.StatusCode, apiErr.Message)
			if len(apiErr.Errors) > 0 {
				for _, e := range apiErr.Errors {
					log.Printf("  - %s", e)
				}
			}
			return nil, fmt.Errorf("API error: %w", apiErr)
		}

		log.Printf("Pasis Unknown Error: %v", err)
		return nil, fmt.Errorf("failed to initiate deposit: %w", err)
	}

	// Convert Pasis transaction to our Transaction type
	transaction := &Transaction{
		ID:       pasisTxn.ID,
		Status:   string(pasisTxn.Status),   // Convert TransactionStatus to string
		Amount:   pasisTxn.Amount,
		Currency: pasisTxn.Currency,
		Type:     string(pasisTxn.Type),     // Convert TransactionType to string
		Provider: pasisTxn.Provider,
	}

	return transaction, nil
}

// InitiateWithdraw processes a withdrawal through Pasis
func (pc *PasisClient) InitiateWithdraw(ctx context.Context, params WithdrawParams) (*Transaction, error) {
	withdrawReq := &pasis.WithdrawRequest{
		Amount:      params.Amount,
		Currency:    params.Currency,
		Provider:    params.Provider,
		Region:      params.Region,
		PhoneNumber: params.PhoneNumber,
		Metadata:    params.Metadata,
	}

	pasisTxn, err := pc.client.Withdraw(ctx, withdrawReq)
	if err != nil {
		// Check for authentication errors
		var authErr *pasis.AuthError
		if errors.As(err, &authErr) {
			log.Printf("Pasis Authentication Error: %v", authErr)
			return nil, fmt.Errorf("authentication failed: %w", authErr)
		}
		
		// Check for validation errors
		var valErr *pasis.ValidationError
		if errors.As(err, &valErr) {
			log.Printf("Pasis Validation Error: %v", valErr)
			return nil, fmt.Errorf("validation failed: %w", valErr)
		}
		
		// Check for API errors
		var apiErr *pasis.APIError
		if errors.As(err, &apiErr) {
			log.Printf("Pasis API Error (status %d): %s", apiErr.StatusCode, apiErr.Message)
			if len(apiErr.Errors) > 0 {
				for _, e := range apiErr.Errors {
					log.Printf("  - %s", e)
				}
			}
			return nil, fmt.Errorf("API error: %w", apiErr)
		}

		log.Printf("Pasis Unknown Error: %v", err)
		return nil, fmt.Errorf("failed to initiate withdrawal: %w", err)
	}

	// Convert Pasis transaction to our Transaction type
	transaction := &Transaction{
		ID:       pasisTxn.ID,
		Status:   string(pasisTxn.Status),   // Convert TransactionStatus to string
		Amount:   pasisTxn.Amount,
		Currency: pasisTxn.Currency,
		Type:     string(pasisTxn.Type),     // Convert TransactionType to string
		Provider: pasisTxn.Provider,
	}

	return transaction, nil
}

// GetTransactionStatus retrieves the status of a transaction
func (pc *PasisClient) GetTransactionStatus(ctx context.Context, transactionID string) (*Transaction, error) {
	pasisTxn, err := pc.client.GetTransaction(ctx, transactionID)
	if err != nil {
		// Check for authentication errors
		var authErr *pasis.AuthError
		if errors.As(err, &authErr) {
			log.Printf("Pasis Authentication Error: %v", authErr)
			return nil, fmt.Errorf("authentication failed: %w", authErr)
		}
		
		// Check for API errors
		var apiErr *pasis.APIError
		if errors.As(err, &apiErr) {
			log.Printf("Pasis API Error (status %d): %s", apiErr.StatusCode, apiErr.Message)
			return nil, fmt.Errorf("API error: %w", apiErr)
		}

		log.Printf("Pasis Unknown Error: %v", err)
		return nil, fmt.Errorf("failed to get transaction status: %w", err)
	}

	// Convert Pasis transaction to our Transaction type
	transaction := &Transaction{
		ID:       pasisTxn.ID,
		Status:   string(pasisTxn.Status),   
		Amount:   pasisTxn.Amount,
		Currency: pasisTxn.Currency,
		Type:     string(pasisTxn.Type),     // Convert TransactionType to string
		Provider: pasisTxn.Provider,
	}

	return transaction, nil
}

// GetWalletBalance retrieves the current wallet balance
func (pc *PasisClient) GetWalletBalance(ctx context.Context) (*Wallet, error) {
	pasisWallet, err := pc.client.GetWallet(ctx)
	if err != nil {
		// Check for authentication errors
		var authErr *pasis.AuthError
		if errors.As(err, &authErr) {
			log.Printf("Pasis Authentication Error: %v", authErr)
			return nil, fmt.Errorf("authentication failed: %w", authErr)
		}
		
		// Check for API errors
		var apiErr *pasis.APIError
		if errors.As(err, &apiErr) {
			log.Printf("Pasis API Error (status %d): %s", apiErr.StatusCode, apiErr.Message)
			return nil, fmt.Errorf("API error: %w", apiErr)
		}

		log.Printf("Pasis Unknown Error: %v", err)
		return nil, fmt.Errorf("failed to get wallet balance: %w", err)
	}

	// Convert Pasis wallet to our Wallet type
	wallet := &Wallet{
		ID:               pasisWallet.ID,
		Balance:          pasisWallet.Balance,
		Currency:         pasisWallet.Currency,
		AvailableBalance: pasisWallet.Balance, 
	}

	return wallet, nil
}

// ListTransactions retrieves paginated list of transactions
func (pc *PasisClient) ListTransactions(ctx context.Context, page, pageSize int) (*TransactionList, error) {
	pasisList, err := pc.client.ListTransactions(ctx, page, pageSize)
	if err != nil {
		// Check for authentication errors
		var authErr *pasis.AuthError
		if errors.As(err, &authErr) {
			log.Printf("Pasis Authentication Error: %v", authErr)
			return nil, fmt.Errorf("authentication failed: %w", authErr)
		}
		
		// Check for API errors
		var apiErr *pasis.APIError
		if errors.As(err, &apiErr) {
			log.Printf("Pasis API Error (status %d): %s", apiErr.StatusCode, apiErr.Message)
			return nil, fmt.Errorf("API error: %w", apiErr)
		}

		log.Printf("Pasis Unknown Error: %v", err)
		return nil, fmt.Errorf("failed to list transactions: %w", err)
	}

	// Convert Pasis transactions to our Transaction type
	transactions := make([]Transaction, len(pasisList.Data))
	for i, pasisTxn := range pasisList.Data {
		transactions[i] = Transaction{
			ID:       pasisTxn.ID,
			Status:   string(pasisTxn.Status),   // Convert TransactionStatus to string
			Amount:   pasisTxn.Amount,
			Currency: pasisTxn.Currency,
			Type:     string(pasisTxn.Type),     // Convert TransactionType to string
			Provider: pasisTxn.Provider,
		}
	}

	list := &TransactionList{
		Data: transactions,
	}
	list.Pagination.Total = pasisList.Pagination.Total
	list.Pagination.CurrentPage = pasisList.Pagination.Page  // Use Page field
	list.Pagination.PerPage = pasisList.Pagination.PerPage

	return list, nil
}

// GetMerchantProfile retrieves the merchant profile information
func (pc *PasisClient) GetMerchantProfile(ctx context.Context) (*MerchantProfile, error) {
	pasisProfile, err := pc.client.GetMerchantProfile(ctx)
	if err != nil {
		// Check for authentication errors
		var authErr *pasis.AuthError
		if errors.As(err, &authErr) {
			log.Printf("Pasis Authentication Error: %v", authErr)
			return nil, fmt.Errorf("authentication failed: %w", authErr)
		}
		
		// Check for API errors
		var apiErr *pasis.APIError
		if errors.As(err, &apiErr) {
			log.Printf("Pasis API Error (status %d): %s", apiErr.StatusCode, apiErr.Message)
			return nil, fmt.Errorf("API error: %w", apiErr)
		}

		log.Printf("Pasis Unknown Error: %v", err)
		return nil, fmt.Errorf("failed to get merchant profile: %w", err)
	}

	// Convert Pasis profile to our MerchantProfile type
	profile := &MerchantProfile{
		FirstName:   pasisProfile.FirstName,
		LastName:    pasisProfile.LastName,
		Email:       pasisProfile.Email,
		PhoneNumber: pasisProfile.PhoneNumber,
	}

	return profile, nil
}