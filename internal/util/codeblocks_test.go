package util

import (
	"testing"
	"time"
)

func TestExtractCodeBlock(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Simple code block",
			input:    "Here is some code:\n```\nfunc main() {\n\tfmt.Println(\"Hello, World!\")\n}\n```",
			expected: "func main() {\n\tfmt.Println(\"Hello, World!\")\n}",
		},
		{
			name:     "Code block with language",
			input:    "Here is some Go code:\n```go\nfunc main() {\n\tfmt.Println(\"Hello, World!\")\n}\n```",
			expected: "func main() {\n\tfmt.Println(\"Hello, World!\")\n}",
		},
		{
			name:     "Multiple code blocks (extracts first)",
			input:    "First block:\n```\nfirst\n```\nSecond block:\n```\nsecond\n```",
			expected: "first",
		},
		{
			name:     "No code block",
			input:    "This is just plain text with no code block.",
			expected: "This is just plain text with no code block.",
		},
		{
			name:     "Empty code block",
			input:    "```\n```",
			expected: "",
		},
		{
			name:     "Empty code block with preceding text",
			input:    "Here is an empty code block:\n```\n```",
			expected: "",
		},
		{
			name:     "Empty code block with following text",
			input:    "```\n```\nThis follows an empty code block.",
			expected: "",
		},
		{
			name:     "Empty code block with surrounding text",
			input:    "Before the block\n```\n```\nAfter the block",
			expected: "",
		},
		{
			name:     "Empty code block without newline before closing",
			input:    "```\n```",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractCodeBlock(tt.input)
			if result != tt.expected {
				t.Errorf("ExtractCodeBlock() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestExtractCodeBlockStream(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name: "Simple streaming code block",
			input: []string{
				"Here is some ",
				"code:\n```\n",
				"func main() {\n",
				"\tfmt.Println(\"Hello, World!\")\n",
				"}\n```",
			},
			expected: []string{
				"func main() {\n\tfmt.Println(\"Hello, World!\")\n}",
			},
		},
		{
			name: "Code block with language",
			input: []string{
				"Here is some Go code:\n```go\n",
				"func main() {\n",
				"\tfmt.Println(\"Hello, World!\")\n",
				"}\n```",
			},
			expected: []string{
				"func main() {\n\tfmt.Println(\"Hello, World!\")\n}",
			},
		},
		{
			name: "No code block",
			input: []string{
				"This is just plain text ",
				"with no code block.",
			},
			expected: []string{
				"This is just plain text with no code block.",
			},
		},
		{
			name: "Empty code block",
			input: []string{
				"```\n",
				"```",
			},
			expected: []string{
				"",
			},
		},
		{
			name: "Empty code block with preceding text",
			input: []string{
				"Here is an empty ",
				"code block:\n```\n",
				"```",
			},
			expected: []string{
				"",
			},
		},
		{
			name: "Empty code block without newline before closing",
			input: []string{
				"```\n```",
			},
			expected: []string{
				"",
			},
		},
		{
			name: "Split backticks across chunks",
			input: []string{
				"Here is some Python code: `",
				"``",
				"python\n",
				"def hello_world():\n",
				"    print('Hello, World!')\n",
				"\n`",
				"``",
			},
			expected: []string{
				"def hello_world():\n    print('Hello, World!')\n",
			},
		},
		{
			name: "Split backticks with empty code block",
			input: []string{
				"`",
				"``\n",
				"`",
				"``",
			},
			expected: []string{
				"",
			},
		},
		{
			name: "Complex split with partial content",
			input: []string{
				"Let me show you `",     // First backtick
				"`",                     // Second backtick
				"`",                     // Third backtick
				"js\nconst ",            // Language identifier and start of content
				"greeting = '",          // Middle of content
				"Hello",                 // More content
				", World!';",            // More content
				"\nconsole.log",         // More content
				"(greeting",             // More content
				");\n`",                 // End of content and first closing backtick
				"`",                     // Second closing backtick
				"`",                     // Third closing backtick
				" That was JavaScript!", // Text after code block
			},
			expected: []string{
				"const greeting = 'Hello, World!';\nconsole.log(greeting);",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create input channel
			inputChan := make(chan string)
			go func() {
				defer close(inputChan)
				for _, part := range tt.input {
					inputChan <- part
					// Add a small delay to simulate streaming
					time.Sleep(10 * time.Millisecond)
				}
			}()

			// Get output channel
			outputChan := ExtractCodeBlockStream(inputChan)

			// Collect results
			var results []string
			for part := range outputChan {
				results = append(results, part)
			}

			// Check results
			if len(results) != len(tt.expected) {
				t.Errorf("ExtractCodeBlockStream() returned %d parts, expected %d", len(results), len(tt.expected))
				return
			}

			for i, result := range results {
				if result != tt.expected[i] {
					t.Errorf("ExtractCodeBlockStream() part %d = %q, want %q", i, result, tt.expected[i])
				}
			}
		})
	}
}
