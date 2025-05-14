package test_test

import (
	"testing"

	"github.com/H3199/doggodb/internal/data"
)

func TestInMemoryStorage(t *testing.T) {
	// Create an in-memory storage instance
	storage := data.NewInMemoryStorage()

	// Test CreateTable operation
	_, err := storage.CreateTable("users")
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	// Try creating a table with the same name again (should fail)
	_, err = storage.CreateTable("users")
	if err == nil {
		t.Fatalf("Expected error when creating table 'users' again, but got nil")
	}

	// Test Insert operation
	row1 := data.CreateRow(map[string]interface{}{"id": 1, "name": "Alice"}) // Using CreateRow with map
	err = storage.Insert("users", row1)                                      // Inserting row into the table
	if err != nil {
		t.Fatalf("Insert failed: %v", err)
	}

	// Insert another row
	row2 := data.CreateRow(map[string]interface{}{"id": 2, "name": "Bob"})
	err = storage.Insert("users", row2)
	if err != nil {
		t.Fatalf("Insert failed: %v", err)
	}

	// Test Query operation: Select rows where id == 1
	rows, err := storage.Query("users", func(r *data.Row) bool {
		id, _ := r.GetValue("id") // column 'id' now accessed by name
		return id == 1
	})
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}

	// Check if the correct row was found
	if len(rows) != 1 {
		t.Errorf("Expected 1 row, got %d", len(rows))
	}
	if name, _ := rows[0].GetValue("name"); name != "Alice" {
		t.Errorf("Expected name 'Alice', got '%v'", name)
	}

	// Test Update operation: Update name of user where id == 2
	err = storage.Update("users", map[string]interface{}{"name": "Charlie"}, func(r *data.Row) bool {
		id, _ := r.GetValue("id")
		return id == 2
	})
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	// Query again after the update
	rows, err = storage.Query("users", func(r *data.Row) bool {
		id, _ := r.GetValue("id")
		return id == 2
	})
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}

	// Check if the name was updated correctly
	if len(rows) != 1 {
		t.Errorf("Expected 1 row, got %d", len(rows))
	}
	if name, _ := rows[0].GetValue("name"); name != "Charlie" {
		t.Errorf("Expected name 'Charlie', got '%v'", name)
	}

	// Test Delete operation: Delete user where id == 1
	err = storage.Delete("users", 0) // delete the first row (id == 1)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Query again after deletion
	rows, err = storage.Query("users", func(r *data.Row) bool {
		id, _ := r.GetValue("id")
		return id == 1
	})
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}

	// Check if the row was deleted
	if len(rows) != 0 {
		t.Errorf("Expected 0 rows, got %d", len(rows))
	}
}
