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
	s_get_by_token_query = `SELECT users.id, users.username, users.password_hash
		FROM sessions
		INNER JOIN users ON users.id = sessions.user_id
		WHERE sessions.token_hash = $1`
	s_delete_session = `DELETE FROM sessions WHERE token_hash = $1`
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

func (sr *PostgresSessionRepository) CreateNew(user_id int, hashed_token string) (*models.Session, error) {
	result := models.Session{UserId: user_id, TokenHash: hashed_token}
	row := sr.db.QueryRow(s_create_new_query, user_id, hashed_token)
	err := row.Scan(&result.Id)
	if err != nil {
		return nil, fmt.Errorf("PostgresSessionRepository.CreateSession failed : %w", err)
	}
	return &result, nil
}

func (sr *PostgresSessionRepository) GetUserByToken(hashed_token string) (*models.User, error) {
	result := models.User{}
	row := sr.db.QueryRow(s_get_by_token_query, hashed_token)
	err := row.Scan(&result.Id, &result.Username, &result.Password_hash)
	if err != nil {
		return nil, fmt.Errorf("PostgresSessionRepository.GetUserByToken : %w", err)
	}
	return &result, nil
}

func (sr *PostgresSessionRepository) DeleteSession(hashed_token string) error {
	_, err := sr.db.Exec(s_delete_session, hashed_token)
	if err != nil {
		return fmt.Errorf("PostgresSessionRepository.DeleteSession error : %w", err)
	}
	return nil
}
