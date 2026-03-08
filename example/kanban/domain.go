package main

import (
	"errors"
	"time"
)

// Task represents a task in the Kanban board.
type Task struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	Status      string    `json:"status"` // E.g., "TODO", "DOING", "DONE"
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ==== Standardized Domain Errors ====
// These are errors that originate purely from the domain/business layer.
// They know nothing about HTTP, Status Codes, or RFC 7807.

var (
	// ErrTaskNotFound is returned when a task cannot be found in the store.
	ErrTaskNotFound = errors.New("task not found in the kanban board")

	// ErrInvalidStatus is returned when a task is moved to an invalid board column.
	ErrInvalidStatus = errors.New("invalid task status transition")

	// ErrTitleCannotBeEmpty is a domain rule.
	ErrTitleCannotBeEmpty = errors.New("a task title must have at least 3 characters")
)

// Repository defines how the application interacts with the data store.
// In Clean Architecture, the usecase/service layer interacts with this.
type Repository interface {
	Create(task Task) error
	GetByID(id string) (Task, error)
	Update(task Task) error
	Delete(id string) error
	List() ([]Task, error)
}

// Service is our primary use-case layer for business logic.
type Service interface {
	CreateTask(title, description string) (Task, error)
	GetTask(id string) (Task, error)
	MoveTask(id string, newStatus string) (Task, error)
	ListTasks() ([]Task, error)
}
