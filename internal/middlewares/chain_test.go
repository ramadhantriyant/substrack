package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestChain(t *testing.T) {
	// Track the order of middleware execution
	var executionOrder []string

	// Create test middlewares
	middleware1 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			executionOrder = append(executionOrder, "middleware1-before")
			next.ServeHTTP(w, r)
			executionOrder = append(executionOrder, "middleware1-after")
		})
	}

	middleware2 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			executionOrder = append(executionOrder, "middleware2-before")
			next.ServeHTTP(w, r)
			executionOrder = append(executionOrder, "middleware2-after")
		})
	}

	middleware3 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			executionOrder = append(executionOrder, "middleware3-before")
			next.ServeHTTP(w, r)
			executionOrder = append(executionOrder, "middleware3-after")
		})
	}

	// Create a test handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		executionOrder = append(executionOrder, "handler")
		w.WriteHeader(http.StatusOK)
	})

	// Chain the middlewares
	chained := Chain(handler, middleware1, middleware2, middleware3)

	// Execute the request
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder()
	chained.ServeHTTP(rr, req)

	// Check the execution order
	expectedOrder := []string{
		"middleware1-before",
		"middleware2-before",
		"middleware3-before",
		"handler",
		"middleware3-after",
		"middleware2-after",
		"middleware1-after",
	}

	if len(executionOrder) != len(expectedOrder) {
		t.Errorf("Execution order length = %v, want %v", len(executionOrder), len(expectedOrder))
	}

	for i, expected := range expectedOrder {
		if i >= len(executionOrder) {
			t.Errorf("Missing execution step at index %d: expected %s", i, expected)
			continue
		}
		if executionOrder[i] != expected {
			t.Errorf("Execution order[%d] = %v, want %v", i, executionOrder[i], expected)
		}
	}

	if rr.Code != http.StatusOK {
		t.Errorf("Status code = %v, want %v", rr.Code, http.StatusOK)
	}
}

func TestChainNoMiddleware(t *testing.T) {
	// Test that Chain works with no middlewares
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("OK"))
	})

	chained := Chain(handler)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder()
	chained.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("Status code = %v, want %v", rr.Code, http.StatusCreated)
	}

	if rr.Body.String() != "OK" {
		t.Errorf("Body = %v, want OK", rr.Body.String())
	}
}

func TestChainSingleMiddleware(t *testing.T) {
	var called bool

	middleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			next.ServeHTTP(w, r)
		})
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	chained := Chain(handler, middleware)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder()
	chained.ServeHTTP(rr, req)

	if !called {
		t.Error("Middleware was not called")
	}

	if rr.Code != http.StatusOK {
		t.Errorf("Status code = %v, want %v", rr.Code, http.StatusOK)
	}
}

func TestChainModifiesRequest(t *testing.T) {
	// Test that middlewares can modify the request
	addHeaderMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Header.Set("X-Test-Header", "test-value")
			next.ServeHTTP(w, r)
		})
	}

	var receivedHeader string
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedHeader = r.Header.Get("X-Test-Header")
		w.WriteHeader(http.StatusOK)
	})

	chained := Chain(handler, addHeaderMiddleware)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder()
	chained.ServeHTTP(rr, req)

	if receivedHeader != "test-value" {
		t.Errorf("Received header = %v, want test-value", receivedHeader)
	}
}
