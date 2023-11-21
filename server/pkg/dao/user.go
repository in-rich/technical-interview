package dao

import (
	"context"
	"errors"
	"technical-interview/pkg/models"

	"cloud.google.com/go/firestore"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrUserNotFound  = errors.New("user not found")
	ErrParseDocument = errors.New("error while parsing document")
	ErrEmailTaken    = errors.New("email already taken")
)

type UserRepository interface {
	Create(ctx context.Context, email string, password string, username string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	UpdateEmail(ctx context.Context, id string, email string) error
}

func NewUserRepository(collection *firestore.CollectionRef) UserRepository {
	return &userRepositoryImpl{
		collection: collection,
	}
}

type userRepositoryImpl struct {
	collection *firestore.CollectionRef
}

func (repository *userRepositoryImpl) Create(ctx context.Context, email string, password string, username string) (*models.User, error) {
	id := uuid.New()

	// Verify email is available.
	_, err := repository.GetUserByEmail(ctx, email)
	if err == nil {
		return nil, ErrEmailTaken
	}
	if !errors.Is(err, ErrUserNotFound) {
		return nil, err
	}

	// Hash the password so it doesn't get exposed in case of data leak.
	passwordHashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	output := &models.User{
		ID:       id.String(),
		Email:    email,
		Username: username,
		Password: string(passwordHashed),
	}

	_, err = repository.collection.Doc(id.String()).Set(ctx, output)
	if err != nil {
		return nil, err
	}

	return output, nil
}

func (repository *userRepositoryImpl) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	output := new(models.User)

	doc, err := repository.collection.Where("email", "==", email).Limit(1).Documents(ctx).Next()
	if err != nil {
		return nil, lo.Ternary(err == iterator.Done, ErrUserNotFound, err)
	}

	if err := doc.DataTo(output); err != nil {
		return nil, errors.Join(ErrParseDocument, err)
	}

	return output, nil
}

func (repository *userRepositoryImpl) UpdateEmail(ctx context.Context, id string, email string) error {
	// Verify email is available.
	_, err := repository.GetUserByEmail(ctx, email)
	if err == nil {
		return ErrEmailTaken
	}
	if !errors.Is(err, ErrUserNotFound) {
		return err
	}

	// Verify user exists.
	_, err = repository.collection.Doc(id).Get(ctx)
	if err != nil {
		return lo.Ternary(status.Code(err) == codes.NotFound, ErrUserNotFound, err)
	}

	_, err = repository.collection.Doc(id).
		Set(ctx, map[string]interface{}{"email": email}, firestore.Merge(firestore.FieldPath{"email"}))
	if err != nil {
		return err
	}

	return nil
}
