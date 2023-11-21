package services

import (
	"context"
	"errors"
	"technical-interview/pkg/dao"
	"technical-interview/pkg/models"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidEntity   = errors.New("invalid entity")
	ErrMissingEmail    = errors.New("missing email")
	ErrMissingPassword = errors.New("missing password")
	ErrMissingUsername = errors.New("missing username")
)

type RegisterService interface {
	Exec(ctx context.Context, email string, password string, username string) (*models.User, *models.TokenIntrospection, error)
}

func NewRegisterService(repository dao.UserRepository, generateToken GenerateTokenService) RegisterService {
	return &registerServiceImpl{
		repository:    repository,
		generateToken: generateToken,
	}
}

type registerServiceImpl struct {
	repository    dao.UserRepository
	generateToken GenerateTokenService
}

func (s *registerServiceImpl) Exec(ctx context.Context, email string, password string, username string) (*models.User, *models.TokenIntrospection, error) {
	if email == "" {
		return nil, nil, errors.Join(ErrInvalidEntity, ErrMissingEmail)
	}
	if password == "" {
		return nil, nil, errors.Join(ErrInvalidEntity, ErrMissingPassword)
	}
	if username == "" {
		return nil, nil, errors.Join(ErrInvalidEntity, ErrMissingUsername)
	}

	user, err := s.repository.Create(ctx, email, password, username)
	if err != nil {
		return nil, nil, err
	}

	token, err := s.generateToken.GenerateToken(models.UserTokenPayload{ID: user.ID}, uuid.New(), time.Now())
	if err != nil {
		return nil, nil, err
	}

	return user, token, nil
}
