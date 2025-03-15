package parsing

import (
	"regexp"
)

var (
	// JSON keywords
	jsonKeywords = map[string]bool{
		"true":  true,
		"false": true,
		"null":  true,
	}

	// Regular expressions for JSON tokens
	jsonNumberRegex     = regexp.MustCompile(`^-?(?:0|[1-9]\d*)(?:\.\d+)?(?:[eE][+-]?\d+)?`)
	jsonIdentifierRegex = regexp.MustCompile(`^[a-zA-Z_$][a-zA-Z0-9_$]*`)
	jsonWhitespaceRegex = regexp.MustCompile(`^[ \t\r\n]+`)
)

// JSONParser implements the Parser interface for JSON
type JSONParser struct{}

// Parse implements the Parser interface for JSON
func (p *JSONParser) Parse(code string) (TokenSequence, error) {
	return ParseJSON(code)
}

// isJSONStringStart checks if the code starts with a string delimiter
func isJSONStringStart(code string) bool {
	return len(code) > 0 && code[0] == '"'
}

// findJSONStringEnd finds the end of a string literal
func findJSONStringEnd(code string) int {
	if len(code) < 2 {
		return -1
	}

	for i := 1; i < len(code); i++ {
		if code[i] == '\\' && i+1 < len(code) {
			// Skip escaped character
			i++
			continue
		}
		if code[i] == '"' {
			return i + 1
		}
	}

	return -1
}

// ParseJSON parses JSON code and returns a sequence of tokens
func ParseJSON(code string) (TokenSequence, error) {
	var tokens TokenSequence

	for len(code) > 0 {
		// Try to match whitespace first to preserve formatting
		if match := jsonWhitespaceRegex.FindString(code); match != "" {
			tokens = append(tokens, Token{Type: TokenWhitespace, Text: match})
			code = code[len(match):]
			continue
		}

		// Try to match a string literal
		if isJSONStringStart(code) {
			end := findJSONStringEnd(code)
			if end > 0 {
				text := code[:end]
				tokenType := TokenLiteral

				// If this string is followed by a colon (ignoring whitespace),
				// it's an object key and should be treated as an identifier
				remaining := code[end:]
				for len(remaining) > 0 {
					if match := jsonWhitespaceRegex.FindString(remaining); match != "" {
						remaining = remaining[len(match):]
						continue
					}
					if len(remaining) > 0 && remaining[0] == ':' {
						tokenType = TokenIdentifier
					}
					break
				}

				tokens = append(tokens, Token{Type: tokenType, Text: text})
				code = code[end:]
				continue
			}
		}

		// Try to match a number
		if match := jsonNumberRegex.FindString(code); match != "" {
			tokens = append(tokens, Token{Type: TokenLiteral, Text: match})
			code = code[len(match):]
			continue
		}

		// Try to match keywords (true, false, null)
		if match := jsonIdentifierRegex.FindString(code); match != "" {
			if jsonKeywords[match] {
				tokens = append(tokens, Token{Type: TokenKeyword, Text: match})
			} else {
				tokens = append(tokens, Token{Type: TokenIdentifier, Text: match})
			}
			code = code[len(match):]
			continue
		}

		// If none of the above matched, it's an "other" token (punctuation, etc.)
		tokens = append(tokens, Token{Type: TokenOther, Text: string(code[0])})
		code = code[1:]
	}

	return tokens, nil
}
