package handler

import (
	"WebMessanger/internal/adapter/handler/mail"
	"WebMessanger/internal/model"
	"WebMessanger/internal/usecase"
	"WebMessanger/pkg/middleware"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

type AuthHandler struct {
	uc usecase.UserUsecase
}

func NewAuthHandler(uc usecase.UserUsecase) *AuthHandler {
	return &AuthHandler{uc: uc}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var input model.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if input.Name == "" || input.Password == "" || input.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	if len(input.Name) < 3 || len(input.Name) > 16 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid name"})
		return
	}
	if len(input.Password) < 5 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password too short"})
		return
	}

	user, err := h.uc.Register(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	verificationToken := middleware.GenerateToken(user.ID, user.Name, time.Now().Add(24*time.Hour))
	verifyURL := "http://localhost:8080/verify?token=" + verificationToken

	emailBody := fmt.Sprintf(`
		<h1>Welcome to LinkUp, %s!</h1>
		<p>Please click the button below to verify your email address:</p>
		<a href="%s" style="padding: 10px 20px; background-color: #3498db; color: white; border-radius: 4px; text-decoration: none;">Verify Email</a>
		<p>Or paste this link in your browser:</p>
		<p>%s</p>
	`, user.Name, verifyURL, verifyURL)

	if err := mail.SendEmail(user.Email, "Confirm your LinkUp account", emailBody); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send verification email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Registration successful! Please check your email."})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var input model.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	user, err := h.uc.Login(input.Name, input.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	expiration := time.Now().Add(7 * 24 * time.Hour)
	token := middleware.GenerateToken(user.ID, user.Name, expiration)

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	tokenString := c.Query("token")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	claims := &model.Claims{}
	jwtKey := middleware.JwtKey
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired verification token"})
		return
	}

	if err := h.uc.VerifyUserEmail(claims.UserID); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired verification token"})
		return
	}

	c.Redirect(http.StatusFound, "http://localhost:5173/email-verified")
}
