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
	codeBlockFlag := pflag.BoolP("codeblock", "c", false, "Extract code block from response")
	streamFlag := pflag.BoolP("stream", "s", false, "Stream completions from the AI model")
	prettyFlag := pflag.BoolP("pretty", "p", false, "Enable pretty printing with colors and formatting")
	reasoningFlag := pflag.BoolP("reasoning", "r", false, "Use reasoning model")
	fastFlag := pflag.BoolP("fast", "f", false, "Use fast model")
	thinkingFlag := pflag.BoolP("thinking", "t", false, "Show thinking process")

	// Parse command line flags - pflag allows flags to be placed anywhere
	pflag.Parse()

	// Combine short and long flags
	isCodeBlock := *codeBlockFlag
	isStream := *streamFlag
	isPretty := *prettyFlag
	isReasoning := *reasoningFlag
	isFast := *fastFlag
	showThinking := *thinkingFlag
	// Get prompt from command line arguments
	var argPrompt string
	if pflag.NArg() > 0 {
		argPrompt = strings.Join(pflag.Args(), " ")
	}

	// Run the AI query
	err := runAIQuery(isCodeBlock, isStream, isPretty, isReasoning, isFast, showThinking, argPrompt)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func runAIQuery(isCodeBlock, isStream, isPretty, isReasoning, isFast, showThinking bool, argPrompt string) error {
	// Check for mutually exclusive options

	if isReasoning && isFast {
		return fmt.Errorf("the --reasoning and --fast options cannot be used together")
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
	if isFast {
		model = llm.ModelTypeFast
	}

	// Create LLM client
	config := &llm.Config{
		APIEndpoint:    apiConfig.APIEndpoint,
		APIToken:       apiConfig.APIToken,
		IsCodeBlock:    isCodeBlock,
		IsStream:       isStream,
		ModelType:      model,
		DefaultModel:   apiConfig.DefaultModel,
		FastModel:      apiConfig.FastModel,
		ReasoningModel: apiConfig.ReasoningModel,
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
		stream := client.CreateCompletionStream(prompt)
		if !showThinking {
			stream = util.StripThinkTagsStream(stream)
		}
		if isCodeBlock {
			codeBlockStream := util.ExtractCodeBlockStream(stream)

			if isPretty {
				printer := display.NewPrettyPrinter()
				defer printer.Close()

				for result := range codeBlockStream {
					if result.Type != "" {
						printer.SetCodeBlockState(result.Type)
					}
					printer.Print(result.Text)
				}

				// Make sure to flush any remaining content before closing
				printer.Flush()
			} else {
				for result := range codeBlockStream {
					fmt.Print(result.Text)
				}
				// Add a newline if the last part doesn't end with one
				fmt.Println()
			}
		} else {
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
		}
	} else {
		var response string
		var err error

		response, err = client.CreateCompletion(prompt)
		if err != nil {
			return err
		}

		if !showThinking {
			response = util.StripThinkTags(response)
		}

		if isCodeBlock {
			result := util.ExtractCodeBlock(response)

			if isPretty {
				printer := display.NewPrettyPrinter()
				defer printer.Close()
				if result.Type != "" {
					printer.SetCodeBlockState(result.Type)
				}
				printer.Print(result.Text)
				printer.Flush()
			} else {
				fmt.Println(result.Text)
			}
		} else {
			if isPretty {
				printer := display.NewPrettyPrinter()
				defer printer.Close()
				printer.Print(response)
				printer.Flush()
			} else {
				fmt.Println(response)
			}
		}
	}

	return nil
}
