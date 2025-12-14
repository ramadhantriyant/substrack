package middlewares

import "net/http"

// Middleware is a function that wraps an http.Handler
type Middleware func(http.Handler) http.Handler

// Chain applies middlewares in the order they are provided
// The first middleware in the slice will be the outermost (executed first)
func Chain(handler http.Handler, middlewares ...Middleware) http.Handler {
	// Apply middlewares in reverse order so the first one wraps everything
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}
