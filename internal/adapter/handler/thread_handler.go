package handler

import (
	"WebMessanger/internal/app/service"
	"WebMessanger/internal/model"
	"WebMessanger/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type ThreadHandler struct {
	uc usecase.ThreadUsecase
}

func NewThreadHandler(uc *service.ThreadService) *ThreadHandler {
	return &ThreadHandler{uc: uc}
}
func (h *ThreadHandler) Create(c *gin.Context) {
	var thread model.Thread
	if err := c.ShouldBindJSON(&thread); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authorized"})
		return
	}
	thread.UserID = userID.(uuid.UUID)
	created, err := h.uc.Create(&thread)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusCreated, created)
}

func (h *ThreadHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid thread id"})
		return
	}

	var body struct {
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid or empty content"})
		return
	}

	existing, err := h.uc.GetThreadById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "thread not found"})
		return
	}

	existing.Content = body.Content

	updated, err := h.uc.Update(existing)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updated)
}

func (h *ThreadHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.uc.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Thread has been deleted"})
}
func (h *ThreadHandler) GetAll(c *gin.Context) {
	threads, err := h.uc.GetAllThreads()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, threads)
}
func (h *ThreadHandler) GetById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	thread, err := h.uc.GetThreadById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, thread)
}
func (h *ThreadHandler) GetByUser(c *gin.Context) {
	UserIDStr := c.Param("user_id")
	id, err := uuid.Parse(UserIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	threads, err := h.uc.GetByUser(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, threads)
}
