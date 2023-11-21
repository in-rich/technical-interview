package services

import (
	"context"
	"errors"
	"technical-interview/pkg/dao"
	"technical-interview/pkg/models"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidPassword = errors.New("invalid password")
)

type LoginService interface {
	Exec(ctx context.Context, email string, password string) (*models.User, *models.TokenIntrospection, error)
}

func NewLoginService(repository dao.UserRepository, generateToken GenerateTokenService) LoginService {
	return &loginServiceImpl{
		repository:    repository,
		generateToken: generateToken,
	}
}

type loginServiceImpl struct {
	repository    dao.UserRepository
	generateToken GenerateTokenService
}

func (s *loginServiceImpl) Exec(ctx context.Context, email string, password string) (*models.User, *models.TokenIntrospection, error) {
	user, err := s.repository.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return nil, nil, ErrInvalidPassword
		}

		return nil, nil, err
	}

	token, err := s.generateToken.GenerateToken(models.UserTokenPayload{ID: user.ID}, uuid.New(), time.Now())
	if err != nil {
		return nil, nil, err
	}

	return user, token, nil
}
