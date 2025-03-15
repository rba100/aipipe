package parsing

import (
	"testing"
)

func TestJSONParser(t *testing.T) {
	parser := &JSONParser{}

	testCases := []struct {
		name     string
		input    string
		expected int // Expected number of tokens
	}{
		{
			name:     "Empty object",
			input:    "{}",
			expected: 2, // "{", "}"
		},
		{
			name:     "Simple object",
			input:    "{\"key\": \"value\"}",
			expected: 6, // "{", "\"key\"", ":", " ", "\"value\"", "}"
		},
		{
			name:     "Object with multiple properties",
			input:    "{\"key1\": \"value1\", \"key2\": 42}",
			expected: 12, // "{", "\"key1\"", ":", " ", "\"value1\"", ",", " ", "\"key2\"", ":", " ", "42", "}"
		},
		{
			name:     "Nested object",
			input:    "{\"outer\": {\"inner\": \"value\"}}",
			expected: 11, // "{", "\"outer\"", ":", " ", "{", "\"inner\"", ":", " ", "\"value\"", "}", "}"
		},
		{
			name:     "Array",
			input:    "[1, 2, 3]",
			expected: 9, // "[", "1", ",", " ", "2", ",", " ", "3", "]"
		},
		{
			name:     "Object with array",
			input:    "{\"items\": [1, 2, 3]}",
			expected: 14, // "{", "\"items\"", ":", " ", "[", "1", ",", " ", "2", ",", " ", "3", "]", "}"
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tokens, err := parser.Parse(tc.input)
			if err != nil {
				t.Fatalf("Error parsing JSON: %v", err)
			}

			if len(tokens) != tc.expected {
				t.Errorf("Expected %d tokens, got %d", tc.expected, len(tokens))
				for i, token := range tokens {
					t.Logf("Token %d: Type=%d, Text=%q", i, token.Type, token.Text)
				}
			}
		})
	}
}

func TestJSONObjectKeyIdentification(t *testing.T) {
	testCases := []struct {
		name  string
		input string
		keys  []string // Expected object keys
	}{
		{
			name:  "Simple object",
			input: "{\"key\": \"value\"}",
			keys:  []string{"\"key\""},
		},
		{
			name:  "Multiple keys",
			input: "{\"key1\": \"value1\", \"key2\": 42}",
			keys:  []string{"\"key1\"", "\"key2\""},
		},
		{
			name:  "Nested object",
			input: "{\"outer\": {\"inner\": \"value\"}}",
			keys:  []string{"\"outer\"", "\"inner\""},
		},
		{
			name:  "Object with array",
			input: "{\"items\": [1, 2, 3]}",
			keys:  []string{"\"items\""},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tokens, err := ParseJSON(tc.input)
			if err != nil {
				t.Fatalf("Error parsing JSON: %v", err)
			}

			// Count how many keys we found
			foundKeys := 0
			for _, token := range tokens {
				if token.Type == TokenIdentifier {
					foundKeys++
					// Check if this key is in our expected list
					found := false
					for _, expectedKey := range tc.keys {
						if token.Text == expectedKey {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("Unexpected key identified: %q", token.Text)
					}
				}
			}

			if foundKeys != len(tc.keys) {
				t.Errorf("Expected %d keys, found %d", len(tc.keys), foundKeys)
				for i, token := range tokens {
					t.Logf("Token %d: Type=%d, Text=%q", i, token.Type, token.Text)
				}
			}
		})
	}
}

func TestJSONStringLiteralVsObjectKey(t *testing.T) {
	input := "{\"key\": \"value\"}"
	tokens, err := ParseJSON(input)
	if err != nil {
		t.Fatalf("Error parsing JSON: %v", err)
	}

	// Find all string tokens
	var keyToken, valueToken Token
	for _, token := range tokens {
		if token.Text == "\"key\"" {
			keyToken = token
		} else if token.Text == "\"value\"" {
			valueToken = token
		}
	}

	// Verify that "key" is identified as an identifier
	if keyToken.Type != TokenIdentifier {
		t.Errorf("Expected \"key\" to be TokenIdentifier (type %d), got type %d", TokenIdentifier, keyToken.Type)
	}

	// Verify that "value" is identified as a literal
	if valueToken.Type != TokenLiteral {
		t.Errorf("Expected \"value\" to be TokenLiteral (type %d), got type %d", TokenLiteral, valueToken.Type)
	}
}
