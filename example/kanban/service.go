package main

import (
	"time"

	"github.com/google/uuid"
)

type kanbanService struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &kanbanService{repo: repo}
}

func (s *kanbanService) CreateTask(title, description string) (Task, error) {
	// Domain Business Rule validation
	if len(title) < 3 {
		return Task{}, ErrTitleCannotBeEmpty
	}

	task := Task{
		ID:          uuid.NewString(), // Generates an infrastructure-layer unique ID
		Title:       title,
		Description: description,
		Status:      "TODO",
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	if err := s.repo.Create(task); err != nil {
		return Task{}, err
	}

	return task, nil
}

func (s *kanbanService) GetTask(id string) (Task, error) {
	return s.repo.GetByID(id)
}

func (s *kanbanService) MoveTask(id string, newStatus string) (Task, error) {
	validStatuses := map[string]bool{"TODO": true, "DOING": true, "DONE": true}
	if !validStatuses[newStatus] {
		// Example of a domain error violation
		return Task{}, ErrInvalidStatus
	}

	task, err := s.repo.GetByID(id)
	if err != nil {
		return Task{}, err
	}

	task.Status = newStatus
	task.UpdatedAt = time.Now().UTC()

	if err := s.repo.Update(task); err != nil {
		return Task{}, err
	}

	return task, nil
}

func (s *kanbanService) ListTasks() ([]Task, error) {
	return s.repo.List()
}
