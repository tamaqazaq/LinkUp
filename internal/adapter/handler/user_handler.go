package handler

import (
	"WebMessanger/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type UserHandler struct {
	uc usecase.UserUsecase
}

func NewUserHandler(uc usecase.UserUsecase) *UserHandler {
	return &UserHandler{uc: uc}
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
		return
	}
	user, err := h.uc.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user_id": user.ID, "name": user.Name, "email": user.Email, "avatar_url": user.AvatarURL,
		"bio": user.Bio, "location": user.Location, "social_links": user.SocialLinks, "created_at": user.CreatedAt})
}

func (h *UserHandler) GetMe(c *gin.Context) {
	rawID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	id := rawID.(uuid.UUID)
	user, err := h.uc.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id":      user.ID,
		"name":         user.Name,
		"email":        user.Email,
		"avatar_url":   user.AvatarURL,
		"bio":          user.Bio,
		"location":     user.Location,
		"social_links": user.SocialLinks,
		"created_at":   user.CreatedAt,
	})
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	rawID, _ := c.Get("user_id")
	id := rawID.(uuid.UUID)
	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}
	currentUser, err := h.uc.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if val, ok := body["name"].(string); ok && len(val) < 3 {
		c.JSON(400, gin.H{"error": "Name too short"})
		return
	}
	if val, ok := body["name"]; ok {
		currentUser.Name = val.(string)
	}
	if val, ok := body["bio"]; ok {
		currentUser.Bio = val.(string)
	}
	if val, ok := body["social_links"].(string); ok && val != "" {
		currentUser.SocialLinks = val
	}
	if val, ok := body["location"]; ok {
		currentUser.Location = val.(string)
	}
	if val, ok := body["avatar_url"].(string); ok && val != "" {
		currentUser.AvatarURL = val
	}
	if err := h.uc.UpdateProfile(currentUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})

}
func (h *UserHandler) SearchUsers(c *gin.Context) {
	query := c.Query("query")
	user, err := h.uc.SearchUsers(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"users": user})
}
