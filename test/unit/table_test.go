package data_test

import (
	"testing"

	"github.com/H3199/doggodb/internal/data"
)

// This is almost pure vibe code at the moment.

func TestTableCreation(t *testing.T) {
	tableName := "TestTable"
	table := data.NewTable(tableName)

	if table.Name != tableName {
		t.Errorf("expected table name %s, got %s", tableName, table.Name)
	}

	if len(table.Rows) != 0 {
		t.Errorf("expected no rows in the table, got %d", len(table.Rows))
	}
}

func TestInsert(t *testing.T) {
	table := data.NewTable("InsertTest")

	row1 := &data.Row{Values: []interface{}{"Alice", 30}}
	row2 := &data.Row{Values: []interface{}{"Bob", 25}}

	table.Insert(row1)
	table.Insert(row2)

	if len(table.Rows) != 2 {
		t.Errorf("expected 2 rows, got %d", len(table.Rows))
	}

	if table.Rows[0] != row1 {
		t.Errorf("expected first row to be %+v, got %+v", row1, table.Rows[0])
	}

	if table.Rows[1] != row2 {
		t.Errorf("expected second row to be %+v, got %+v", row2, table.Rows[1])
	}
}

func TestDelete(t *testing.T) {
	table := data.NewTable("DeleteTest")

	row1 := &data.Row{Values: []interface{}{"Alice", 30}}
	row2 := &data.Row{Values: []interface{}{"Bob", 25}}
	table.Insert(row1)
	table.Insert(row2)

	err := table.Delete(0)
	if err != nil {
		t.Errorf("unexpected error when deleting row: %v", err)
	}

	if len(table.Rows) != 1 {
		t.Errorf("expected 1 row after deletion, got %d", len(table.Rows))
	}

	if table.Rows[0] != row2 {
		t.Errorf("expected remaining row to be %+v, got %+v", row2, table.Rows[0])
	}

	err = table.Delete(5) // Invalid index
	if err == nil {
		t.Errorf("expected error for out-of-bounds index, got nil")
	}
}

func TestQuery(t *testing.T) {
	table := data.NewTable("QueryTest")

	row1 := &data.Row{Values: []interface{}{"Alice", 30}}
	row2 := &data.Row{Values: []interface{}{"Bob", 25}}
	row3 := &data.Row{Values: []interface{}{"Charlie", 35}}
	table.Insert(row1)
	table.Insert(row2)
	table.Insert(row3)

	result := table.Query(func(r *data.Row) bool {
		return r.Values[1] == 25 // Query for age == 25
	})

	if len(result) != 1 {
		t.Errorf("expected 1 row to match query, got %d", len(result))
	}

	if result[0] != row2 {
		t.Errorf("expected matching row to be %+v, got %+v", row2, result[0])
	}
}
