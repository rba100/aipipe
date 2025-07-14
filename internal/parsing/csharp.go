package parsing

import (
	"regexp"
)

var (
	// C# keywords
	csharpKeywords = map[string]bool{
		"abstract":    true,
		"as":          true,
		"base":        true,
		"bool":        true,
		"break":       true,
		"byte":        true,
		"case":        true,
		"catch":       true,
		"char":        true,
		"checked":     true,
		"class":       true,
		"const":       true,
		"continue":    true,
		"decimal":     true,
		"default":     true,
		"delegate":    true,
		"do":          true,
		"double":      true,
		"else":        true,
		"enum":        true,
		"event":       true,
		"explicit":    true,
		"extern":      true,
		"false":       true,
		"finally":     true,
		"fixed":       true,
		"float":       true,
		"for":         true,
		"foreach":     true,
		"goto":        true,
		"if":          true,
		"implicit":    true,
		"in":          true,
		"int":         true,
		"interface":   true,
		"internal":    true,
		"is":          true,
		"lock":        true,
		"long":        true,
		"namespace":   true,
		"new":         true,
		"null":        true,
		"object":      true,
		"operator":    true,
		"out":         true,
		"override":    true,
		"params":      true,
		"private":     true,
		"protected":   true,
		"public":      true,
		"readonly":    true,
		"ref":         true,
		"return":      true,
		"sbyte":       true,
		"sealed":      true,
		"short":       true,
		"sizeof":      true,
		"stackalloc":  true,
		"static":      true,
		"string":      true,
		"struct":      true,
		"switch":      true,
		"this":        true,
		"throw":       true,
		"true":        true,
		"try":         true,
		"typeof":      true,
		"uint":        true,
		"ulong":       true,
		"unchecked":   true,
		"unsafe":      true,
		"ushort":      true,
		"using":       true,
		"virtual":     true,
		"void":        true,
		"volatile":    true,
		"while":       true,
		"yield":       true,
		// Contextual keywords
		"add":         true,
		"async":       true,
		"await":       true,
		"dynamic":     true,
		"get":         true,
		"global":      true,
		"partial":     true,
		"remove":      true,
		"set":         true,
		"value":       true,
		"var":         true,
		"where":       true,
	}

	// Regular expressions for C# tokens
	csharpNumberRegex     = regexp.MustCompile(`^(0[xX][0-9a-fA-F]+[ULul]*|0[bB][01]+[ULul]*|[0-9]+(\.[0-9]+)?([eE][+-]?[0-9]+)?[fFdDmMULul]*)`)
	csharpIdentifierRegex = regexp.MustCompile(`^[a-zA-Z_@][a-zA-Z0-9_]*`)
	csharpCommentRegex    = regexp.MustCompile(`^(//.*|/\*[\s\S]*?\*/)`)
	csharpWhitespaceRegex = regexp.MustCompile(`^[ \t\r\n]+`)
	csharpVerbatimRegex   = regexp.MustCompile(`^@"(?:[^"]|"")*"`)
)

// CsharpParser implements the Parser interface for C# code
type CsharpParser struct{}

// Parse parses C# code and returns a sequence of tokens
func (p *CsharpParser) Parse(code string) (TokenSequence, error) {
	return ParseCsharp(code)
}

// isCsharpStringStart checks if the code starts with a string delimiter
func isCsharpStringStart(code string) bool {
	return len(code) > 0 && (code[0] == '"' || code[0] == '\'')
}

// findCsharpStringEnd finds the end of a string literal
func findCsharpStringEnd(code string) int {
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

// ParseCsharp parses C# code and returns a sequence of tokens
func ParseCsharp(code string) (TokenSequence, error) {
	var tokens TokenSequence

	// Process the code without trimming whitespace
	// This preserves indentation

	for len(code) > 0 {
		// Try to match whitespace first to preserve indentation
		if match := csharpWhitespaceRegex.FindString(code); match != "" {
			tokens = append(tokens, Token{Type: TokenWhitespace, Text: match})
			code = code[len(match):]
			continue
		}

		// Try to match a comment
		if match := csharpCommentRegex.FindString(code); match != "" {
			tokens = append(tokens, Token{Type: TokenComment, Text: match})
			code = code[len(match):]
			continue
		}

		// Try to match a verbatim string (@"...")
		if match := csharpVerbatimRegex.FindString(code); match != "" {
			tokens = append(tokens, Token{Type: TokenLiteral, Text: match})
			code = code[len(match):]
			continue
		}

		// Try to match a string literal
		if isCsharpStringStart(code) {
			end := findCsharpStringEnd(code)
			if end > 0 {
				tokens = append(tokens, Token{Type: TokenLiteral, Text: code[:end]})
				code = code[end:]
				continue
			}
		}

		// Try to match a number
		if match := csharpNumberRegex.FindString(code); match != "" {
			tokens = append(tokens, Token{Type: TokenLiteral, Text: match})
			code = code[len(match):]
			continue
		}

		// Try to match an identifier or keyword
		if match := csharpIdentifierRegex.FindString(code); match != "" {
			if csharpKeywords[match] {
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