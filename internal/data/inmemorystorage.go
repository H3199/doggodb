package data

import (
	"fmt"
)

// InMemoryStorage implements the Storage interface for in-memory tables.
type InMemoryStorage struct {
	tables map[string]*Table
}

// NewInMemoryStorage creates a new instance of InMemoryStorage.
func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		tables: make(map[string]*Table),
	}
}

// CreateTable creates a new table with the specified name.
func (s *InMemoryStorage) CreateTable(tableName string) (*Table, error) {
	if _, exists := s.tables[tableName]; exists {
		return nil, fmt.Errorf("table %s already exists", tableName)
	}

	table := NewTable(tableName)
	s.tables[tableName] = table
	return table, nil
}

// GetTable retrieves a table by its name.
func (s *InMemoryStorage) GetTable(tableName string) (*Table, error) {
	table, exists := s.tables[tableName]
	if !exists {
		return nil, fmt.Errorf("table %s not found", tableName)
	}
	return table, nil
}

// Insert inserts a row into the specified table.
func (s *InMemoryStorage) Insert(tableName string, row *Row) error {
	table, err := s.GetTable(tableName)
	if err != nil {
		return err
	}
	table.Insert(row)
	return nil
}

// Query performs a SELECT query on the specified table and returns the result set.
func (s *InMemoryStorage) Query(tableName string, condition func(*Row) bool) ([]*Row, error) {
	table, err := s.GetTable(tableName)
	if err != nil {
		return nil, err
	}
	return table.Query(condition), nil
}

// Update updates rows in the specified table based on the given condition and assignments.
func (s *InMemoryStorage) Update(tableName string, assignments map[string]interface{}, condition func(*Row) bool) error {
	table, err := s.GetTable(tableName)
	if err != nil {
		return err
	}
	return table.Update(assignments, condition)
}

// Delete deletes a row from the specified table by its index.
func (s *InMemoryStorage) Delete(tableName string, index int) error {
	table, err := s.GetTable(tableName)
	if err != nil {
		return err
	}
	return table.Delete(index)
}
