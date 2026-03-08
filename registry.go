package problem

import (
	"net/http"
)

// BaseURI is an optional prefix that can be globally modified
// to set the base URI for all automatically generated Problem Types.
// However, since this is a library, hardcoding isn't ideal,
// so standard problem types typically use "about:blank" or a common schema.
// We'll use "about:blank" for standard HTTP errors to keep it generic,
// but users can define their own Types.
var BaseURI = "about:blank"

// TypeTemplate defines a reusable structure for a specific type of problem.
type TypeTemplate struct {
	Type   string
	Title  string
	Status int
}

// New creates a new Problem instance based on a provided TypeTemplate,
// applying any functional Options.
func New(template TypeTemplate, opts ...Option) *Problem {
	p := &Problem{
		Type:   template.Type,
		Title:  template.Title,
		Status: template.Status,
	}
	for _, opt := range opts {
		opt(p)
	}

	// Default fallback: if Type is empty but we're creating from a standard HTTP status template,
	// we use a generic "about:blank" as recommended by RFC 7807 for simple HTTP codes.
	if p.Type == "" {
		p.Type = "about:blank"
	}

	return p
}

// Wrap is a convenience function that creates a new Problem from a template
// and automatically wraps the provided error.
// It is equivalent to New(template, WithErr(err)) with additional options.
func Wrap(err error, template TypeTemplate, opts ...Option) *Problem {
	allOpts := make([]Option, 0, len(opts)+1)
	allOpts = append(allOpts, WithErr(err))
	allOpts = append(allOpts, opts...)
	return New(template, allOpts...)
}

// ==== Standard HTTP Problem Templates ====
// These are pre-defined templates covering common HTTP 4xx and 5xx responses.

var (
	// BadRequest (HTTP 400)
	BadRequest = TypeTemplate{
		Type:   "about:blank",
		Title:  "Bad Request",
		Status: http.StatusBadRequest,
	}

	// Unauthorized (HTTP 401)
	Unauthorized = TypeTemplate{
		Type:   "about:blank",
		Title:  "Unauthorized",
		Status: http.StatusUnauthorized,
	}

	// PaymentRequired (HTTP 402)
	PaymentRequired = TypeTemplate{
		Type:   "about:blank",
		Title:  "Payment Required",
		Status: http.StatusPaymentRequired,
	}

	// Forbidden (HTTP 403)
	Forbidden = TypeTemplate{
		Type:   "about:blank",
		Title:  "Forbidden",
		Status: http.StatusForbidden,
	}

	// NotFound (HTTP 404)
	NotFound = TypeTemplate{
		Type:   "about:blank",
		Title:  "Not Found",
		Status: http.StatusNotFound,
	}

	// MethodNotAllowed (HTTP 405)
	MethodNotAllowed = TypeTemplate{
		Type:   "about:blank",
		Title:  "Method Not Allowed",
		Status: http.StatusMethodNotAllowed,
	}

	// NotAcceptable (HTTP 406)
	NotAcceptable = TypeTemplate{
		Type:   "about:blank",
		Title:  "Not Acceptable",
		Status: http.StatusNotAcceptable,
	}

	// RequestTimeout (HTTP 408)
	RequestTimeout = TypeTemplate{
		Type:   "about:blank",
		Title:  "Request Timeout",
		Status: http.StatusRequestTimeout,
	}

	// Conflict (HTTP 409)
	Conflict = TypeTemplate{
		Type:   "about:blank",
		Title:  "Conflict",
		Status: http.StatusConflict,
	}

	// Gone (HTTP 410)
	Gone = TypeTemplate{
		Type:   "about:blank",
		Title:  "Gone",
		Status: http.StatusGone,
	}

	// LengthRequired (HTTP 411)
	LengthRequired = TypeTemplate{
		Type:   "about:blank",
		Title:  "Length Required",
		Status: http.StatusLengthRequired,
	}

	// PayloadTooLarge (HTTP 413)
	PayloadTooLarge = TypeTemplate{
		Type:   "about:blank",
		Title:  "Payload Too Large",
		Status: http.StatusRequestEntityTooLarge,
	}

	// URITooLong (HTTP 414)
	URITooLong = TypeTemplate{
		Type:   "about:blank",
		Title:  "URI Too Long",
		Status: http.StatusRequestURITooLong,
	}

	// UnsupportedMediaType (HTTP 415)
	UnsupportedMediaType = TypeTemplate{
		Type:   "about:blank",
		Title:  "Unsupported Media Type",
		Status: http.StatusUnsupportedMediaType,
	}

	// UnprocessableEntity (HTTP 422)
	UnprocessableEntity = TypeTemplate{
		Type:   "about:blank",
		Title:  "Unprocessable Entity",
		Status: http.StatusUnprocessableEntity,
	}

	// TooManyRequests (HTTP 429)
	TooManyRequests = TypeTemplate{
		Type:   "about:blank",
		Title:  "Too Many Requests",
		Status: http.StatusTooManyRequests,
	}

	// InternalServerError (HTTP 500)
	InternalServerError = TypeTemplate{
		Type:   "about:blank",
		Title:  "Internal Server Error",
		Status: http.StatusInternalServerError,
	}

	// NotImplemented (HTTP 501)
	NotImplemented = TypeTemplate{
		Type:   "about:blank",
		Title:  "Not Implemented",
		Status: http.StatusNotImplemented,
	}

	// BadGateway (HTTP 502)
	BadGateway = TypeTemplate{
		Type:   "about:blank",
		Title:  "Bad Gateway",
		Status: http.StatusBadGateway,
	}

	// ServiceUnavailable (HTTP 503)
	ServiceUnavailable = TypeTemplate{
		Type:   "about:blank",
		Title:  "Service Unavailable",
		Status: http.StatusServiceUnavailable,
	}

	// GatewayTimeout (HTTP 504)
	GatewayTimeout = TypeTemplate{
		Type:   "about:blank",
		Title:  "Gateway Timeout",
		Status: http.StatusGatewayTimeout,
	}
)
