package problem

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Problem represents an RFC 7807 Problem Details object.
type Problem struct {
	// Type is a URI reference that identifies the problem type.
	Type string `json:"type"`

	// Title is a short, human-readable summary of the problem type.
	Title string `json:"title"`

	// Status is the HTTP status code for this occurrence of the problem.
	Status int `json:"status"`

	// Detail is a human-readable explanation specific to this occurrence of the problem.
	Detail string `json:"detail,omitempty"`

	// Instance is a URI reference that identifies the specific occurrence of the problem.
	Instance string `json:"instance,omitempty"`

	// Extensions contains additional properties beyond the standard RFC 7807 fields.
	Extensions map[string]any `json:"-"`

	// err is the underlying error that caused this problem, if any.
	// This makes Problem compatible with Go 1.13+ error wrapping.
	err error
}

// MarshalJSON implements custom JSON marshaling to flatten the Extensions
// into the top-level object, as required by RFC 7807.
func (p Problem) MarshalJSON() ([]byte, error) {
	// Start with standard fields
	base := map[string]any{
		"type":   p.Type,
		"title":  p.Title,
		"status": p.Status,
	}

	if p.Detail != "" {
		base["detail"] = p.Detail
	}
	if p.Instance != "" {
		base["instance"] = p.Instance
	}

	// Merge extensions to the top level, avoiding overwriting standard fields
	protectedFields := map[string]bool{
		"type": true, "title": true, "status": true,
		"detail": true, "instance": true,
	}

	for k, v := range p.Extensions {
		if !protectedFields[k] {
			base[k] = v
		}
	}

	return json.Marshal(base)
}

// UnmarshalJSON implements custom JSON unmarshaling to extract non-standard
// fields into the Extensions map.
func (p *Problem) UnmarshalJSON(data []byte) error {
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	p.Extensions = make(map[string]any)
	protectedFields := map[string]bool{
		"type": true, "title": true, "status": true,
		"detail": true, "instance": true,
	}

	for k, v := range raw {
		switch k {
		case "type":
			if s, ok := v.(string); ok {
				p.Type = s
			}
		case "title":
			if s, ok := v.(string); ok {
				p.Title = s
			}
		case "status":
			if f, ok := v.(float64); ok { // JSON numbers decode to float64 by default
				p.Status = int(f)
			}
		case "detail":
			if s, ok := v.(string); ok {
				p.Detail = s
			}
		case "instance":
			if s, ok := v.(string); ok {
				p.Instance = s
			}
		default:
			if !protectedFields[k] {
				p.Extensions[k] = v
			}
		}
	}
	return nil
}

// Error implements the standard Go `error` interface.
func (p *Problem) Error() string {
	if p.err != nil {
		return fmt.Sprintf("[%d] %s: %s: %v", p.Status, p.Type, p.Title, p.err)
	}
	if p.Detail != "" {
		return fmt.Sprintf("[%d] %s: %s", p.Status, p.Type, p.Detail)
	}
	return fmt.Sprintf("[%d] %s: %s", p.Status, p.Type, p.Title)
}

// Unwrap makes Problem compatible with Go 1.13+ error wrapping.
// It returns the underlying error if one was optionally provided.
func (p *Problem) Unwrap() error {
	return p.err
}

// IsProblem checks if an error is or wraps a *Problem details object.
func IsProblem(err error) (*Problem, bool) {
	var p *Problem
	if err != nil && errors.As(err, &p) {
		return p, true
	}
	return nil, false
}
