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
		if testType == "csharp" || testType == "cs" || testType == "c#" {
			testCsharpHighlighting()
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
		fmt.Println("\nRunning C# test...")
		testCsharpHighlighting()
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

func testCsharpHighlighting() {
	// Create a sample markdown text with C# code
	markdownText := `# C# Syntax Highlighting Test

This is a test of syntax highlighting for C# code:

` + "```csharp" + `
// This is a C# comment
/* This is a multi-line
   C# comment */

// Using statements
using System;
using System.Collections.Generic;
using System.Linq;

// Namespace declaration
namespace MyNamespace
{
    // Class definition
    public class Person
    {
        // Properties
        public string Name { get; set; }
        public int Age { get; private set; }
        
        // Constructor
        public Person(string name, int age)
        {
            Name = name;
            Age = age;
        }
        
        // Method
        public virtual void Greet()
        {
            Console.WriteLine($"Hello, my name is {Name} and I'm {Age} years old.");
        }
    }
    
    // Interface
    public interface IRepository<T> where T : class
    {
        Task<T> GetByIdAsync(int id);
        Task<IEnumerable<T>> GetAllAsync();
        Task SaveAsync(T entity);
    }
    
    // Static class
    public static class Utils
    {
        // Static method with different number formats
        public static void DisplayNumbers()
        {
            int decimal = 42;
            float floatVal = 3.14f;
            double doubleVal = 3.14159;
            decimal decimalVal = 123.45m;
            long longVal = 1234567890L;
            
            // Hexadecimal and binary
            int hex = 0xFF;
            int binary = 0b1010;
            
            // String literals
            string regular = "Hello, World!";
            string verbatim = @"C:\Users\John\Documents";
            char character = 'A';
            
            // Boolean and null
            bool isValid = true;
            bool isEmpty = false;
            string nullValue = null;
            
            Console.WriteLine($"Decimal: {decimal}, Float: {floatVal}, Double: {doubleVal}");
        }
    }
    
    // Generic class
    public class Repository<T> : IRepository<T> where T : class
    {
        private readonly List<T> _items = new List<T>();
        
        public async Task<T> GetByIdAsync(int id)
        {
            // Simulated async operation
            await Task.Delay(100);
            return _items.FirstOrDefault();
        }
        
        public async Task<IEnumerable<T>> GetAllAsync()
        {
            await Task.Delay(50);
            return _items;
        }
        
        public async Task SaveAsync(T entity)
        {
            if (entity == null)
                throw new ArgumentNullException(nameof(entity));
                
            _items.Add(entity);
            await Task.Delay(10);
        }
    }
    
    // Main program
    public class Program
    {
        public static async Task Main(string[] args)
        {
            var person = new Person("John Doe", 30);
            person.Greet();
            
            // Exception handling
            try
            {
                var repository = new Repository<Person>();
                await repository.SaveAsync(person);
                
                var allPeople = await repository.GetAllAsync();
                foreach (var p in allPeople)
                {
                    p.Greet();
                }
            }
            catch (Exception ex)
            {
                Console.WriteLine($"Error: {ex.Message}");
            }
            finally
            {
                Console.WriteLine("Cleanup completed.");
            }
        }
    }
}
` + "```" + `

And here's some text after the code block.
`

	// Create a pretty printer
	printer := display.NewPrettyPrinter()
	defer printer.Close()

	// Print the markdown text with syntax highlighting
	printer.Print(markdownText)
	printer.Flush()

	fmt.Println("\nC# syntax highlighting test complete!")
}
