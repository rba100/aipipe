package parsing

import (
	"testing"
)

func TestCsharpParser(t *testing.T) {
	parser := &CsharpParser{}

	testCases := []struct {
		name     string
		input    string
		expected int // Expected number of tokens
	}{
		{
			name:     "Empty class",
			input:    "class MyClass { }",
			expected: 7, // "class", " ", "MyClass", " ", "{", " ", "}"
		},
		{
			name:     "Simple method",
			input:    "public void Method() { }",
			expected: 11, // "public", " ", "void", " ", "Method", "(", ")", " ", "{", " ", "}"
		},
		{
			name:     "Variable declaration",
			input:    "int x = 42;",
			expected: 8, // "int", " ", "x", " ", "=", " ", "42", ";"
		},
		{
			name:     "String literal",
			input:    "string msg = \"Hello, World!\";",
			expected: 8, // "string", " ", "msg", " ", "=", " ", "\"Hello, World!\"", ";"
		},
		{
			name:     "Verbatim string",
			input:    "string path = @\"C:\\Users\\John\";",
			expected: 8, // "string", " ", "path", " ", "=", " ", "@\"C:\\Users\\John\"", ";"
		},
		{
			name:     "Comment",
			input:    "// This is a comment",
			expected: 1, // "// This is a comment"
		},
		{
			name:     "Multi-line comment",
			input:    "/* This is a\nmulti-line comment */",
			expected: 1, // "/* This is a\nmulti-line comment */"
		},
		{
			name:     "If statement",
			input:    "if (x > 0) { return true; }",
			expected: 18, // "if", " ", "(", "x", " ", ">", " ", "0", ")", " ", "{", " ", "return", " ", "true", ";", " ", "}"
		},
		{
			name:     "Using statement",
			input:    "using System;",
			expected: 4, // "using", " ", "System", ";"
		},
		{
			name:     "Namespace",
			input:    "namespace MyNamespace { }",
			expected: 7, // "namespace", " ", "MyNamespace", " ", "{", " ", "}"
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tokens, err := parser.Parse(tc.input)
			if err != nil {
				t.Fatalf("Error parsing C#: %v", err)
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

func TestCsharpKeywordIdentification(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		keywords []string // Expected keywords
	}{
		{
			name:     "Class declaration",
			input:    "public class MyClass { }",
			keywords: []string{"public", "class"},
		},
		{
			name:     "Method declaration",
			input:    "private static void Main(string[] args) { }",
			keywords: []string{"private", "static", "void", "string"},
		},
		{
			name:     "If statement",
			input:    "if (x > 0) return true; else return false;",
			keywords: []string{"if", "return", "true", "else", "return", "false"},
		},
		{
			name:     "Using and namespace",
			input:    "using System; namespace MyApp { }",
			keywords: []string{"using", "namespace"},
		},
		{
			name:     "Variable types",
			input:    "int x; bool flag; string text; double val;",
			keywords: []string{"int", "bool", "string", "double"},
		},
		{
			name:     "Async method",
			input:    "public async Task<bool> GetDataAsync() { }",
			keywords: []string{"public", "async", "bool"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tokens, err := ParseCsharp(tc.input)
			if err != nil {
				t.Fatalf("Error parsing C#: %v", err)
			}

			// Count keywords found
			var foundKeywords []string
			for _, token := range tokens {
				if token.Type == TokenKeyword {
					foundKeywords = append(foundKeywords, token.Text)
				}
			}

			if len(foundKeywords) != len(tc.keywords) {
				t.Errorf("Expected %d keywords, found %d", len(tc.keywords), len(foundKeywords))
				t.Logf("Expected: %v", tc.keywords)
				t.Logf("Found: %v", foundKeywords)
			}

			// Check each expected keyword is found
			for _, expectedKeyword := range tc.keywords {
				found := false
				for _, foundKeyword := range foundKeywords {
					if foundKeyword == expectedKeyword {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected keyword %q not found", expectedKeyword)
				}
			}
		})
	}
}

func TestCsharpStringLiterals(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		literals []string // Expected string literals
	}{
		{
			name:     "Regular string",
			input:    "string msg = \"Hello, World!\";",
			literals: []string{"\"Hello, World!\""},
		},
		{
			name:     "Verbatim string",
			input:    "string path = @\"C:\\Users\\John\";",
			literals: []string{"@\"C:\\Users\\John\""},
		},
		{
			name:     "Character literal",
			input:    "char c = 'A';",
			literals: []string{"'A'"},
		},
		{
			name:     "Multiple strings",
			input:    "string a = \"first\"; string b = \"second\";",
			literals: []string{"\"first\"", "\"second\""},
		},
		{
			name:     "Escaped string",
			input:    "string msg = \"Hello\\nWorld\";",
			literals: []string{"\"Hello\\nWorld\""},
		},
		{
			name:     "Verbatim string with quotes",
			input:    "string quote = @\"He said \"\"Hello\"\"\";",
			literals: []string{"@\"He said \"\"Hello\"\"\""},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tokens, err := ParseCsharp(tc.input)
			if err != nil {
				t.Fatalf("Error parsing C#: %v", err)
			}

			// Find string literals
			var foundLiterals []string
			for _, token := range tokens {
				if token.Type == TokenLiteral && (token.Text[0] == '"' || token.Text[0] == '\'' || (len(token.Text) > 1 && token.Text[0] == '@' && token.Text[1] == '"')) {
					foundLiterals = append(foundLiterals, token.Text)
				}
			}

			if len(foundLiterals) != len(tc.literals) {
				t.Errorf("Expected %d string literals, found %d", len(tc.literals), len(foundLiterals))
				t.Logf("Expected: %v", tc.literals)
				t.Logf("Found: %v", foundLiterals)
			}

			// Check each expected literal is found
			for _, expectedLiteral := range tc.literals {
				found := false
				for _, foundLiteral := range foundLiterals {
					if foundLiteral == expectedLiteral {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected literal %q not found", expectedLiteral)
				}
			}
		})
	}
}

func TestCsharpNumbers(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		numbers  []string // Expected number literals
	}{
		{
			name:     "Integer",
			input:    "int x = 42;",
			numbers:  []string{"42"},
		},
		{
			name:     "Float",
			input:    "float f = 3.14f;",
			numbers:  []string{"3.14f"},
		},
		{
			name:     "Double",
			input:    "double d = 3.14159;",
			numbers:  []string{"3.14159"},
		},
		{
			name:     "Decimal",
			input:    "decimal m = 123.45m;",
			numbers:  []string{"123.45m"},
		},
		{
			name:     "Hexadecimal",
			input:    "int hex = 0xFF;",
			numbers:  []string{"0xFF"},
		},
		{
			name:     "Binary",
			input:    "int bin = 0b1010;",
			numbers:  []string{"0b1010"},
		},
		{
			name:     "Scientific notation",
			input:    "double sci = 1.23e4;",
			numbers:  []string{"1.23e4"},
		},
		{
			name:     "Long",
			input:    "long big = 123456789L;",
			numbers:  []string{"123456789L"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tokens, err := ParseCsharp(tc.input)
			if err != nil {
				t.Fatalf("Error parsing C#: %v", err)
			}

			// Find number literals
			var foundNumbers []string
			for _, token := range tokens {
				if token.Type == TokenLiteral && (token.Text[0] >= '0' && token.Text[0] <= '9') {
					foundNumbers = append(foundNumbers, token.Text)
				}
			}

			if len(foundNumbers) != len(tc.numbers) {
				t.Errorf("Expected %d number literals, found %d", len(tc.numbers), len(foundNumbers))
				t.Logf("Expected: %v", tc.numbers)
				t.Logf("Found: %v", foundNumbers)
			}

			// Check each expected number is found
			for _, expectedNumber := range tc.numbers {
				found := false
				for _, foundNumber := range foundNumbers {
					if foundNumber == expectedNumber {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected number %q not found", expectedNumber)
				}
			}
		})
	}
}

func TestCsharpComments(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		comments []string // Expected comments
	}{
		{
			name:     "Single line comment",
			input:    "// This is a comment",
			comments: []string{"// This is a comment"},
		},
		{
			name:     "Multi-line comment",
			input:    "/* This is a multi-line comment */",
			comments: []string{"/* This is a multi-line comment */"},
		},
		{
			name:     "Comment with code",
			input:    "int x = 5; // Initialize x",
			comments: []string{"// Initialize x"},
		},
		{
			name:     "Multiple comments",
			input:    "// First comment\n/* Second comment */",
			comments: []string{"// First comment", "/* Second comment */"},
		},
		{
			name:     "XML documentation comment",
			input:    "/// <summary>This is a summary</summary>",
			comments: []string{"/// <summary>This is a summary</summary>"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tokens, err := ParseCsharp(tc.input)
			if err != nil {
				t.Fatalf("Error parsing C#: %v", err)
			}

			// Find comments
			var foundComments []string
			for _, token := range tokens {
				if token.Type == TokenComment {
					foundComments = append(foundComments, token.Text)
				}
			}

			if len(foundComments) != len(tc.comments) {
				t.Errorf("Expected %d comments, found %d", len(tc.comments), len(foundComments))
				t.Logf("Expected: %v", tc.comments)
				t.Logf("Found: %v", foundComments)
			}

			// Check each expected comment is found
			for _, expectedComment := range tc.comments {
				found := false
				for _, foundComment := range foundComments {
					if foundComment == expectedComment {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected comment %q not found", expectedComment)
				}
			}
		})
	}
}

func TestCsharpIdentifiers(t *testing.T) {
	testCases := []struct {
		name        string
		input       string
		identifiers []string // Expected identifiers
	}{
		{
			name:        "Simple identifier",
			input:       "int myVariable = 5;",
			identifiers: []string{"myVariable"},
		},
		{
			name:        "Method call",
			input:       "Console.WriteLine(\"Hello\");",
			identifiers: []string{"Console", "WriteLine"},
		},
		{
			name:        "Class name",
			input:       "class MyClass { }",
			identifiers: []string{"MyClass"},
		},
		{
			name:        "Namespace",
			input:       "using System.Collections.Generic;",
			identifiers: []string{"System", "Collections", "Generic"},
		},
		{
			name:        "Escaped identifier",
			input:       "int @class = 5;",
			identifiers: []string{"@class"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tokens, err := ParseCsharp(tc.input)
			if err != nil {
				t.Fatalf("Error parsing C#: %v", err)
			}

			// Find identifiers
			var foundIdentifiers []string
			for _, token := range tokens {
				if token.Type == TokenIdentifier {
					foundIdentifiers = append(foundIdentifiers, token.Text)
				}
			}

			if len(foundIdentifiers) != len(tc.identifiers) {
				t.Errorf("Expected %d identifiers, found %d", len(tc.identifiers), len(foundIdentifiers))
				t.Logf("Expected: %v", tc.identifiers)
				t.Logf("Found: %v", foundIdentifiers)
			}

			// Check each expected identifier is found
			for _, expectedIdentifier := range tc.identifiers {
				found := false
				for _, foundIdentifier := range foundIdentifiers {
					if foundIdentifier == expectedIdentifier {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected identifier %q not found", expectedIdentifier)
				}
			}
		})
	}
}