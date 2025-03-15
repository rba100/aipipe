# AIPipe

A command-line tool for transforming text or otherwise making adhoc LLM calls.

## Usage

Use the pipe operator (`|`) to send text through the `aipipe` command:

```bash
some-app | aipipe "instruction"
```

...or just call it by itself for shell based one-off llm calls.

### Example

Simple reformatting

Input:
```bash
echo "Robin Anderson 1 High Street CB1 1AA" | aipipe "format as JSON" --cb
```

Output:
```json
{
    "name": "Robin Anderson",
    "address": {
        "street": "1 High Street",
        "postcode": "CB1 1AA"
    }
}
```

The 'code block' flag `--cb` is best for when you want something specifically formatted, rather than just for you to read yourself. Without it the LLM might write "Sure, here's your thing..." which you might not want to pipe into another application.

### Options

- `-c / --cb`: outputs only the first code block emitted by the LLM, discarding all other output. Otherwise all output is emitted to std out.
- `-p / --pretty`: use console colours to highlight markdown.
- `-s / --stream`: stream the output for faster perceived response
- `-r / --reasoning`: use a reasoning model instead, for extra oomph.
- `-f / --fast`: use a fast-but-thick model instead, for extra speed.
- `-m / --mic`: microphone input for the instruction prompt (probably Windows only)

## Installation

`dotnet publish aipipe.csproj -c Release -o bld --self-contained true -p:PublishSingleFile=true -p:PublishTrimmed=true -p:DebugType=None`

copy the .\bld\aipipe(.exe) to your bin folder.

Set env vars
```
GROQ_API_KEY
```

## AI Providers

Supports https://groq.com/ using the OpenAI client interface with a custom base URL.

# LLM Client Implementation

This project provides a Go implementation for interacting with LLM providers using the OpenAI client interface.

## Structure

- `go/internal/llm/openaiclient.go`: Contains the complete LLM client implementation
- `go/internal/util/codeblocks.go`: Contains utilities for extracting code blocks from LLM responses
- `go/cmd/main.go`: Example usage of the LLM client

## Implementation Details

The implementation uses the OpenAI client interface to interact with Groq's API by overriding the base URL. This approach allows us to use the same client for different LLM providers that support the OpenAI API format.

Key features:
- Configurable API endpoint and token
- Support for different model types (fast, default, reasoning)
- Streaming and non-streaming completions
- Code block extraction utilities

## Usage

```go
import (
    "../internal/llm"
    "../internal/util"
)

func main() {
    // Create a configuration for using Groq with the OpenAI client
    config := &llm.Config{
        APIEndpoint:    "https://api.groq.com/v1", // Groq API endpoint
        APIToken:       "your-api-key",
        DefaultModel:   "llama3-70b-8192",         // Default Groq model
        FastModel:      "gemma-7b-it",             // Fast Groq model
        ReasoningModel: "mixtral-8x7b-32768",      // Reasoning Groq model
        IsCodeBlock:    true,
        IsStream:       false,
        ModelType:      llm.ModelTypeDefault,
    }

    // Create a client
    client, err := llm.NewClient(config)
    if err != nil {
        panic(err)
    }

    // Get a completion
    response, err := client.CreateCompletion("Write a hello world program in Go")
    if err != nil {
        panic(err)
    }

    // Extract code block if needed
    codeBlock := util.ExtractCodeBlock(response)
    
    // Or use streaming
    if config.IsStream {
        stream := client.CreateCompletionStream("Write a hello world program in Go")
        for chunk := range stream {
            // Process each chunk
            fmt.Print(chunk)
        }
    }
}
```

## Notes

- The implementation uses standard Go HTTP client to make API requests to Groq's API.
- The client is designed to work with any API that follows the OpenAI API format.
- The code block extraction utilities are provided to help extract code blocks from LLM responses.