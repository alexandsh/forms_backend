package handler

import (
	"context"
	"forms/internal/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FormService interface {
	CreateForm(ctx context.Context, userID int, req model.CreateFormRequest) (int, error)
	GetForm(ctx context.Context, id int) (*model.GetFormResponse, error)
	GetForms(ctx context.Context, userID int) (*model.GetFormsResponse, error)
	UpdateForm(ctx context.Context, userID int, formID int, req model.UpdateFormRequest) error
	DeleteForm(ctx context.Context, userID int, formID int) error
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

func (h *FormHandler) GetForm(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	form, err := h.service.GetForm(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, form)
}

func (h *FormHandler) GetForms(c *gin.Context) {
	userID := c.GetInt("user_id")

	form, err := h.service.GetForms(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, form)
}

func (h *FormHandler) UpdateForm(c *gin.Context) {
	userID := c.GetInt("user_id")

	formID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var req model.UpdateFormRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = h.service.UpdateForm(c.Request.Context(), userID, formID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *FormHandler) DeleteForm(c *gin.Context) {
	userID := c.GetInt("user_id")

	formID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = h.service.DeleteForm(c.Request.Context(), userID, formID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}
