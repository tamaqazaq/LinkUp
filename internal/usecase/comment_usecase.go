package usecase

import (
	"WebMessanger/internal/model"
	"github.com/google/uuid"
)

type CommentUsecase interface {
	Create(comment *model.Comment) (*model.Comment, error)
	Delete(id uuid.UUID) error
	GetByThread(threadID uuid.UUID) ([]*model.Comment, error)
}
type CommentRepository interface {
	Create(comment *model.Comment) (*model.Comment, error)
	Delete(id uuid.UUID) error
	GetByThread(threadID uuid.UUID) ([]*model.Comment, error)
}
