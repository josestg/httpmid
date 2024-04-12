package httpmid

import (
	"net/http"
	"testing"
)

func TestStackAndApply(t *testing.T) {
	trace := make([]int, 0)
	m1 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			trace = append(trace, -1)
			defer func() {
				trace = append(trace, 1)
			}()
			next.ServeHTTP(w, r)
		})
	}

	m2 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			trace = append(trace, -2)
			defer func() {
				trace = append(trace, 2)
			}()
			next.ServeHTTP(w, r)
		})
	}

	m3 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			trace = append(trace, -3)
			defer func() {
				trace = append(trace, 3)
			}()
			next.ServeHTTP(w, r)
		})
	}

	h := func(w http.ResponseWriter, r *http.Request) {
		trace = append(trace, 0)
	}

	expectedTrace := []int{-3, -2, -1, 0, 1, 2, 3}
	Stack(m3, m2, m1).Then(http.HandlerFunc(h)).ServeHTTP(nil, nil)
	if len(trace) != len(expectedTrace) {
		t.Fatalf("expected trace to be %v, but got %v", expectedTrace, trace)
	}

	trace = make([]int, 0)
	Apply(http.HandlerFunc(h), m3, m2, m1).ServeHTTP(nil, nil)
	if len(trace) != len(expectedTrace) {
		t.Fatalf("expected trace to be %v, but got %v", expectedTrace, trace)
	}
}
