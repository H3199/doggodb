package data

import (
	"fmt"
)

// Row represents a row in a table, holding values for each column.
type Row struct {
	Columns map[string]interface{} // Mapping of column names to values
}

// CreateRow creates a new Row with the specified column values.
func CreateRow(columns map[string]interface{}) *Row {
	return &Row{
		Columns: columns,
	}
}

// GetValue retrieves a column value by its name.
func (r *Row) GetValue(columnName string) (interface{}, error) {
	value, exists := r.Columns[columnName]
	if !exists {
		return nil, fmt.Errorf("column '%s' not found", columnName)
	}
	return value, nil
}

// SetValue sets a column value by its name.
func (r *Row) SetValue(columnName string, value interface{}) error {
	if _, exists := r.Columns[columnName]; !exists {
		return fmt.Errorf("column '%s' not found", columnName)
	}
	r.Columns[columnName] = value
	return nil
}
