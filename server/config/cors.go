package config

import (
	"github.com/gin-contrib/cors"
	"time"
)

var Cors = cors.Config{
	AllowOrigins: []string{"*"},
	AllowMethods: cors.DefaultConfig().AllowMethods,
	AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
	ExposeHeaders: []string{
		"Content-Type",
		"Content-Length",
		"Access-Control-Allow-Origin",
	},
	AllowCredentials: false,
	MaxAge:           12 * time.Hour,
}
