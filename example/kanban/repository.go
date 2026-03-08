package main

import (
	"sync"
)

// memoryRepo implements the domain Repository interface
// using a concurrent-safe in-memory map.
type memoryRepo struct {
	tasks map[string]Task
	mu    sync.RWMutex
}

func NewMemoryRepository() Repository {
	return &memoryRepo{
		tasks: make(map[string]Task),
	}
}

func (r *memoryRepo) Create(task Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tasks[task.ID] = task
	return nil
}

func (r *memoryRepo) GetByID(id string) (Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	task, exists := r.tasks[id]
	if !exists {
		// Return our clear domain error.
		return Task{}, ErrTaskNotFound
	}
	return task, nil
}

func (r *memoryRepo) Update(task Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.tasks[task.ID]; !exists {
		return ErrTaskNotFound
	}

	r.tasks[task.ID] = task
	return nil
}

func (r *memoryRepo) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.tasks[id]; !exists {
		return ErrTaskNotFound
	}

	delete(r.tasks, id)
	return nil
}

func (r *memoryRepo) List() ([]Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	list := make([]Task, 0, len(r.tasks))
	for _, t := range r.tasks {
		list = append(list, t)
	}
	return list, nil
}
