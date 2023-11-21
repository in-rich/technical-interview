package services

import (
	"context"
	"technical-interview/pkg/dao"
	"technical-interview/pkg/models"
)

type GetUserService interface {
	Exec(ctx context.Context, email string) (*models.User, error)
}

func NewGetUserService(repository dao.UserRepository) GetUserService {
	return &getUserServiceImpl{
		repository: repository,
	}
}

type getUserServiceImpl struct {
	repository dao.UserRepository
}

func (s *getUserServiceImpl) Exec(ctx context.Context, email string) (*models.User, error) {
	user, err := s.repository.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return user, nil
}
