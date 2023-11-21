package models

import (
	"crypto/ed25519"
)

var (
	JWTPublicKey  ed25519.PublicKey
	JWTPrivateKey ed25519.PrivateKey
)

func init() {
	var err error

	JWTPublicKey, JWTPrivateKey, err = ed25519.GenerateKey(nil)
	if err != nil {
		panic(err)
	}
}
