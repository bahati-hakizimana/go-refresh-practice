package types

import "time"


type ApartmentStore interface{
	GetApartments() ([]Apartment, error)
	GetApartmentByCode(code string) (*Apartment, error)
	CreateApartment(apartment Apartment) (Apartment, error)
}

type Apartment struct {

	ID		  int     `json:"id"`
	Name     string  `json:"name"`
	Code	 string  `json:"code"`
	Rooms	 int     `json:"rooms"`
	Description string  `json:"description"`
	Price     float64 `json:"price"`
	Status	string  `json:"status"`
	CreatedAt time.Time `json:"createdAt"`

}

type CreateApartmentPayload struct {
	Name     string  `json:"name" validate:"required"`
	Code	 string  `json:"code" validate:"required"`
	Rooms	 int     `json:"rooms" validate:"required,min=1"`
	Description string  `json:"description" validate:"required"`
	Price     float64 `json:"price" validate:"required,gt=0"`
}

type ApartmentImagesStore interface {
    GetAllImages() ([]ApartmentImage, error)                     
    GetImagesByApartmentID(apartmentID int) ([]ApartmentImage, error)
    CreateApartmentImage(image ApartmentImage) (ApartmentImage, error)
}


type ApartmentImage struct {
	ID          int       `json:"id"`
	ApartmentID int       `json:"apartmentId"`
	ImageURL    string    `json:"imageUrl"`
	Caption    string     `json:"caption"`
	CreatedAt   time.Time `json:"createdAt"`
}

type ApartmentImagePayload struct {
	ImageURL string `json:"imageUrl" validate:"required,url"`
	Caption  string `json:"caption" validate:"required"`
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