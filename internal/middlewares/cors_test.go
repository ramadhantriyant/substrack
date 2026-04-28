package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCORS(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	corsHandler := CORS(handler)

	tests := []struct {
		name         string
		method       string
		wantStatus   int
		checkHeaders bool
		wantHeaders  map[string]string
	}{
		{
			name:         "GET request",
			method:       http.MethodGet,
			wantStatus:   http.StatusOK,
			checkHeaders: true,
			wantHeaders: map[string]string{
				"Access-Control-Allow-Origin":      "*",
				"Access-Control-Allow-Methods":     "GET, POST, PUT, DELETE, OPTIONS, PATCH",
				"Access-Control-Allow-Headers":     "Content-Type, Authorization",
				"Access-Control-Allow-Credentials": "true",
			},
		},
		{
			name:         "POST request",
			method:       http.MethodPost,
			wantStatus:   http.StatusOK,
			checkHeaders: true,
			wantHeaders: map[string]string{
				"Access-Control-Allow-Origin":      "*",
				"Access-Control-Allow-Methods":     "GET, POST, PUT, DELETE, OPTIONS, PATCH",
				"Access-Control-Allow-Headers":     "Content-Type, Authorization",
				"Access-Control-Allow-Credentials": "true",
			},
		},
		{
			name:         "OPTIONS request (preflight)",
			method:       http.MethodOptions,
			wantStatus:   http.StatusOK,
			checkHeaders: true,
			wantHeaders: map[string]string{
				"Access-Control-Allow-Origin":      "*",
				"Access-Control-Allow-Methods":     "GET, POST, PUT, DELETE, OPTIONS, PATCH",
				"Access-Control-Allow-Headers":     "Content-Type, Authorization",
				"Access-Control-Allow-Credentials": "true",
			},
		},
		{
			name:         "PUT request",
			method:       http.MethodPut,
			wantStatus:   http.StatusOK,
			checkHeaders: true,
			wantHeaders: map[string]string{
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Methods": "GET, POST, PUT, DELETE, OPTIONS, PATCH",
			},
		},
		{
			name:         "DELETE request",
			method:       http.MethodDelete,
			wantStatus:   http.StatusOK,
			checkHeaders: true,
			wantHeaders: map[string]string{
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Methods": "GET, POST, PUT, DELETE, OPTIONS, PATCH",
			},
		},
		{
			name:         "PATCH request",
			method:       http.MethodPatch,
			wantStatus:   http.StatusOK,
			checkHeaders: true,
			wantHeaders: map[string]string{
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Methods": "GET, POST, PUT, DELETE, OPTIONS, PATCH",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/test", nil)
			rr := httptest.NewRecorder()

			corsHandler.ServeHTTP(rr, req)

			if rr.Code != tt.wantStatus {
				t.Errorf("Status code = %v, want %v", rr.Code, tt.wantStatus)
			}

			if tt.checkHeaders {
				for header, wantValue := range tt.wantHeaders {
					gotValue := rr.Header().Get(header)
					if gotValue != wantValue {
						t.Errorf("Header %s = %v, want %v", header, gotValue, wantValue)
					}
				}
			}
		})
	}
}

func TestCORSOPTIONSRequest(t *testing.T) {
	// Test that OPTIONS request returns early without calling the handler
	handlerCalled := false
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		w.WriteHeader(http.StatusOK)
	})

	corsHandler := CORS(handler)

	req := httptest.NewRequest(http.MethodOptions, "/test", nil)
	rr := httptest.NewRecorder()

	corsHandler.ServeHTTP(rr, req)

	if handlerCalled {
		t.Error("Handler was called for OPTIONS request - should return early")
	}

	if rr.Code != http.StatusOK {
		t.Errorf("Status code = %v, want %v", rr.Code, http.StatusOK)
	}

	// Check CORS headers are set even for OPTIONS
	if origin := rr.Header().Get("Access-Control-Allow-Origin"); origin != "*" {
		t.Errorf("Access-Control-Allow-Origin = %v, want *", origin)
	}
}

func TestCORSPassesToHandler(t *testing.T) {
	// Test that non-OPTIONS requests pass through to handler
	handlerCalled := false
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("response body"))
	})

	corsHandler := CORS(handler)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder()

	corsHandler.ServeHTTP(rr, req)

	if !handlerCalled {
		t.Error("Handler was not called for GET request")
	}

	if rr.Code != http.StatusCreated {
		t.Errorf("Status code = %v, want %v", rr.Code, http.StatusCreated)
	}

	if rr.Body.String() != "response body" {
		t.Errorf("Body = %v, want response body", rr.Body.String())
	}
}
