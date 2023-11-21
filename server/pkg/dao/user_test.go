package dao_test

import (
	"context"
	"technical-interview/config"
	"technical-interview/pkg/dao"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

const UsersTestCollection = "test-users"

func TestUserCreate(t *testing.T) {
	firestoreClient := config.FirestoreClient
	repository := dao.NewUserRepository(firestoreClient.Collection(UsersTestCollection))

	fixtures := map[string]interface{}{
		"01010101-0101-0101-0101-010101010101": map[string]interface{}{
			"id":       "01010101-0101-0101-0101-010101010101",
			"email":    "user1@gmail.com",
			"password": "safely-hashed-password",
			"username": "user1",
		},
	}

	data := []struct {
		name string

		email    string
		password string
		username string

		expectErr error
	}{
		{
			name:     "Success",
			email:    "user2@gmail.com",
			password: "1234",
			username: "user2",
		},
		{
			name:      "EmailTaken",
			email:     "user1@gmail.com",
			password:  "1234",
			username:  "user2",
			expectErr: dao.ErrEmailTaken,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			defer func() {
				require.NoError(t, CleanFirestore(firestoreClient))
			}()

			for id, note := range fixtures {
				_, err := firestoreClient.Collection(UsersTestCollection).Doc(id).Set(context.Background(), note)
				require.NoError(t, err)
			}

			res, err := repository.Create(context.Background(), d.email, d.password, d.username)
			require.ErrorIs(t, err, d.expectErr)

			if err == nil {
				require.Equal(t, d.email, res.Email)
				require.Equal(t, d.username, res.Username)

				err = bcrypt.CompareHashAndPassword([]byte(res.Password), []byte(d.password))
				require.NoError(t, err)
			}
		})
	}
}

func TestGetUser(t *testing.T) {
	firestoreClient := config.FirestoreClient
	repository := dao.NewUserRepository(firestoreClient.Collection(UsersTestCollection))

	fixtures := map[string]interface{}{
		"01010101-0101-0101-0101-010101010101": map[string]interface{}{
			"id":       "01010101-0101-0101-0101-010101010101",
			"email":    "user1@gmail.com",
			"password": "safely-hashed-password",
			"username": "user1",
		},
	}

	data := []struct {
		name string

		email string

		expectErr error
	}{
		{
			name:  "Success",
			email: "user1@gmail.com",
		},
		{
			name:      "UserNotFound",
			email:     "user2@gmail.com",
			expectErr: dao.ErrUserNotFound,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			defer func() {
				require.NoError(t, CleanFirestore(firestoreClient))
			}()

			for id, note := range fixtures {
				_, err := firestoreClient.Collection(UsersTestCollection).Doc(id).Set(context.Background(), note)
				require.NoError(t, err)
			}

			res, err := repository.GetUserByEmail(context.Background(), d.email)
			require.ErrorIs(t, err, d.expectErr)

			if err == nil {
				require.Equal(t, d.email, res.Email)
			}
		})
	}
}

func TestUpdateEmail(t *testing.T) {
	firestoreClient := config.FirestoreClient
	repository := dao.NewUserRepository(firestoreClient.Collection(UsersTestCollection))

	fixtures := map[string]interface{}{
		"01010101-0101-0101-0101-010101010101": map[string]interface{}{
			"id":       "01010101-0101-0101-0101-010101010101",
			"email":    "user1@gmail.com",
			"password": "safely-hashed-password",
			"username": "user1",
		},
		"02020202-0202-0202-0202-020202020202": map[string]interface{}{
			"id":       "02020202-0202-0202-0202-020202020202",
			"email":    "user2@gmail.com",
			"password": "safely-hashed-password",
			"username": "user2",
		},
	}

	data := []struct {
		name string

		email string
		id    string

		expectErr error
	}{
		{
			name:  "Success",
			id:    "01010101-0101-0101-0101-010101010101",
			email: "user3@gmail.com",
		},
		{
			name:      "UserNotFound",
			id:        "03030303-0303-0303-0303-030303030303",
			email:     "user3@gmail.com",
			expectErr: dao.ErrUserNotFound,
		},
		{
			name:      "EmailTaken",
			id:        "01010101-0101-0101-0101-010101010101",
			email:     "user2@gmail.com",
			expectErr: dao.ErrEmailTaken,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			defer func() {
				require.NoError(t, CleanFirestore(firestoreClient))
			}()

			for id, note := range fixtures {
				_, err := firestoreClient.Collection(UsersTestCollection).Doc(id).Set(context.Background(), note)
				require.NoError(t, err)
			}

			err := repository.UpdateEmail(context.Background(), d.id, d.email)
			require.ErrorIs(t, err, d.expectErr)

			if err == nil {
				res, err := repository.GetUserByEmail(context.Background(), d.email)
				require.NoError(t, err)
				require.Equal(t, d.email, res.Email)
			}
		})
	}
}
