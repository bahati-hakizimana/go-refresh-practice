package aprtment

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

func (s *Store) GetApartments() ([]types.Apartment, error) {
	rows, err := s.db.Query("SELECT * FROM apartments")
	if err != nil {
		return  nil, err
	}

	apartments := make([]types.Apartment, 0)
	for rows.Next() {
		apt, err:= scanRowsIntoApartment(rows)
		if err != nil {
			return nil, err
		}

		apartments = append(apartments, *apt)
	}

	return apartments, nil
}


func (s *Store) CreateApartment(ap types.Apartment) (types.Apartment, error) {
    // Insert into Postgres and return the generated ID
    err := s.db.QueryRow(`
        INSERT INTO apartments (code, name, rooms, description, price)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id
    `, ap.Code, ap.Name, ap.Rooms, ap.Description, ap.Price).Scan(&ap.ID)

    if err != nil {
        return types.Apartment{}, err
    }

    return ap, nil
}

func scanRowsIntoApartment(rows *sql.Rows) (*types.Apartment, error) {
	apartment := new(types.Apartment)

	err := rows.Scan(

		&apartment.ID,
		&apartment.Name,
		&apartment.Code,
		&apartment.Rooms,
		&apartment.Description,
		&apartment.Price,
		&apartment.CreatedAt,

	)

	if err != nil {
		return nil, err
	}

	return  apartment, nil
}


