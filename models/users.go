package models

type User struct {
	Id            int    `json:"id"`
	Username      string `json:"username"`
	Password_hash string `json:"password_hash"`
	Email         string `json:"email"`
	Date_created  string `json:"date_created"`
	Date_updated  string `json:"date_updated"`
	User_type     int    `json:"user_type"`
	Password      string
}
