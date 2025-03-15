package parsing

// Parser defines an interface for code parsers
type Parser interface {
	// Parse parses code and returns a sequence of tokens
	Parse(code string) (TokenSequence, error)
}

// GetParser returns a parser for the specified language
func GetParser(language string) Parser {
	switch language {
	case "python", "py":
		return &PythonParser{}
	case "typescript", "ts", "javascript", "js":
		return &TypeScriptParser{}
	case "bash", "sh", "shell":
		return &BashParser{}
	case "json":
		return &JSONParser{}
	default:
		return nil
	}
}
