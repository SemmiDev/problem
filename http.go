package problem

import (
	"encoding/json"
	"net/http"
)

// ContentType is the standard media type for Problem Details over HTTP.
const ContentType = "application/problem+json; charset=utf-8"

// Write automatically writes the Problem Details as a JSON HTTP response
// with the correct Application/Problem+JSON content type.
func (p *Problem) Write(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", ContentType)
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(p.Status)

	enc := json.NewEncoder(w)
	// Optionally prevent escaping HTML inside JSON to keep URI strings clean
	enc.SetEscapeHTML(false)

	if err := enc.Encode(p); err != nil {
		return err
	}
	return nil
}

// JSON returns the JSON encoding of the Problem.
// This is useful for frameworks like Gin, Echo, or Fiber where you might
// want to write the raw JSON bytes directly and set the Content-Type manually.
func (p *Problem) JSON() []byte {
	b, _ := json.Marshal(p)
	return b
}

// Headers returns a map of HTTP headers that should be set when returning this problem.
// This is convenient for translating to other framework's header structures.
func (p *Problem) Headers() map[string]string {
	return map[string]string{
		"Content-Type":           ContentType,
		"X-Content-Type-Options": "nosniff",
	}
}

// Map converts the Problem to a generic map map[string]any.
// This is particularly useful when returning Problem Details via frameworks
// that prefer or require maps for custom JSON serialization overrides.
func (p *Problem) Map() map[string]any {
	var m map[string]any
	// Reusing Marshal/Unmarshal avoids duplicating MarshalJSON logic
	// and preserves consistency with the struct's existing serialization.
	b, _ := json.Marshal(p)
	_ = json.Unmarshal(b, &m)
	return m
}
