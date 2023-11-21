package models

import (
	"github.com/google/uuid"
	"time"
)

// TokenIntrospection is the result of a token introspection.
type TokenIntrospection struct {
	// OK is true if the token is valid.
	OK bool `json:"ok"`
	// Expired is true if the token is past expiration date.
	Expired bool `json:"expired"`
	// NotIssued is true if the token has an issuedAt date in the future.
	NotIssued bool `json:"notIssued"`
	// Malformed is true if the token is not a valid JWT.
	Malformed bool `json:"malformed"`
	// Token contains the decoded token, if decoding was successful.
	Token *UserToken `json:"token,omitempty"`
	// TokenRaw is the original token sent in the headers.
	TokenRaw string `json:"tokenRaw,omitempty"`
}

type UserTokenHeader struct {
	// IAT (issuedAt) sets the date when the token starts to become valid.
	IAT time.Time `json:"iat"`
	// EXP (expiration) sets the date when the token becomes invalid.
	EXP time.Time `json:"exp"`
	// ID is a unique identifier for this token, that guarantees a unique encoded string.
	ID uuid.UUID `json:"id"`
}

type UserTokenPayload struct {
	// ID of the user who owns this token.
	ID string `json:"id"`
}

// UserToken represents the token issued to a user, for authentication.
type UserToken struct {
	Header  UserTokenHeader  `json:"header"`
	Payload UserTokenPayload `json:"payload"`
}
