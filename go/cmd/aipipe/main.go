package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/rba100/aipipe/internal/display"
	"github.com/rba100/aipipe/internal/llm"
	"github.com/rba100/aipipe/internal/util"
	"github.com/spf13/pflag"
)

func main() {
	// Define command line flags
	codeBlockFlag := pflag.Bool("cb", false, "Extract code block from response")
	codeBlockFlagShort := pflag.BoolP("c", "c", false, "Extract code block from response (shorthand)")
	streamFlag := pflag.Bool("stream", false, "Stream completions from the AI model")
	streamFlagShort := pflag.BoolP("s", "s", false, "Stream completions from the AI model (shorthand)")
	prettyFlag := pflag.Bool("pretty", false, "Enable pretty printing with colors and formatting")
	prettyFlagShort := pflag.BoolP("p", "p", false, "Enable pretty printing with colors and formatting (shorthand)")
	reasoningFlag := pflag.Bool("reasoning", false, "Use reasoning model")
	reasoningFlagShort := pflag.BoolP("r", "r", false, "Use reasoning model (shorthand)")

	// Parse command line flags - pflag allows flags to be placed anywhere
	pflag.Parse()

	// Combine short and long flags
	isCodeBlock := *codeBlockFlag || *codeBlockFlagShort
	isStream := *streamFlag || *streamFlagShort
	isPretty := *prettyFlag || *prettyFlagShort
	isReasoning := *reasoningFlag || *reasoningFlagShort

	// Get prompt from command line arguments
	var argPrompt string
	if pflag.NArg() > 0 {
		argPrompt = strings.Join(pflag.Args(), " ")
	}

	// Run the AI query
	err := runAIQuery(isCodeBlock, isStream, isPretty, isReasoning, argPrompt)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func runAIQuery(isCodeBlock, isStream, isPretty, isReasoning bool, argPrompt string) error {
	// Check for mutually exclusive options
	if isCodeBlock && isPretty {
		return fmt.Errorf("the --cb and --pretty options cannot be used together")
	}

	// Get API configuration from environment variables
	apiConfig, err := util.GetAPIConfig()
	if err != nil {
		return err
	}

	model := llm.ModelTypeDefault
	if isReasoning {
		model = llm.ModelTypeReasoning
	}

	// Create LLM client
	config := &llm.Config{
		APIEndpoint:    apiConfig.APIEndpoint,
		APIToken:       apiConfig.APIToken,
		IsCodeBlock:    isCodeBlock,
		IsStream:       isStream,
		ModelType:      model,
		DefaultModel:   "llama-3.3-70b-versatile",
		FastModel:      "llama-3.1-8b-instant",
		ReasoningModel: "qwen-2.5-32b",
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
