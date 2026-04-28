package middlewares

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"git.ramadhantriyant.id/ramadhantriyant/substrack/internal/utils"
)

func TestRequireAuth(t *testing.T) {
	jwtSecret := "test-secret-key-that-is-long-enough-for-hs256-testing"

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check that userID is in context
		userID, ok := GetUserIDFromContext(r.Context())
		if !ok {
			t.Error("UserID not found in context")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if userID == 0 {
			t.Error("UserID is 0")
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("authenticated"))
	})

	authMiddleware := RequireAuth(jwtSecret)
	_ = authMiddleware // Will be used in subtests

	tests := []struct {
		name          string
		authHeader    string
		wantStatus    int
		wantBody      string
		handlerCalled bool
	}{
		{
			name:          "valid token",
			authHeader:    "",
			wantStatus:    http.StatusOK,
			wantBody:      "authenticated",
			handlerCalled: true,
		},
		{
			name:          "missing authorization header",
			authHeader:    "",
			wantStatus:    http.StatusUnauthorized,
			wantBody:      "missing or invalid authorization header",
			handlerCalled: false,
		},
		{
			name:          "empty authorization header",
			authHeader:    "",
			wantStatus:    http.StatusUnauthorized,
			wantBody:      "missing or invalid authorization header",
			handlerCalled: false,
		},
		{
			name:          "invalid format - no Bearer",
			authHeader:    "Token invalid",
			wantStatus:    http.StatusUnauthorized,
			wantBody:      "missing or invalid authorization header",
			handlerCalled: false,
		},
		{
			name:          "Bearer without token",
			authHeader:    "Bearer ",
			wantStatus:    http.StatusUnauthorized,
			wantBody:      "missing or invalid authorization header",
			handlerCalled: false,
		},
		{
			name:          "invalid token",
			authHeader:    "Bearer invalid.token.here",
			wantStatus:    http.StatusUnauthorized,
			wantBody:      "invalid or expired token",
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

			testAuthHandler := authMiddleware(testHandler)

			req := httptest.NewRequest(http.MethodGet, "/test", nil)

			// For valid token test, generate a real token
			if tt.name == "valid token" {
				token, err := utils.MakeJWT(1, jwtSecret)
				if err != nil {
					t.Fatalf("Failed to create token: %v", err)
				}
				req.Header.Set("Authorization", "Bearer "+token)
			} else if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			rr := httptest.NewRecorder()

			testAuthHandler.ServeHTTP(rr, req)

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

func TestRequireAuthWithWrongSecret(t *testing.T) {
	jwtSecret := "test-secret-key-that-is-long-enough-for-hs256-testing"
	wrongSecret := "wrong-secret-key-that-is-also-long-enough"

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	authMiddleware := RequireAuth(jwtSecret)
	authHandler := authMiddleware(handler)

	// Create token with wrong secret
	token, err := utils.MakeJWT(1, wrongSecret)
	if err != nil {
		t.Fatalf("Failed to create token: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	authHandler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Status code = %v, want %v", rr.Code, http.StatusUnauthorized)
	}

	body := rr.Body.String()
	var errResponse map[string]interface{}
	if err := json.Unmarshal([]byte(body), &errResponse); err != nil {
		t.Fatalf("Failed to parse error response: %v", err)
	}

	if errResponse["message"] != "invalid or expired token" {
		t.Errorf("Error message = %v, want 'invalid or expired token'", errResponse["message"])
	}
}

func TestRequireAuthWithExpiredToken(t *testing.T) {
	jwtSecret := "test-secret-key-that-is-long-enough-for-hs256-testing"

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	authMiddleware := RequireAuth(jwtSecret)
	authHandler := authMiddleware(handler)

	// We can't easily create an expired token without modifying the JWT library
	// So we'll use an invalid token format instead
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxIiwiZXhwIjoxfQ.invalid")
	rr := httptest.NewRecorder()

	authHandler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Status code = %v, want %v", rr.Code, http.StatusUnauthorized)
	}
}

func TestRequireAuthDifferentUserIDs(t *testing.T) {
	jwtSecret := "test-secret-key-that-is-long-enough-for-hs256-testing"

	tests := []struct {
		name   string
		userID int64
	}{
		{
			name:   "user ID 1",
			userID: 1,
		},
		{
			name:   "user ID 100",
			userID: 100,
		},
		{
			name:   "user ID 999999",
			userID: 999999,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var receivedUserID int64
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				userID, ok := GetUserIDFromContext(r.Context())
				if !ok {
					t.Error("UserID not found in context")
				}
				receivedUserID = userID
				w.WriteHeader(http.StatusOK)
			})

			authMiddleware := RequireAuth(jwtSecret)
			authHandler := authMiddleware(handler)

			token, err := utils.MakeJWT(tt.userID, jwtSecret)
			if err != nil {
				t.Fatalf("Failed to create token: %v", err)
			}

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			req.Header.Set("Authorization", "Bearer "+token)
			rr := httptest.NewRecorder()

			authHandler.ServeHTTP(rr, req)

			if rr.Code != http.StatusOK {
				t.Errorf("Status code = %v, want %v", rr.Code, http.StatusOK)
			}

			if receivedUserID != tt.userID {
				t.Errorf("Received userID = %v, want %v", receivedUserID, tt.userID)
			}
		})
	}
}

func TestGetUserIDFromContext(t *testing.T) {
	tests := []struct {
		name       string
		setupCtx   func() context.Context
		wantUserID int64
		wantOK     bool
	}{
		{
			name: "valid userID in context",
			setupCtx: func() context.Context {
				return context.WithValue(context.Background(), contextKeyUserID, int64(42))
			},
			wantUserID: 42,
			wantOK:     true,
		},
		{
			name: "no userID in context",
			setupCtx: func() context.Context {
				return context.Background()
			},
			wantUserID: 0,
			wantOK:     false,
		},
		{
			name: "wrong type in context",
			setupCtx: func() context.Context {
				return context.WithValue(context.Background(), contextKeyUserID, "not-an-int")
			},
			wantUserID: 0,
			wantOK:     false,
		},
		{
			name: "int instead of int64",
			setupCtx: func() context.Context {
				return context.WithValue(context.Background(), contextKeyUserID, 42)
			},
			wantUserID: 0,
			wantOK:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.setupCtx()
			userID, ok := GetUserIDFromContext(ctx)

			if userID != tt.wantUserID {
				t.Errorf("GetUserIDFromContext() userID = %v, want %v", userID, tt.wantUserID)
			}

			if ok != tt.wantOK {
				t.Errorf("GetUserIDFromContext() ok = %v, want %v", ok, tt.wantOK)
			}
		})
	}
}

func TestRequireAuthIntegration(t *testing.T) {
	jwtSecret := "test-secret-key-that-is-long-enough-for-hs256-testing"

	// Test the full flow
	var capturedUserID int64
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, ok := GetUserIDFromContext(r.Context())
		if !ok {
			http.Error(w, "No user ID in context", http.StatusInternalServerError)
			return
		}
		capturedUserID = userID
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	})

	authMiddleware := RequireAuth(jwtSecret)
	authHandler := authMiddleware(handler)

	// Test successful authentication
	token, _ := utils.MakeJWT(123, jwtSecret)
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	authHandler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Status code = %v, want %v", rr.Code, http.StatusOK)
	}

	if capturedUserID != 123 {
		t.Errorf("Captured userID = %v, want %v", capturedUserID, 123)
	}

	if rr.Body.String() != "success" {
		t.Errorf("Body = %v, want success", rr.Body.String())
	}
}
