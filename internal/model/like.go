package model

import (
	"github.com/google/uuid"
	"time"
)

type Like struct {
	ID        uuid.UUID
	ThreadID  uuid.UUID
	UserID    uuid.UUID
	UserName  string
	CreatedAt time.Time
}
