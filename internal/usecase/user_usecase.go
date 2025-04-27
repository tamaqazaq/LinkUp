package usecase

import (
	"WebMessanger/internal/model"
	"github.com/google/uuid"
)

type UserUsecase interface {
	Register(user *model.User) (*model.User, error)
	Login(name, password string) (*model.User, error)
	GetUserByID(id uuid.UUID) (*model.User, error)
	UpdateProfile(user *model.User) error
	SearchUsers(query string) ([]*model.PublicUser, error)
	VerifyUserEmail(userID uuid.UUID) error // ✅ добавлено
}

type UserRepository interface {
	Create(user *model.User) error
	GetByName(name string) (*model.User, error)
	GetByID(id uuid.UUID) (*model.User, error)
	UpdateProfile(user *model.User) error
	Search(query string) ([]*model.User, error)
	VerifyUserEmail(userID uuid.UUID) error
	GetByEmail(email string) (*model.User, error)
}
