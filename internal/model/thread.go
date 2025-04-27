package model

import (
	"github.com/google/uuid"
	"time"
)

type Thread struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	UserName  string    `json:"user_name"`
	AvatarURL string    `json:"avatar_url"`
	Content   string    `json:"content"`   // текст поста
	MediaURL  string    `json:"media_url"` // (опционально) фото/видео
	CreatedAt time.Time `json:"created_at"`
}
