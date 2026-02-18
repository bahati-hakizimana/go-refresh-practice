// package aprtment

// import (
// 	"database/sql"
// 	"fmt"

// 	"github.com/go-refresh-practice/go-refresh-course/types"
// )

// type Store struct {
// 	db *sql.DB
// }

// func NewStore(db *sql.DB) *Store {
// 	return &Store{db: db}
// }

// // GetPublicApartments returns only available apartments for public viewing
// func (s *Store) GetPublicApartments() ([]types.Apartment, error) {
// 	rows, err := s.db.Query(`
// 		SELECT id, code, name, rooms, description, price, status, created_at
// 		FROM apartments
// 		WHERE status = 'available'
// 		ORDER BY created_at DESC
// 	`)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	apartments := make([]types.Apartment, 0)
// 	for rows.Next() {
// 		apt, err := scanRowsIntoApartment(rows)
// 		if err != nil {
// 			return nil, err
// 		}
// 		apartments = append(apartments, *apt)
// 	}

// 	return apartments, nil
// }

// // GetApartments returns all apartments (admin view)
// func (s *Store) GetApartments() ([]types.Apartment, error) {
// 	rows, err := s.db.Query(`
// 		SELECT id, code, name, rooms, description, price, status, created_at
// 		FROM apartments
// 		ORDER BY created_at DESC
// 	`)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	apartments := make([]types.Apartment, 0)
// 	for rows.Next() {
// 		apt, err := scanRowsIntoApartment(rows)
// 		if err != nil {
// 			return nil, err
// 		}
// 		apartments = append(apartments, *apt)
// 	}

// 	return apartments, nil
// }

// func (s *Store) GetApartmentByCode(code string) (*types.Apartment, error) {
// 	rows, err := s.db.Query(`
// 		SELECT id, code, name, rooms, description, price, status, created_at
// 		FROM apartments
// 		WHERE code = $1
// 	`, code)

// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	apt := new(types.Apartment)
// 	for rows.Next() {
// 		apt, err = scanRowsIntoApartment(rows)
// 		if err != nil {
// 			return nil, err
// 		}
// 	}

// 	if apt.ID == 0 {
// 		return nil, fmt.Errorf("Apartment not found")
// 	}

// 	return apt, nil
// }

// func (s *Store) CreateApartment(ap types.Apartment) (types.Apartment, error) {
// 	err := s.db.QueryRow(`
// 		INSERT INTO apartments (code, name, rooms, description, price)
// 		VALUES ($1, $2, $3, $4, $5)
// 		RETURNING id
// 	`, ap.Code, ap.Name, ap.Rooms, ap.Description, ap.Price).Scan(&ap.ID)

// 	if err != nil {
// 		return types.Apartment{}, err
// 	}

// 	return ap, nil
// }

// func (s *Store) DeleteApartment(id int) (types.Apartment, error) {
// 	var apartment types.Apartment

// 	// First, get the apartment info before deleting
// 	err := s.db.QueryRow(`
// 		SELECT id, code, name, rooms, description, price, status, created_at
// 		FROM apartments
// 		WHERE id = $1
// 	`, id).Scan(
// 		&apartment.ID,
// 		&apartment.Code,
// 		&apartment.Name,
// 		&apartment.Rooms,
// 		&apartment.Description,
// 		&apartment.Price,
// 		&apartment.Status,
// 		&apartment.CreatedAt,
// 	)

// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return types.Apartment{}, fmt.Errorf("apartment not found")
// 		}
// 		return types.Apartment{}, err
// 	}

// 	// Delete the apartment
// 	_, err = s.db.Exec(`DELETE FROM apartments WHERE id = $1`, id)
// 	if err != nil {
// 		return types.Apartment{}, err
// 	}

// 	return apartment, nil
// }

// func scanRowsIntoApartment(rows *sql.Rows) (*types.Apartment, error) {
// 	apartment := new(types.Apartment)

// 	err := rows.Scan(
// 		&apartment.ID,
// 		&apartment.Code,
// 		&apartment.Name,
// 		&apartment.Rooms,
// 		&apartment.Description,
// 		&apartment.Price,
// 		&apartment.Status,
// 		&apartment.CreatedAt,
// 	)

// 	if err != nil {
// 		return nil, err
// 	}

// 	return apartment, nil
// }

package aprtment

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

// Get apartments publicly

