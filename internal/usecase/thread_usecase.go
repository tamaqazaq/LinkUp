package usecase

import (
	"WebMessanger/internal/model"
	"github.com/google/uuid"
)

type ThreadUsecase interface {
	Create(thread *model.Thread) (*model.Thread, error)
	GetAllThreads() ([]*model.Thread, error)
	GetThreadById(id uuid.UUID) (*model.Thread, error)
	Update(thread *model.Thread) (*model.Thread, error)
	Delete(id uuid.UUID) error
	GetByUser(userID uuid.UUID) ([]*model.Thread, error)
}

type ThreadRepository interface {
	Create(thread *model.Thread) (*model.Thread, error)
	GetAllThreads() ([]*model.Thread, error)
	GetThreadById(id uuid.UUID) (*model.Thread, error)
	Update(thread *model.Thread) (*model.Thread, error)
	Delete(id uuid.UUID) error
	GetByUser(userID uuid.UUID) ([]*model.Thread, error)
}
