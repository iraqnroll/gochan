package repos

import (
	"database/sql"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/iraqnroll/gochan/models"
)

type PostgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	if db == nil {
		panic("Missing db")
	}

	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) dbInstance() *goqu.Database {
	return pgDialect.DB(r.db)
}

func (r *PostgresUserRepository) CreateNew(username, password_hash, email string, user_type int) (*models.User, error) {
	var result models.User

	_, err := r.dbInstance().Insert("users").
		Cols("username", "password_hash", "email", "user_type").
		Vals([]interface{}{username, password_hash, email, user_type}).
		Returning("id").
		Executor().
		ScanStruct(&result)
	if err != nil {
		return nil, fmt.Errorf("PostgresUserRepository.CreateNew error: %w", err)
	}

	result.Username = username
	result.Password_hash = password_hash
	result.Email = email
	result.User_type = user_type
	return &result, nil
}

func (r *PostgresUserRepository) Delete(user_id int) error {
	_, err := r.dbInstance().Delete("users").
		Where(goqu.C("id").Eq(user_id)).
		Executor().Exec()
	if err != nil {
		return fmt.Errorf("PostgresUserRepository.Delete error: %w", err)
	}
	return nil
}

func (r *PostgresUserRepository) GetAll() ([]models.User, error) {
	var result []models.User

	err := r.dbInstance().From("users").
		Select(
			"id",
			"username",
			"email",
			goqu.L("COALESCE(to_char(date_created, 'YYYY-MM-DD HH24:MI:SS'), 'Never')").As("date_created"),
			goqu.L("COALESCE(to_char(date_updated, 'YYYY-MM-DD HH24:MI:SS'), 'Never')").As("date_updated"),
			"user_type",
		).
		Order(goqu.C("date_updated").Desc()).
		ScanStructs(&result)
	if err != nil {
		return nil, fmt.Errorf("PostgresUserRepository.GetAll error: %w", err)
	}

	return result, nil
}

func (r *PostgresUserRepository) GetUserById(user_id int) (*models.User, error) {
	var result models.User

	found, err := r.dbInstance().From("users").
		Select(
			"username",
			"password_hash",
			"email",
			goqu.L("COALESCE(to_char(date_created, 'YYYY-MM-DD HH24:MI:SS'), 'Never')").As("date_created"),
			goqu.L("COALESCE(to_char(date_updated, 'YYYY-MM-DD HH24:MI:SS'), 'Never')").As("date_updated"),
			"user_type",
		).
		Where(goqu.C("id").Eq(user_id)).
		ScanStruct(&result)
	if err != nil {
		return nil, fmt.Errorf("PostgresUserRepository.GetUserById error: %w", err)
	}
	if !found {
		return nil, fmt.Errorf("PostgresUserRepository.GetUserById: user not found")
	}

	result.Id = user_id
	return &result, nil
}

func (r *PostgresUserRepository) GetPwHashByUsername(username string) (*models.User, error) {
	var result models.User

	found, err := r.dbInstance().From("users").
		Select("id", "password_hash").
		Where(goqu.C("username").Eq(username)).
		ScanStruct(&result)
	if err != nil {
		return nil, fmt.Errorf("PostgresUserRepository.GetPwHashByUsername error: %w", err)
	}
	if !found {
		return nil, fmt.Errorf("PostgresUserRepository.GetPwHashByUsername: user not found")
	}

	result.Username = username
	return &result, nil
}

func (r *PostgresUserRepository) UpdateUser(user_id, user_type int, username, password_hash, email string) (*models.User, error) {
	_, err := r.dbInstance().Update("users").
		Set(goqu.Record{
			"username":      username,
			"password_hash": password_hash,
			"email":         email,
			"date_updated":  goqu.L("NOW()"),
			"user_type":     user_type,
		}).
		Where(goqu.C("id").Eq(user_id)).
		Executor().Exec()
	if err != nil {
		return nil, fmt.Errorf("PostgresUserRepository.UpdateUser error: %w", err)
	}

	result := &models.User{
		Id:            user_id,
		Username:      username,
		Password_hash: password_hash,
		Email:         email,
		User_type:     user_type,
	}
	return result, nil
}
