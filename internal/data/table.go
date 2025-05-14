package data

import (
	"errors"
	"sync"
)

// Table represents a table in the database, which contains rows.
type Table struct {
	Name  string
	Rows  []*Row
	mutex sync.Mutex
}

// NewTable creates a new empty table with the given name.
func NewTable(name string) *Table {
	return &Table{
		Name: name,
		Rows: []*Row{},
	}
}

// Insert adds a row to the table.
func (t *Table) Insert(row *Row) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.Rows = append(t.Rows, row)
}

// Delete removes a row by its index.
func (t *Table) Delete(index int) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if index < 0 || index >= len(t.Rows) {
		return errors.New("index out of bounds")
	}

	t.Rows = append(t.Rows[:index], t.Rows[index+1:]...)
	return nil
}

// Query retrieves rows that satisfy a condition function.
func (t *Table) Query(condition func(*Row) bool) []*Row {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	var result []*Row
	for _, row := range t.Rows {
		if condition(row) {
			result = append(result, row)
		}
	}
	return result
}

// Update modifies rows that satisfy the given condition and apply column assignments.
func (t *Table) Update(assignments map[string]interface{}, condition func(*Row) bool) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	for _, row := range t.Rows {
		if condition(row) {
			for column, value := range assignments {
				if err := row.SetValue(column, value); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
