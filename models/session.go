package models

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"

	"github.com/iraqnroll/gochan/rand"
)

const (
	MIN_BYTES_PER_TOKEN = 32
)

type Session struct {
	Id     int
	UserId int

	//WARNING: Raw token value will only be set and accessible on the creation of a new active session.
	//When looking up an already existing session, only the hashed token will be returned and Token will be empty.
	//We only store a hash of the token in the database due to security reasons.
	Token     string
	TokenHash string
}

type SessionService struct {
	DB            *sql.DB
	BytesPerToken int
}

func (ss *SessionService) CreateSession(userId int) (*Session, error) {
	bytesPerToken := ss.BytesPerToken
	if bytesPerToken < MIN_BYTES_PER_TOKEN {
		bytesPerToken = MIN_BYTES_PER_TOKEN
	}

	token, err := rand.String(bytesPerToken)
	if err != nil {
		return nil, fmt.Errorf("session.CreateSession error : %w", err)
	}

	session := Session{
		UserId:    userId,
		Token:     token,
		TokenHash: ss.HashToken(token),
	}
	//Check if user has any existing active sessions, if he does - update current record with a new token.
	//TODO: Add session expiration logic, we dont want to keep sessions active indefinitely...

	row := ss.DB.QueryRow(`
		INSERT INTO sessions (user_id, token_hash)
		VALUES ($1, $2) ON CONFLICT (user_id) DO UPDATE
		SET token_hash = $2
		WHERE sessions.user_id = $1
		RETURNING id;`, session.UserId, session.TokenHash)

	err = row.Scan(&session.Id)
	if err != nil {
		return nil, fmt.Errorf("SessionService.Create error : %w", err)
	}

	return &session, nil
}

func (ss *SessionService) GetUser(token string) (*User, error) {
	tokenHash := ss.HashToken(token)
	var user User
	row := ss.DB.QueryRow(`
		SELECT users.id,
			users.username,
			users.password_hash
		FROM sessions
		INNER JOIN users ON users.id = sessions.user_id
		WHERE sessions.token_hash = $1`, tokenHash)

	err := row.Scan(&user.Id, &user.Username, &user.password_hash)
	if err != nil {
		return nil, fmt.Errorf("SessionService.GetUser error : %w", err)
	}

	return &user, nil
}

func (ss *SessionService) HashToken(token string) string {
	tokenHash := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(tokenHash[:])
}

func (ss *SessionService) DeleteSession(token string) error {
	tokenHash := ss.HashToken(token)

	//TODO: Refactor to not filter by NVARCHAR/TEXT columns (it's slow)
	_, err := ss.DB.Exec(`
		DELETE FROM sessions
		WHERE token_hash = $1;`, tokenHash)
	if err != nil {
		return fmt.Errorf("SessionService.DeleteService error : %w", err)
	}
	return nil
}
