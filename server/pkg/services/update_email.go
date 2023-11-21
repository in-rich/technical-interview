package services

import (
	"context"
	"errors"
	"technical-interview/pkg/dao"
	"time"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type UpdateEmailService interface {
	UpdateEmail(ctx context.Context, tokenRaw string, email string) error
}

func NewUpdateEmailService(repository dao.UserRepository, getTokenStatus GetTokenStatusService) UpdateEmailService {
	return &updateEmailServiceImpl{
		getTokenStatus: getTokenStatus,
		repository:     repository,
	}
}

type updateEmailServiceImpl struct {
	getTokenStatus GetTokenStatusService
	repository     dao.UserRepository
}

func (s *updateEmailServiceImpl) UpdateEmail(ctx context.Context, tokenRaw string, email string) error {
	tokenStatus, err := s.getTokenStatus.GetTokenStatus(tokenRaw, time.Now())
	if err != nil {
		return err
	}
	if !tokenStatus.OK {
		return ErrInvalidCredentials
	}

	if err := s.repository.UpdateEmail(ctx, tokenStatus.Token.Payload.ID, email); err != nil {
		return err
	}

	return nil
}