func (s *Store) GetPublicApartments() ([]types.Apartment, error) {
	rows, err := s.db.Query(getPublicApartmentsQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanApartmentRows(rows)
}

func (s *Store) GetApartments() ([]types.Apartment, error) {
	rows, err := s.db.Query(getApartmentsQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanApartmentRows(rows)
}

func (s *Store) GetApartmentByCode(code string) (*types.Apartment, error) {
	rows, err := s.db.Query(getApartmentByCodeQuery, code)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		return scanRowsIntoApartment(rows)
	}
	return nil, fmt.Errorf("apartment not found")
}

func (s *Store) GetApartmentByID(id int) (*types.Apartment, error) {
	rows, err := s.db.Query(getApartmentByIDQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		return scanRowsIntoApartment(rows)
	}
	return nil, fmt.Errorf("apartment not found")
}

func (s *Store) CreateApartment(ap types.Apartment) (types.Apartment, error) {
	var created types.Apartment
	err := s.db.QueryRow(createApartmentQuery,
		ap.Code, ap.Name, ap.Rooms, ap.Description, ap.Price,
	).Scan(
		&created.ID,
		&created.Code,
		&created.Name,
		&created.Rooms,
		&created.Description,
		&created.Price,
		&created.Status,
		&created.CreatedAt,
	)
	if err != nil {
		return types.Apartment{}, err
	}
	return created, nil
}

// UpdateApartment applies only the non-zero fields from the payload.
// It builds the SET clause dynamically so callers can do partial updates.
func (s *Store) UpdateApartment(id int, p types.UpdateApartmentPayload) (types.Apartment, error) {
	// Ensure the apartment exists first
	existing, err := s.GetApartmentByID(id)
	if err != nil {
		return types.Apartment{}, err
	}

	// Merge: keep existing values when payload field is zero-value
	if p.Name != "" {
		existing.Name = p.Name
	}
	if p.Rooms > 0 {
		existing.Rooms = p.Rooms
	}
	if p.Description != "" {
		existing.Description = p.Description
	}
	if p.Price > 0 {
		existing.Price = p.Price
	}
	if p.Status != "" {
		existing.Status = p.Status
	}

	const updateApartmentQuery = `
		UPDATE apartments
		SET name = $1, rooms = $2, description = $3, price = $4, status = $5
		WHERE id = $6
		RETURNING id, code, name, rooms, description, price, status, created_at
	`

	var updated types.Apartment
	err = s.db.QueryRow(updateApartmentQuery,
		existing.Name,
		existing.Rooms,
		existing.Description,
		existing.Price,
		existing.Status,
		id,
	).Scan(
		&updated.ID,
		&updated.Code,
		&updated.Name,
		&updated.Rooms,
		&updated.Description,
		&updated.Price,
		&updated.Status,
		&updated.CreatedAt,
	)
	if err != nil {
		return types.Apartment{}, err
	}
	return updated, nil
}

func (s *Store) DeleteApartment(id int) (types.Apartment, error) {
	existing, err := s.GetApartmentByID(id)
	if err != nil {
		return types.Apartment{}, err
	}

	_, err = s.db.Exec(deleteApartmentQuery, id)
	if err != nil {
		return types.Apartment{}, err
	}
	return *existing, nil
}





// Helper
func scanApartmentRows(rows *sql.Rows) ([]types.Apartment, error) {
	apartments := make([]types.Apartment, 0)
	for rows.Next() {
		apt, err := scanRowsIntoApartment(rows)
		if err != nil {
			return nil, err
		}
		apartments = append(apartments, *apt)
	}
	return apartments, nil
}

func scanRowsIntoApartment(rows *sql.Rows) (*types.Apartment, error) {
	a := new(types.Apartment)
	err := rows.Scan(
		&a.ID,
		&a.Code,
		&a.Name,
		&a.Rooms,
		&a.Description,
		&a.Price,
		&a.Status,
		&a.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return a, nil
}

// Get all apartments publicly
const getPublicApartmentsQuery = ` 
   SELECT id, code, name, rooms, description, price, status, created_at
   FROM apartments
   WHERE status = 'available'
   ORDER BY created_at DESC
 `
// Get all apartments for admin
 const getApartmentsQuery = ` 
     SELECT id, code, name, rooms, description, price, status, created_at
	 FROM apartments
	 ORDER BY created_at DEC
  `

//   get apartment by code

const getApartmentByCodeQuery = ` 
   SELECT id, code, name, rooms, description, price, status, created_at
   FROM apartments
   WHERE code = $1
 `

//  Get apartment by id query

const getApartmentByIDQuery = `
   SELECT id, code, name, rooms, description, price, status, created_at
   FROM apartments
   WHERE id = $1 
 `
const createApartmentQuery = ` 
  INSERT INTO apartments (code, name, rooms, description, price)
  VALUES (S1, $2, $3, $4, $5)
  RETURNING id, code, name, rooms, decsription, price, status, created_at
 `

 const deleteApartmentQuery = ` 
 DELETE FROM apartments WHERE id = $1
  `
