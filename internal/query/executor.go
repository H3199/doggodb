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
	case *InsertStatement:
		return e.executeInsert(s)
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
