package parsing

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

// TokenSequence is a sequence of tokens that represents parsed code
type TokenSequence []Token
