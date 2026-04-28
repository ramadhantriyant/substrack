package middlewares

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestShouldJSON(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	jsonHandler := ShouldJSON(handler)

	_ = jsonHandler // Will be used in subtests

	tests := []struct {
		name          string
		method        string
		contentType   string
		wantStatus    int
		wantBody      string
		handlerCalled bool
	}{
		{
			name:          "POST with application/json",
			method:        http.MethodPost,
			contentType:   "application/json",
			wantStatus:    http.StatusOK,
			wantBody:      "OK",
			handlerCalled: true,
		},
		{
			name:          "PUT with application/json",
			method:        http.MethodPut,
			contentType:   "application/json",
			wantStatus:    http.StatusOK,
			wantBody:      "OK",
			handlerCalled: true,
		},
		{
			name:          "GET without content type",
			method:        http.MethodGet,
			contentType:   "",
			wantStatus:    http.StatusOK,
			wantBody:      "OK",
			handlerCalled: true,
		},
		{
			name:          "GET with content type",
			method:        http.MethodGet,
			contentType:   "application/json",
			wantStatus:    http.StatusOK,
			wantBody:      "OK",
			handlerCalled: true,
		},
		{
			name:          "PATCH without content type",
			method:        http.MethodPatch,
			contentType:   "",
			wantStatus:    http.StatusOK,
			wantBody:      "OK",
			handlerCalled: true,
		},
		{
			name:          "DELETE without content type",
			method:        http.MethodDelete,
			contentType:   "",
			wantStatus:    http.StatusOK,
			wantBody:      "OK",
			handlerCalled: true,
		},
		{
			name:          "POST without content type",
			method:        http.MethodPost,
			contentType:   "",
			wantStatus:    http.StatusUnsupportedMediaType,
			wantBody:      "unsupported media type",
			handlerCalled: false,
		},
		{
			name:          "PUT with text/plain",
			method:        http.MethodPut,
			contentType:   "text/plain",
			wantStatus:    http.StatusUnsupportedMediaType,
			wantBody:      "unsupported media type",
			handlerCalled: false,
		},
		{
			name:          "POST with application/xml",
			method:        http.MethodPost,
			contentType:   "application/xml",
			wantStatus:    http.StatusUnsupportedMediaType,
			wantBody:      "unsupported media type",
			handlerCalled: false,
		},
		{
			name:          "POST with multipart/form-data",
			method:        http.MethodPost,
			contentType:   "multipart/form-data",
			wantStatus:    http.StatusUnsupportedMediaType,
			wantBody:      "unsupported media type",
			handlerCalled: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlerCalled := false
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				handlerCalled = true
				handler.ServeHTTP(w, r)
			})

			wrappedHandler := ShouldJSON(testHandler)

			req := httptest.NewRequest(tt.method, "/test", nil)
			if tt.contentType != "" {
				req.Header.Set("Content-Type", tt.contentType)
			}
			rr := httptest.NewRecorder()

			wrappedHandler.ServeHTTP(rr, req)

			if rr.Code != tt.wantStatus {
				t.Errorf("Status code = %v, want %v", rr.Code, tt.wantStatus)
			}

			if !strings.Contains(rr.Body.String(), tt.wantBody) {
				t.Errorf("Body = %v, want to contain %v", rr.Body.String(), tt.wantBody)
			}

			if handlerCalled != tt.handlerCalled {
				t.Errorf("Handler called = %v, want %v", handlerCalled, tt.handlerCalled)
			}
		})
	}
}

func TestShouldJSONContentTypeCase(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	jsonHandler := ShouldJSON(handler)

	tests := []struct {
		name        string
		contentType string
		wantStatus  int
	}{
		{
			name:        "lowercase application/json",
			contentType: "application/json",
			wantStatus:  http.StatusOK,
		},
		{
			name:        "uppercase APPLICATION/JSON",
			contentType: "APPLICATION/JSON",
			wantStatus:  http.StatusUnsupportedMediaType,
		},
		{
			name:        "mixed case Application/Json",
			contentType: "Application/Json",
			wantStatus:  http.StatusUnsupportedMediaType,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/test", nil)
			req.Header.Set("Content-Type", tt.contentType)
			rr := httptest.NewRecorder()

			jsonHandler.ServeHTTP(rr, req)

			if rr.Code != tt.wantStatus {
				t.Errorf("Status code = %v, want %v", rr.Code, tt.wantStatus)
			}
		})
	}
}

func TestShouldJSONWithCharset(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	jsonHandler := ShouldJSON(handler)

	// Content-Type with charset should not match exactly
	req := httptest.NewRequest(http.MethodPost, "/test", nil)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	rr := httptest.NewRecorder()

	jsonHandler.ServeHTTP(rr, req)

	// This will fail because the check is for exact match
	if rr.Code != http.StatusUnsupportedMediaType {
		t.Errorf("Status code = %v, want %v (Note: exact match required)", rr.Code, http.StatusUnsupportedMediaType)
	}
}
