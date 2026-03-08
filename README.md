# problem

[![Go Reference](https://pkg.go.dev/badge/github.com/semmidev/problem.svg)](https://pkg.go.dev/github.com/semmidev/problem)

A comprehensive, idiomatic, and robust Go library for implementing **RFC 7807 (Problem Details for HTTP APIs)**.

This library provides a standard way to return machine-readable errors from your HTTP APIs, ensuring consistency across your microservices and APIs.

## Features

- **Standard RFC 7807 Compliance**: Supports all standard fields (`type`, `title`, `status`, `detail`, `instance`).
- **Custom Extensions**: Easily add custom arbitrary fields (extension members) that automatically flatten into the JSON root.
- **Go 1.13+ Error Wrapping**: Implements `Error()` and `Unwrap()` making it fully compatible with `errors.Is` and `errors.As`.
- **Pre-defined HTTP Templates**: Built-in templates for all common HTTP 4xx and 5xx errors.
- **Fluent Options API**: Clean, chainable API for building problem details.
- **Zero Dependencies**: Relies solely on the Go standard library.

## Installation

```bash
go get github.com/semmidev/problem
```

## Basic Usage

The easiest way to use the library is with the pre-defined templates in your HTTP handlers.

```go
package main

import (
	"net/http"
	"github.com/semmidev/problem"
)

func myHandler(w http.ResponseWriter, r *http.Request) {
	// Create a problem detail
	p := problem.New(
		problem.NotFound,
		problem.WithDetail("The requested user with id 12345 was not found."),
		problem.WithInstance(r.URL.Path),
	)

	// Write directly to the HTTP response
	p.Write(w)
}
```

This automatically sets the `Content-Type: application/problem+json` header and writes:

```json
{
  "type": "about:blank",
  "title": "Not Found",
  "status": 404,
  "detail": "The requested user with id 12345 was not found.",
  "instance": "/users/12345"
}
```

## Adding Custom Extensions

RFC 7807 allows you to add custom fields to provide domain-specific context.

```go
p := problem.New(
	problem.UnprocessableEntity,
	problem.WithDetail("You do not have enough credit."),
	problem.WithExtension("balance", 30),
	problem.WithExtensions(map[string]any{
		"currency": "USD",
		"account": "/account/12345",
	}),
)
```

Generates:

```json
{
  "type": "about:blank",
  "title": "Unprocessable Entity",
  "status": 422,
  "detail": "You do not have enough credit.",
  "balance": 30,
  "currency": "USD",
  "account": "/account/12345"
}
```

## Creating Custom Problem Types

You are encouraged to define your own problem types for domain-specific errors.

```go
var OutOfCredit = problem.TypeTemplate{
	Type:   "https://example.com/probs/out-of-credit",
	Title:  "You do not have enough credit.",
	Status: http.StatusForbidden,
}

// Later in a handler:
p := problem.New(OutOfCredit, problem.WithDetail("Current balance is 30, but that costs 50."))
```

## Error Wrapping & Unwrapping

Because `*problem.Problem` implements the `error` interface, it can be seamlessly passed through your service layer and wrapped.

```go
// In your repository/service layer
func queryDB() error {
	err := db.Query(...) // let's imagine this fails

	// Wrap the internal error inside a Problem
	return problem.Wrap(err, problem.InternalServerError,
		problem.WithDetail("Database timeout"),
		problem.WithExtension("trace_id", "req-789"),
	)
}

// In your HTTP handler
func handler(w http.ResponseWriter, r *http.Request) {
	err := queryDB()
	if err != nil {
		// Use errors.As or the provided IsProblem helper to extract it
		if p, ok := problem.IsProblem(err); ok {
			p.Write(w)
			return
		}

		// Fallback for unknown errors
		problem.New(problem.InternalServerError).Write(w)
	}
}
```

## Example Application

Check out the [example/kanban](./example/kanban) directory for a complete, runnable example of using the `problem` library.
It features:
- **Clean Architecture**: Clear separation of domain rules, repositories, and HTTP handlers.
- **Clear Error Boundaries**: Mapping standard domain errors to RFC 7807 problem details.
- **Go Chi & Govalidator**: Demonstrates using the library with widely used community packages, mapping validation errors to `422 Unprocessable Entity` problem extensions.

