package data_test

import (
    "github.com/H3199/doggodb/internal/data"
    "testing"
)

func TestRow(t *testing.T) {
    row := data.CreateRow([]interface{}{1, "Aliisa", true})

    value, err := row.GetValue(1)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    if value != "Aliisa" {
        t.Fatalf("expected 'Aliisa', got %v", value)
    }

    err = row.SetValue(2, false)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    value, _ = row.GetValue(2)
    if value != false {
        t.Fatalf("expected false, got %v", value)
    }
}
