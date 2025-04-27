package model

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID             uuid.UUID
	Name           string
	Password       string
	HashedPassword string
	Email          string
	Bio            string
	Location       string
	SocialLinks    string
	AvatarURL      string
	CreatedAt      time.Time
	IsVerified     bool
}

type Claims struct {
	UserID uuid.UUID
	Name   string `json:"name"`
	jwt.RegisteredClaims
}

type Response struct {
	Message string    `json:"message"`
	UserID  uuid.UUID `json:"user_id"`
	Name    string    `json:"name"`
	Email   string    `json:"email"`
}
type PublicUser struct {
	ID        uuid.UUID `json:"user_id"`
	Name      string    `json:"name"`
	AvatarURL string    `json:"avatar_url"`
}
