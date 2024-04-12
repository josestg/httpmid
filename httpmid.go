package httpmid

import (
	"net/http"
)

// Middleware is the signature of the middleware.
// The middleware decorates the given handler with additional behavior.
type Middleware func(http.Handler) http.Handler

// Then is a syntactic sugar for applying the middleware to the handler.
func (m Middleware) Then(h http.Handler) http.Handler { return m(h) }

// Apply applies middleware stack to the handler.
// Equivalent with Stack(xs...).Then(h)
func Apply(h http.Handler, xs ...Middleware) http.Handler {
	return fold(xs).Then(h)
}

// Stack stacks the middleware in the order they are passed.
// The first middleware passed will be the outermost layer of the middleware stack, and the last one will be the
// innermost. For example:
//
//	Stack(m1, m2, m3).Then(handler)
//
// will be equivalent with following middleware stack:
//
//	m1 begin:
//		m2 begin:
//			m3 begin:
//				handler
//			m3 end
//		m2 end
//	m1 end
func Stack(xs ...Middleware) Middleware { return fold(xs) }

// fold folds set of middleware into a single middleware.
func fold(xs []Middleware) Middleware {
	return func(h http.Handler) http.Handler {
		f := h
		for i := len(xs) - 1; i >= 0; i-- {
			f = xs[i](f)
		}
		return f
	}
}
