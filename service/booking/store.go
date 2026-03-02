package booking

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

func (s *Store) GetBookings() ([]types.Booking, error) {
	rows, err := s.db.Query(getBookingsQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	bookings := make([]types.Booking, 0)

	for rows.Next() {
		b, err := scanRowsIntoBooking(rows)
		if err != nil {
			return nil, err
		}
		bookings = append(bookings, *b)
	}

	return bookings, nil
}

func (s *Store) GetBookingByID(id int) (*types.Booking, error) {
	row := s.db.QueryRow(getBookingByIDQuery, id)

	booking := new(types.Booking)

	err := row.Scan(
		&booking.ID,
		&booking.ApartmentID,
		&booking.FirstName,
		&booking.LastName,
		&booking.Email,
		&booking.PhoneNumber,
		&booking.CheckinDate,
		&booking.CheckoutDate,
		&booking.GuestNumber,
		&booking.TotalPrice,
		&booking.BookingAmount,
		&booking.BalanceAmount,
		&booking.Currency,
		&booking.BookingStatus,
		&booking.CreatedAt,
		&booking.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("booking with id %d not found", id)
		}
		return nil, err
	}

	return booking, nil
}

func (s *Store) CreateBooking(b types.Booking) (types.Booking, error) {
	row := s.db.QueryRow(
		createBookingQuery,
		b.ApartmentID,
		b.FirstName,
		b.LastName,
		b.Email,
		b.PhoneNumber,
		b.CheckinDate,
		b.CheckoutDate,
		b.GuestNumber,
		b.TotalPrice,
		b.BookingAmount,
		b.BalanceAmount,
		b.Currency,
	)

	err := row.Scan(&b.ID)
	if err != nil {
		return types.Booking{}, fmt.Errorf("failed to create booking: %w", err)
	}

	return b, nil
}


// =========================
// Helper
// =========================


func scanRowsIntoBooking(rows *sql.Rows) (*types.Booking, error) {
	booking := new(types.Booking)

	err := rows.Scan(
		&booking.ID,
		&booking.ApartmentID,
		&booking.FirstName,
		&booking.LastName,
		&booking.Email,
		&booking.PhoneNumber,
		&booking.CheckinDate,
		&booking.CheckoutDate,
		&booking.GuestNumber,
		&booking.TotalPrice,
		&booking.BookingAmount,
		&booking.BalanceAmount,
		&booking.Currency,
		&booking.BookingStatus,
		&booking.CreatedAt,
		&booking.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return booking, nil
}




// =========================
// Queries (Professional Way)
// =========================

const getBookingsQuery = `
	SELECT id, apartment_id, first_name, last_name, email, phone_number,
	       checkin_date, checkout_date, guest_number, total_price,
	       booking_amount, balance_amount, currency, booking_status,
	       created_at, updated_at
	FROM bookings
`

const getBookingByIDQuery = `
	SELECT id, apartment_id, first_name, last_name, email, phone_number,
	       checkin_date, checkout_date, guest_number, total_price,
	       booking_amount, balance_amount, currency, booking_status,
	       created_at, updated_at
	FROM bookings
	WHERE id = $1
`

const createBookingQuery = `
	INSERT INTO bookings
	(apartment_id, first_name, last_name, email, phone_number,
	 checkin_date, checkout_date, guest_number, total_price,
	 booking_amount, balance_amount, currency,
	 booking_status, created_at, updated_at)
	VALUES
	($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,'pending', NOW(), NOW())
	RETURNING id
`


// =========================
// Methods
// =========================

