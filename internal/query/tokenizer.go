package query

import (
	"errors"
	"strings"
)

// Tokenize splits a query into tokens.
func Tokenize(query string) ([]Token, error) {
	parts := strings.Fields(query)
	if len(parts) != 4 {
		return nil, errors.New("invalid query format")
	}

	// Simplistic tokenizer for SELECT * FROM <table>
	return []Token{
		{Type: SELECT, Literal: parts[0]},
		{Type: ASTERISK, Literal: parts[1]},
		{Type: FROM, Literal: parts[2]},
		{Type: IDENTIFIER, Literal: parts[3]},
	}, nil
}
