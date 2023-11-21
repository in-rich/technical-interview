package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"technical-interview/pkg/services"
)

type loginForm struct {
	Email    string `json:"email" form:"email" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}

type LoginHandler interface {
	Handle(c *gin.Context)
}

func NewLoginHandler(service services.LoginService) LoginHandler {
	return &loginHandlerImpl{
		service: service,
	}
}

type loginHandlerImpl struct {
	service services.LoginService
}

func (h *loginHandlerImpl) Handle(c *gin.Context) {
	form := new(loginForm)

	if err := c.ShouldBind(form); err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	user, token, err := h.service.Exec(c, form.Email, form.Password)

	if err != nil {
		if errors.Is(err, services.ErrInvalidPassword) {
			_ = c.AbortWithError(http.StatusForbidden, err)
			return
		}
		if errors.Is(err, services.ErrInvalidEntity) {
			_ = c.AbortWithError(http.StatusUnprocessableEntity, err)
			return
		}

		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user":  user,
		"token": token,
	})
}
