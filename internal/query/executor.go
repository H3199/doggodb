package query

import (
	"fmt"

	"github.com/H3199/doggodb/internal/data"
)

// Executor handles the execution of SQL queries.
type Executor struct {
	storage data.InMemoryStorage
}

// NewExecutor creates a new Executor with the provided storage.
func NewExecutor(storage data.InMemoryStorage) *Executor {
	return &Executor{storage: storage}
}

// Execute executes the given statement.
func (e *Executor) Execute(stmt Statement) (interface{}, error) {
	switch s := stmt.(type) {
	case *CreateTableStatement:
		return e.executeCreateTable(s)
	case *InsertStatement:
		return e.executeInsert(s)
	case *SelectStatement:
		return e.executeSelect(s)
	case *UpdateStatement:
		return e.executeUpdate(s)
	case *DeleteStatement:
		return e.executeDelete(s)
	default:
		return nil, fmt.Errorf("unsupported statement type")
	}
}

// executeInsert handles INSERT statements.
func (e *Executor) executeInsert(stmt *InsertStatement) (interface{}, error) {
	// Prepare values as a map from column name to value.
	values := make(map[string]interface{})
	for i, col := range stmt.Columns {
		values[col] = stmt.Values[i]
	}

	// Create a row from the values.
	row := data.CreateRow(values)

	// Insert the row into the storage.
	if err := e.storage.Insert(stmt.Table, row); err != nil {
		return nil, fmt.Errorf("failed to execute INSERT: %v", err)
	}

	// Return success with no result.
	return nil, nil
}

func (e *Executor) executeSelect(stmt *SelectStatement) (interface{}, error) {
	// Retrieve the table.
	table, err := e.storage.GetTable(stmt.Table)
	if err != nil {
		return nil, fmt.Errorf("failed to execute SELECT: %v", err)
	}

	// Parse the condition into a function.
	var conditionFunc func(*data.Row) bool
	if stmt.Conditions != "" {
		conditionFunc, err = parseCondition(stmt.Conditions)
		if err != nil {
			return nil, fmt.Errorf("invalid condition: %v", err)
		}
	}

	// Query the table with the condition.
	filteredRows := table.Query(func(row *data.Row) bool {
		if conditionFunc == nil {
			return true // No condition means include all rows.
		}
		return conditionFunc(row)
	})

	// If selecting all columns (SELECT *), return the filtered rows directly:
	if len(stmt.Columns) == 1 && stmt.Columns[0] == "*" {
		return filteredRows, nil
	}

	// Otherwise, create new rows with only selected columns:
	var result []*data.Row
	for _, row := range filteredRows {
		columns := make(map[string]interface{})
		for _, col := range stmt.Columns {
			if val, ok := row.Columns[col]; ok {
				columns[col] = val
			} else {
				columns[col] = nil
			}
		}
		result = append(result, data.CreateRow(columns))
	}

	return result, nil
}

func (e *Executor) executeUpdate(stmt *UpdateStatement) (interface{}, error) {
	// Retrieve the table
	table, err := e.storage.GetTable(stmt.Table)
	if err != nil {
		return nil, fmt.Errorf("failed to execute UPDATE: %v", err)
	}

	// Prevent empty WHERE clause from nuking the entire table.
	if stmt.Conditions == "" {
		return nil, fmt.Errorf("UPDATE query must have a WHERE clause to prevent unintentional updates")
	}

	// Parse the condition into a function
	var conditionFunc func(*data.Row) bool
	if stmt.Conditions != "" {
		conditionFunc, err = parseCondition(stmt.Conditions)
		if err != nil {
			return nil, fmt.Errorf("invalid condition: %v", err)
		}
	}

	// Track rows updated
	rowsUpdated := 0

	// Iterate over the rows and apply the update
	for _, row := range table.Rows {
		if conditionFunc == nil || conditionFunc(row) {
			for col, newValue := range stmt.Assignments {
				if err := row.SetValue(col, newValue); err != nil {
					return nil, fmt.Errorf("failed to set value for column '%s': %v", col, err)
				}
			}
			rowsUpdated++
		}
	}

	return rowsUpdated, nil // Return the count of rows updated
}

func (e *Executor) executeDelete(stmt *DeleteStatement) (interface{}, error) {
	// Retrieve the table.
	table, err := e.storage.GetTable(stmt.Table)
	if err != nil {
		return nil, fmt.Errorf("failed to execute DELETE: %v", err)
	}

	// Parse the condition into a function.
	var conditionFunc func(*data.Row) bool
	if stmt.Conditions != "" {
		conditionFunc, err = parseCondition(stmt.Conditions)
		if err != nil {
			return nil, fmt.Errorf("invalid condition: %v", err)
		}
	}

	// Find rows to delete using table.Query and the condition function.
	rowsToDelete := table.Query(func(row *data.Row) bool {
		if conditionFunc == nil {
			return true // No condition means delete all rows.
		}
		return conditionFunc(row)
	})

	// Delete the filtered rows from the table.
	rowsDeleted := 0
	for _, delRow := range rowsToDelete {
		err := table.DeleteRow(delRow)
		if err != nil {
			return nil, fmt.Errorf("failed to delete row: %v", err)
		}
		rowsDeleted++
	}

	return rowsDeleted, nil
}

func (e *Executor) executeCreateTable(stmt *CreateTableStatement) (interface{}, error) {
	_, err := e.storage.CreateTable(stmt.Table)
	if err != nil {
		return nil, fmt.Errorf("failed to create table: %v", err)
	}
	return nil, nil
}
