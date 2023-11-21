package config

import "os"

var ENV = os.Getenv("ENV")

const (
	DevENV  = "dev"
	ProdENV = "prod"

	// DefaultENV is used to target all environments at once. This value should not be used as the actual content
	// of the ENV variable.
	DefaultENV = "all"
)
