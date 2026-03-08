package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/semmidev/problem"
)

func main() {
	// 1. Setup Data Store (In-Memory)
	repo := NewMemoryRepository()

	// 2. Setup Business Logic Layer
	service := NewService(repo)

	// 3. Setup HTTP Handler
	handler := NewHandler(service)

	// 4. Setup Router using Chi
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Injecting an application-level panic recovery middleware
	// returning an RFC 7807 500 error on panic
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			defer func() {
				if rvr := recover(); rvr != nil {
					p := problem.New(
						problem.InternalServerError,
						problem.WithDetail("A critical panic occurred internally."),
						problem.WithInstance(req.URL.Path),
					)
					p.Write(w)
				}
			}()
			next.ServeHTTP(w, req)
		})
	})

	r.Route("/tasks", func(r chi.Router) {
		r.Get("/", handler.ListTasks)
		r.Post("/", handler.CreateTask)

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", handler.GetTask)
			r.Put("/status", handler.MoveTask)
		})
	})

	// Setup a standard 404 handler returning Problem details
	r.NotFound(func(w http.ResponseWriter, req *http.Request) {
		problem.New(
			problem.NotFound,
			problem.WithDetail("The requested endpoint does not exist."),
			problem.WithInstance(req.URL.Path),
		).Write(w)
	})

	// Setup a standard 405 Method Not Allowed
	r.MethodNotAllowed(func(w http.ResponseWriter, req *http.Request) {
		problem.New(
			problem.MethodNotAllowed,
			problem.WithDetail("The requested HTTP method is not allowed on this endpoint."),
			problem.WithInstance(req.URL.Path),
		).Write(w)
	})

	log.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
