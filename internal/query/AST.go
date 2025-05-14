package query

import "strings"

// Node is the interface that all AST nodes implement.
type Node interface {
	String() string
}

// Statement is a high-level node that represents a complete query.
type Statement interface {
	Node
	statementNode()
}

// SelectStatement represents a SELECT query in the AST.
type SelectStatement struct {
	Table      string
	Columns    []string
	Conditions string // Optional WHERE clause
}

func (s *SelectStatement) statementNode() {}

// String returns a string representation of the SelectStatement.
func (s *SelectStatement) String() string {
	return "SELECT * FROM " + s.Table
}

// InsertStatement represents an INSERT query in the AST.
type InsertStatement struct {
	Table   string   // The name of the table being inserted into.
	Columns []string // The list of column names.
	Values  []string // The corresponding list of values.
}

func (i *InsertStatement) statementNode() {}

// String returns a string representation of the InsertStatement.
func (i *InsertStatement) String() string {
	columns := strings.Join(i.Columns, ", ")
	values := strings.Join(i.Values, ", ")
	return "INSERT INTO " + i.Table + " (" + columns + ") VALUES (" + values + ")"
}

type UpdateStatement struct {
	Table       string            // The table to update
	Assignments map[string]string // Column-value pairs to update
	Conditions  string            // WHERE clause (string for now, could be more structured later)
}

func (i *UpdateStatement) statementNode() {} // I have no idea why this is needed.

// Implement the `string` method for UpdateStatement
func (u *UpdateStatement) String() string {
	assignments := []string{}
	for col, val := range u.Assignments {
		assignments = append(assignments, col+"="+val)
	}
	assignmentStr := strings.Join(assignments, ", ")

	whereClause := ""
	if u.Conditions != "" {
		whereClause = " WHERE " + u.Conditions
	}

	return "UPDATE " + u.Table + " SET " + assignmentStr + whereClause
}
