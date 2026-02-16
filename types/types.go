package types

import "time"

type ApartmentStore interface {
	GetApartments() ([]Apartment, error)
	GetApartmentByCode(code string) (*Apartment, error)
	CreateApartment(apartment Apartment) (Apartment, error)
	GetPublicApartments() ([]Apartment, error)
	DeleteApartment(id int) (Apartment, error)
}

type Apartment struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	Rooms       int       `json:"rooms"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
}

type CreateApartmentPayload struct {
	Name        string  `json:"name" validate:"required"`
	Code        string  `json:"code" validate:"required"`
	Rooms       int     `json:"rooms" validate:"required,min=1"`
	Description string  `json:"description" validate:"required"`
	Price       float64 `json:"price" validate:"required,gt=0"`
}

type ApartmentImagesStore interface {
	GetAllImages() ([]ApartmentImage, error)
	GetImagesByApartmentID(apartmentID int) ([]ApartmentImage, error)
	CreateApartmentImage(image ApartmentImage) (ApartmentImage, error)
	DeleteApartmentImage(imageID int) (ApartmentImage, error) // declare only
}


type ApartmentImage struct {
	ID          int       `json:"id"`
	ApartmentID int       `json:"apartmentId"`
	ImageURL    string    `json:"imageUrl"`
	Caption     string    `json:"caption"`
	CreatedAt   time.Time `json:"createdAt"`
}

type ApartmentImagePayload struct {
	ImageURL string `json:"imageUrl" validate:"required,url"`
	Caption  string `json:"caption" validate:"required"`
}

// Starting  bookings struct and interface

type BookingStore interface {
	GetBookings() ([]Booking, error)
	GetBookingByID(id int) (*Booking, error)
	CreateBooking(booking Booking) (Booking, error)
}

type Booking struct {
	ID            int       `json:"id"`
	ApartmentID   int       `json:"apartment_id"`
	FirstName     string    `json:"firstName"`
	LastName      string    `json:"lastName"`
	Email         string    `json:"email"`
	PhoneNumber   string    `json:"phoneNumber"`
	CheckinDate   time.Time `json:"checkinDate"`
	CheckoutDate  time.Time `json:"checkoutDate"`
	GuestNumber   int       `json:"guestNumber"`
	TotalPrice    float64   `json:"totalPrice"`
	BookingAmount float64   `json:"bookingAmount"`
	BalanceAmount float64   `json:"balanceAmount"`
	Currency      string    `json:"currency"`
	BookingStatus string    `json:"bookingStatus"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

// BookingPayload is used for incoming booking requests
type BookingPayload struct {
	ApartmentID   int     `json:"apartment_id" validate:"required,gt=0"`
	FirstName     string  `json:"firstName" validate:"required"`
	LastName      string  `json:"lastName" validate:"required"`
	Email         string  `json:"email" validate:"required,email"`
	PhoneNumber   string  `json:"phoneNumber" validate:"required"`
	CheckinDate   string  `json:"checkinDate" validate:"required,datetime=2006-01-02"`
	CheckoutDate  string  `json:"checkoutDate" validate:"required,datetime=2006-01-02"`
	GuestNumber   int     `json:"guestNumber" validate:"required,gt=0"`
	TotalPrice    float64 `json:"totalPrice" validate:"required,gt=0"`
	BookingAmount float64 `json:"bookingAmount" validate:"required,gt=0"`
	BalanceAmount float64 `json:"balanceAmount" validate:"required,gte=0"`
	Currency      string  `json:"currency" validate:"omitempty"`
}

// payment struct and interface

type PaymentStore interface {
	GetPayments() ([]Payment, error)
	GetPaymentByID(id int) (*Payment, error)
	CreatePayment(payment Payment) (Payment, error)
}

type Payment struct {
	ID                   int       `json:"id"`
	BookingID            int       `json:"booking_id"`
	PaymentType          string    `json:"payment_type"`
	Amount               float64   `json:"amount"`
	Currency             string    `json:"currency"`
	PaymentStatus        string    `json:"payment_status"`
	PaymentMethod        string    `json:"payment_method"`
	TransactionReference string    `json:"transaction_reference"`
	PaidAt               time.Time `json:"paid_at"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

type PaymentPayload struct {
	BookingID     int     `json:"booking_id" validate:"required,gt=0"`
	PaymentType   string  `json:"payment_type" validate:"required,oneof=deposit balance full refund"`
	Amount        float64 `json:"amount" validate:"required,gt=0"`
	Currency      string  `json:"currency" validate:"omitempty"`
	PaymentMethod string  `json:"payment_method" validate:"required,oneof=mtn card paypack bank_transfer cash"`
	PaidAt        string  `json:"paidAt" validate:"required,datetime=2006-01-02"`
}

type PaypackTxReq struct {
	Amount float64
	Phone  string
	Mode   string
}

type UserStore interface {
	GetUserByEmail(email string) (*User, error)
	GetUserById(id int) (*User, error)
	CreateUser(User) error
}

type User struct {
	ID        int       `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
}

type RegisterUserPayload struct {
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=3,max=130"`
}

type LoginUserPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
