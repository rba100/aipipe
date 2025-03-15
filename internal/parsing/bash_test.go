package parsing

import (
	"testing"
)

func TestBashParser(t *testing.T) {
	parser := &BashParser{}

	testCases := []struct {
		name     string
		input    string
		expected int // Expected number of tokens
	}{
		{
			name:     "Simple command",
			input:    "echo hello world",
			expected: 5, // "echo", " ", "hello", " ", "world"
		},
		{
			name:     "Command with variable",
			input:    "echo $HOME",
			expected: 3, // "echo", " ", "$HOME"
		},
		{
			name:     "Command with string",
			input:    "echo \"hello world\"",
			expected: 3, // "echo", " ", "\"hello world\""
		},
		{
			name:     "Command with single quotes",
			input:    "echo 'hello world'",
			expected: 3, // "echo", " ", "'hello world'"
		},
		{
			name:     "Command with backticks",
			input:    "echo `date`",
			expected: 3, // "echo", " ", "`date`"
		},
		{
			name:     "Command with redirection",
			input:    "ls > output.txt",
			expected: 7, // "ls", " ", ">", " ", "output", ".", "txt"
		},
		{
			name:     "Command with pipe",
			input:    "ls | grep file",
			expected: 7, // "ls", " ", "|", " ", "grep", " ", "file"
		},
		{
			name:     "Command with here-document",
			input:    "cat <<EOF\nHello\nEOF",
			expected: 7, // "cat", " ", "<<EOF", "\n", "Hello", "\n", "EOF"
		},
		{
			name:     "If statement",
			input:    "if [ $a -eq 5 ]; then echo \"equal\"; fi",
			expected: 22, // Complex parsing with multiple tokens
		},
		{
			name:     "Command with comment",
			input:    "echo hello # This is a comment",
			expected: 5, // "echo", " ", "hello", " ", "# This is a comment"
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tokens, err := parser.Parse(tc.input)
			if err != nil {
				t.Fatalf("Error parsing bash code: %v", err)
			}

			if len(tokens) != tc.expected {
				t.Errorf("Expected %d tokens, got %d. Tokens: %v", tc.expected, len(tokens), tokens)
			}
		})
	}
}

func TestBashKeywordIdentification(t *testing.T) {
	parser := &BashParser{}

	code := "if [ $a -eq 5 ]; then echo \"equal\"; fi"
	tokens, err := parser.Parse(code)
	if err != nil {
		t.Fatalf("Error parsing bash code: %v", err)
	}

	// Print tokens for debugging
	for i, token := range tokens {
		t.Logf("Token %d: Type=%v, Text=%q", i, token.Type, token.Text)
	}

	// Check that "if", "then", and "fi" are identified as keywords
	// Based on the actual token positions in the parsed output
	expectedKeywords := map[string]bool{
		"if":   true,
		"then": true,
		"fi":   true,
	}

	for i, token := range tokens {
		if expectedKeywords[token.Text] && token.Type != TokenKeyword {
			t.Errorf("Expected token at index %d with text %q to be a keyword, got %v", i, token.Text, token.Type)
		}
	}
}
