package display

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/rba100/aipipe/internal/parsing"
)

// TokenType represents the type of a token
type TokenType int

const (
	// TokenOther represents miscellaneous tokens like operators, punctuation, etc.
	TokenOther TokenType = iota
	// TokenKeyword represents language keywords
	TokenKeyword
	// TokenIdentifier represents variable names, function names, etc.
	TokenIdentifier
	// TokenLiteral represents string, number, and other literals
	TokenLiteral
	// TokenComment represents comments
	TokenComment
	// TokenWhitespace represents spaces, tabs, newlines
	TokenWhitespace
)

// Token represents a single token in the parsed code
type Token struct {
	// Type is the type of the token
	Type TokenType
	// Text is the actual text content of the token
	Text string
}

// ANSI color codes for syntax highlighting
const (
	KeywordColor    = "\033[35m" // Magenta for keywords
	IdentifierColor = "\033[37m" // White for identifiers
	LiteralColor    = "\033[32m" // Green for literals
	CommentColor    = "\033[90m" // Grey for comments
	OtherColor      = "\033[36m" // Cyan for other tokens
)

// SyntaxHighlighter handles syntax highlighting for code blocks
type SyntaxHighlighter struct {
	languageRegex   *regexp.Regexp
	currentLanguage string

	// Python specific regexes
	pythonKeywords        map[string]bool
	pythonNumberRegex     *regexp.Regexp
	pythonIdentifierRegex *regexp.Regexp
	pythonCommentRegex    *regexp.Regexp
	pythonWhitespaceRegex *regexp.Regexp
}

// NewSyntaxHighlighter creates a new syntax highlighter
func NewSyntaxHighlighter() *SyntaxHighlighter {
	h := &SyntaxHighlighter{
		languageRegex: regexp.MustCompile(`^\s*\x60\x60\x60(\w+)`),

		// Python specific regexes
		pythonKeywords: map[string]bool{
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
		},
		pythonNumberRegex:     regexp.MustCompile(`^[0-9]+(\.[0-9]+)?([eE][+-]?[0-9]+)?`),
		pythonIdentifierRegex: regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*`),
		pythonCommentRegex:    regexp.MustCompile(`^#.*`),
		pythonWhitespaceRegex: regexp.MustCompile(`^[ \t\r\n]+`),
	}

	return h
}

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

// parsePython parses Python code and returns a sequence of tokens
func (h *SyntaxHighlighter) parsePython(code string) []Token {
	var tokens []Token

	// Process the code without trimming whitespace
	// This preserves indentation

	for len(code) > 0 {
		// Try to match whitespace first to preserve indentation
		if match := h.pythonWhitespaceRegex.FindString(code); match != "" {
			tokens = append(tokens, Token{Type: TokenWhitespace, Text: match})
			code = code[len(match):]
			continue
		}

		// Try to match a comment
		if match := h.pythonCommentRegex.FindString(code); match != "" {
			tokens = append(tokens, Token{Type: TokenComment, Text: match})
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
		if match := h.pythonNumberRegex.FindString(code); match != "" {
			tokens = append(tokens, Token{Type: TokenLiteral, Text: match})
			code = code[len(match):]
			continue
		}

		// Try to match an identifier or keyword
		if match := h.pythonIdentifierRegex.FindString(code); match != "" {
			if h.pythonKeywords[match] {
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

	return tokens
}

// parseTypeScript parses TypeScript/JavaScript code and returns a sequence of tokens
func (h *SyntaxHighlighter) parseTypeScript(code string) []Token {
	// Use the parsing package's TypeScript parser
	parsingTokens, err := parsing.ParseTypeScript(code)
	if err != nil {
		// If there's an error, return a single token with the original code
		return []Token{{Type: TokenOther, Text: code}}
	}

	// Convert parsing.Token to display.Token
	var tokens []Token
	for _, token := range parsingTokens {
		tokens = append(tokens, Token{
			Type: TokenType(token.Type), // TokenType enums match between packages
			Text: token.Text,
		})
	}

	return tokens
}

// HighlightCode highlights code based on the language identifier
func (h *SyntaxHighlighter) HighlightCode(code string, language string) string {
	var tokens []Token

	// Select the appropriate parser based on the language
	switch language {
	case "python", "py":
		tokens = h.parsePython(code)
	case "typescript", "ts", "javascript", "js":
		tokens = h.parseTypeScript(code)
	default:
		// For unsupported languages, just return the code as is
		return code
	}

	// Build the highlighted code
	var highlighted strings.Builder
	for _, token := range tokens {
		switch token.Type {
		case TokenKeyword:
			highlighted.WriteString(KeywordColor + token.Text + Reset)
		case TokenIdentifier:
			highlighted.WriteString(IdentifierColor + token.Text + Reset)
		case TokenLiteral:
			highlighted.WriteString(LiteralColor + token.Text + Reset)
		case TokenComment:
			highlighted.WriteString(CommentColor + token.Text + Reset)
		case TokenWhitespace:
			highlighted.WriteString(token.Text)
		default:
			highlighted.WriteString(OtherColor + token.Text + Reset)
		}
	}

	return highlighted.String()
}

// ExtractLanguage extracts the language identifier from a code block start line
func (h *SyntaxHighlighter) ExtractLanguage(line string) string {
	matches := h.languageRegex.FindStringSubmatch(line)
	if len(matches) >= 2 {
		return strings.ToLower(matches[1])
	}
	return ""
}

// ProcessCodeLine processes a line of code with syntax highlighting
func (h *SyntaxHighlighter) ProcessCodeLine(line string) {
	// Check if this is a code block start line with a language identifier
	if h.languageRegex.MatchString(line) {
		h.currentLanguage = h.ExtractLanguage(line)
		fmt.Print(Cyan)
		fmt.Print(line)
		return
	}

	// Check if this is a code block end line
	if regexp.MustCompile(`^\s*\x60\x60\x60\s*$`).MatchString(line) {
		h.currentLanguage = ""
		fmt.Print(Cyan)
		fmt.Print(line)
		return
	}

	// This is a line of code within the code block
	if h.currentLanguage != "" {
		highlightedLine := h.HighlightCode(line, h.currentLanguage)
		fmt.Print(highlightedLine)
	} else {
		// If we don't know the language, just print in cyan
		fmt.Print(Cyan)
		fmt.Print(line)
	}
}
