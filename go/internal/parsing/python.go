package parsing

import (
	"regexp"
	"strings"
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
	pythonStringRegex     = regexp.MustCompile(`^(["'])(?:\\.|[^\\])*?\1`)
	pythonNumberRegex     = regexp.MustCompile(`^[0-9]+(\.[0-9]+)?([eE][+-]?[0-9]+)?`)
	pythonIdentifierRegex = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*`)
	pythonCommentRegex    = regexp.MustCompile(`^#.*`)
	pythonWhitespaceRegex = regexp.MustCompile(`^[ \t\r\n]+`)
)

// ParsePython parses Python code and returns a sequence of tokens
func ParsePython(code string) (TokenSequence, error) {
	var tokens TokenSequence

	// Trim any BOM or other markers
	code = strings.TrimSpace(code)

	for len(code) > 0 {
		// Try to match a string literal
		if match := pythonStringRegex.FindString(code); match != "" {
			tokens = append(tokens, Token{Type: TokenLiteral, Text: match})
			code = code[len(match):]
			continue
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

		// Try to match whitespace
		if match := pythonWhitespaceRegex.FindString(code); match != "" {
			tokens = append(tokens, Token{Type: TokenWhitespace, Text: match})
			code = code[len(match):]
			continue
		}

		// If none of the above matched, it's an "other" token (operator, punctuation, etc.)
		tokens = append(tokens, Token{Type: TokenOther, Text: string(code[0])})
		code = code[1:]
	}

	return tokens, nil
}
