package test

import (
	"testing"

	"github.com/H3199/doggodb/internal/data"
	"github.com/H3199/doggodb/internal/query"
)

func TestExecutorInsert(t *testing.T) {
	// Step 1: Create an in-memory storage instance.
	storage := data.NewInMemoryStorage()

	// Step 2: Create a new executor with the storage.
	executor := query.NewExecutor(*storage)

	// Step 3: Create a new table in the storage.
	tableName := "users"
	_, err := storage.CreateTable(tableName)
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	// Step 4: Define the INSERT statement (Values as []string).
	insertStmt := &query.InsertStatement{
		Table:   tableName,
		Columns: []string{"id", "name", "email"},
		Values:  []string{"1", "Alice", "alice@example.com"}, // String values
	}

	// Step 5: Execute the INSERT statement.
	_, err = executor.Execute(insertStmt)
	if err != nil {
		t.Fatalf("ExecuteInsert failed: %v", err)
	}

	// Step 6: Verify the row was inserted into the table.
	table, err := storage.GetTable(tableName)
	if err != nil {
		t.Fatalf("Failed to retrieve table: %v", err)
	}

	if len(table.Rows) != 1 {
		t.Fatalf("Expected 1 row, got %d", len(table.Rows))
	}

	row := table.Rows[0]

	// Verify the row's columns and values.
	expectedValues := map[string]interface{}{
		"id":    "1",                 // String value
		"name":  "Alice",             // String value
		"email": "alice@example.com", // String value
	}
	for column, expected := range expectedValues {
		actual, err := row.GetValue(column)
		if err != nil {
			t.Errorf("Column '%s' not found in row: %v", column, err)
			continue
		}
		if actual != expected {
			t.Errorf("Column '%s' mismatch: expected %v, got %v", column, expected, actual)
		}
	}
}
