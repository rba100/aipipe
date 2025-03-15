# aipipe - Golang CLI for LLM Interaction

A command-line tool for interacting with LLM providers, currently supporting Groq.

## Features

- Send prompts to LLM providers via command line or piped input
- Stream responses in real-time
- Extract code blocks from responses
- Pretty print responses with syntax highlighting

## Installation

```bash
# Clone the repository
git clone https://github.com/rba100/aipipe.git
cd aipipe

# Build the binary
go build -o aipipe ./cmd/aipipe
```

## Usage

```bash
# Set your API key
export GROQ_API_KEY=your-api-key

# Basic usage with a prompt
./aipipe "What is the capital of France?"

# Pipe content to the tool
cat file.txt | ./aipipe

# Stream the response
./aipipe --stream "Explain quantum computing"

# Extract code blocks
./aipipe --cb "Write a Python function to calculate Fibonacci numbers"

# Pretty print the response
./aipipe --pretty "Explain the benefits of Go over C++"
```

## Command Line Options

- `--cb`, `-c`: Extract code block from response
- `--stream`, `-s`: Stream completions from the AI model
- `--pretty`, `-p`: Enable pretty printing with colors and formatting

## Environment Variables

- `GROQ_API_KEY`: Your Groq API key
- `GROQ_ENDPOINT`: Groq API endpoint (defaults to "https://api.groq.com/openai/v1")

## License

MIT 