package repos

import (
	"database/sql"
	"fmt"

	"github.com/iraqnroll/gochan/models"
)

const (
	s_create_new_query = `INSERT INTO sessions (user_id, token_hash)
		VALUES ($1, $2) ON CONFLICT (user_id) DO UPDATE
		SET token_hash = $2
		WHERE sessions.user_id = $1 RETURNING id;`
)

type PostgresSessionRepository struct {
	db *sql.DB
}

func NewPostgresSessionRepository(db *sql.DB) *PostgresSessionRepository {
	if db == nil {
		panic("Missing db")
	}

	return &PostgresSessionRepository{db: db}
}

func (sr *PostgresSessionRepository) CreateSession(user_id int, hashed_token string) (*models.Session, error) {
	result := models.Session{UserId: user_id, TokenHash: hashed_token}
	row := sr.db.QueryRow(s_create_new_query, user_id, hashed_token)
	err := row.Scan(&result.Id)
	if err != nil {
		return nil, fmt.Errorf("PostgresSessionRepository.CreateSession failed : %w", err)
	}
	return &result, nil
}
