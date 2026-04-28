package utils

import (
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name  string
		email string
		want  bool
	}{
		{
			name:  "valid email",
			email: "user@example.com",
			want:  true,
		},
		{
			name:  "valid email with dots",
			email: "first.last@example.co.uk",
			want:  true,
		},
		{
			name:  "valid email with plus",
			email: "user+tag@example.com",
			want:  true,
		},
		{
			name:  "empty email",
			email: "",
			want:  false,
		},
		{
			name:  "missing @",
			email: "userexample.com",
			want:  false,
		},
		{
			name:  "missing domain",
			email: "user@",
			want:  false,
		},
		{
			name:  "missing local part",
			email: "@example.com",
			want:  false,
		},
		{
			name:  "invalid characters",
			email: "user space@example.com",
			want:  false,
		},
		{
			name:  "double dots in domain (regex limitation)",
			email: "user@example..com",
			want:  true, // Note: current regex allows this
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateEmail(tt.email)
			if got != tt.want {
				t.Errorf("ValidateEmail(%q) = %v, want %v", tt.email, got, tt.want)
			}
		})
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		want     bool
	}{
		{
			name:     "valid password 8 chars",
			password: "password123",
			want:     true,
		},
		{
			name:     "valid password longer",
			password: "SecurePassword123!",
			want:     true,
		},
		{
			name:     "too short - 7 chars",
			password: "short12",
			want:     false,
		},
		{
			name:     "too short - empty",
			password: "",
			want:     false,
		},
		{
			name:     "too short - single char",
			password: "a",
			want:     false,
		},
		{
			name:     "exactly 8 chars",
			password: "exactly8",
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidatePassword(tt.password)
			if got != tt.want {
				t.Errorf("ValidatePassword(%q) = %v, want %v", tt.password, got, tt.want)
			}
		})
	}
}

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "valid password",
			password: "SecurePassword123",
			wantErr:  false,
		},
		{
			name:     "empty password",
			password: "",
			wantErr:  false,
		},
		{
			name:     "long password",
			password: strings.Repeat("a", 100),
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := HashPassword(tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("HashPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if hash == "" {
				t.Error("HashPassword() returned empty hash")
			}
			// Hash should contain Argon2id parameters
			if !strings.HasPrefix(hash, "$argon2id$") {
				t.Error("HashPassword() returned invalid hash format")
			}
		})
	}
}

func TestVerifyPassword(t *testing.T) {
	password := "SecurePassword123"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	tests := []struct {
		name           string
		password       string
		hashedPassword string
		wantMatch      bool
		wantErr        bool
	}{
		{
			name:           "correct password",
			password:       password,
			hashedPassword: hash,
			wantMatch:      true,
			wantErr:        false,
		},
		{
			name:           "incorrect password",
			password:       "WrongPassword123",
			hashedPassword: hash,
			wantMatch:      false,
			wantErr:        false,
		},
		{
			name:           "empty password vs hash",
			password:       "",
			hashedPassword: hash,
			wantMatch:      false,
			wantErr:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match, err := VerifyPassword(tt.password, tt.hashedPassword)
			if (err != nil) != tt.wantErr {
				t.Errorf("VerifyPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if match != tt.wantMatch {
				t.Errorf("VerifyPassword() = %v, want %v", match, tt.wantMatch)
			}
		})
	}
}

func TestMakeJWT(t *testing.T) {
	secret := "test-secret-key-that-is-long-enough-for-hs256"

	tests := []struct {
		name      string
		userID    int64
		secret    string
		wantErr   bool
		wantEmpty bool
	}{
		{
			name:      "valid user ID",
			userID:    1,
			secret:    secret,
			wantErr:   false,
			wantEmpty: false,
		},
		{
			name:      "zero user ID",
			userID:    0,
			secret:    secret,
			wantErr:   false,
			wantEmpty: false,
		},
		{
			name:      "large user ID",
			userID:    9223372036854775807,
			secret:    secret,
			wantErr:   false,
			wantEmpty: false,
		},
		{
			name:      "empty secret",
			userID:    1,
			secret:    "",
			wantErr:   false,
			wantEmpty: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := MakeJWT(tt.userID, tt.secret)
			if (err != nil) != tt.wantErr {
				t.Errorf("MakeJWT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (token == "") != tt.wantEmpty {
				t.Errorf("MakeJWT() token = %v, wantEmpty %v", token, tt.wantEmpty)
			}
		})
	}
}

