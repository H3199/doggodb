package test_test

import (
	"testing"

	"fmt"

	"github.com/H3199/doggodb/internal/data"
	"github.com/H3199/doggodb/internal/query"
)

func TestExecutorInsert(t *testing.T) {
	// Step 1: Create an in-memory storage instance.
	storage := data.NewInMemoryStorage()

	// Step 2: Create a new executor with the storage.
	executor := query.NewExecutor(*storage)

	// Step 3: Create a new table in the storage.
	fmt.Print("Create a new table in the storage directly with storage layer.\n")
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

func TestExecutorSelect(t *testing.T) {
	// Step 1: Create an in-memory storage instance.
	storage := data.NewInMemoryStorage()

	// Step 2: Create a new executor with the storage.
	executor := query.NewExecutor(*storage)

	// Step 3: Create a new table in the storage.
	fmt.Println("Creating table...")
	tableName := "users"
	_, err := storage.CreateTable(tableName)
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}
	fmt.Println("Table created.")

	// Step 4: Insert some rows into the table.
	fmt.Println("Inserting rows...")
	insertStmt1 := &query.InsertStatement{
		Table:   tableName,
		Columns: []string{"id", "name", "email"},
		Values:  []string{"1", "Alice", "alice@example.com"},
	}
	_, err = executor.Execute(insertStmt1)
	if err != nil {
		t.Fatalf("Failed to insert row: %v", err)
	}

	insertStmt2 := &query.InsertStatement{
		Table:   tableName,
		Columns: []string{"id", "name", "email"},
		Values:  []string{"2", "Bob", "bob@example.com"},
	}
	_, err = executor.Execute(insertStmt2)
	if err != nil {
		t.Fatalf("Failed to insert row: %v", err)
	}

	fmt.Println("Rows inserted.")
	// Step 5: Define the SELECT statement.
	selectStmt := &query.SelectStatement{
		Table:      tableName,
		Columns:    []string{"id", "name", "email"}, // Selecting all columns
		Conditions: "id = 1",                        // Filter: WHERE id = 1
	}

	// Step 6: Execute the SELECT statement.
	fmt.Println("Executing SELECT statement...")
	result, err := executor.Execute(selectStmt)
	if err != nil {
		t.Fatalf("ExecuteSelect failed: %v", err)
	}

	// Step 7: Verify the result.
	fmt.Println("Verifying result...")
	rows, ok := result.([]*data.Row)
	//fmt.Print("Here are the rows:")
	//fmt.Print(rows)
	if !ok {
		t.Fatalf("Expected result to be []*data.Row, got %T", result)
	}

	if len(rows) != 1 {
		t.Fatalf("Expected 1 row, got %d", len(rows))
	}

	row := rows[0]
	expectedValues := map[string]interface{}{
		"id":    "1",
		"name":  "Alice",
		"email": "alice@example.com",
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

	// SELECT * test
	fmt.Println("SELECT * test")
	selectAllStmt := &query.SelectStatement{
		Table:      tableName,
		Columns:    []string{"*"},
		Conditions: "", // No WHERE clause
	}

	resultAll, err := executor.Execute(selectAllStmt)
	if err != nil {
		t.Fatalf("ExecuteSelect * failed: %v", err)
	}

	rowsAll, ok := resultAll.([]*data.Row)

	//	fmt.Println("Here are the rows:")
	//	fmt.Println(rowsAll)
	//	for i, r := range rowsAll {
	//		fmt.Printf("Row %d columns: %+v\n", i, r.Columns)
	//	}

	if !ok {
		t.Fatalf("Expected result to be []*data.Row, got %T", resultAll)
	}

	if len(rowsAll) != 2 {
		t.Fatalf("Expected 2 rows, got %d", len(rowsAll))
	}

	expectedRows := []map[string]interface{}{
		{"id": "1", "name": "Alice", "email": "alice@example.com"},
		{"id": "2", "name": "Bob", "email": "bob@example.com"},
	}

	for i, expected := range expectedRows {
		row := rowsAll[i]
		for column, expVal := range expected {
			actual, err := row.GetValue(column)
			if err != nil {
				t.Errorf("Row %d: column '%s' not found: %v", i, column, err)
				continue
			}
			if actual != expVal {
				t.Errorf("Row %d: column '%s' mismatch: expected %v, got %v", i, column, expVal, actual)
			}
		}
	}
}

func TestExecutorUpdate(t *testing.T) {
	fmt.Println("Step 1: Creating an in-memory storage instance...")
	storage := data.NewInMemoryStorage()

	fmt.Println("Step 2: Creating an executor...")
	executor := query.NewExecutor(*storage)

	fmt.Println("Step 3: Creating a new table in the storage...")
	tableName := "users"
	_, err := storage.CreateTable(tableName)
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}
	fmt.Println("Table created successfully.")

	fmt.Println("Step 4: Inserting rows into the table...")
	insertStmt1 := &query.InsertStatement{
		Table:   tableName,
		Columns: []string{"id", "name", "email"},
		Values:  []string{"1", "Alice", "alice@example.com"},
	}
	_, err = executor.Execute(insertStmt1)
	if err != nil {
		t.Fatalf("Failed to insert row: %v", err)
	}
	fmt.Println("Inserted first row.")

	insertStmt2 := &query.InsertStatement{
		Table:   tableName,
		Columns: []string{"id", "name", "email"},
		Values:  []string{"2", "Bob", "bob@example.com"},
	}
	_, err = executor.Execute(insertStmt2)
	if err != nil {
		t.Fatalf("Failed to insert row: %v", err)
	}
	fmt.Println("Inserted second row.")

	fmt.Println("Step 5: Defining the UPDATE statement...")
	updateStmt := &query.UpdateStatement{
		Table: tableName,
		Assignments: map[string]string{
			"name":  "Updated Alice",
			"email": "updated_alice@example.com",
		},
		Conditions: "id = 1",
	}
	fmt.Println("UPDATE statement defined:", updateStmt)

	fmt.Println("Step 6: Executing the UPDATE statement...")
	result, err := executor.Execute(updateStmt)
	if err != nil {
		t.Fatalf("ExecuteUpdate failed: %v", err)
	}

	rowsUpdated, ok := result.(int)
	if !ok {
		t.Fatalf("Expected result to be int, got %T", result)
	}

	fmt.Printf("Rows updated: %d\n", rowsUpdated)
	if rowsUpdated != 1 {
		t.Fatalf("Expected 1 row to be updated, got %d", rowsUpdated)
	}

	fmt.Println("Step 7: Verifying the updated row...")
	table, err := storage.GetTable(tableName)
	if err != nil {
		t.Fatalf("Failed to retrieve table: %v", err)
	}

	var updatedRow *data.Row
	for _, row := range table.Rows {
		id, err := row.GetValue("id")
		if err == nil && id == "1" {
			updatedRow = row
			break
		}
	}

	if updatedRow == nil {
		t.Fatalf("Updated row not found")
	}
	fmt.Println("Updated row found:", updatedRow)

	expectedValues := map[string]interface{}{
		"id":    "1",
		"name":  "Updated Alice",
		"email": "updated_alice@example.com",
	}

	fmt.Println("Verifying updated row values...")
	for column, expected := range expectedValues {
		actual, err := updatedRow.GetValue(column)
		if err != nil {
			t.Errorf("Column '%s' not found in row: %v", column, err)
			continue
		}
		if actual != expected {
			t.Errorf("Column '%s' mismatch: expected %v, got %v", column, expected, actual)
		} else {
			fmt.Printf("Column '%s' value verified: %v\n", column, actual)
		}
	}

	fmt.Println("Step 8: Verifying the other row remains unchanged...")
	var otherRow *data.Row
	for _, row := range table.Rows {
		id, err := row.GetValue("id")
		if err == nil && id == "2" {
			otherRow = row
			break
		}
	}

	if otherRow == nil {
		t.Fatalf("Other row not found")
	}
	fmt.Println("Other row found:", otherRow)

	expectedOtherValues := map[string]interface{}{
		"id":    "2",
		"name":  "Bob",
		"email": "bob@example.com",
	}

	fmt.Println("Verifying other row values...")
	for column, expected := range expectedOtherValues {
		actual, err := otherRow.GetValue(column)
		if err != nil {
			t.Errorf("Column '%s' not found in row: %v", column, err)
			continue
		}
		if actual != expected {
			t.Errorf("Column '%s' mismatch: expected %v, got %v", column, expected, actual)
		} else {
			fmt.Printf("Column '%s' value verified: %v\n", column, actual)
		}
	}

	fmt.Println("TestExecutorUpdate completed successfully.")
}

