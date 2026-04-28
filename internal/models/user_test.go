package models

import (
	"encoding/json"
	"testing"
	"time"

	"git.ramadhantriyant.id/ramadhantriyant/substrack/internal/database"
)

func TestUserToResponse(t *testing.T) {
	now := time.Now()
	past := now.Add(-24 * time.Hour)

	tests := []struct {
		name string
		user database.User
		want UserResponse
	}{
		{
			name: "complete user",
			user: database.User{
				ID:        1,
				Email:     "test@example.com",
				Password:  "hashed_password_should_not_appear",
				Name:      "Test User",
				CreatedAt: &past,
				UpdatedAt: &now,
			},
			want: UserResponse{
				ID:        1,
				Email:     "test@example.com",
				Name:      "Test User",
				CreatedAt: &past,
				UpdatedAt: &now,
			},
		},
		{
			name: "user without timestamps",
			user: database.User{
				ID:        2,
				Email:     "minimal@example.com",
				Password:  "secret",
				Name:      "Minimal User",
				CreatedAt: nil,
				UpdatedAt: nil,
			},
			want: UserResponse{
				ID:        2,
				Email:     "minimal@example.com",
				Name:      "Minimal User",
				CreatedAt: nil,
				UpdatedAt: nil,
			},
		},
		{
			name: "zero ID user",
			user: database.User{
				ID:        0,
				Email:     "zero@example.com",
				Password:  "password",
				Name:      "Zero ID",
				CreatedAt: nil,
				UpdatedAt: nil,
			},
			want: UserResponse{
				ID:        0,
				Email:     "zero@example.com",
				Name:      "Zero ID",
				CreatedAt: nil,
				UpdatedAt: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := UserToResponse(tt.user)

			if got.ID != tt.want.ID {
				t.Errorf("UserToResponse() ID = %v, want %v", got.ID, tt.want.ID)
			}

			if got.Email != tt.want.Email {
				t.Errorf("UserToResponse() Email = %v, want %v", got.Email, tt.want.Email)
			}

			if got.Name != tt.want.Name {
				t.Errorf("UserToResponse() Name = %v, want %v", got.Name, tt.want.Name)
			}

			// Verify password is NOT in response
			// We can't check this directly since UserResponse doesn't have Password field
			// But we can verify by marshaling to JSON
			jsonBytes, err := json.Marshal(got)
			if err != nil {
				t.Fatalf("Failed to marshal UserResponse: %v", err)
			}

			if string(jsonBytes) == "" {
				t.Error("UserToResponse() produced empty JSON")
			}

			// Ensure password hash is not in the response
			if string(jsonBytes) != "" && len(jsonBytes) > 0 {
				var result map[string]interface{}
				if err := json.Unmarshal(jsonBytes, &result); err != nil {
					t.Fatalf("Failed to unmarshal JSON: %v", err)
				}
				if _, exists := result["password"]; exists {
					t.Error("UserToResponse() response contains password field - should be omitted")
				}
			}

			if got.CreatedAt != tt.want.CreatedAt {
				if tt.want.CreatedAt != nil && got.CreatedAt != nil {
					if !got.CreatedAt.Equal(*tt.want.CreatedAt) {
						t.Errorf("UserToResponse() CreatedAt = %v, want %v", got.CreatedAt, tt.want.CreatedAt)
					}
				} else if !(got.CreatedAt == nil && tt.want.CreatedAt == nil) {
					t.Errorf("UserToResponse() CreatedAt = %v, want %v", got.CreatedAt, tt.want.CreatedAt)
				}
			}

			if got.UpdatedAt != tt.want.UpdatedAt {
				if tt.want.UpdatedAt != nil && got.UpdatedAt != nil {
					if !got.UpdatedAt.Equal(*tt.want.UpdatedAt) {
						t.Errorf("UserToResponse() UpdatedAt = %v, want %v", got.UpdatedAt, tt.want.UpdatedAt)
					}
				} else if !(got.UpdatedAt == nil && tt.want.UpdatedAt == nil) {
					t.Errorf("UserToResponse() UpdatedAt = %v, want %v", got.UpdatedAt, tt.want.UpdatedAt)
				}
			}
		})
	}
}

func TestUserResponseJSONStructure(t *testing.T) {
	now := time.Now()
	response := UserResponse{
		ID:        1,
		Email:     "test@example.com",
		Name:      "Test User",
		CreatedAt: &now,
		UpdatedAt: &now,
	}

	jsonBytes, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal UserResponse: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &result); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	// Check required fields exist
	requiredFields := []string{"id", "email", "name", "created_at", "updated_at"}
	for _, field := range requiredFields {
		if _, exists := result[field]; !exists {
			t.Errorf("UserResponse JSON missing field: %s", field)
		}
	}

	// Check password is NOT in response
	if _, exists := result["password"]; exists {
		t.Error("UserResponse JSON should not contain password field")
	}
}

