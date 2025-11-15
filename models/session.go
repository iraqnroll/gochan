package models

type Session struct {
	Id     int `json:"id"`
	UserId int `json:"user_id"`
	//WARNING: Raw token value will only be set and accessible on the creation of a new active session.
	//When looking up an already existing session, only the hashed token will be returned and Token will be empty.
	//We only store a hash of the token in the database due to security reasons.
	Token     string `json:"token"`
	TokenHash string `json:"token_hash"`
}
