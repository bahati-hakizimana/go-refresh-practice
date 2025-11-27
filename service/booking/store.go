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

func(s *Store)GetBookings() ([]types.Booking, error) {
	rows, err := s.db.Query(`
	SELECT id, apartment_id, first_name, last_name, email, phone_number, checkin_date, checkout_date,
	       guest_number, total_price, booking_amount, balance_amount, currency, booking_status,
	       created_at, updated_at
	FROM bookings
`)


	if err != nil {
		return nil, err
	}

	bookings := make([]types.Booking, 0)
	for rows.Next() {
		bking, err := scanRowsIntoBooking(rows)
		if err != nil {
			return nil, err
		}

		bookings = append(bookings, *bking)
	}

	return  bookings, nil
}

// Get booking by id

func(s *Store)GetBookingByID(id int) (*types.Booking, error) {

	rows, err :=s.db.Query("SELECT, id, apartment_id, first_name, last_name, email, phoneNumber, checkinDate, checkoutDate, guestNumber, totalPrice, bookingAmount, balanceAmount, currency, bookingStatus, createdAt, updatedAt FROM bookings WHERE id = $1", id)

	if err != nil {
		return nil, err
	}

	bking := new(types.Booking)
	for rows.Next(){
		bking, err = scanRowsIntoBooking(rows)

		if err != nil {
			return nil, err
		}
	}

	if bking.ID == 0 {
			return  nil, fmt.Errorf("Booking with this id %v not found", id)
		}

		return  bking, nil


}



// Create booking

func (s *Store) CreateBooking(bkg types.Booking) (types.Booking, error) {
	// Use QueryRow with INSERT and RETURNING id
	row := s.db.QueryRow(`
		INSERT INTO bookings 
		(apartment_id, first_name, last_name, email, phone_number, checkin_date, checkout_date, guest_number, total_price, booking_amount, balance_amount, currency, booking_status, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,'pending', NOW(), NOW())
		RETURNING id
	`,
		bkg.ApartmentID,
		bkg.FirstName,
		bkg.LastName,
		bkg.Email,
		bkg.PhoneNumber,
		bkg.CheckinDate,
		bkg.CheckoutDate,
		bkg.GuestNumber,
		bkg.TotalPrice,
		bkg.BookingAmount,
		bkg.BalanceAmount,
		bkg.Currency,
	)

	// Scan returned ID into bkg
	err := row.Scan(&bkg.ID)
	if err != nil {
		return types.Booking{}, fmt.Errorf("failed to create booking: %w", err)
	}

	return bkg, nil
}




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
