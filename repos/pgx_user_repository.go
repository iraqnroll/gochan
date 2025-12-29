package repos

import (
	"database/sql"
	"fmt"

	"github.com/iraqnroll/gochan/models"
)

const (
	u_create_new_query = `INSERT INTO users (username, password_hash, email, user_type) VALUES ($1, $2, $3, $4) RETURNING id`
	u_delete_query     = `DELETE FROM users WHERE id = $1`
	u_get_all_query    = `SELECT 
		id,
		username,
		email,
		COALESCE(to_char(date_created, 'YYYY-MM-DD HH24:MI:SS'), 'Never') AS date_created,
		COALESCE(to_char(date_updated, 'YYYY-MM-DD HH24:MI:SS'), 'Never') AS date_updated,
		user_type FROM users`
	u_get_pw_hash_query = `SELECT id, password_hash FROM users WHERE username = $1`
	u_get_by_id_query   = `SELECT 
		username,
		password_hash,
		email,
		COALESCE(to_char(date_created, 'YYYY-MM-DD HH24:MI:SS'), 'Never') AS date_created,
		COALESCE(to_char(date_updated, 'YYYY-MM-DD HH24:MI:SS'), 'Never') AS date_updated,
		usertype FROM users WHERE id = $1`
	u_update_query = `UPDATE users SET
		username = $2,
		password_hash = $3,
		email = $4,
		user_type = $5
		WHERE id = $1`
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

func (r *PostgresUserRepository) CreateNew(username, password_hash, email string, user_type int) (*models.User, error) {
	result := models.User{Username: username, Email: email, User_type: user_type, Password_hash: password_hash}
	row := r.db.QueryRow(u_create_new_query, result.Username, password_hash, result.Email, result.User_type)
	err := row.Scan(&result.Id)
	if err != nil {
		return nil, fmt.Errorf("PostgresUserRepository.CreateNew failed : %w", err)
	}

	return &result, nil
}

func (r *PostgresUserRepository) Delete(user_id int) error {
	_, err := r.db.Exec(u_delete_query, user_id)
	if err != nil {
		return fmt.Errorf("PostgresUserRepository.Delete failed : %w", err)
	}
	return nil
}

func (r *PostgresUserRepository) GetAll() ([]models.User, error) {
	var result []models.User
	rows, err := r.db.Query(u_get_all_query)
	if err != nil {
		return nil, fmt.Errorf("PostgresUserRepository.GetAll failed : %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var usr models.User
		err := rows.Scan(&usr.Id, &usr.Username, &usr.Email, &usr.Date_created, &usr.Date_updated, &usr.User_type)
		if err != nil {
			return result, fmt.Errorf("PostgresUserRepository.GetAll failed : %w", err)
		}
		result = append(result, usr)
	}

	return result, nil
}

func (r *PostgresUserRepository) GetUserById(user_id int) (*models.User, error) {
	result := models.User{Id: user_id}
	row := r.db.QueryRow(u_get_by_id_query, user_id)
	err := row.Scan(&result.Username, &result.Password_hash, &result.Email, &result.Date_created, &result.Date_updated, &result.User_type)
	if err != nil {
		return nil, fmt.Errorf("PostgresUserRepository.GetById failed : %w", err)
	}
	return &result, nil
}

func (r *PostgresUserRepository) GetPwHashByUsername(username string) (*models.User, error) {
	result := models.User{Username: username}
	row := r.db.QueryRow(u_get_pw_hash_query, username)
	err := row.Scan(&result.Id, &result.Password_hash)
	if err != nil {
		return nil, fmt.Errorf("PostgresUserRepository.GetPwHashByUsername failed : %w", err)
	}

	return &result, nil
}

func (r *PostgresUserRepository) UpdateUser(user_id, user_type int, username, password_hash, email string) (*models.User, error) {
	result := models.User{
		Id:            user_id,
		Username:      username,
		Password_hash: password_hash,
		Email:         email,
		User_type:     user_type,
	}

	_, err := r.db.Exec(u_update_query, result.Id, result.Username, result.Password_hash, result.Email, result.User_type)
	if err != nil {
		return nil, fmt.Errorf("PostgresUserRepository.UpdateUser failed : %w", err)
	}
	return &result, nil
}
