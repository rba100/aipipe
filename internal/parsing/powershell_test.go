package parsing

import (
	"testing"
)

func TestPowerShellParser(t *testing.T) {
	parser := &PowerShellParser{}

	testCases := []struct {
		name     string
		input    string
		expected int // Expected number of tokens
	}{
		{
			name:     "Simple command",
			input:    "Get-Process",
			expected: 1, // "Get-Process"
		},
		{
			name:     "Command with parameter",
			input:    "Get-Process -Name powershell",
			expected: 6, // "Get-Process", " ", "-", "Name", " ", "powershell"
		},
		{
			name:     "Command with variable",
			input:    "Write-Host $env:USERNAME",
			expected: 5, // "Write-Host", " ", "$env", ":", "USERNAME"
		},
		{
			name:     "Command with string",
			input:    "Write-Host \"Hello World\"",
			expected: 3, // "Write-Host", " ", "\"Hello World\""
		},
		{
			name:     "Command with single quotes",
			input:    "Write-Host 'Hello World'",
			expected: 3, // "Write-Host", " ", "'Hello World'"
		},
		{
			name:     "Variable assignment",
			input:    "$name = \"John\"",
			expected: 5, // "$name", " ", "=", " ", "\"John\""
		},
		{
			name:     "If statement",
			input:    "if ($true) { Write-Host \"Yes\" }",
			expected: 13, // "if", " ", "(", "$true", ")", " ", "{", " ", "Write-Host", " ", "\"Yes\"", " ", "}"
		},
		{
			name:     "PowerShell operators",
			input:    "$a -eq $b",
			expected: 5, // "$a", " ", "-eq", " ", "$b"
		},
		{
			name:     "Comment",
			input:    "Get-Process # Get running processes",
			expected: 3, // "Get-Process", " ", "# Get running processes"
		},
		{
			name:     "Multi-line comment",
			input:    "<# This is a\n   multi-line comment #>",
			expected: 1, // "<# This is a\n   multi-line comment #>"
		},
		{
			name:     "Boolean literals",
			input:    "$true $false $null",
			expected: 5, // "$true", " ", "$false", " ", "$null"
		},
		{
			name:     "Numbers",
			input:    "123 3.14 0xFF",
			expected: 5, // "123", " ", "3.14", " ", "0xFF"
		},
		{
			name:     "Foreach loop",
			input:    "foreach ($item in $array) { Write-Host $item }",
			expected: 17, // Complex parsing with multiple tokens
		},
		{
			name:     "Here-string",
			input:    "@\"\nHello\nWorld\n\"@",
			expected: 1, // "@\"\nHello\nWorld\n\"@"
		},
		{
			name:     "Pipeline",
			input:    "Get-Process | Where-Object { $_.Name -like '*powershell*' }",
			expected: 17, // Complex pipeline with multiple tokens
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tokens, err := parser.Parse(tc.input)
			if err != nil {
				t.Fatalf("Error parsing PowerShell code: %v", err)
			}

			if len(tokens) != tc.expected {
				t.Errorf("Expected %d tokens, got %d. Tokens: %v", tc.expected, len(tokens), tokens)
			}
		})
	}
}

func TestPowerShellKeywordIdentification(t *testing.T) {
	parser := &PowerShellParser{}

	code := "if ($true) { Write-Host \"Hello\" } else { Write-Error \"Error\" }"
	tokens, err := parser.Parse(code)
	if err != nil {
		t.Fatalf("Error parsing PowerShell code: %v", err)
	}

	// Print tokens for debugging
	for i, token := range tokens {
		t.Logf("Token %d: Type=%v, Text=%q", i, token.Type, token.Text)
	}

	// Check that "if", "else", and "Write-Host" are identified as keywords
	expectedKeywords := map[string]bool{
		"if":         true,
		"else":       true,
		"Write-Host": true,
	}

	for i, token := range tokens {
		if expectedKeywords[token.Text] && token.Type != TokenKeyword {
			t.Errorf("Expected token at index %d with text %q to be a keyword, got %v", i, token.Text, token.Type)
		}
	}
}

func TestPowerShellVariables(t *testing.T) {
	parser := &PowerShellParser{}

	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Simple variable",
			input:    "$name",
			expected: "$name",
		},
		{
			name:     "Environment variable",
			input:    "$env:PATH",
			expected: "$env:PATH",
		},
		{
			name:     "Global variable",
			input:    "$global:variable",
			expected: "$global:variable",
		},
		{
			name:     "Complex variable",
			input:    "${complex variable}",
			expected: "${complex variable}",
		},
		{
			name:     "Special variables",
			input:    "$$ $? $^",
			expected: "$$ $? $^",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tokens, err := parser.Parse(tc.input)
			if err != nil {
				t.Fatalf("Error parsing PowerShell code: %v", err)
			}

			// Find variable tokens
			var foundVariables []string
			for _, token := range tokens {
				if token.Type == TokenIdentifier && token.Text[0] == '$' {
					foundVariables = append(foundVariables, token.Text)
				}
			}

			if len(foundVariables) == 0 {
				t.Errorf("No variables found in: %s", tc.input)
			}
		})
	}
}

func TestPowerShellOperators(t *testing.T) {
	parser := &PowerShellParser{}

	code := "$a -eq $b -and $c -ne $d"
	tokens, err := parser.Parse(code)
	if err != nil {
		t.Fatalf("Error parsing PowerShell code: %v", err)
	}

	// Check that operators are identified
	expectedOperators := map[string]bool{
		"-eq":  true,
		"-and": true,
		"-ne":  true,
	}

	for i, token := range tokens {
		if expectedOperators[token.Text] && token.Type != TokenOther {
			t.Errorf("Expected token at index %d with text %q to be an operator, got %v", i, token.Text, token.Type)
		}
	}
}

func TestGetParserPowerShell(t *testing.T) {
	// Test all PowerShell variations
	testCases := []string{"powershell", "ps1", "ps"}
	
	for _, lang := range testCases {
		t.Run(lang, func(t *testing.T) {
			parser := GetParser(lang)
			if parser == nil {
				t.Errorf("Expected parser for language '%s', got nil", lang)
				return
			}
			
			// Verify it's a PowerShell parser
			if _, ok := parser.(*PowerShellParser); !ok {
				t.Errorf("Expected PowerShellParser for language '%s', got %T", lang, parser)
				return
			}
			
			// Test parsing some PowerShell code
			code := `# PowerShell comment
$name = "test"
Write-Host $name`
			
			tokens, err := parser.Parse(code)
			if err != nil {
				t.Errorf("Failed to parse PowerShell code with language '%s': %v", lang, err)
				return
			}
			
			if len(tokens) == 0 {
				t.Errorf("Expected tokens for PowerShell code with language '%s', got 0", lang)
			}
		})
	}
}