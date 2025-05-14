package test_test

import (
	"reflect"
	"testing"

	"github.com/H3199/doggodb/internal/query"
)

func TestSimpleSelect(t *testing.T) {
	queryString := "SELECT * FROM users"

	// Expected tokens
	expectedTokens := []query.Token{
		{Type: query.SELECT, Literal: "SELECT"},
		{Type: query.ASTERISK, Literal: "*"},
		{Type: query.FROM, Literal: "FROM"},
		{Type: query.IDENTIFIER, Literal: "users"},
	}

	// Expected AST
	expectedAST := &query.SelectStatement{
		Table:   "users",
		Columns: []string{"*"}, // Include the columns for completeness
	}

	// Test the tokenizer
	tokens, err := query.Tokenize(queryString)
	if err != nil {
		t.Fatalf("Tokenizer failed: %v", err)
	}
	if !reflect.DeepEqual(tokens, expectedTokens) {
		t.Errorf("Unexpected tokens:\nExpected: %v\nGot: %v", expectedTokens, tokens)
	}

	// Test the parser
	stmt, err := query.Parse(tokens)
	if err != nil {
		t.Fatalf("Parser failed: %v", err)
	}

	// Ensure stmt is of type *query.SelectStatement
	ast, ok := stmt.(*query.SelectStatement)
	if !ok {
		t.Fatalf("Expected *query.SelectStatement, got %T", stmt)
	}

	// Compare the AST fields
	if expectedAST.Table != ast.Table {
		t.Errorf("Table mismatch: expected %s, got %s", expectedAST.Table, ast.Table)
	}
	if !reflect.DeepEqual(expectedAST.Columns, ast.Columns) {
		t.Errorf("Columns mismatch: expected %v, got %v", expectedAST.Columns, ast.Columns)
	}
}

func TestAdvancedSelect(t *testing.T) {
	queryString := "SELECT name, age FROM users WHERE id = 1"
	tokens, err := query.Tokenize(queryString)
	if err != nil {
		t.Fatalf("Tokenization failed: %v", err)
	}

	stmt, err := query.Parse(tokens)
	if err != nil {
		t.Fatalf("Parsing failed: %v", err)
	}

	selectStmt, ok := stmt.(*query.SelectStatement)
	if !ok {
		t.Fatalf("Expected SelectStatement, got %T", stmt)
	}

	expectedColumns := []string{"name", "age"}
	if !reflect.DeepEqual(selectStmt.Columns, expectedColumns) {
		t.Errorf("Expected columns %v, got %v", expectedColumns, selectStmt.Columns)
	}

	if selectStmt.Table != "users" {
		t.Errorf("Expected table 'users', got %s", selectStmt.Table)
	}

	if selectStmt.Conditions != "id = 1" {
		t.Errorf("Expected condition 'id = 1', got %s", selectStmt.Conditions)
	}
}

func TestInsert(t *testing.T) {
	queryString := "INSERT INTO users (name, age) VALUES ('Alice', 30)"

	expectedTokens := []query.Token{
		{Type: query.INSERT, Literal: "INSERT"},
		{Type: query.INTO, Literal: "INTO"},
		{Type: query.IDENTIFIER, Literal: "users"},
		{Type: query.LEFT_PAREN, Literal: "("},
		{Type: query.IDENTIFIER, Literal: "name"},
		{Type: query.COMMA, Literal: ","},
		{Type: query.IDENTIFIER, Literal: "age"},
		{Type: query.RIGHT_PAREN, Literal: ")"},
		{Type: query.VALUES, Literal: "VALUES"},
		{Type: query.LEFT_PAREN, Literal: "("},
		{Type: query.STRING, Literal: "'Alice'"},
		{Type: query.COMMA, Literal: ","},
		{Type: query.NUMBER, Literal: "30"},
		{Type: query.RIGHT_PAREN, Literal: ")"},
	}

	expectedAST := &query.InsertStatement{
		Table:   "users",
		Columns: []string{"name", "age"},
		Values:  []string{"'Alice'", "30"},
	}

	tokens, err := query.Tokenize(queryString)
	if err != nil {
		t.Fatalf("Tokenize failed: %v", err)
	}
	if !reflect.DeepEqual(tokens, expectedTokens) {
		t.Errorf("Tokens do not match. Expected %v, got %v", expectedTokens, tokens)
	}

	ast, err := query.Parse(tokens)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	if !reflect.DeepEqual(ast, expectedAST) {
		t.Errorf("AST does not match. Expected %v, got %v", expectedAST, ast)
	}
}

func TestUpdateParsing(t *testing.T) {
	queryString := "UPDATE users SET name = 'Alice', age = 25 WHERE id = 1"
	tokens, err := query.Tokenize(queryString)
	if err != nil {
		t.Fatalf("Tokenization failed: %v", err)
	}

	stmt, err := query.Parse(tokens)
	updateStmt, ok := stmt.(*query.UpdateStatement)
	if !ok {
		t.Fatalf("Expected UpdateStatement, got %T", stmt)
	}

	if updateStmt.Table != "users" {
		t.Errorf("Expected table 'users', got %s", updateStmt.Table)
	}
	if err != nil {
		t.Fatalf("Parsing failed: %v", err)
	}

	expectedAssignments := map[string]string{"name": "'Alice'", "age": "25"}
	for col, val := range expectedAssignments {
		if updateStmt.Assignments[col] != val {
			t.Errorf("Expected %s = %s, got %s", col, val, updateStmt.Assignments[col])
		}
	}

	if updateStmt.Conditions != "id = 1" {
		t.Errorf("Expected condition 'id = 1', got %s", updateStmt.Conditions)
	}
}
