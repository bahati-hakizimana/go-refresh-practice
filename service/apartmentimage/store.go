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

// ----------------------------------------------------
// Add image to apartment (simple version)
// ----------------------------------------------------

func (s *Store) AddImageToApartment(apartmentID int, imageURL string) (types.ApartmentImage, error) {
	img := types.ApartmentImage{
		ApartmentID: apartmentID,
		ImageURL:    imageURL,
	}

	err := s.db.QueryRow(`
        INSERT INTO apartment_images (apartment_id, image_url)
        VALUES ($1, $2)
        RETURNING id, apartment_id, image_url, caption, created_at
    `, apartmentID, imageURL).Scan(
		&img.ID,
		&img.ApartmentID,
		&img.ImageURL,
		&img.Caption,
		&img.CreatedAt,
	)

	if err != nil {
		return types.ApartmentImage{}, err
	}

	return img, nil
}

// ----------------------------------------------------
// Main CREATE method (with caption)
// ----------------------------------------------------

func (s *Store) CreateApartmentImage(img types.ApartmentImage) (types.ApartmentImage, error) {
	err := s.db.QueryRow(`
        INSERT INTO apartment_images (apartment_id, image_url, caption)
        VALUES ($1, $2, $3)
        RETURNING id, apartment_id, image_url, caption, created_at
    `, img.ApartmentID, img.ImageURL, img.Caption).Scan(
		&img.ID,
		&img.ApartmentID,
		&img.ImageURL,
		&img.Caption,
		&img.CreatedAt,
	)

	if err != nil {
		return types.ApartmentImage{}, err
	}

	return img, nil
}

// ----------------------------------------------------
// GET images by apartment id
// ----------------------------------------------------

func (s *Store) GetImagesByApartmentID(apartmentID int) ([]types.ApartmentImage, error) {
	rows, err := s.db.Query(`
        SELECT id, apartment_id, image_url, caption, created_at
        FROM apartment_images
        WHERE apartment_id = $1
        ORDER BY created_at DESC
    `, apartmentID)

	if err != nil {
		return nil, err
	}

	images := []types.ApartmentImage{}
	for rows.Next() {
		img, err := ScanRowsIntoApartmentImage(rows)
		if err != nil {
			return nil, err
		}
		images = append(images, *img)
	}

	return images, nil
}

func (s *Store) GetAllImages() ([]types.ApartmentImage, error) {
    rows, err := s.db.Query("SELECT id, apartment_id, image_url, caption, created_at FROM apartment_images")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    images := make([]types.ApartmentImage, 0)
    for rows.Next() {
        img, err := ScanRowsIntoApartmentImage(rows)
        if err != nil {
            return nil, err
        }
        images = append(images, *img)
    }

    return images, nil
}


// ----------------------------------------------------
// SCAN rows into struct
// ----------------------------------------------------

func ScanRowsIntoApartmentImage(rows *sql.Rows) (*types.ApartmentImage, error) {
	img := new(types.ApartmentImage)

	err := rows.Scan(
		&img.ID,
		&img.ApartmentID,
		&img.ImageURL,
		&img.Caption,
		&img.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return img, nil
}
