package types

import "time"


type ApartmentStore interface{
	GetApartments() ([]Apartment, error)
	CreateApartment(apartment Apartment) (Apartment, error)
}

type Apartment struct {

	ID		  int     `json:"id"`
	Name     string  `json:"name"`
	Code	 string  `json:"code"`
	Rooms	 int     `json:"rooms"`
	Description string  `json:"description"`
	Price     float64 `json:"price"`
	CreatedAt time.Time `json:"createdAt"`

}

type UserStore interface{
	GetUserByEmail(email string) (*User, error)
	GetUserById(id int) (*User, error)
	CreateUser(User) error
}


type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
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
	Email string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}