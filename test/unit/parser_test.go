package query_test

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
		Table: "users",
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
	ast, err := query.Parse(tokens)
	if err != nil {
		t.Fatalf("Parser failed: %v", err)
	}
	if !reflect.DeepEqual(ast, expectedAST) {
		t.Errorf("Unexpected AST:\nExpected: %v\nGot: %v", expectedAST, ast)
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
