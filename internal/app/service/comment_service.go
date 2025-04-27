package service

import (
	"WebMessanger/internal/model"
	"WebMessanger/internal/usecase"
	"github.com/google/uuid"
	"time"
)

type CommentService struct {
	repo usecase.CommentRepository
}

func NewCommentService(repo usecase.CommentRepository) *CommentService {
	return &CommentService{repo: repo}
}
func (s *CommentService) Create(comment *model.Comment) (*model.Comment, error) {
	comment.ID = uuid.New()
	comment.CreatedAt = time.Now()
	createdComment, err := s.repo.Create(comment)
	if err != nil {
		return nil, err
	}
	return createdComment, nil
}
func (s *CommentService) Delete(id uuid.UUID) error {
	return s.repo.Delete(id)
}
func (s *CommentService) GetByThread(threadID uuid.UUID) ([]*model.Comment, error) {
	return s.repo.GetByThread(threadID)
}
