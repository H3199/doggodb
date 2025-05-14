package query

import (
	"errors"
	"strconv"
	"strings"
)

// Tokenize splits a query into tokens.
func Tokenize(query string) ([]Token, error) {
	var tokens []Token
	//valuesMode := false // Track whether we are inside the VALUES clause
	var current string // Accumulator for building tokens

	flushCurrent := func() {
		if current == "" {
			return
		}
		upperCurrent := strings.ToUpper(current)
		switch {
		case upperCurrent == "SELECT":
			tokens = append(tokens, Token{Type: SELECT, Literal: current})
		case upperCurrent == "INSERT":
			tokens = append(tokens, Token{Type: INSERT, Literal: current})
		case upperCurrent == "INTO":
			tokens = append(tokens, Token{Type: INTO, Literal: current})
		case upperCurrent == "VALUES":
			tokens = append(tokens, Token{Type: VALUES, Literal: current})
			//valuesMode = true // Enter VALUES mode
		case upperCurrent == "FROM":
			tokens = append(tokens, Token{Type: FROM, Literal: current})
		case upperCurrent == "UPDATE":
			tokens = append(tokens, Token{Type: UPDATE, Literal: current})
		case upperCurrent == "SET":
			tokens = append(tokens, Token{Type: SET, Literal: current})
		case upperCurrent == "WHERE":
			tokens = append(tokens, Token{Type: WHERE, Literal: current})
		case upperCurrent == "=":
			tokens = append(tokens, Token{Type: EQUALS, Literal: current})
		case upperCurrent == "*":
			tokens = append(tokens, Token{Type: ASTERISK, Literal: current})
		default:
			/*
				if valuesMode {
					// Handle literals within VALUES
					if _, err := strconv.Atoi(current); err == nil {
						tokens = append(tokens, Token{Type: NUMBER, Literal: current})
					} else if strings.HasPrefix(current, "'") && strings.HasSuffix(current, "'") {
						tokens = append(tokens, Token{Type: STRING, Literal: current})
					} else {
						tokens = append(tokens, Token{Type: IDENTIFIER, Literal: current})
					}
				}
			*/
			if _, err := strconv.Atoi(current); err == nil {
				tokens = append(tokens, Token{Type: NUMBER, Literal: current})
			} else if strings.HasPrefix(current, "'") && strings.HasSuffix(current, "'") {
				tokens = append(tokens, Token{Type: STRING, Literal: current})
			} else {
				// Default case for identifiers
				tokens = append(tokens, Token{Type: IDENTIFIER, Literal: current})
			}
		}
		current = ""
	}

	for _, char := range query {
		switch char {
		case ' ', '\t', '\n': // Handle whitespace as token separators
			flushCurrent()
		case '(':
			flushCurrent()
			tokens = append(tokens, Token{Type: LEFT_PAREN, Literal: string(char)})
		case ')':
			flushCurrent()
			tokens = append(tokens, Token{Type: RIGHT_PAREN, Literal: string(char)})
		case ',':
			flushCurrent()
			tokens = append(tokens, Token{Type: COMMA, Literal: string(char)})
		case '=':
			flushCurrent()
			tokens = append(tokens, Token{Type: EQUALS, Literal: string(char)})
		case '\'':
			// Handle quoted strings
			if strings.HasPrefix(current, "'") {
				current += string(char) // Closing quote
				flushCurrent()
			} else {
				flushCurrent()
				current = "'" // Start a new quoted string
			}
		default:
			current += string(char)
		}
	}

	flushCurrent() // Add any remaining token

	if len(tokens) == 0 {
		return nil, errors.New("query is empty or could not be tokenized")
	}

	return tokens, nil
}
