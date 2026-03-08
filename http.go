package problem

import (
	"encoding/json"
	"net/http"
)

// Write automatically writes the Problem Details as a JSON HTTP response
// with the correct Application/Problem+JSON content type.
func (p *Problem) Write(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/problem+json; charset=utf-8")
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

// Handler is a generic middleware interface or handler concept, but realistically
// you just need `problem.Write(w)`.
// We can provide a functional wrapper if you wanted, but usually `Write` is enough.
