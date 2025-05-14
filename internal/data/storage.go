package data

// Storage represents the interface for interacting with the database.
type Storage interface {
	// Insert adds a new row to the table.
	Insert(tableName string, columns map[string]interface{}) error

	// Select retrieves rows from the table based on columns and conditions.
	Select(tableName string, columns []string, condition string) ([]*Row, error)

	// Update modifies rows in the table based on conditions and column assignments.
	Update(tableName string, assignments map[string]interface{}, condition string) error

	// CreateTable creates a new table in the storage.
	CreateTable(name string) error

	// GetTable retrieves a table by its name.
	GetTable(name string) (*Table, error)
}
