package display

import (
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
// These are now defined in colours.go
// const (
// 	KeywordColor    = "\033[35m" // Magenta for keywords
// 	IdentifierColor = "\033[37m" // White for identifiers
// 	LiteralColor    = "\033[32m" // Green for literals
// 	CommentColor    = "\033[90m" // Grey for comments
// 	OtherColor      = "\033[36m" // Cyan for other tokens
// )

// SyntaxHighlighter handles syntax highlighting for code blocks
type SyntaxHighlighter struct {
	languageRegex   *regexp.Regexp
	currentLanguage string
}

// NewSyntaxHighlighter creates a new syntax highlighter
func NewSyntaxHighlighter() *SyntaxHighlighter {
	h := &SyntaxHighlighter{
		languageRegex: regexp.MustCompile(`^\s*\x60\x60\x60(\w+)`),
	}

	return h
}

// parsePython parses Python code and returns a sequence of tokens
func (h *SyntaxHighlighter) parsePython(code string) []Token {
	// Use the parsing package's Python parser
	parsingTokens, err := parsing.ParsePython(code)
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
			highlighted.WriteString(TokenKeywordColor + token.Text + ResetFormat)
		case TokenIdentifier:
			highlighted.WriteString(TokenIdentifierColor + token.Text + ResetFormat)
		case TokenLiteral:
			highlighted.WriteString(TokenLiteralColor + token.Text + ResetFormat)
		case TokenComment:
			highlighted.WriteString(TokenCommentColor + token.Text + ResetFormat)
		case TokenWhitespace:
			highlighted.WriteString(token.Text)
		default:
			highlighted.WriteString(TokenOtherColor + token.Text + ResetFormat)
		}
	}

	return highlighted.String()
}

// ExtractLanguage extracts the language identifier from a code block start line
func (h *SyntaxHighlighter) ExtractLanguage(line string) string {
	matches := h.languageRegex.FindStringSubmatch(line)
	if len(matches) >= 2 {
		h.currentLanguage = strings.ToLower(matches[1])
		return h.currentLanguage
	}
	return ""
}

// ProcessCodeLine processes a line of code for syntax highlighting
func (h *SyntaxHighlighter) ProcessCodeLine(line string) {
	// This method is not used in the current implementation
	// It's here for compatibility with potential future extensions
}
