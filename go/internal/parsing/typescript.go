package parsing

import (
	"regexp"
	"strings"
)

var (
	// TypeScript/JavaScript keywords
	typescriptKeywords = map[string]bool{
		"abstract":    true,
		"any":         true,
		"as":          true,
		"async":       true,
		"await":       true,
		"boolean":     true,
		"break":       true,
		"case":        true,
		"catch":       true,
		"class":       true,
		"const":       true,
		"constructor": true,
		"continue":    true,
		"debugger":    true,
		"declare":     true,
		"default":     true,
		"delete":      true,
		"do":          true,
		"else":        true,
		"enum":        true,
		"export":      true,
		"extends":     true,
		"false":       true,
		"finally":     true,
		"for":         true,
		"from":        true,
		"function":    true,
		"get":         true,
		"if":          true,
		"implements":  true,
		"import":      true,
		"in":          true,
		"instanceof":  true,
		"interface":   true,
		"is":          true,
		"keyof":       true,
		"let":         true,
		"module":      true,
		"namespace":   true,
		"new":         true,
		"null":        true,
		"number":      true,
		"object":      true,
		"of":          true,
		"package":     true,
		"private":     true,
		"protected":   true,
		"public":      true,
		"readonly":    true,
		"require":     true,
		"return":      true,
		"set":         true,
		"static":      true,
		"string":      true,
		"super":       true,
		"switch":      true,
		"symbol":      true,
		"this":        true,
		"throw":       true,
		"true":        true,
		"try":         true,
		"type":        true,
		"typeof":      true,
		"undefined":   true,
		"unique":      true,
		"unknown":     true,
		"var":         true,
		"void":        true,
		"while":       true,
		"with":        true,
		"yield":       true,
	}

	// Regular expressions for TypeScript/JavaScript tokens
	typescriptNumberRegex     = regexp.MustCompile(`^(0[xX][0-9a-fA-F]+|0[oO][0-7]+|0[bB][01]+|[0-9]+(\.[0-9]+)?([eE][+-]?[0-9]+)?)`)
	typescriptIdentifierRegex = regexp.MustCompile(`^[a-zA-Z_$][a-zA-Z0-9_$]*`)
	typescriptCommentRegex    = regexp.MustCompile(`^(//.*|/\*[\s\S]*?\*/)`)
	typescriptWhitespaceRegex = regexp.MustCompile(`^[ \t\r\n]+`)
	typescriptTemplateRegex   = regexp.MustCompile("^`(?:\\\\.|[^`\\\\])*?`")
)

// isStringStart checks if the code starts with a string delimiter
func isStringStart(code string) bool {
	return len(code) > 0 && (code[0] == '"' || code[0] == '\'')
}

// findStringEnd finds the end of a string literal
func findStringEnd(code string) int {
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

// ParseTypeScript parses TypeScript/JavaScript code and returns a sequence of tokens
func ParseTypeScript(code string) (TokenSequence, error) {
	var tokens TokenSequence

	// Trim any BOM or other markers
	code = strings.TrimSpace(code)

	for len(code) > 0 {
		// Try to match a comment
		if match := typescriptCommentRegex.FindString(code); match != "" {
			tokens = append(tokens, Token{Type: TokenComment, Text: match})
			code = code[len(match):]
			continue
		}

		// Try to match a template string
		if match := typescriptTemplateRegex.FindString(code); match != "" {
			tokens = append(tokens, Token{Type: TokenLiteral, Text: match})
			code = code[len(match):]
			continue
		}

		// Try to match a string literal
		if isStringStart(code) {
			end := findStringEnd(code)
			if end > 0 {
				tokens = append(tokens, Token{Type: TokenLiteral, Text: code[:end]})
				code = code[end:]
				continue
			}
		}

		// Try to match a number
		if match := typescriptNumberRegex.FindString(code); match != "" {
			tokens = append(tokens, Token{Type: TokenLiteral, Text: match})
			code = code[len(match):]
			continue
		}

		// Try to match an identifier or keyword
		if match := typescriptIdentifierRegex.FindString(code); match != "" {
			if typescriptKeywords[match] {
				tokens = append(tokens, Token{Type: TokenKeyword, Text: match})
			} else {
				tokens = append(tokens, Token{Type: TokenIdentifier, Text: match})
			}
			code = code[len(match):]
			continue
		}

		// Try to match whitespace
		if match := typescriptWhitespaceRegex.FindString(code); match != "" {
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
