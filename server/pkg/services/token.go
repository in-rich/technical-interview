package services

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"technical-interview/pkg/models"
	"time"

	"github.com/google/uuid"
)

var (
	ErrEncodeTokenHeader  = errors.New("unable to encode token header")
	ErrEncodeTokenPayload = errors.New("unable to encode token payload")
	ErrInvalidToken       = errors.New("invalid token")
	ErrBadSignature       = errors.New("bad signature")
)

type GenerateTokenService interface {
	// GenerateToken creates a new, valid token for the given user.
	GenerateToken(data models.UserTokenPayload, id uuid.UUID, now time.Time) (*models.TokenIntrospection, error)
}

func NewGenerateTokenService(tokenTTL time.Duration) GenerateTokenService {
	return &generateTokenServiceImpl{
		tokenTTL: tokenTTL,
	}
}

type generateTokenServiceImpl struct {
	tokenTTL time.Duration
}

func (s *generateTokenServiceImpl) GenerateToken(data models.UserTokenPayload, id uuid.UUID, now time.Time) (*models.TokenIntrospection, error) {
	// Create the content of the token.
	source := models.UserToken{
		Header:  models.UserTokenHeader{IAT: now, EXP: now.Add(s.tokenTTL), ID: id},
		Payload: data,
	}

	// Marshal token header into a base64 string.
	mrshHeader, err := json.Marshal(source.Header)
	if err != nil {
		return nil, errors.Join(ErrEncodeTokenHeader, err)
	}
	header := base64.RawURLEncoding.EncodeToString(mrshHeader)

	// Marshal token payload into a base64 string.
	mrshPayload, err := json.Marshal(source.Payload)
	if err != nil {
		return nil, errors.Join(ErrEncodeTokenPayload, err)
	}
	payload := base64.RawURLEncoding.EncodeToString(mrshPayload)

	// Merge together header and payload strings to create the unsigned version of the token.
	unsigned := fmt.Sprintf("%s.%s", header, payload)
	// Generate a signature to prevent data tampering.
	signature := base64.RawURLEncoding.EncodeToString(ed25519.Sign(models.JWTPrivateKey, []byte(unsigned)))

	return &models.TokenIntrospection{
		OK:       true,
		Token:    &source,
		TokenRaw: fmt.Sprintf("%s.%s", unsigned, signature),
	}, nil
}

type GetTokenStatusService interface {
	// GetTokenStatus introspects a token and returns the results.
	GetTokenStatus(token string, now time.Time) (*models.TokenIntrospection, error)
}

func NewGetTokenStatusService() GetTokenStatusService {
	return &getTokenStatusServiceImpl{}
}

type getTokenStatusServiceImpl struct{}

func (s *getTokenStatusServiceImpl) splitToken(token string) (string, string, string, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return "", "", "", ErrInvalidToken
	}

	header := parts[0]
	payload := parts[1]
	signature := parts[2]

	return header, payload, signature, nil
}

func (s *getTokenStatusServiceImpl) decodeToken(header, payload, signature string) ([]byte, []byte, []byte, error) {
	decodedSignature, err := base64.RawURLEncoding.DecodeString(signature)
	if err != nil {
		return nil, nil, nil, errors.Join(ErrInvalidToken, err)
	}

	decodedHeader, err := base64.RawURLEncoding.DecodeString(header)
	if err != nil {
		return nil, nil, nil, errors.Join(ErrInvalidToken, err)
	}

	decodedPayload, err := base64.RawURLEncoding.DecodeString(payload)
	if err != nil {
		return nil, nil, nil, errors.Join(ErrInvalidToken, err)
	}

	return decodedHeader, decodedPayload, decodedSignature, nil
}

func (s *getTokenStatusServiceImpl) validateToken(header, payload string, decodedSignature []byte) error {
	ok := ed25519.Verify(
		// We know for sure the public type is correct, because we read it from the private key.
		models.JWTPublicKey,
		[]byte(fmt.Sprintf("%s.%s", header, payload)),
		decodedSignature,
	)

	if !ok {
		return ErrBadSignature
	}

	return nil
}

func (s *getTokenStatusServiceImpl) GetTokenStatus(token string, now time.Time) (*models.TokenIntrospection, error) {
	status := &models.TokenIntrospection{TokenRaw: token}

	if token == "" {
		return status, nil
	}

	status.TokenRaw = token

	header, payload, signature, err := s.splitToken(token)
	if err != nil {
		if errors.Is(err, ErrInvalidToken) {
			status.Malformed = true
			return status, nil
		}

		return nil, err
	}

	decodedHeader, decodedPayload, decodedSignature, err := s.decodeToken(header, payload, signature)
	if err != nil {
		if errors.Is(err, ErrInvalidToken) {
			status.Malformed = true
			return status, nil
		}

		return nil, err
	}

	if err := s.validateToken(header, payload, decodedSignature); err != nil {
		if errors.Is(err, ErrBadSignature) {
			status.Expired = true
		} else {
			return nil, err
		}
	}

	parsedToken := new(models.UserToken)

	if err := json.Unmarshal(decodedHeader, &parsedToken.Header); err != nil {
		status.Malformed = true
		return status, nil
	}
	if err := json.Unmarshal(decodedPayload, &parsedToken.Payload); err != nil {
		status.Malformed = true
		return status, nil
	}

	status.Token = parsedToken

	if !status.Expired {
		if parsedToken.Header.ID == uuid.Nil {
			status.Malformed = true
			return status, nil
		}
		if parsedToken.Header.IAT.After(now) {
			status.NotIssued = true
			return status, nil
		}
		if parsedToken.Header.EXP.Before(now) {
			status.Expired = true
			return status, nil
		}

		status.OK = true
	}

	return status, nil
}
