package util

import (
	"testing"
	"time"
)

func TestExtractCodeBlock(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expectedText string
		expectedType string
	}{
		{
			name:         "Simple code block",
			input:        "Here is some code:\n```\nfunc main() {\n\tfmt.Println(\"Hello, World!\")\n}\n```",
			expectedText: "func main() {\n\tfmt.Println(\"Hello, World!\")\n}",
			expectedType: "",
		},
		{
			name:         "Code block with language",
			input:        "Here is some Go code:\n```go\nfunc main() {\n\tfmt.Println(\"Hello, World!\")\n}\n```",
			expectedText: "func main() {\n\tfmt.Println(\"Hello, World!\")\n}",
			expectedType: "go",
		},
		{
			name:         "Multiple code blocks (extracts first)",
			input:        "First block:\n```\nfirst\n```\nSecond block:\n```\nsecond\n```",
			expectedText: "first",
			expectedType: "",
		},
		{
			name:         "No code block",
			input:        "This is just plain text with no code block.",
			expectedText: "This is just plain text with no code block.",
			expectedType: "",
		},
		{
			name:         "Empty code block",
			input:        "```\n```",
			expectedText: "",
			expectedType: "",
		},
		{
			name:         "Empty code block with preceding text",
			input:        "Here is an empty code block:\n```\n```",
			expectedText: "",
			expectedType: "",
		},
		{
			name:         "Empty code block with following text",
			input:        "```\n```\nThis follows an empty code block.",
			expectedText: "",
			expectedType: "",
		},
		{
			name:         "Empty code block with surrounding text",
			input:        "Before the block\n```\n```\nAfter the block",
			expectedText: "",
			expectedType: "",
		},
		{
			name:         "Empty code block without newline before closing",
			input:        "```\n```",
			expectedText: "",
			expectedType: "",
		},
		{
			name:         "Code block with language type",
			input:        "Here is some Python code:\n```python\nprint('Hello, World!')\n```",
			expectedText: "print('Hello, World!')",
			expectedType: "python",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractCodeBlock(tt.input)
			if result.Text != tt.expectedText {
				t.Errorf("ExtractCodeBlock().Text = %q, want %q", result.Text, tt.expectedText)
			}
			if result.Type != tt.expectedType {
				t.Errorf("ExtractCodeBlock().Type = %q, want %q", result.Type, tt.expectedType)
			}
		})
	}
}

func TestExtractCodeBlockStream(t *testing.T) {
	tests := []struct {
		name         string
		input        []string
		expectedText []string
		expectedType []string
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
			expectedText: []string{
				"func main() {\n\tfmt.Println(\"Hello, World!\")\n}",
			},
			expectedType: []string{
				"",
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
			expectedText: []string{
				"func main() {\n\tfmt.Println(\"Hello, World!\")\n}",
			},
			expectedType: []string{
				"go",
			},
		},
		{
			name: "No code block",
			input: []string{
				"This is just plain text ",
				"with no code block.",
			},
			expectedText: []string{
				"This is just plain text with no code block.",
			},
			expectedType: []string{
				"",
			},
		},
		{
			name: "Empty code block",
			input: []string{
				"```\n",
				"```",
			},
			expectedText: []string{
				"",
			},
			expectedType: []string{
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
			expectedText: []string{
				"",
			},
			expectedType: []string{
				"",
			},
		},
		{
			name: "Empty code block without newline before closing",
			input: []string{
				"```\n```",
			},
			expectedText: []string{
				"",
			},
			expectedType: []string{
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
			expectedText: []string{
				"def hello_world():\n    print('Hello, World!')\n",
			},
			expectedType: []string{
				"python",
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
			expectedText: []string{
				"",
			},
			expectedType: []string{
				"",
			},
		},
		{
			name: "Complex split with partial content",
			input: []string{
				"Let me show you\n`",    // First backtick
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
			expectedText: []string{
				"const greeting = 'Hello, World!';\nconsole.log(greeting);",
			},
			expectedType: []string{
				"js",
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
			var results []CodeBlockResult
			for part := range outputChan {
				results = append(results, part)
			}

			aggregatedText := ""
			for _, result := range results {
				aggregatedText += result.Text
			}

			if aggregatedText != tt.expectedText[0] {
				t.Errorf("ExtractCodeBlockStream() returned %q, want %q", aggregatedText, tt.expectedText[0])
			}
		})
	}
}
