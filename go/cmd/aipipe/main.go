package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/rba100/aipipe/internal/display"
	"github.com/rba100/aipipe/internal/llm"
	"github.com/rba100/aipipe/internal/util"
)

func main() {
	// Define command line flags
	codeBlockFlag := flag.Bool("cb", false, "Extract code block from response")
	codeBlockFlagShort := flag.Bool("c", false, "Extract code block from response (shorthand)")
	streamFlag := flag.Bool("stream", false, "Stream completions from the AI model")
	streamFlagShort := flag.Bool("s", false, "Stream completions from the AI model (shorthand)")
	prettyFlag := flag.Bool("pretty", false, "Enable pretty printing with colors and formatting")
	prettyFlagShort := flag.Bool("p", false, "Enable pretty printing with colors and formatting (shorthand)")

	// Parse command line flags
	flag.Parse()

	// Combine short and long flags
	isCodeBlock := *codeBlockFlag || *codeBlockFlagShort
	isStream := *streamFlag || *streamFlagShort
	isPretty := *prettyFlag || *prettyFlagShort

	// Get prompt from command line arguments
	var argPrompt string
	if flag.NArg() > 0 {
		argPrompt = strings.Join(flag.Args(), " ")
	}

	// Run the AI query
	err := runAIQuery(isCodeBlock, isStream, isPretty, argPrompt)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func runAIQuery(isCodeBlock, isStream, isPretty bool, argPrompt string) error {
	// Check for mutually exclusive options
	if isCodeBlock && isPretty {
		return fmt.Errorf("the --cb and --pretty options cannot be used together")
	}

	// Create LLM client
	groqToken := os.Getenv("GROQ_API_KEY")
	if groqToken == "" {
		return fmt.Errorf("GROQ_API_KEY environment variable is not set")
	}

	groqEndpoint := os.Getenv("GROQ_ENDPOINT")
	if groqEndpoint == "" {
		groqEndpoint = "https://api.groq.com/openai/v1"
	}

	config := &llm.Config{
		APIEndpoint:  groqEndpoint,
		APIToken:     groqToken,
		IsCodeBlock:  isCodeBlock,
		IsStream:     isStream,
		ModelType:    llm.ModelTypeDefault,
		DefaultModel: "llama-3.3-70b-versatile",
		FastModel:    "llama-3.1-8b-instant",
	}

	client, err := llm.NewClient(config)
	if err != nil {
		return err
	}

	// Build prompt from stdin and/or command line argument
	promptBuilder := strings.Builder{}

	// Check if there's input from stdin
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		// Read from stdin
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			promptBuilder.WriteString(scanner.Text())
			promptBuilder.WriteString("\n")
		}

		if err := scanner.Err(); err != nil {
			return fmt.Errorf("error reading from stdin: %v", err)
		}
	}

	// Add command line argument if provided
	if argPrompt != "" {
		if promptBuilder.Len() > 0 {
			promptBuilder.WriteString("-----\n")
		}
		promptBuilder.WriteString(argPrompt)
	}

	// Check if we have any input
	if promptBuilder.Len() == 0 {
		return fmt.Errorf("no input provided")
	}

	prompt := promptBuilder.String()

	// Process the prompt with the LLM
	if isStream {
		var stream <-chan string
		if isCodeBlock {
			stream = util.ExtractCodeBlockStream(client.CreateCompletionStream(prompt))
		} else {
			stream = client.CreateCompletionStream(prompt)
		}

		if isPretty {
			printer := display.NewPrettyPrinter()
			defer printer.Close()

			for part := range stream {
				printer.Print(part)
			}

			// Make sure to flush any remaining content before closing
			printer.Flush()
		} else {
			for part := range stream {
				fmt.Print(part)
			}
			// Add a newline if the last part doesn't end with one
			fmt.Println()
		}
	} else {
		var response string
		var err error

		response, err = client.CreateCompletion(prompt)
		if err != nil {
			return err
		}

		if isCodeBlock {
			response = util.ExtractCodeBlock(response)
		}

		if isPretty {
			printer := display.NewPrettyPrinter()
			defer printer.Close()
			printer.Print(response)
			printer.Flush()
		} else {
			fmt.Println(response)
		}
	}

	return nil
}
