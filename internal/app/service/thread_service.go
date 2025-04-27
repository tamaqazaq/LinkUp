package service

import (
	"WebMessanger/internal/model"
	"WebMessanger/internal/usecase"
	"github.com/google/uuid"
	"time"
)

type ThreadService struct {
	repo usecase.ThreadRepository
}

func NewThreadService(repo usecase.ThreadRepository) *ThreadService {
	return &ThreadService{repo: repo}
}
func (s *ThreadService) Create(thread *model.Thread) (*model.Thread, error) {
	thread.ID = uuid.New()
	thread.CreatedAt = time.Now()
	return s.repo.Create(thread)
}
func (s *ThreadService) Update(thread *model.Thread) (*model.Thread, error) {
	return s.repo.Update(thread)
}
func (s *ThreadService) Delete(id uuid.UUID) error {
	return s.repo.Delete(id)
}
func (s *ThreadService) GetAllThreads() ([]*model.Thread, error) {
	return s.repo.GetAllThreads()
}
func (s *ThreadService) GetThreadById(id uuid.UUID) (*model.Thread, error) {
	return s.repo.GetThreadById(id)
}
func (s *ThreadService) GetByUser(userID uuid.UUID) ([]*model.Thread, error) {
	return s.repo.GetByUser(userID)
}
