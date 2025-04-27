package usecase

import (
	"WebMessanger/internal/model"
	"github.com/google/uuid"
)

type LikeUsecase interface {
	AddLike(like *model.Like) (*model.Like, error)
	RemoveLike(threadID uuid.UUID, userID uuid.UUID) error
	GetLikesByThread(threadID uuid.UUID) ([]*model.Like, error)
	GetLikesByUser(userID uuid.UUID) ([]*model.Like, error)
}

type LikeRepository interface {
	AddLike(like *model.Like) error
	RemoveLike(threadID uuid.UUID, userID uuid.UUID) error
	GetLikesByThread(threadID uuid.UUID) ([]*model.Like, error)
	GetLikesByUser(userID uuid.UUID) ([]*model.Like, error)
}
