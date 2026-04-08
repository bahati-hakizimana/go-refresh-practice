package comment

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

//
// =====================
// Queries
// =====================
//

const getCommentsQuery = `
	SELECT id, name, comment, created_at
	FROM comments
	ORDER BY created_at DESC
`

const getCommentByIDQuery = `
	SELECT id, name, comment, created_at
	FROM comments
	WHERE id = $1
`

const createCommentQuery = `
	INSERT INTO comments (name, comment, created_at)
	VALUES ($1, $2, NOW())
	RETURNING id, created_at
`

const deleteCommentQuery = `
	DELETE FROM comments
	WHERE id = $1
`

//
// =====================
// Methods
// =====================
//

func (s *Store) GetComments() ([]types.Comment, error) {
	rows, err := s.db.Query(getCommentsQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := make([]types.Comment, 0)

	for rows.Next() {
		var c types.Comment

		err := rows.Scan(
			&c.ID,
			&c.Name,
			&c.Comment,
			&c.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		comments = append(comments, c)
	}

	return comments, nil
}

func (s *Store) GetCommentByID(id int) (*types.Comment, error) {
	row := s.db.QueryRow(getCommentByIDQuery, id)

	var c types.Comment

	err := row.Scan(
		&c.ID,
		&c.Name,
		&c.Comment,
		&c.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("comment not found")
		}
		return nil, err
	}

	return &c, nil
}

func (s *Store) CreateComment(comment types.Comment) (types.Comment, error) {

	err := s.db.QueryRow(
		createCommentQuery,
		comment.Name,
		comment.Comment,
	).Scan(&comment.ID, &comment.CreatedAt)

	if err != nil {
		return types.Comment{}, fmt.Errorf("failed to create comment: %w", err)
	}

	return comment, nil
}


func (s *Store) DeleteComment(id int) error {
	result, err := s.db.Exec(deleteCommentQuery, id)
	if err != nil {
		return fmt.Errorf("failed to delete comment: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("comment not found")
	}

	return nil
}