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
	email         string
	date_created  string
	date_updated  string
	user_type     int
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
		email:         email,
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
