package test_test

import (
	"testing"

	"github.com/H3199/doggodb/internal/data"
)

func TestRow(t *testing.T) {
	row := data.CreateRow(map[string]interface{}{"id": 1, "name": "Aliisa"})
	value, err := row.GetValue("name")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if value != "Aliisa" {
		t.Fatalf("expected 'Aliisa', got %v", value)
	}

	err = row.SetValue("name", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	value, _ = row.GetValue("name")
	if value != false {
		t.Fatalf("expected false, got %v", value)
	}
}
