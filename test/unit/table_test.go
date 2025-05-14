package test_test

import (
	"testing"

	"github.com/H3199/doggodb/internal/data"
)

// This is pure vibe code.

func TestTableOperations(t *testing.T) {
	// Create a new table
	table := data.NewTable("users")

	// Insert a row into the table
	row := data.CreateRow(map[string]interface{}{"id": 1, "name": "Alice"})
	table.Insert(row)

	// Query the table for rows where id = 1
	rows := table.Query(func(r *data.Row) bool {
		id, err := r.GetValue("id")
		if err != nil {
			t.Errorf("Error retrieving column 'id': %v", err)
		}
		return id == 1
	})

	if len(rows) != 1 {
		t.Errorf("Expected 1 row, got %d", len(rows))
	}

	// Verify that the row's name is "Alice"
	name, err := rows[0].GetValue("name")
	if err != nil {
		t.Errorf("Error retrieving column 'name': %v", err)
	}
	if name != "Alice" {
		t.Errorf("Expected name 'Alice', got %s", name)
	}

	// Update the row's name
	err = table.Update(map[string]interface{}{"name": "Bob"}, func(r *data.Row) bool {
		id, err := r.GetValue("id")
		if err != nil {
			t.Errorf("Error retrieving column 'id': %v", err)
		}
		return id == 1
	})

	if err != nil {
		t.Fatalf("Error updating row: %v", err)
	}

	// Verify that the name was updated
	rows = table.Query(func(r *data.Row) bool {
		id, err := r.GetValue("id")
		if err != nil {
			t.Errorf("Error retrieving column 'id': %v", err)
		}
		return id == 1
	})

	if len(rows) != 1 {
		t.Errorf("Expected 1 row, got %d", len(rows))
	}

	name, err = rows[0].GetValue("name")
	if err != nil {
		t.Errorf("Error retrieving column 'name': %v", err)
	}
	if name != "Bob" {
		t.Errorf("Expected name 'Bob', got %s", name)
	}
}
