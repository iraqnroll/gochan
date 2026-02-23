package models

type User struct {
	Id            int    `db:"id" json:"id" schema:"id"`
	Username      string `db:"username" json:"username" schema:"username"`
	Password_hash string `db:"password_hash" json:"password_hash" schema:"password_hash"`
	Email         string `db:"email" json:"email" schema:"email"`
	Date_created  string `db:"date_created" json:"date_created" schema:"date_created"`
	Date_updated  string `db:"date_updated" json:"date_updated" schema:"date_updated"`
	User_type     int    `db:"user_type" json:"user_type" schema:"user_type"`
	Password      string `schema:"password"`
}
