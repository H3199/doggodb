package query

import (
	"errors"
)

//
// Token Definitions
//

type TokenType string

const (
	INSERT      TokenType = "INSERT"
	SELECT      TokenType = "SELECT"
	ASTERISK    TokenType = "ASTERISK"
	FROM        TokenType = "FROM"
	INTO        TokenType = "INTO"
	VALUES      TokenType = "VALUES"
	COMMA       TokenType = "COMMA"
	LEFT_PAREN  TokenType = "LEFT_PAREN"
	RIGHT_PAREN TokenType = "RIGHT_PAREN"
	STRING      TokenType = "STRING" // For literal strings
	IDENTIFIER  TokenType = "IDENTIFIER"
	NUMBER      TokenType = "NUMBER"
)

type Token struct {
	Type    TokenType
	Literal string
}

// Parse takes a slice of tokens and converts them into an AST.
func Parse(tokens []Token) (Statement, error) {
	if len(tokens) < 1 {
		return nil, errors.New("invalid query: no tokens provided")
	}

	switch tokens[0].Type {
	case SELECT:
		return parseSelect(tokens)
	case INSERT:
		return parseInsert(tokens)
	default:
		return nil, errors.New("unsupported query type")
	}
}

func parseSelect(tokens []Token) (*SelectStatement, error) {
	if len(tokens) < 4 {
		return nil, errors.New("invalid query: insufficient tokens")
	}

	if tokens[0].Type != SELECT || tokens[1].Type != ASTERISK || tokens[2].Type != FROM || tokens[3].Type != IDENTIFIER {
		return nil, errors.New("invalid SELECT query format")
	}

	return &SelectStatement{
		Table: tokens[3].Literal,
	}, nil
}

func parseInsert(tokens []Token) (*InsertStatement, error) {
	if len(tokens) < 6 {
		return nil, errors.New("invalid query: insufficient tokens for INSERT")
	}

	if tokens[0].Type != INSERT || tokens[1].Type != INTO || tokens[2].Type != IDENTIFIER {
		return nil, errors.New("invalid INSERT query format")
	}

	table := tokens[2].Literal
	var columns []string
	var values []string

	// Extract columns
	i := 3
	if tokens[i].Type == LEFT_PAREN {
		i++
		for i < len(tokens) && tokens[i].Type != RIGHT_PAREN {
			if tokens[i].Type == IDENTIFIER {
				columns = append(columns, tokens[i].Literal)
			} else if tokens[i].Type != COMMA {
				return nil, errors.New("unexpected token in column list")
			}
			i++
		}
		if i >= len(tokens) || tokens[i].Type != RIGHT_PAREN {
			return nil, errors.New("expected ')' after column names")
		}
		i++ // Move past ')'
	}

	// Now expect VALUES
	if i >= len(tokens) || tokens[i].Type != VALUES {
		return nil, errors.New("expected VALUES after columns")
	}
	i++ // Move past 'VALUES'

	// Extract values
	if i >= len(tokens) || tokens[i].Type != LEFT_PAREN {
		return nil, errors.New("expected '(' after VALUES")
	}
	i++ // Skip '(' token

	for i < len(tokens) && tokens[i].Type != RIGHT_PAREN {
		if tokens[i].Type == STRING || tokens[i].Type == NUMBER {
			values = append(values, tokens[i].Literal)
		} else if tokens[i].Type != COMMA {
			return nil, errors.New("unexpected token in values list")
		}
		i++
	}
	if i >= len(tokens) || tokens[i].Type != RIGHT_PAREN {
		return nil, errors.New("expected ')' after values")
	}
	i++ // Move past ')'

	if len(columns) == 0 || len(values) == 0 {
		return nil, errors.New("no columns or values found")
	}

	return &InsertStatement{
		Table:   table,
		Columns: columns,
		Values:  values,
	}, nil
}
