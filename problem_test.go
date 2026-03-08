package problem_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/semmidev/problem"
)

func TestProblemSerialization(t *testing.T) {
	t.Run("marshaling", func(t *testing.T) {
		p := problem.New(
			problem.BadRequest,
			problem.WithDetail("Invalid input provided"),
			problem.WithInstance("/api/v1/users"),
			problem.WithExtension("trace_id", "req-1234"),
			problem.WithExtension("balance", 5000),
		)

		data, err := json.Marshal(p)
		if err != nil {
			t.Fatalf("failed to marshal: %v", err)
		}

		var raw map[string]any
		if err := json.Unmarshal(data, &raw); err != nil {
			t.Fatalf("failed to unmarshal back to map: %v", err)
		}

		// Ensure top-level fields match expectations
		if raw["type"] != "about:blank" {
			t.Errorf("expected type 'about:blank', got %v", raw["type"])
		}
		if raw["title"] != "Bad Request" {
			t.Errorf("expected title 'Bad Request', got %v", raw["title"])
		}
		if raw["status"].(float64) != float64(http.StatusBadRequest) {
			t.Errorf("expected status %v, got %v", http.StatusBadRequest, raw["status"])
		}
		if raw["detail"] != "Invalid input provided" {
			t.Errorf("expected detail, got %v", raw["detail"])
		}
		if raw["instance"] != "/api/v1/users" {
			t.Errorf("expected instance, got %v", raw["instance"])
		}
		if raw["trace_id"] != "req-1234" {
			t.Errorf("expected trace_id extension, got %v", raw["trace_id"])
		}
		if raw["balance"].(float64) != 5000 {
			t.Errorf("expected balance extension, got %v", raw["balance"])
		}
	})

	t.Run("unmarshaling", func(t *testing.T) {
		rawJSON := []byte(`{
			"type": "https://example.com/probs/out-of-credit",
			"title": "You do not have enough credit.",
			"status": 403,
			"detail": "Your current balance is 30, but that costs 50.",
			"instance": "/account/12345/msgs/abc",
			"balance": 30,
			"accounts": ["/account/12345", "/account/67890"]
		}`)

		var p problem.Problem
		if err := json.Unmarshal(rawJSON, &p); err != nil {
			t.Fatalf("failed to unmarshal: %v", err)
		}

		if p.Type != "https://example.com/probs/out-of-credit" {
			t.Errorf("bad type: %s", p.Type)
		}
		if p.Title != "You do not have enough credit." {
			t.Errorf("bad title: %s", p.Title)
		}
		if p.Status != 403 {
			t.Errorf("bad status: %d", p.Status)
		}
		if p.Detail != "Your current balance is 30, but that costs 50." {
			t.Errorf("bad detail: %s", p.Detail)
		}
		if p.Instance != "/account/12345/msgs/abc" {
			t.Errorf("bad instance: %s", p.Instance)
		}

		// Extensions
		if p.Extensions == nil {
			t.Fatal("expected extensions to be initialized")
		}
		if p.Extensions["balance"].(float64) != 30 {
			t.Errorf("bad extension balance: %v", p.Extensions["balance"])
		}

		// check accounts list length
		accounts, ok := p.Extensions["accounts"].([]any)
		if !ok || len(accounts) != 2 {
			t.Errorf("bad extension accounts array: %v", p.Extensions["accounts"])
		}
	})
}

func TestProblemErrorWrapping(t *testing.T) {
	origErr := errors.New("underlying generic network timeout")

	p := problem.Wrap(origErr, problem.GatewayTimeout, problem.WithDetail("Upstream service failed"))

	// Implement Standard Error()
	if p.Error() == "" {
		t.Error("error string should not be empty")
	}

	// Unwrap
	if unwrapped := p.Unwrap(); unwrapped != origErr {
		t.Errorf("expected to unwrap origErr, got: %v", unwrapped)
	}

	// errors.Is check
	if !errors.Is(p, origErr) {
		t.Error("errors.Is should return true for underlying error")
	}

	// errors.As check (getting Problem back out of an error)
	var asErr *problem.Problem
	if !errors.As(p, &asErr) {
		t.Error("errors.As should extract Problem from itself")
	}

	// Double-wrapping scenario (simulating fmt.Errorf)
	wrap2 := fmt.Errorf("adding more context: %w", p)
	if !errors.As(wrap2, &asErr) {
		t.Error("errors.As should extract Problem from a wrapped generic error")
	}

	// IsProblem convenience
	if extracted, ok := problem.IsProblem(wrap2); !ok || extracted == nil {
		t.Error("IsProblem should successfully extract the Problem from a wrapped error")
	}
}

func TestProblemHTTPWrite(t *testing.T) {
	p := problem.New(problem.InternalServerError, problem.WithDetail("DB connection failed"))

	rec := httptest.NewRecorder()
	if err := p.Write(rec); err != nil {
		t.Fatalf("Write should not return error: %v", err)
	}

	res := rec.Result()
	if res.StatusCode != 500 {
		t.Errorf("bad status code, expected 500 got %d", res.StatusCode)
	}

	contentType := res.Header.Get("Content-Type")
	if contentType != "application/problem+json; charset=utf-8" {
		t.Errorf("bad content type: %s", contentType)
	}

	if snif := res.Header.Get("X-Content-Type-Options"); snif != "nosniff" {
		t.Errorf("expected nosniff, got %v", snif)
	}
}
