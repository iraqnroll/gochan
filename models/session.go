package models

type Session struct {
	Id        int    `db:"id" json:"id" schema:"id"`
	UserId    int    `db:"user_id" json:"user_id" schema:"user_id"`
	Token     string `db:"token" json:"token" schema:"token"`
	TokenHash string `db:"token_hash" json:"token_hash" schema:"token_hash"`
}
