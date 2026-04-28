package models

import (
	"encoding/json"
	"testing"
	"time"

	"git.ramadhantriyant.id/ramadhantriyant/substrack/internal/database"
)

func TestCategoryRequest(t *testing.T) {
	request := CategoryRequest{
		Name:        "Entertainment",
		Description: "Movies, TV shows, and streaming services",
	}

	jsonBytes, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("Failed to marshal CategoryRequest: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &result); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if result["name"] != request.Name {
		t.Errorf("CategoryRequest name = %v, want %v", result["name"], request.Name)
	}

	if result["description"] != request.Description {
		t.Errorf("CategoryRequest description = %v, want %v", result["description"], request.Description)
	}

	// Test unmarshaling
	var unmarshaled CategoryRequest
	if err := json.Unmarshal(jsonBytes, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal CategoryRequest: %v", err)
	}

	if unmarshaled.Name != request.Name {
		t.Errorf("Unmarshaled name = %v, want %v", unmarshaled.Name, request.Name)
	}

	if unmarshaled.Description != request.Description {
		t.Errorf("Unmarshaled description = %v, want %v", unmarshaled.Description, request.Description)
	}
}

func TestCategoryRequestEmptyDescription(t *testing.T) {
	request := CategoryRequest{
		Name:        "Uncategorized",
		Description: "",
	}

	jsonBytes, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("Failed to marshal CategoryRequest with empty description: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &result); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if result["name"] != request.Name {
		t.Errorf("CategoryRequest name = %v, want %v", result["name"], request.Name)
	}

	if result["description"] != "" {
		t.Errorf("CategoryRequest description = %v, want empty string", result["description"])
	}
}

func TestCategoryList(t *testing.T) {
	now := time.Now()
	categories := []database.Category{
		{
			ID:          1,
			Name:        "Entertainment",
			Description: nil,
			CreatedAt:   &now,
			UpdatedAt:   &now,
		},
		{
			ID:          2,
			Name:        "Productivity",
			Description: nil,
			CreatedAt:   &now,
			UpdatedAt:   &now,
		},
	}

	list := CategoryList{
		Total:      len(categories),
		Categories: categories,
	}

	jsonBytes, err := json.Marshal(list)
	if err != nil {
		t.Fatalf("Failed to marshal CategoryList: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &result); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if result["total"] != float64(list.Total) {
		t.Errorf("CategoryList total = %v, want %v", result["total"], list.Total)
	}

	if cats, ok := result["categories"].([]interface{}); !ok {
		t.Error("CategoryList categories should be an array")
	} else if len(cats) != list.Total {
		t.Errorf("CategoryList categories length = %v, want %v", len(cats), list.Total)
	}
}

func TestCategoryListEmpty(t *testing.T) {
	list := CategoryList{
		Total:      0,
		Categories: []database.Category{},
	}

	jsonBytes, err := json.Marshal(list)
	if err != nil {
		t.Fatalf("Failed to marshal empty CategoryList: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &result); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if result["total"] != float64(0) {
		t.Errorf("Empty CategoryList total = %v, want 0", result["total"])
	}

	if cats, ok := result["categories"].([]interface{}); !ok {
		t.Error("Empty CategoryList categories should be an array")
	} else if len(cats) != 0 {
		t.Errorf("Empty CategoryList categories length = %v, want 0", len(cats))
	}
}

func TestCategoryListSingleCategory(t *testing.T) {
	now := time.Now()
	description := "Software subscriptions"
	categories := []database.Category{
		{
			ID:          1,
			Name:        "Software",
			Description: &description,
			CreatedAt:   &now,
			UpdatedAt:   &now,
		},
	}

	list := CategoryList{
		Total:      1,
		Categories: categories,
	}

	jsonBytes, err := json.Marshal(list)
	if err != nil {
		t.Fatalf("Failed to marshal CategoryList with single category: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &result); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if result["total"] != float64(1) {
		t.Errorf("CategoryList total = %v, want 1", result["total"])
	}

	if cats, ok := result["categories"].([]interface{}); !ok {
		t.Fatal("CategoryList categories should be an array")
	} else if len(cats) != 1 {
		t.Errorf("CategoryList categories length = %v, want 1", len(cats))
	} else {
		cat := cats[0].(map[string]interface{})
		if cat["name"] != "Software" {
			t.Errorf("Category name = %v, want Software", cat["name"])
		}
	}
}