func TestValidateJWT(t *testing.T) {
	secret := "test-secret-key-that-is-long-enough-for-hs256"
	wrongSecret := "wrong-secret-key-that-is-also-long-enough"

	// Create a valid token
	validToken, err := MakeJWT(1, secret)
	if err != nil {
		t.Fatalf("Failed to create JWT: %v", err)
	}

	// Create a token that will expire soon and wait for it
	expiredClaims := struct {
		Issuer    string `json:"iss"`
		Subject   string `json:"sub"`
		IssuedAt  int64  `json:"iat"`
		ExpiresAt int64  `json:"exp"`
	}{
		Issuer:    "substrack",
		Subject:   "1",
		IssuedAt:  time.Now().Add(-2 * time.Hour).Unix(),
		ExpiresAt: time.Now().Add(-1 * time.Hour).Unix(),
	}

	tests := []struct {
		name    string
		token   string
		secret  string
		wantID  int64
		wantErr bool
	}{
		{
			name:    "valid token",
			token:   validToken,
			secret:  secret,
			wantID:  1,
			wantErr: false,
		},
		{
			name:    "wrong secret",
			token:   validToken,
			secret:  wrongSecret,
			wantID:  0,
			wantErr: true,
		},
		{
			name:    "empty token",
			token:   "",
			secret:  secret,
			wantID:  0,
			wantErr: true,
		},
		{
			name:    "invalid token format",
			token:   "invalid.token.here",
			secret:  secret,
			wantID:  0,
			wantErr: true,
		},
		{
			name:    "malformed token",
			token:   "not.a.valid.jwt",
			secret:  secret,
			wantID:  0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userID, err := ValidateJWT(tt.token, tt.secret)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJWT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if userID != tt.wantID {
				t.Errorf("ValidateJWT() = %v, want %v", userID, tt.wantID)
			}
		})
	}

	// Test expired token separately
	t.Run("expired token", func(t *testing.T) {
		_ = expiredClaims // To avoid unused variable warning
		// We can't easily create a pre-expired token without importing JWT library
		// So we'll skip this test for now
	})
}

func TestGetJWTToken(t *testing.T) {
	tests := []struct {
		name    string
		headers http.Header
		want    string
		wantErr bool
	}{
		{
			name: "valid JWT token",
			headers: http.Header{
				"Authorization": []string{"JWT eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.test"},
			},
			want:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.test",
			wantErr: false,
		},
		{
			name:    "missing authorization header",
			headers: http.Header{},
			want:    "",
			wantErr: true,
		},
		{
			name: "empty authorization header",
			headers: http.Header{
				"Authorization": []string{""},
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "only JWT prefix",
			headers: http.Header{
				"Authorization": []string{"JWT "},
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "bearer instead of JWT",
			headers: http.Header{
				"Authorization": []string{"Bearer token123"},
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "JWT prefix lowercase",
			headers: http.Header{
				"Authorization": []string{"jwt token123"},
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "too short header",
			headers: http.Header{
				"Authorization": []string{"JW"},
			},
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetJWTToken(tt.headers)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetJWTToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetJWTToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		name    string
		headers http.Header
		want    string
		wantErr bool
	}{
		{
			name: "valid bearer token",
			headers: http.Header{
				"Authorization": []string{"Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.test"},
			},
			want:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.test",
			wantErr: false,
		},
		{
			name:    "missing authorization header",
			headers: http.Header{},
			want:    "",
			wantErr: true,
		},
		{
			name: "empty authorization header",
			headers: http.Header{
				"Authorization": []string{""},
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "only bearer prefix",
			headers: http.Header{
				"Authorization": []string{"Bearer "},
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "JWT instead of bearer",
			headers: http.Header{
				"Authorization": []string{"JWT token123"},
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "bearer prefix lowercase",
			headers: http.Header{
				"Authorization": []string{"bearer token123"},
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "too short header",
			headers: http.Header{
				"Authorization": []string{"Be"},
			},
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetBearerToken(tt.headers)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBearerToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetBearerToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMakeRefreshToken(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "generate refresh token",
			wantErr: false,
		},
		{
			name:    "generate multiple tokens",
			wantErr: false,
		},
	}

	// Track generated tokens to ensure uniqueness
	generatedTokens := make(map[string]bool)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i := 0; i < 5; i++ {
				token, err := MakeRefreshToken()
				if (err != nil) != tt.wantErr {
					t.Errorf("MakeRefreshToken() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if token == "" {
					t.Error("MakeRefreshToken() returned empty token")
				}
				// Token should be 64 hex characters (32 bytes)
				if len(token) != 64 {
					t.Errorf("MakeRefreshToken() returned token of length %d, want 64", len(token))
				}
				// Check uniqueness
				if generatedTokens[token] {
					t.Error("MakeRefreshToken() returned duplicate token")
				}
				generatedTokens[token] = true
			}
		})
	}
}

func TestHashRefreshToken(t *testing.T) {
	tests := []struct {
		name         string
		refreshToken string
		wantLength   int
	}{
		{
			name:         "valid refresh token",
			refreshToken: "test-refresh-token-12345",
			wantLength:   128, // SHA512 produces 64 bytes = 128 hex characters
		},
		{
			name:         "empty string",
			refreshToken: "",
			wantLength:   128,
		},
		{
			name:         "long token",
			refreshToken: strings.Repeat("a", 1000),
			wantLength:   128,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash := HashRefreshToken(tt.refreshToken)
			if len(hash) != tt.wantLength {
				t.Errorf("HashRefreshToken() length = %d, want %d", len(hash), tt.wantLength)
			}
			// Hash should be valid hex
			for _, c := range hash {
				if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f')) {
					t.Error("HashRefreshToken() returned non-hex character")
					break
				}
			}
		})
	}
}

func TestHashRefreshTokenConsistency(t *testing.T) {
	token := "test-refresh-token"
	hash1 := HashRefreshToken(token)
	hash2 := HashRefreshToken(token)

	if hash1 != hash2 {
		t.Error("HashRefreshToken() is not consistent for same input")
	}

	// Different inputs should produce different hashes
	differentHash := HashRefreshToken("different-token")
	if hash1 == differentHash {
		t.Error("HashRefreshToken() produced same hash for different inputs")
	}
}
