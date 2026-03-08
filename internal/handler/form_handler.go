package handler

import (
	"context"
	"forms/internal/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

type FormService interface {
	CreateForm(ctx context.Context, userID int, req model.CreateFormRequest) (int, error)
}

type FormHandler struct {
	service FormService
}

func NewFormHandler(service FormService) *FormHandler {
	return &FormHandler{service: service}
}

func (h *FormHandler) CreateForm(c *gin.Context) {
	userID := c.GetInt("user_id")

	var req model.CreateFormRequest

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error:": err.Error(),
		})
		return
	}

	formID, err := h.service.CreateForm(c, userID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error:": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id": formID,
	})

}
