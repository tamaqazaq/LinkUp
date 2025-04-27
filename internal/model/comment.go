package model

import (
	"github.com/google/uuid"
	"time"
)

type Comment struct {
	ID        uuid.UUID `json:"id"`
	ThreadID  uuid.UUID `json:"thread_id"`
	UserID    uuid.UUID `json:"user_id"`
	UserName  string    `json:"user_name"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}
