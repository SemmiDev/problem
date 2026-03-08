package problem

import "fmt"

// Option is a functional option for customizing a Problem details object.
type Option func(*Problem)

// WithDetail sets the human-readable explanation specific to this occurrence of the problem.
func WithDetail(detail string) Option {
	return func(p *Problem) {
		p.Detail = detail
	}
}

// WithDetailf formats and sets the detail message.
func WithDetailf(format string, args ...any) Option {
	return func(p *Problem) {
		p.Detail = fmt.Sprintf(format, args...)
	}
}

// WithInstance sets the URI reference that identifies the specific occurrence of the problem.
func WithInstance(instance string) Option {
	return func(p *Problem) {
		p.Instance = instance
	}
}

// WithExtension sets a single extension field.
// According to RFC 7807, extension members are additional members within the Problem Details object.
func WithExtension(key string, value any) Option {
	return func(p *Problem) {
		if p.Extensions == nil {
			p.Extensions = make(map[string]any)
		}
		p.Extensions[key] = value
	}
}

// WithExtensions merges multiple extension fields at once.
func WithExtensions(ext map[string]any) Option {
	return func(p *Problem) {
		if p.Extensions == nil {
			p.Extensions = make(map[string]any)
		}
		for k, v := range ext {
			p.Extensions[k] = v
		}
	}
}

// WithErr sets the underlying error that caused this problem.
// This allows the Problem to wrap another error, supporting standard Go wrapping (errors.Is/As).
func WithErr(err error) Option {
	return func(p *Problem) {
		p.err = err
	}
}
