package query

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/H3199/doggodb/internal/data"
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
	STRING      TokenType = "STRING"
	IDENTIFIER  TokenType = "IDENTIFIER"
	NUMBER      TokenType = "NUMBER"
	UPDATE      TokenType = "UPDATE"
	EQUALS      TokenType = "EQUALS"
	WHERE       TokenType = "WHERE"
	SET         TokenType = "SET"
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
	case UPDATE:
		return parseUpdate(tokens)
	default:
		return nil, errors.New("unsupported query type")
	}
}

func parseSelect(tokens []Token) (*SelectStatement, error) {
	if len(tokens) < 4 {
		return nil, errors.New("invalid query: insufficient tokens")
	}

	// Ensure the query starts with SELECT
	if tokens[0].Type != SELECT {
		return nil, errors.New("invalid SELECT query format: missing SELECT")
	}

	var columns []string
	i := 1

	// Parse columns
	if tokens[i].Type == ASTERISK {
		columns = append(columns, "*")
		i++
	} else {
		for i < len(tokens) && tokens[i].Type != FROM {
			if tokens[i].Type == IDENTIFIER {
				columns = append(columns, tokens[i].Literal)
			} else if tokens[i].Type != COMMA {
				return nil, errors.New("unexpected token in column list")
			}
			i++
		}
		if len(columns) == 0 {
			return nil, errors.New("no columns specified in SELECT")
		}
	}

	// Ensure FROM keyword
	if i >= len(tokens) || tokens[i].Type != FROM {
		return nil, errors.New("expected FROM after column list")
	}
	i++ // Skip 'FROM'

	// Parse table name
	if i >= len(tokens) || tokens[i].Type != IDENTIFIER {
		return nil, errors.New("expected table name after FROM")
	}
	table := tokens[i].Literal
	i++ // Skip table name

	// Parse optional WHERE clause
	var conditions string
	if i < len(tokens) && tokens[i].Type == WHERE {
		i++ // Skip 'WHERE'
		var whereParts []string
		for i < len(tokens) {
			whereParts = append(whereParts, tokens[i].Literal)
			i++
		}
		conditions = strings.Join(whereParts, " ")
	}

	return &SelectStatement{
		Table:      table,
		Columns:    columns,
		Conditions: conditions,
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

func parseUpdate(tokens []Token) (*UpdateStatement, error) {
	if len(tokens) < 4 {
		return nil, errors.New("invalid query: insufficient tokens for UPDATE")
	}

	if tokens[0].Type != UPDATE {
		//fmt.Println("DEBUG: parseUpdate: tokens[0] is not UPDATE")
		return nil, errors.New("invalid UPDATE query format")
	}

	table := tokens[1].Literal
	// fmt.Println("DEBUG: parseUpdate: Found table:", table)

	if tokens[2].Type != SET {
		//	fmt.Print("DEBUG: expected SET after table name in UPDATE")
		return nil, errors.New("expected SET after table name in UPDATE")
	}

	assignments := make(map[string]string)
	i := 3

	// Parse assignments (SET clause)
	for i < len(tokens) && tokens[i].Type != WHERE {
		if tokens[i].Type == IDENTIFIER {
			column := tokens[i].Literal
			i++
			if i >= len(tokens) || tokens[i].Type != EQUALS {
				//	fmt.Println("DEBUG: parseUpdate: expected '=' after column name in SET clause")
				return nil, errors.New("expected '=' after column name in SET clause")
			}
			i++
			if i >= len(tokens) || (tokens[i].Type != STRING && tokens[i].Type != NUMBER) {
				//	fmt.Printf("DEBUG: tokens[%d]: %+v\n", i, tokens[i])
				//	fmt.Println("DEBUG:expected '=' after column name in SET clause II")
				return nil, errors.New("expected value after '=' in SET clause")
			}
			value := tokens[i].Literal
			assignments[column] = value
			i++

			if i < len(tokens) && tokens[i].Type == COMMA {
				i++ // Skip comma
			}
		} else {
			//	fmt.Println("DEBUG: parseUpdate: invalid token in SET clause")
			return nil, errors.New("invalid token in SET clause")
		}
	}

	// Parse WHERE clause (optional)
	var conditions string
	if i < len(tokens) && tokens[i].Type == WHERE {
		i++
		var whereParts []string
		for i < len(tokens) {
			whereParts = append(whereParts, tokens[i].Literal)
			i++
		}
		conditions = strings.Join(whereParts, " ")
	}

	// Debug output to check the flow
	// fmt.Printf("Parsed UPDATE: table=%s, assignments=%v, conditions=%s\n", table, assignments, conditions)

	return &UpdateStatement{
		Table:       table,
		Assignments: assignments,
		Conditions:  conditions,
	}, nil
}

func parseCondition(condition string) (func(*data.Row) bool, error) {
	// Example: "age > 30"
	// Note: This is a basic implementation. For complex conditions, a full parser is needed.
	parts := strings.Split(condition, " ")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid condition format")
	}

	column, operator, value := parts[0], parts[1], parts[2]

	// Convert the value from the condition to a float if needed
	valNum, _ := toFloat(value)

	return func(row *data.Row) bool {
		colValue, exists := row.Columns[column]
		if !exists {
			return false
		}

		switch operator {
		case "=":
			// Handle equality
			return fmt.Sprintf("%v", colValue) == value

		case ">", "<":
			// Ensure colValue is numeric
			colNum, err := toFloat(fmt.Sprintf("%v", colValue))
			if err != nil {
				return false
			}

			// Perform the comparison
			if operator == ">" {
				return colNum > valNum
			} else if operator == "<" {
				return colNum < valNum
			}

		default:
			// Unsupported operator
			return false
		}

		return false
	}, nil
}

func toFloat(value string) (float64, error) {
	return strconv.ParseFloat(value, 64)
}
