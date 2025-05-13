package query

import (
	"errors"
)

//
// Token Definitions
//

type TokenType string

const (
	SELECT     TokenType = "SELECT"
	ASTERISK   TokenType = "ASTERISK"
	FROM       TokenType = "FROM"
	IDENTIFIER TokenType = "IDENTIFIER"
	EOF        TokenType = "EOF"
)

type Token struct {
	Type    TokenType
	Literal string
}

// Parse takes a slice of tokens and converts them into an AST.
func Parse(tokens []Token) (*SelectStatement, error) {
	// Validate the number of tokens.
	if len(tokens) < 4 {
		return nil, errors.New("invalid query: insufficient tokens")
	}

	// Check for the required SELECT * FROM <table> format.
	if tokens[0].Type != SELECT {
		return nil, errors.New("invalid query: expected SELECT at the beginning")
	}
	if tokens[1].Type != ASTERISK {
		return nil, errors.New("invalid query: expected * after SELECT")
	}
	if tokens[2].Type != FROM {
		return nil, errors.New("invalid query: expected FROM after *")
	}
	if tokens[3].Type != IDENTIFIER {
		return nil, errors.New("invalid query: expected table name after FROM")
	}

	// Construct the AST using the SelectStatement from AST.go.
	return &SelectStatement{
		Table: tokens[3].Literal,
	}, nil
}
