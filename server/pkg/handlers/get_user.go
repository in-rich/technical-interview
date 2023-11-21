package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"technical-interview/pkg/dao"
	"technical-interview/pkg/services"
)

type GetUserHandler interface {
	Handle(c *gin.Context)
}

func NewGetUserHandler(service services.GetUserService) GetUserHandler {
	return &getUserHandlerImpl{
		service: service,
	}
}

type getUserHandlerImpl struct {
	service services.GetUserService
}

func (h *getUserHandlerImpl) Handle(c *gin.Context) {
	res, err := h.service.Exec(c, c.Query("email"))

	if err != nil {
		if errors.Is(err, dao.ErrUserNotFound) {
			_ = c.AbortWithError(http.StatusNotFound, err)
			return
		}

		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, res)
}
