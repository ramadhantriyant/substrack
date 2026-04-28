package utils

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWriteJSON(t *testing.T) {
	tests := []struct {
		name        string
		status      int
		data        any
		wantStatus  int
		wantContent string
		wantErr     bool
	}{
		{
			name:        "simple object",
			status:      http.StatusOK,
			data:        map[string]string{"message": "success"},
			wantStatus:  http.StatusOK,
			wantContent: `{"message":"success"}`,
			wantErr:     false,
		},
		{
			name:        "array data",
			status:      http.StatusOK,
			data:        []string{"item1", "item2"},
			wantStatus:  http.StatusOK,
			wantContent: `["item1","item2"]`,
			wantErr:     false,
		},
		{
			name:        "nested object",
			status:      http.StatusCreated,
			data:        map[string]any{"user": map[string]string{"name": "John"}},
			wantStatus:  http.StatusCreated,
			wantContent: `{"user":{"name":"John"}}`,
			wantErr:     false,
		},
		{
			name:        "empty object",
			status:      http.StatusOK,
			data:        map[string]string{},
			wantStatus:  http.StatusOK,
			wantContent: `{}`,
			wantErr:     false,
		},
		{
			name:        "nil data",
			status:      http.StatusOK,
			data:        nil,
			wantStatus:  http.StatusOK,
			wantContent: `null`,
			wantErr:     false,
		},
		{
			name:        "error status",
			status:      http.StatusBadRequest,
			data:        map[string]string{"error": "bad request"},
			wantStatus:  http.StatusBadRequest,
			wantContent: `{"error":"bad request"}`,
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			err := WriteJSON(w, tt.status, tt.data)

			if (err != nil) != tt.wantErr {
				t.Errorf("WriteJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if w.Code != tt.wantStatus {
				t.Errorf("WriteJSON() status = %v, want %v", w.Code, tt.wantStatus)
			}

			contentType := w.Header().Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("WriteJSON() Content-Type = %v, want application/json", contentType)
			}

			body := w.Body.String()
			if body != tt.wantContent+"\n" {
				t.Errorf("WriteJSON() body = %v, want %v", body, tt.wantContent+"\n")
			}
		})
	}
}

func TestWriteJSONError(t *testing.T) {
	tests := []struct {
		name        string
		status      int
		errMessage  string
		wantStatus  int
		wantMessage string
	}{
		{
			name:        "bad request error",
			status:      http.StatusBadRequest,
			errMessage:  "invalid input",
			wantStatus:  http.StatusBadRequest,
			wantMessage: "invalid input",
		},
		{
			name:        "unauthorized error",
			status:      http.StatusUnauthorized,
			errMessage:  "invalid credentials",
			wantStatus:  http.StatusUnauthorized,
			wantMessage: "invalid credentials",
		},
		{
			name:        "not found error",
			status:      http.StatusNotFound,
			errMessage:  "resource not found",
			wantStatus:  http.StatusNotFound,
			wantMessage: "resource not found",
		},
		{
			name:        "internal server error",
			status:      http.StatusInternalServerError,
			errMessage:  "something went wrong",
			wantStatus:  http.StatusInternalServerError,
			wantMessage: "something went wrong",
		},
		{
			name:        "empty error message",
			status:      http.StatusBadRequest,
			errMessage:  "",
			wantStatus:  http.StatusBadRequest,
			wantMessage: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			WriteJSONError(w, tt.status, tt.errMessage)

			if w.Code != tt.wantStatus {
				t.Errorf("WriteJSONError() status = %v, want %v", w.Code, tt.wantStatus)
			}

			contentType := w.Header().Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("WriteJSONError() Content-Type = %v, want application/json", contentType)
			}

			var response ErrorResponse
			if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
				t.Fatalf("Failed to unmarshal response: %v", err)
			}

			if response.Status != tt.wantStatus {
				t.Errorf("WriteJSONError() response.Status = %v, want %v", response.Status, tt.wantStatus)
			}

			if response.Message != tt.wantMessage {
				t.Errorf("WriteJSONError() response.Message = %v, want %v", response.Message, tt.wantMessage)
			}
		})
	}
}

func TestErrorResponseStructure(t *testing.T) {
	response := ErrorResponse{
		Status:  400,
		Message: "test error",
	}

	jsonBytes, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal ErrorResponse: %v", err)
	}

	expected := `{"status":400,"message":"test error"}`
	if string(jsonBytes) != expected {
		t.Errorf("ErrorResponse JSON = %v, want %v", string(jsonBytes), expected)
	}

	// Test unmarshaling
	var unmarshaled ErrorResponse
	if err := json.Unmarshal(jsonBytes, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal ErrorResponse: %v", err)
	}

	if unmarshaled.Status != response.Status {
		t.Errorf("Unmarshaled Status = %v, want %v", unmarshaled.Status, response.Status)
	}

	if unmarshaled.Message != response.Message {
		t.Errorf("Unmarshaled Message = %v, want %v", unmarshaled.Message, response.Message)
	}
}
