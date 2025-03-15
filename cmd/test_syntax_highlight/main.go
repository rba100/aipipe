package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/rba100/aipipe/internal/display"
)

func main() {
	// Check if we should run a specific test
	if len(os.Args) > 1 {
		testType := strings.ToLower(os.Args[1])
		if testType == "typescript" || testType == "ts" {
			testTypeScriptHighlighting()
			return
		}
		// Default to Python test
	}

	// Run Python test by default
	testPythonHighlighting()

	// Run TypeScript test if no specific test was requested
	if len(os.Args) <= 1 {
		fmt.Println("\nRunning TypeScript test...")
		testTypeScriptHighlighting()
	}
}

func testPythonHighlighting() {
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

	fmt.Println("\nPython syntax highlighting test complete!")
}

func testTypeScriptHighlighting() {
	// Create a sample markdown text with TypeScript code
	markdownText := `# TypeScript Syntax Highlighting Test

This is a test of syntax highlighting for TypeScript code:

` + "```typescript" + `
// This is a TypeScript comment
/* This is a multi-line
   TypeScript comment */

// Interface definition
interface User {
    id: number;
    name: string;
    isActive: boolean;
}

// Class definition
class UserService {
    private users: User[] = [];

    constructor() {
        console.log("UserService initialized");
    }

    public addUser(user: User): void {
        this.users.push(user);
    }

    public getUser(id: number): User | undefined {
        return this.users.find(user => user.id === id);
    }
}

// Function with arrow syntax
const calculateTotal = (items: number[]): number => {
    return items.reduce((total, item) => total + item, 0);
};

// String concatenation
const name = "John";
const greeting = "Hello, " + name + "! Today is " + new Date().toLocaleDateString();

// Numeric literals
const decimal = 6;
const hex = 0xf00d;
const binary = 0b1010;
const octal = 0o744;

// Boolean and null
const isValid: boolean = true;
const isEmpty: boolean = false;
const nothing = null;
const notDefined = undefined;

// Async/await
async function fetchData(): Promise<User[]> {
    try {
        const response = await fetch('https://api.example.com/users');
        return await response.json();
    } catch (error) {
        console.error("Error fetching data:", error);
        return [];
    }
}
` + "```" + `

And here's some JavaScript code:

` + "```javascript" + `
// This is a JavaScript comment
function calculateSum(a, b) {
    return a + b;
}

// ES6 class
class Person {
    constructor(name, age) {
        this.name = name;
        this.age = age;
    }

    greet() {
        console.log("Hello, my name is " + this.name);
    }
}

// Arrow function
const multiply = (x, y) => x * y;

// Destructuring
const { name, age } = new Person("John", 30);

// Spread operator
const numbers = [1, 2, 3];
const moreNumbers = [...numbers, 4, 5];

// Promises
fetch('https://api.example.com/data')
    .then(response => response.json())
    .then(data => console.log(data))
    .catch(error => console.error(error));
` + "```" + `

And here's some text after the code blocks.
`

	// Create a pretty printer
	printer := display.NewPrettyPrinter()
	defer printer.Close()

	// Print the markdown text with syntax highlighting
	printer.Print(markdownText)
	printer.Flush()

	fmt.Println("\nTypeScript/JavaScript syntax highlighting test complete!")
}
