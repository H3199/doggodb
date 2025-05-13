package data

import (
	"errors"
	"sync"
)

// A Table is a slice of rows
type Table struct {
	Name  string
	Rows  []*Row
	mutex sync.Mutex
}

func NewTable(name string) *Table {
	return &Table{
		Name: name,
		Rows: []*Row{},
	}
}

func (t *Table) Insert(row *Row) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.Rows = append(t.Rows, row)
}

func (t *Table) Delete(index int) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if index < 0 || index >= len(t.Rows) {
		return errors.New("index out of bounds")
	}

	t.Rows = append(t.Rows[:index], t.Rows[index+1:]...)
	return nil
}

func (t *Table) Query(condition func(*Row) bool) []*Row {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	result := []*Row{}
	for _, row := range t.Rows {
		if condition(row) {
			result = append(result, row)
		}
	}
	return result
}
