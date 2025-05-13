package query

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
	Table string // The name of the table being queried.
}

func (s *SelectStatement) statementNode() {}

// String returns a string representation of the SelectStatement.
func (s *SelectStatement) String() string {
	return "SELECT * FROM " + s.Table
}