func TestExecutorDelete(t *testing.T) {
	fmt.Println("Starting TestExecutorDelete")

	// Step 1: Create in-memory storage and executor
	storage := data.NewInMemoryStorage()
	executor := query.NewExecutor(*storage)

	// Step 2: Create table
	tableName := "users"
	fmt.Println("Creating table:", tableName)
	_, err := storage.CreateTable(tableName)
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	// Step 3: Insert rows
	fmt.Println("Inserting rows")
	rowsToInsert := []map[string]string{
		{"id": "1", "name": "Alice", "email": "alice@example.com"},
		{"id": "2", "name": "Bob", "email": "bob@example.com"},
	}
	for _, rowData := range rowsToInsert {
		insertStmt := &query.InsertStatement{
			Table:   tableName,
			Columns: []string{"id", "name", "email"},
			Values:  []string{rowData["id"], rowData["name"], rowData["email"]},
		}
		_, err := executor.Execute(insertStmt)
		if err != nil {
			t.Fatalf("Failed to insert row: %v", err)
		}
	}
	fmt.Println("Rows inserted successfully")

	// Step 4: Define DELETE statement with condition to delete Bob (id=2)
	fmt.Println("Preparing DELETE statement: DELETE FROM users WHERE id = 2")
	deleteStmt := &query.DeleteStatement{
		Table:      tableName,
		Conditions: "id = 2",
	}

	// Step 5: Execute DELETE
	fmt.Println("Executing DELETE statement...")
	result, err := executor.Execute(deleteStmt)
	if err != nil {
		t.Fatalf("ExecuteDelete failed: %v", err)
	}
	fmt.Println("Delete result:", result)

	// Step 6: Verify that Bob's row is deleted, only Alice remains
	table, err := storage.GetTable(tableName)
	if err != nil {
		t.Fatalf("Failed to get table: %v", err)
	}

	fmt.Printf("Table now has %d rows\n", len(table.Rows))
	if len(table.Rows) != 1 {
		t.Fatalf("Expected 1 row after delete, got %d", len(table.Rows))
	}

	remainingRow := table.Rows[0]
	name, err := remainingRow.GetValue("name")
	if err != nil {
		t.Fatalf("Failed to get 'name' from remaining row: %v", err)
	}
	if name != "Alice" {
		t.Errorf("Expected remaining row to be 'Alice', got %v", name)
	}

	fmt.Println("TestExecutorDelete completed successfully")
}

func TestExecutorCreateTable(t *testing.T) {
	storage := data.NewInMemoryStorage()
	executor := query.NewExecutor(*storage)

	stmt := &query.CreateTableStatement{Table: "test_table"}

	_, err := executor.Execute(stmt)
	if err != nil {
		t.Fatalf("CreateTableStatement execution failed: %v", err)
	}

	_, err = storage.GetTable("test_table")
	if err != nil {
		t.Fatalf("Table was not created in storage: %v", err)
	}
}
