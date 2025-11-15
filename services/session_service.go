package services

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"

	"github.com/iraqnroll/gochan/models"
	"github.com/iraqnroll/gochan/rand"
)

const (
	MIN_BYTES_PER_TOKEN = 32
)

type SessionRepository interface {
	CreateNew(user_id int, hashed_token string) (*models.Session, error)
	GetUserByToken(hashed_token string) (*models.User, error)
	DeleteSession(hashed_token string) error
}

type SessionService struct {
	SessionRepo   SessionRepository
	BytesPerToken int
}

func NewSessionService(repo SessionRepository, bytesPerToken int) *SessionService {
	return &SessionService{SessionRepo: repo, BytesPerToken: bytesPerToken}
}

func (ss *SessionService) CreateNew(user_id int) (*models.Session, error) {
	bytesPerToken := ss.BytesPerToken
	if bytesPerToken < MIN_BYTES_PER_TOKEN {
		bytesPerToken = MIN_BYTES_PER_TOKEN
	}

	token, err := rand.String(bytesPerToken)
	if err != nil {
		return nil, fmt.Errorf("SessionService.CreateNew error : %w", err)
	}
	result, err := ss.SessionRepo.CreateNew(user_id, ss.HashToken(token))
	if err != nil {
		return nil, fmt.Errorf("SessionService.CreateNew error : %w", err)
	}
	return result, nil
}

func (ss *SessionService) GetUser(token string) (*models.User, error) {
	result, err := ss.SessionRepo.GetUserByToken(ss.HashToken(token))
	if err != nil {
		return nil, fmt.Errorf("SessionService.GetUserByToken error : %w", err)
	}
	return result, nil
}

func (ss *SessionService) DeleteSession(token string) error {
	err := ss.SessionRepo.DeleteSession(ss.HashToken(token))
	if err != nil {
		return fmt.Errorf("SessionService.DeleteSession error : %w", err)
	}
	return nil
}

func (ss *SessionService) HashToken(token string) string {
	tokenHash := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(tokenHash[:])
}
