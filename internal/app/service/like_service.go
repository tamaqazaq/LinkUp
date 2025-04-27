package service

import (
	"WebMessanger/internal/model"
	"WebMessanger/internal/usecase"
	"github.com/google/uuid"
	"time"
)

type LikeService struct {
	repo usecase.LikeRepository
}

func NewLikeService(repo usecase.LikeRepository) *LikeService {
	return &LikeService{repo: repo}
}
func (s *LikeService) AddLike(like *model.Like) (*model.Like, error) {
	like.ID = uuid.New()
	like.CreatedAt = time.Now()

	err := s.repo.AddLike(like)
	if err != nil {
		return nil, err
	}
	return like, nil
}

func (s *LikeService) RemoveLike(threadID uuid.UUID, userID uuid.UUID) error {
	return s.repo.RemoveLike(threadID, userID)
}

func (s *LikeService) GetLikesByThread(threadID uuid.UUID) ([]*model.Like, error) {
	return s.repo.GetLikesByThread(threadID)
}

func (s *LikeService) GetLikesByUser(userID uuid.UUID) ([]*model.Like, error) {
	return s.repo.GetLikesByUser(userID)
}
