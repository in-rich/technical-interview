package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"technical-interview/pkg/dao"
	"technical-interview/pkg/services"
)

type registerForm struct {
	Email    string `json:"email" form:"email" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
	Username string `json:"username" form:"username" binding:"required"`
}

type RegisterHandler interface {
	Handle(c *gin.Context)
}

func NewRegisterHandler(service services.RegisterService) RegisterHandler {
	return &registerHandlerImpl{
		service: service,
	}
}

type registerHandlerImpl struct {
	service services.RegisterService
}

func (h *registerHandlerImpl) Handle(c *gin.Context) {
	form := new(registerForm)

	if err := c.ShouldBind(form); err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	user, token, err := h.service.Exec(c, form.Email, form.Password, form.Username)

	if err != nil {
		if errors.Is(err, dao.ErrEmailTaken) {
			_ = c.AbortWithError(http.StatusConflict, err)
			return
		}

		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"user":  user,
		"token": token,
	})
}
