package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/go-chi/chi/v5"
	"github.com/semmidev/problem"
)

type KanbanHandler struct {
	service Service
}

func NewHandler(service Service) *KanbanHandler {
	return &KanbanHandler{service: service}
}

// mapErrorToProblem is the central place where domain errors and infrastructure errors
// are translated into RFC 7807 Problem Details.
func mapErrorToProblem(err error) *problem.Problem {
	// 1. Not Found Domain Error
	if errors.Is(err, ErrTaskNotFound) {
		return problem.Wrap(err, problem.NotFound, problem.WithDetail(err.Error()))
	}
	// 2. Business Logic Validation Errors
	if errors.Is(err, ErrInvalidStatus) || errors.Is(err, ErrTitleCannotBeEmpty) {
		return problem.Wrap(err, problem.UnprocessableEntity, problem.WithDetail(err.Error()))
	}

	// 3. Fallback to 500 Internal Server error for anything unhandled
	// In a real app we'd log the original err securely here and mask details to the client
	return problem.Wrap(err, problem.InternalServerError, problem.WithDetail("An unexpected internal error occurred."))
}

// writeError Helper
func writeError(w http.ResponseWriter, r *http.Request, err error) {
	p := mapErrorToProblem(err)
	// Optionally attach instance URI
	p.Instance = r.URL.Path
	p.Write(w)
}

// ==== Request Models ====

type CreateTaskRequest struct {
	Title       string `json:"title" valid:"required,stringlength(3|100)"`
	Description string `json:"description" valid:"type(string)"`
}

type MoveTaskRequest struct {
	Status string `json:"status" valid:"in(TODO|DOING|DONE),required"`
}

// ==== HTTP Handlers ====

func (h *KanbanHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var req CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		problem.Wrap(err, problem.BadRequest, problem.WithDetail("Invalid or malformed JSON payload")).Write(w)
		return
	}

	if _, err := govalidator.ValidateStruct(req); err != nil {
		var validationErrors []map[string]string
		if errs, ok := err.(govalidator.Errors); ok {
			for _, e := range errs {
				if valErr, isValErr := e.(govalidator.Error); isValErr {
					validationErrors = append(validationErrors, map[string]string{
						"field":   valErr.Name,
						"message": valErr.Err.Error(),
					})
				} else {
					validationErrors = append(validationErrors, map[string]string{"message": e.Error()})
				}
			}
		} else {
			validationErrors = append(validationErrors, map[string]string{"message": err.Error()})
		}

		p := problem.New(
			problem.UnprocessableEntity,
			problem.WithDetail("request parameters failed validation"),
			problem.WithInstance(r.URL.Path),
			problem.WithExtension("invalid_params", validationErrors),
		)
		p.Write(w)
		return
	}

	task, err := h.service.CreateTask(req.Title, req.Description)
	if err != nil {
		writeError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

func (h *KanbanHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	task, err := h.service.GetTask(id)
	if err != nil {
		writeError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func (h *KanbanHandler) MoveTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req MoveTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		problem.Wrap(err, problem.BadRequest, problem.WithDetail("Invalid JSON")).Write(w)
		return
	}

	if _, err := govalidator.ValidateStruct(req); err != nil {
		var validationErrors []map[string]string
		if errs, ok := err.(govalidator.Errors); ok {
			for _, e := range errs {
				if valErr, isValErr := e.(govalidator.Error); isValErr {
					validationErrors = append(validationErrors, map[string]string{
						"field":   valErr.Name,
						"message": valErr.Err.Error(),
					})
				} else {
					validationErrors = append(validationErrors, map[string]string{"message": e.Error()})
				}
			}
		} else {
			validationErrors = append(validationErrors, map[string]string{"message": err.Error()})
		}

		p := problem.New(
			problem.UnprocessableEntity,
			problem.WithDetail("validation failed for status update"),
			problem.WithExtension("invalid_params", validationErrors),
		)
		p.Write(w)
		return
	}

	task, err := h.service.MoveTask(id, req.Status)
	if err != nil {
		writeError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func (h *KanbanHandler) ListTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.service.ListTasks()
	if err != nil {
		writeError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}
