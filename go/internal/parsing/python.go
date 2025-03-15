package parsing

import (
	"regexp"
)

var (
	// Python keywords
	pythonKeywords = map[string]bool{
		"and":      true,
		"as":       true,
		"assert":   true,
		"async":    true,
		"await":    true,
		"break":    true,
		"class":    true,
		"continue": true,
		"def":      true,
		"del":      true,
		"elif":     true,
		"else":     true,
		"except":   true,
		"False":    true,
		"finally":  true,
		"for":      true,
		"from":     true,
		"global":   true,
		"if":       true,
		"import":   true,
		"in":       true,
		"is":       true,
		"lambda":   true,
		"None":     true,
		"nonlocal": true,
		"not":      true,
		"or":       true,
		"pass":     true,
		"raise":    true,
		"return":   true,
		"True":     true,
		"try":      true,
		"while":    true,
		"with":     true,
		"yield":    true,
	}

	// Regular expressions for Python tokens
	pythonNumberRegex     = regexp.MustCompile(`^[0-9]+(\.[0-9]+)?([eE][+-]?[0-9]+)?`)
	pythonIdentifierRegex = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*`)
	pythonCommentRegex    = regexp.MustCompile(`^#.*`)
	pythonWhitespaceRegex = regexp.MustCompile(`^[ \t\r\n]+`)
)

// isPythonStringStart checks if the code starts with a string delimiter
func isPythonStringStart(code string) bool {
	return len(code) > 0 && (code[0] == '"' || code[0] == '\'')
}

// findPythonStringEnd finds the end of a string literal
func findPythonStringEnd(code string) int {
	if len(code) < 2 {
		return -1
	}

	delimiter := code[0]
	for i := 1; i < len(code); i++ {
		if code[i] == '\\' && i+1 < len(code) {
			// Skip escaped character
			i++
			continue
		}
		if code[i] == delimiter {
			return i + 1
		}
	}

	return -1
}

// ParsePython parses Python code and returns a sequence of tokens
func ParsePython(code string) (TokenSequence, error) {
	var tokens TokenSequence

	// Process the code without trimming whitespace
	// This preserves indentation

	for len(code) > 0 {
		// Try to match whitespace first to preserve indentation
		if match := pythonWhitespaceRegex.FindString(code); match != "" {
			tokens = append(tokens, Token{Type: TokenWhitespace, Text: match})
			code = code[len(match):]
			continue
		}

		// Try to match a string literal
		if isPythonStringStart(code) {
			end := findPythonStringEnd(code)
			if end > 0 {
				tokens = append(tokens, Token{Type: TokenLiteral, Text: code[:end]})
				code = code[end:]
				continue
			}
		}

		// Try to match a comment
		if match := pythonCommentRegex.FindString(code); match != "" {
			tokens = append(tokens, Token{Type: TokenComment, Text: match})
			code = code[len(match):]
			continue
		}

		// Try to match a number
		if match := pythonNumberRegex.FindString(code); match != "" {
			tokens = append(tokens, Token{Type: TokenLiteral, Text: match})
			code = code[len(match):]
			continue
		}

		// Try to match an identifier or keyword
		if match := pythonIdentifierRegex.FindString(code); match != "" {
			if pythonKeywords[match] {
				tokens = append(tokens, Token{Type: TokenKeyword, Text: match})
			} else {
				tokens = append(tokens, Token{Type: TokenIdentifier, Text: match})
			}
			code = code[len(match):]
			continue
		}

		// If none of the above matched, it's an "other" token (operator, punctuation, etc.)
		tokens = append(tokens, Token{Type: TokenOther, Text: string(code[0])})
		code = code[1:]
	}

	return tokens, nil
}
