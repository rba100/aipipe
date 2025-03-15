package parsing

import (
	"regexp"
)

var (
	// Bash keywords
	bashKeywords = map[string]bool{
		"if":       true,
		"then":     true,
		"else":     true,
		"elif":     true,
		"fi":       true,
		"case":     true,
		"esac":     true,
		"for":      true,
		"select":   true,
		"while":    true,
		"until":    true,
		"do":       true,
		"done":     true,
		"in":       true,
		"function": true,
		"time":     true,
		"declare":  true,
		"local":    true,
		"readonly": true,
		"export":   true,
		"alias":    true,
		"unalias":  true,
		"source":   true,
		"let":      true,
		"return":   true,
		"exit":     true,
		"exec":     true,
		"set":      true,
		"unset":    true,
		"trap":     true,
		"break":    true,
		"continue": true,
		"eval":     true,
		"cd":       true,
		"pwd":      true,
		"echo":     true,
		"printf":   true,
		"read":     true,
		"shift":    true,
		"test":     true,
		"[":        true,
		"]":        true,
		"[[":       true,
		"]]":       true,
	}

	// Regular expressions for Bash tokens
	bashVariableRegex     = regexp.MustCompile(`^\$([a-zA-Z_][a-zA-Z0-9_]*|\{[a-zA-Z_][a-zA-Z0-9_]*\}|[0-9])`)
	bashIdentifierRegex   = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*`)
	bashCommentRegex      = regexp.MustCompile(`^#.*`)
	bashOperatorRegex     = regexp.MustCompile(`^(&&|\|\||>>|<<|>=|<=|==|!=|>|<|\+|-|\*|/|=|;|\||&)`)
	bashRedirectionRegex  = regexp.MustCompile(`^([0-9]*>>?|[0-9]*<<?)`)
	bashWhitespaceRegex   = regexp.MustCompile(`^[ \t\r\n]+`)
	bashHeredocStartRegex = regexp.MustCompile(`^<<-?\s*([a-zA-Z_][a-zA-Z0-9_]*|'[^']*'|"[^"]*")`)
	bashProcessSubRegex   = regexp.MustCompile(`^[<>]\(`)
	bashNumberRegex       = regexp.MustCompile(`^[0-9]+(\.[0-9]+)?`)
)

// BashParser implements the Parser interface for Bash shell commands
type BashParser struct{}

// Parse implements the Parser interface for Bash
func (p *BashParser) Parse(code string) (TokenSequence, error) {
	return ParseBash(code)
}

// isBashStringStart checks if the code starts with a string delimiter
func isBashStringStart(code string) (string, bool) {
	if len(code) == 0 {
		return "", false
	}

	if code[0] == '"' {
		return "\"", true
	}
	if code[0] == '\'' {
		return "'", true
	}
	if len(code) >= 1 && code[0] == '`' {
		return "`", true
	}

	return "", false
}

// findBashStringEnd finds the end of a string
func findBashStringEnd(code string, delimiter string) int {
	escaped := false

	for i := 1; i < len(code); i++ {
		if escaped {
			escaped = false
			continue
		}

		if code[i] == '\\' && delimiter != "'" {
			escaped = true
			continue
		}

		if string(code[i]) == delimiter {
			return i
		}
	}

	return len(code) - 1
}

// ParseBash parses Bash shell commands and returns a sequence of tokens
func ParseBash(code string) (TokenSequence, error) {
	tokens := TokenSequence{}

	for len(code) > 0 {
		// Check for whitespace
		if match := bashWhitespaceRegex.FindString(code); match != "" {
			tokens = append(tokens, Token{Type: TokenWhitespace, Text: match})
			code = code[len(match):]
			continue
		}

		// Check for comments
		if match := bashCommentRegex.FindString(code); match != "" {
			tokens = append(tokens, Token{Type: TokenComment, Text: match})
			code = code[len(match):]
			continue
		}

		// Check for here-documents
		if match := bashHeredocStartRegex.FindString(code); match != "" {
			tokens = append(tokens, Token{Type: TokenOther, Text: match})
			code = code[len(match):]
			continue
		}

		// Check for process substitution
		if match := bashProcessSubRegex.FindString(code); match != "" {
			tokens = append(tokens, Token{Type: TokenOther, Text: match})
			code = code[len(match):]
			continue
		}

		// Check for strings
		if delimiter, isString := isBashStringStart(code); isString {
			end := findBashStringEnd(code, delimiter)
			if end < 0 {
				end = len(code) - 1
			}
			tokens = append(tokens, Token{Type: TokenLiteral, Text: code[:end+1]})
			code = code[end+1:]
			continue
		}

		// Check for variables
		if match := bashVariableRegex.FindString(code); match != "" {
			tokens = append(tokens, Token{Type: TokenIdentifier, Text: match})
			code = code[len(match):]
			continue
		}

		// Check for numbers
		if match := bashNumberRegex.FindString(code); match != "" {
			tokens = append(tokens, Token{Type: TokenLiteral, Text: match})
			code = code[len(match):]
			continue
		}

		// Check for redirections
		if match := bashRedirectionRegex.FindString(code); match != "" {
			tokens = append(tokens, Token{Type: TokenOther, Text: match})
			code = code[len(match):]
			continue
		}

		// Check for operators
		if match := bashOperatorRegex.FindString(code); match != "" {
			tokens = append(tokens, Token{Type: TokenOther, Text: match})
			code = code[len(match):]
			continue
		}

		// Check for keywords and identifiers
		if match := bashIdentifierRegex.FindString(code); match != "" {
			if bashKeywords[match] {
				tokens = append(tokens, Token{Type: TokenKeyword, Text: match})
			} else {
				tokens = append(tokens, Token{Type: TokenIdentifier, Text: match})
			}
			code = code[len(match):]
			continue
		}

		// Handle other characters
		tokens = append(tokens, Token{Type: TokenOther, Text: string(code[0])})
		code = code[1:]
	}

	return tokens, nil
}
