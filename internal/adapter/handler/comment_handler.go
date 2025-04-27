package handler

import (
	"WebMessanger/internal/model"
	"WebMessanger/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type CommentHandler struct {
	uc          usecase.CommentUsecase
	userUsecase usecase.UserUsecase
}

func NewCommentHandler(uc usecase.CommentUsecase, userUC usecase.UserUsecase) *CommentHandler {
	return &CommentHandler{uc: uc, userUsecase: userUC}
}
func (h *CommentHandler) Create(c *gin.Context) {
	var comment model.Comment
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	threadIDStr := c.Param("thread_id")
	threadID, err := uuid.Parse(threadIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid thread id"})
		return
	}
	comment.ThreadID = threadID

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authorized"})
		return
	}
	comment.UserID = userID.(uuid.UUID)

	user, err := h.userUsecase.GetUserByID(comment.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	comment.UserName = user.Name

	created, err := h.uc.Create(&comment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, created)
}

func (h *CommentHandler) Delete(c *gin.Context) {
	idStr := c.Param("comment_id")
	id, err := uuid.Parse(idStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.uc.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusNoContent, gin.H{"message": "comment deleted successfully"})
}
func (h *CommentHandler) GetByThread(c *gin.Context) {
	threadIDStr := c.Param("thread_id")
	threadID, err := uuid.Parse(threadIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid thread id"})
		return
	}
	comments, err := h.uc.GetByThread(threadID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, comments)
}
