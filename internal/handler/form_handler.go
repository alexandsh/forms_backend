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
	CreateResponse(ctx context.Context, userID int, formID int, req model.CreateResponseRequest) error
}

type FormHandler struct {
	service FormService
}

func NewFormHandler(service FormService) *FormHandler {
	return &FormHandler{service: service}
}

// CreateForm godoc
// @Summary      Создание формы
// @Description  Создает новую форму с вопросами (text, radio, checkbox)
// @Tags         forms
// @Accept       json
// @Produce      json
// @Param        request  body      model.CreateFormRequest  true  "Данные формы"
// @Success      201      {object}  map[string]int
// @Failure      400      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /api/form [post]
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

// GetForm godoc
// @Summary      Получение формы
// @Description  Возвращает форму по идентификатору
// @Tags         forms
// @Produce      json
// @Param        id   path      int  true  "ID формы"
// @Success      200  {object}  model.GetFormResponse
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/form/{id} [get]
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

// GetForms godoc
// @Summary      Список своих форм
// @Description  Возвращает все формы, созданные текущим пользователем
// @Tags         forms
// @Produce      json
// @Success      200  {object}  model.GetFormsResponse
// @Failure      500  {object}  map[string]string
// @Router       /api/forms [get]
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

// UpdateForm godoc
// @Summary      Обновление формы
// @Description  Обновляет заголовок формы и/или список вопросов
// @Tags         forms
// @Accept       json
// @Produce      json
// @Param        id       path      int                         true  "ID формы"
// @Param        request  body      model.UpdateFormRequest     true  "Данные для обновления"
// @Success      204      "No Content"
// @Failure      400      {object}  map[string]string
// @Failure      403      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /api/form/{id} [patch]
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

// DeleteForm godoc
// @Summary      Удаление формы
// @Description  Удаляет форму по идентификатору
// @Tags         forms
// @Produce      json
// @Param        id   path      int  true  "ID формы"
// @Success      204  "No Content"
// @Failure      400  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/form/{id} [delete]
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

// CreateResponse godoc
// @Summary      Отправка ответа на форму
// @Description  Создает ответ на форму с учетом типов вопросов
// @Tags         forms
// @Accept       json
// @Produce      json
// @Param        id       path      int                          true  "ID формы"
// @Param        request  body      model.CreateResponseRequest  true  "Ответ на форму"
// @Success      201      "Created"
// @Failure      400      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /api/form/{id}/responses [post]
func (h *FormHandler) CreateResponse(c *gin.Context) {
	userID := c.GetInt("user_id")

	formID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var req model.CreateResponseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = h.service.CreateResponse(c.Request.Context(), userID, formID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.Status(http.StatusCreated)
}
