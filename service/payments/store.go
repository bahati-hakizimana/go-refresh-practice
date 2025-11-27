package payments

import (
	"database/sql"
	"fmt"

	"github.com/go-refresh-practice/go-refresh-course/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

/* ------------------ CREATE PAYMENT --------------------- */

func (s *Store) CreatePayment(p types.Payment) (types.Payment, error) {

	err := s.db.QueryRow(`
		INSERT INTO payments
		(booking_id, payment_type, amount, currency, payment_status, payment_method, transaction_reference, paid_at, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`,
		p.BookingID,
		p.PaymentType,
		p.Amount,
		p.Currency,
		p.PaymentStatus,
		p.PaymentMethod,
		p.TransactionReference,
		p.PaidAt,
	).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)

	if err != nil {
		return types.Payment{}, fmt.Errorf("failed to create payment: %w", err)
	}

	return p, nil
}

/* ------------------ GET ALL PAYMENTS --------------------- */

func (s *Store) GetPayments() ([]types.Payment, error) {

	rows, err := s.db.Query(`
		SELECT id, booking_id, payment_type, amount, currency, payment_status, payment_method,
		       transaction_reference, paid_at, created_at, updated_at
		FROM payments
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	payments := []types.Payment{}

	for rows.Next() {
		var p types.Payment

		err := rows.Scan(
			&p.ID,
			&p.BookingID,
			&p.PaymentType,
			&p.Amount,
			&p.Currency,
			&p.PaymentStatus,
			&p.PaymentMethod,
			&p.TransactionReference,
			&p.PaidAt,
			&p.CreatedAt,
			&p.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		payments = append(payments, p)
	}

	return payments, nil
}

/* ------------------ GET PAYMENT BY ID --------------------- */

func (s *Store) GetPaymentByID(id int) (*types.Payment, error) {

	row := s.db.QueryRow(`
		SELECT id, booking_id, payment_type, amount, currency, payment_status, payment_method,
		       transaction_reference, paid_at, created_at, updated_at
		FROM payments WHERE id = $1
	`, id)

	var p types.Payment

	err := row.Scan(
		&p.ID,
		&p.BookingID,
		&p.PaymentType,
		&p.Amount,
		&p.Currency,
		&p.PaymentStatus,
		&p.PaymentMethod,
		&p.TransactionReference,
		&p.PaidAt,
		&p.CreatedAt,
		&p.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("payment not found")
	}

	return &p, nil
}

