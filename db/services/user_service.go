package services

import (
	"fmt"
	"strings"

	"github.com/iraqnroll/gochan/db/models"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	CreateNew(username, password_hash, email string, user_type int) (*models.User, error)
	Delete(user_id int) error
	GetAll() ([]models.User, error)
	GetPwHashByUsername(username string) (*models.User, error)
	GetUserById(user_id int) (*models.User, error)
	UpdateUser(user_id, user_type int, username, password_hash, email string) (*models.User, error)
}

type UserService struct {
	userRepo UserRepository
}

func NewUserService(uRepo UserRepository) *UserService {
	return &UserService{userRepo: uRepo}
}

func (us *UserService) GetAll() ([]models.User, error) {
	result, err := us.userRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("UserService.GetAll failed : %w", err)
	}
	return result, nil
}

func (us *UserService) CreateNew(username, password, email string, user_type int) (*models.User, error) {
	//Postgres is case-sensitive, so convert sensitive strings to lowercase.
	email = strings.ToLower(email)
	username = strings.ToLower(username)

	//Hash the password before storing in DB
	hashedBytes, err := us.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("UserService.CreateNew failed : %w", err)
	}

	passwordHash := string(hashedBytes)
	result, err := us.userRepo.CreateNew(username, passwordHash, email, user_type)
	if err != nil {
		return nil, fmt.Errorf("UserService.CreateNew failed : %w", err)
	}
	return result, nil
}

func (us *UserService) Delete(user_id int) error {
	err := us.userRepo.Delete(user_id)
	if err != nil {
		return fmt.Errorf("UserService.Delete failed : %w", err)
	}
	return nil
}

func (us *UserService) GetUserById(user_id int) (*models.User, error) {
	result, err := us.userRepo.GetUserById(user_id)
	if err != nil {
		return nil, fmt.Errorf("UserService.GetUserById failed : %w", err)
	}
	return result, nil
}

func (us *UserService) Authenticate(username, password string) (*models.User, error) {
	username = strings.ToLower(username)
	result, err := us.userRepo.GetPwHashByUsername(username)
	if err != nil {
		return nil, fmt.Errorf("UserService.Authenticate failed: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(result.Password_hash), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("UserService.Authenticate failed : %w", err)
	}

	return result, nil
}

func (us *UserService) HashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func (us *UserService) UpdateUser(user_id, user_type int, username, password_hash, email string) (*models.User, error) {
	result, err := us.userRepo.UpdateUser(user_id, user_type, username, password_hash, email)
	if err != nil {
		return nil, fmt.Errorf("UserService.UpdateUser failed: %w", err)
	}

	return result, nil
}
