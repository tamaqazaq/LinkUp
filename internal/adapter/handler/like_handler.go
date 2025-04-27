package handler

import (
	"WebMessanger/internal/model"
	"WebMessanger/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"net/http"
)

type LikeHandler struct {
	uc            usecase.LikeUsecase
	userUsecase   usecase.UserUsecase
	threadUsecase usecase.ThreadUsecase
}

func NewLikeHandler(uc usecase.LikeUsecase, userUsecase usecase.UserUsecase, threadUsecase usecase.ThreadUsecase) *LikeHandler {
	return &LikeHandler{uc, userUsecase, threadUsecase}
}
func (h *LikeHandler) Create(c *gin.Context) {

	var like model.Like
	log.Println("Adding like from", like.UserID, "to", like.ThreadID)

	threadIDStr := c.Param("thread_id")
	threadID, err := uuid.Parse(threadIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid thread id"})
		return
	}
	like.ThreadID = threadID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authorized"})
		return
	}
	like.UserID = userID.(uuid.UUID)

	user, err := h.userUsecase.GetUserByID(like.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	like.UserName = user.Name

	created, err := h.uc.AddLike(&like)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, created)
}
func (h *LikeHandler) RemoveLike(c *gin.Context) {
	threadIDStr := c.Param("thread_id")
	threadID, err := uuid.Parse(threadIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authorized"})
		return
	}
	err = h.uc.RemoveLike(threadID, userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, gin.H{"message": "ok"})
}
func (h *LikeHandler) GetLikesByThread(c *gin.Context) {
	threadIDStr := c.Param("thread_id")
	threadID, err := uuid.Parse(threadIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid thread id"})
		return
	}

	likes, err := h.uc.GetLikesByThread(threadID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, likes)
}
func (h *LikeHandler) GetLikesByUser(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	likes, err := h.uc.GetLikesByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, likes)
}
