package apartmentimage

import (
	"database/sql"

	"github.com/go-refresh-practice/go-refresh-course/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

// Get Apartment images by  apartment id

func (s *Store) GetImagesByApartmentID(apartmentID int) ([]types.ApartmentImage, error){

}

// Create Apartment images for a given apartment id


func (s *Store) CreateApartmentImage (img types.ApartmentImage) (types.ApartmentImage, error) {

}


//scan rows into apartment image

func ScanRowsIntoApartmentImage(rows *sql.Rows) (*types.ApartmentImage, error) {

	
}