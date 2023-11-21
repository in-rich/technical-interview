package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"technical-interview/pkg/services"
)

type updateEmailForm struct {
	Email string `json:"email" form:"email" binding:"required"`
}

type UpdateEmailHandler interface {
	Handle(c *gin.Context)
}

func NewUpdateEmailHandler(service services.UpdateEmailService) UpdateEmailHandler {
	return &updateEmailHandlerImpl{
		service: service,
	}
}

type updateEmailHandlerImpl struct {
	service services.UpdateEmailService
}

func (h *updateEmailHandlerImpl) Handle(c *gin.Context) {
	form := new(updateEmailForm)

	if err := c.ShouldBind(form); err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	err := h.service.UpdateEmail(c, c.GetHeader("Authorization"), form.Email)

	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			_ = c.AbortWithError(http.StatusForbidden, err)
			return
		}

		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusCreated)
}
