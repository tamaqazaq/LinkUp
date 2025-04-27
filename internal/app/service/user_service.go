package service

import (
	"WebMessanger/internal/model"
	"WebMessanger/internal/usecase"
	"errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type userService struct {
	repo usecase.UserRepository
}

func NewUserService(repo usecase.UserRepository) usecase.UserUsecase {
	return &userService{repo: repo}
}

func (s *userService) Register(user *model.User) (*model.User, error) {
	if existing, _ := s.repo.GetByName(user.Name); existing != nil {
		return nil, errors.New("user already exists")
	}
	if existingByEmail, _ := s.repo.GetByEmail(user.Email); existingByEmail != nil {
		return nil, errors.New("email already exists")
	}
	hashed, error := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if error != nil {
		return nil, error
	}
	user.CreatedAt = time.Now()
	user.ID = uuid.New()
	user.HashedPassword = string(hashed)
	if err := s.repo.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) Login(name, password string) (*model.User, error) {
	user, err := s.repo.GetByName(name)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}
	if bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password)) != nil {
		return nil, errors.New("invalid credentials")
	}
	if !user.IsVerified {
		return nil, errors.New("please verify your email before logging in")
	}

	return user, nil
}

func (s *userService) GetUserByID(id uuid.UUID) (*model.User, error) {
	return s.repo.GetByID(id)
}

func (s *userService) UpdateProfile(user *model.User) error {
	existing, _ := s.repo.GetByName(user.Name)
	if existing != nil && existing.ID != user.ID {
		return errors.New("this name is already taken")
	}
	return s.repo.UpdateProfile(user)
}
func (s *userService) SearchUsers(query string) ([]*model.PublicUser, error) {
	if query == "" {
		return nil, errors.New("empty query")
	}

	users, err := s.repo.Search(query)
	if err != nil {
		return nil, err
	}
	publicUsers := make([]*model.PublicUser, 0, len(users))
	for _, u := range users {
		publicUsers = append(publicUsers, &model.PublicUser{
			ID:        u.ID,
			Name:      u.Name,
			AvatarURL: u.AvatarURL,
		})
	}

	return publicUsers, nil
}
func (s *userService) VerifyUserEmail(userID uuid.UUID) error {
	return s.repo.VerifyUserEmail(userID)
}
