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
