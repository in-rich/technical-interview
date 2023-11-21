package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"os"
	"technical-interview/config"
	"technical-interview/pkg/api"
	"technical-interview/pkg/dao"
	"technical-interview/pkg/handlers"
	"technical-interview/pkg/services"
	"time"
)

func newLogger() zerolog.Logger {
	return zerolog.
		New(os.Stdout).
		With().
		Dict(
			"application",
			zerolog.Dict().
				Str("name", config.App.Name).
				Str("env", config.ENV).
				Int("port", config.App.Port),
		).
		Logger()
}

func main() {
	logger := newLogger()
	router := gin.New()

	if config.ENV == config.ProdENV {
		gin.SetMode(gin.ReleaseMode)
	}

	routerAPI := router.Use(
		gin.RecoveryWithWriter(logger),
		api.Logger(logger, config.App.ProjectID),
		cors.New(config.Cors),
	)

	userDAO := dao.NewUserRepository(config.FirestoreClient.Collection("users"))

	generateTokenService := services.NewGenerateTokenService(24 * time.Hour)
	introspectTokenService := services.NewGetTokenStatusService()

	getUserService := services.NewGetUserService(userDAO)
	updateEmailService := services.NewUpdateEmailService(userDAO, introspectTokenService)
	loginService := services.NewLoginService(userDAO, generateTokenService)
	registerService := services.NewRegisterService(userDAO, generateTokenService)

	getUserHandler := handlers.NewGetUserHandler(getUserService)
	updateEmailHandler := handlers.NewUpdateEmailHandler(updateEmailService)
	loginHandler := handlers.NewLoginHandler(loginService)
	registerHandler := handlers.NewRegisterHandler(registerService)

	routerAPI.GET("/user", getUserHandler.Handle)
	routerAPI.PUT("/user/email", updateEmailHandler.Handle)
	routerAPI.POST("/user", loginHandler.Handle)
	routerAPI.PUT("/user", registerHandler.Handle)

	if err := router.Run(fmt.Sprintf(":%d", config.App.Port)); err != nil {
		logger.Fatal().Err(err).Msg("a fatal error occurred while running API, and the server had to shut down")
	}
}
