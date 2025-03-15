package main

import (
	"fmt"

	"github.com/rba100/aipipe/internal/display"
)

func main() {
	// Create a sample markdown text with Python code
	markdownText := `# Syntax Highlighting Test

This is a test of syntax highlighting for Python code:

` + "```python" + `
# This is a Python comment
def hello_world():
    """This is a docstring"""
    print("Hello, world!")
    return 42

class MyClass:
    def __init__(self, value):
        self.value = value
        
    def get_value(self):
        return self.value * 2

# Test with some literals
x = 123
y = 3.14
z = "This is a string"
flag = True
none_val = None

if x > 100 and y < 4:
    print(f"x: {x}, y: {y}")
` + "```" + `

And here's some text after the code block.
`

	// Create a pretty printer
	printer := display.NewPrettyPrinter()
	defer printer.Close()

	// Print the markdown text with syntax highlighting
	printer.Print(markdownText)
	printer.Flush()

	fmt.Println("\nSyntax highlighting test complete!")
}
