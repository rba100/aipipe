package parsing

import (
	"regexp"
	"strings"
)

var (
	// PowerShell keywords and common cmdlets
	powershellKeywords = map[string]bool{
		// Control flow keywords
		"if":         true,
		"else":       true,
		"elseif":     true,
		"switch":     true,
		"foreach":    true,
		"for":        true,
		"while":      true,
		"do":         true,
		"until":      true,
		"break":      true,
		"continue":   true,
		"return":     true,
		"exit":       true,
		"throw":      true,
		"try":        true,
		"catch":      true,
		"finally":    true,
		"trap":       true,
		
		// Function and parameter keywords
		"function":   true,
		"filter":     true,
		"param":      true,
		"begin":      true,
		"process":    true,
		"end":        true,
		"class":      true,
		"enum":       true,
		"using":      true,
		"namespace":  true,
		
		// Variable and scope keywords
		"global":     true,
		"local":      true,
		"private":    true,
		"script":     true,
		"static":     true,
		"hidden":     true,
		
		// Common cmdlets (case-insensitive)
		"get-process":      true,
		"get-childitem":    true,
		"get-content":      true,
		"set-content":      true,
		"set-location":     true,
		"get-location":     true,
		"write-host":       true,
		"write-output":     true,
		"write-error":      true,
		"write-warning":    true,
		"write-verbose":    true,
		"write-debug":      true,
		"read-host":        true,
		"select-object":    true,
		"where-object":     true,
		"foreach-object":   true,
		"sort-object":      true,
		"group-object":     true,
		"measure-object":   true,
		"compare-object":   true,
		"out-file":         true,
		"out-string":       true,
		"out-null":         true,
		"invoke-expression": true,
		"invoke-command":   true,
		"start-process":    true,
		"stop-process":     true,
		"get-service":      true,
		"start-service":    true,
		"stop-service":     true,
		"restart-service":  true,
		"new-object":       true,
		"remove-item":      true,
		"copy-item":        true,
		"move-item":        true,
		"rename-item":      true,
		"test-path":        true,
		"join-path":        true,
		"split-path":       true,
		"resolve-path":     true,
		"push-location":    true,
		"pop-location":     true,
		"import-module":    true,
		"export-module":    true,
		"get-module":       true,
		"remove-module":    true,
	}

	// Regular expressions for PowerShell tokens
	powershellVariableRegex     = regexp.MustCompile(`^\$([a-zA-Z_][a-zA-Z0-9_]*|\{[^}]*\}|[0-9]+|\$|\?|\^|args|input|env:[a-zA-Z_][a-zA-Z0-9_]*|global:[a-zA-Z_][a-zA-Z0-9_]*|local:[a-zA-Z_][a-zA-Z0-9_]*|script:[a-zA-Z_][a-zA-Z0-9_]*|private:[a-zA-Z_][a-zA-Z0-9_]*|_|true|false|null)`)
	powershellIdentifierRegex   = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*(-[a-zA-Z_][a-zA-Z0-9_]*)*`)
	powershellCommentRegex      = regexp.MustCompile(`^#.*`)
	powershellMultiCommentRegex = regexp.MustCompile(`^<#[\s\S]*?#>`)
	powershellOperatorRegex     = regexp.MustCompile(`^(-eq|-ne|-lt|-le|-gt|-ge|-like|-notlike|-match|-notmatch|-contains|-notcontains|-in|-notin|-replace|-split|-join|-and|-or|-not|-xor|-band|-bor|-bnot|-bxor|-shl|-shr|>>|<<|&&|\|\||==|!=|<=|>=|>|<|\+\+|--|\+|-|\*|/|%|=|\+=|-=|\*=|/=|%=|\|&)`)
	powershellWhitespaceRegex   = regexp.MustCompile(`^[ \t\r\n]+`)
	powershellNumberRegex       = regexp.MustCompile(`^0x[0-9a-fA-F]+|^[0-9]+(\.[0-9]+)?([eE][+-]?[0-9]+)?[lLdDfFmM]?`)
	powershellBooleanRegex      = regexp.MustCompile(`^\$(true|false|null)`)
	powershellHereStringRegex   = regexp.MustCompile(`^@['"][\s\S]*?['"]@`)
)

// PowerShellParser implements the Parser interface for PowerShell
type PowerShellParser struct{}

// Parse implements the Parser interface for PowerShell
func (p *PowerShellParser) Parse(code string) (TokenSequence, error) {
	return ParsePowerShell(code)
}

// isPowerShellStringStart checks if the code starts with a string delimiter
func isPowerShellStringStart(code string) (string, bool) {
	if len(code) == 0 {
		return "", false
	}

	if code[0] == '"' {
		return "\"", true
	}
	if code[0] == '\'' {
		return "'", true
	}

	return "", false
}

// findPowerShellStringEnd finds the end of a string
func findPowerShellStringEnd(code string, delimiter string) int {
	escaped := false

	for i := 1; i < len(code); i++ {
		if escaped {
			escaped = false
			continue
		}

		if code[i] == '`' {
			escaped = true
			continue
		}

		if string(code[i]) == delimiter {
			return i
		}
	}

	return len(code) - 1
}

// ParsePowerShell parses PowerShell code and returns a sequence of tokens
func ParsePowerShell(code string) (TokenSequence, error) {
	tokens := TokenSequence{}

	for len(code) > 0 {
		// Check for whitespace
		if match := powershellWhitespaceRegex.FindString(code); match != "" {
			tokens = append(tokens, Token{Type: TokenWhitespace, Text: match})
			code = code[len(match):]
			continue
		}

		// Check for multi-line comments first
		if match := powershellMultiCommentRegex.FindString(code); match != "" {
			tokens = append(tokens, Token{Type: TokenComment, Text: match})
			code = code[len(match):]
			continue
		}

		// Check for single-line comments
		if match := powershellCommentRegex.FindString(code); match != "" {
			tokens = append(tokens, Token{Type: TokenComment, Text: match})
			code = code[len(match):]
			continue
		}

		// Check for here-strings
		if match := powershellHereStringRegex.FindString(code); match != "" {
			tokens = append(tokens, Token{Type: TokenLiteral, Text: match})
			code = code[len(match):]
			continue
		}

		// Check for strings
		if delimiter, isString := isPowerShellStringStart(code); isString {
			end := findPowerShellStringEnd(code, delimiter)
			if end < 0 {
				end = len(code) - 1
			}
			tokens = append(tokens, Token{Type: TokenLiteral, Text: code[:end+1]})
			code = code[end+1:]
			continue
		}

		// Check for boolean literals
		if match := powershellBooleanRegex.FindString(code); match != "" {
			tokens = append(tokens, Token{Type: TokenLiteral, Text: match})
			code = code[len(match):]
			continue
		}

		// Check for variables
		if match := powershellVariableRegex.FindString(code); match != "" {
			tokens = append(tokens, Token{Type: TokenIdentifier, Text: match})
			code = code[len(match):]
			continue
		}

		// Check for numbers
		if match := powershellNumberRegex.FindString(code); match != "" {
			tokens = append(tokens, Token{Type: TokenLiteral, Text: match})
			code = code[len(match):]
			continue
		}

		// Check for operators
		if match := powershellOperatorRegex.FindString(code); match != "" {
			tokens = append(tokens, Token{Type: TokenOther, Text: match})
			code = code[len(match):]
			continue
		}

		// Check for keywords and identifiers
		if match := powershellIdentifierRegex.FindString(code); match != "" {
			lowerMatch := strings.ToLower(match)
			if powershellKeywords[lowerMatch] {
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