func TestRegisterRequest(t *testing.T) {
	request := RegisterRequest{
		Email:    "user@example.com",
		Name:     "Test User",
		Password: "securepassword123",
	}

	jsonBytes, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("Failed to marshal RegisterRequest: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &result); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if result["email"] != request.Email {
		t.Errorf("RegisterRequest email = %v, want %v", result["email"], request.Email)
	}

	if result["name"] != request.Name {
		t.Errorf("RegisterRequest name = %v, want %v", result["name"], request.Name)
	}

	if result["password"] != request.Password {
		t.Errorf("RegisterRequest password = %v, want %v", result["password"], request.Password)
	}

	// Test unmarshaling
	var unmarshaled RegisterRequest
	if err := json.Unmarshal(jsonBytes, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal RegisterRequest: %v", err)
	}

	if unmarshaled.Email != request.Email {
		t.Errorf("Unmarshaled email = %v, want %v", unmarshaled.Email, request.Email)
	}
}

func TestLoginRequest(t *testing.T) {
	request := LoginRequest{
		Email:    "user@example.com",
		Password: "securepassword123",
	}

	jsonBytes, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("Failed to marshal LoginRequest: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &result); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if result["email"] != request.Email {
		t.Errorf("LoginRequest email = %v, want %v", result["email"], request.Email)
	}

	if result["password"] != request.Password {
		t.Errorf("LoginRequest password = %v, want %v", result["password"], request.Password)
	}

	// Test unmarshaling
	var unmarshaled LoginRequest
	if err := json.Unmarshal(jsonBytes, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal LoginRequest: %v", err)
	}

	if unmarshaled.Email != request.Email {
		t.Errorf("Unmarshaled email = %v, want %v", unmarshaled.Email, request.Email)
	}
}

func TestUpdateUserRequest(t *testing.T) {
	request := UpdateUserRequest{
		Email: "updated@example.com",
		Name:  "Updated Name",
	}

	jsonBytes, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("Failed to marshal UpdateUserRequest: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &result); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if result["email"] != request.Email {
		t.Errorf("UpdateUserRequest email = %v, want %v", result["email"], request.Email)
	}

	if result["name"] != request.Name {
		t.Errorf("UpdateUserRequest name = %v, want %v", result["name"], request.Name)
	}
}

func TestUpdatePasswordRequest(t *testing.T) {
	request := UpdatePasswordRequest{
		OldPassword: "oldpassword123",
		NewPassword: "newpassword456",
	}

	jsonBytes, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("Failed to marshal UpdatePasswordRequest: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &result); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if result["old_password"] != request.OldPassword {
		t.Errorf("UpdatePasswordRequest old_password = %v, want %v", result["old_password"], request.OldPassword)
	}

	if result["new_password"] != request.NewPassword {
		t.Errorf("UpdatePasswordRequest new_password = %v, want %v", result["new_password"], request.NewPassword)
	}
}

func TestRefreshRequest(t *testing.T) {
	request := RefreshRequest{
		RefreshToken: "abc123refresh",
	}

	jsonBytes, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("Failed to marshal RefreshRequest: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &result); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if result["refresh_token"] != request.RefreshToken {
		t.Errorf("RefreshRequest refresh_token = %v, want %v", result["refresh_token"], request.RefreshToken)
	}
}

func TestLogoutRequest(t *testing.T) {
	request := LogoutRequest{
		RefreshToken: "abc123logout",
	}

	jsonBytes, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("Failed to marshal LogoutRequest: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &result); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if result["refresh_token"] != request.RefreshToken {
		t.Errorf("LogoutRequest refresh_token = %v, want %v", result["refresh_token"], request.RefreshToken)
	}
}

func TestLoginResponse(t *testing.T) {
	response := LoginResponse{
		AccessToken:  "access_token_123",
		RefreshToken: "refresh_token_456",
		TokenType:    "Bearer",
	}

	jsonBytes, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal LoginResponse: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &result); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if result["access_token"] != response.AccessToken {
		t.Errorf("LoginResponse access_token = %v, want %v", result["access_token"], response.AccessToken)
	}

	if result["refresh_token"] != response.RefreshToken {
		t.Errorf("LoginResponse refresh_token = %v, want %v", result["refresh_token"], response.RefreshToken)
	}

	if result["token_type"] != response.TokenType {
		t.Errorf("LoginResponse token_type = %v, want %v", result["token_type"], response.TokenType)
	}
}
