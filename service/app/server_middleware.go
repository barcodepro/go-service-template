package app

import (
	"net/http"
)

// Stack ...
type Stack interface {
	// Adds a middleware to the stack. MWs will be called in the same order that they are added, such that:
	//	Use(Request ID Middleware)
	//  Use(Request Timing Middleware)
	// would result in the request id middleware being the outermost layer, called first, before the timing middleware.
	Use(Middleware)

	// Wraps a given handle with the current stack from the result of Use() calls.
	Wrap(http.Handler) http.Handler
}

// Middleware ...
type Middleware func(http.Handler) http.Handler

type stack struct {
	middlewares []Middleware
}

// NewStack ...
func newMiddlewareStack() *stack {
	return &stack{
		middlewares: []Middleware{},
	}
}

// Use ...
func (s *stack) Use(mw Middleware) {
	s.middlewares = append(s.middlewares, mw)
}

// Wrap ...
func (s *stack) Wrap(fn http.Handler) http.Handler {
	l := len(s.middlewares)
	if l == 0 {
		return fn
	}

	// There is at least one item in the list. Starting with the last item, create the handler to be returned:
	var result http.Handler
	result = s.middlewares[l-1](fn)

	// Reverse through the stack for the remaining elements, and wrap the result with each layer:
	for i := 0; i < (l - 1); i++ {
		result = s.middlewares[l-(2+i)](result)
	}

	return result
}
