package models

import (
	"database/sql"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id            int
	Username      string
	password_hash string
	Email         string
	Date_created  string
	Date_updated  string
	User_type     int
}

type UserService struct {
	DB *sql.DB
}

func (us *UserService) Create(username, password, email string, user_type int) (*User, error) {
	//Postgres is case-sensitive, so convert sensitive strings to lowercase.
	email = strings.ToLower(email)
	username = strings.ToLower(username)

	//Hash the password before storing in DB
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("UserService.Create failed : %w", err)
	}

	passwordHash := string(hashedBytes)

	user := User{
		Email:         email,
		password_hash: passwordHash,
	}

	row := us.DB.QueryRow(`
		INSERT INTO users (username, password_hash, email, user_type)
		VALUES ($1, $2, $3, $4) RETURNING id`, username, passwordHash, email, user_type)

	err = row.Scan(&user.Id)
	if err != nil {
		return nil, fmt.Errorf("UserService.Create failed : %w", err)
	}

	//Return user obj. with a newly created DB record id.
	return &user, nil
}

func (us *UserService) Delete(userId int) error {
	_, err := us.DB.Exec(`DELETE FROM users WHERE id = $1`, userId)
	if err != nil {
		return fmt.Errorf("UserService.Delete failed : %w", err)
	}
	return nil
}

func (us *UserService) Authenticate(username, password string) (*User, error) {
	//Postgres is case-sensitive, so convert sensitive strings to lowercase.
	username = strings.ToLower(username)

	user := User{
		Username: username,
	}

	row := us.DB.QueryRow(`
	SELECT id, password_hash FROM users WHERE username=$1`, username)

	err := row.Scan(&user.Id, &user.password_hash)
	if err != nil {
		return nil, fmt.Errorf("UserService.Authenticate failed : %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.password_hash), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("UserService.Authenticate failed : %w", err)
	}

	return &user, nil
}

func (us *UserService) GetUserList() ([]User, error) {
	//User slice to hold data from returned rows
	var result []User

	rows, err := us.DB.Query(`SELECT 
		id,
		username,
		email,
		COALESCE(to_char(date_created, 'YYYY-MM-DD HH24:MI:SS'), 'Never') AS date_created,
		COALESCE(to_char(date_updated, 'YYYY-MM-DD HH24:MI:SS'), 'Never') AS date_updated,
		user_type FROM users`)
	if err != nil {
		return nil, fmt.Errorf("UserService.GetUserList failed : %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var usr User
		err := rows.Scan(&usr.Id, &usr.Username, &usr.Email, &usr.Date_created, &usr.Date_updated, &usr.User_type)
		if err != nil {
			fmt.Println("UserService.GetUserList loop failed : %w", err)
		}
		result = append(result, usr)
	}

	if err = rows.Err(); err != nil {
		return result, fmt.Errorf("UserService.GetUserList failed : %w", err)
	}

	return result, nil
}
