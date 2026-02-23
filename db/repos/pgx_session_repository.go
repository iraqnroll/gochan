package repos

import (
	"database/sql"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/iraqnroll/gochan/db/models"
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

func (sr *PostgresSessionRepository) dbInstance() *goqu.Database {
	return pgDialect.DB(sr.db)
}

func (sr *PostgresSessionRepository) CreateNew(user_id int, hashed_token string) (*models.Session, error) {
	var result models.Session

	_, err := sr.dbInstance().Insert("sessions").
		Cols("user_id", "token_hash").
		Vals([]interface{}{user_id, hashed_token}).
		OnConflict(goqu.DoUpdate("user_id", goqu.Record{"token_hash": hashed_token})).
		Returning("id").
		Executor().
		ScanStruct(&result)
	if err != nil {
		return nil, fmt.Errorf("PostgresSessionRepository.CreateNew error: %w", err)
	}

	result.UserId = user_id
	result.TokenHash = hashed_token
	return &result, nil
}

func (sr *PostgresSessionRepository) GetUserByToken(hashed_token string) (*models.User, error) {
	var result models.User

	found, err := sr.dbInstance().From("sessions").
		Select(goqu.I("users.id"), goqu.I("users.username"), goqu.I("users.password_hash")).
		Join(goqu.T("users"), goqu.On(goqu.I("users.id").Eq(goqu.I("sessions.user_id")))).
		Where(goqu.I("sessions.token_hash").Eq(hashed_token)).
		ScanStruct(&result)
	if err != nil {
		return nil, fmt.Errorf("PostgresSessionRepository.GetUserByToken error: %w", err)
	}
	if !found {
		return nil, fmt.Errorf("PostgresSessionRepository.GetUserByToken: session not found")
	}
	return &result, nil
}

func (sr *PostgresSessionRepository) DeleteSession(hashed_token string) error {
	_, err := sr.dbInstance().Delete("sessions").
		Where(goqu.C("token_hash").Eq(hashed_token)).
		Executor().Exec()
	if err != nil {
		return fmt.Errorf("PostgresSessionRepository.DeleteSession error: %w", err)
	}
	return nil
}
