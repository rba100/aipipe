package util

import (
	"regexp"
	"strings"
)

// CodeBlockResult represents the result of extracting a code block
type CodeBlockResult struct {
	Text string
	Type string
}

// ExtractCodeBlock extracts a code block from a string
func ExtractCodeBlock(input string) CodeBlockResult {
	// Use a regex pattern that can handle empty code blocks and capture the language type
	re := regexp.MustCompile("```([a-zA-Z0-9.]*)(?:\n)?([\\s\\S]*?)(?:\n```|```)")
	matches := re.FindStringSubmatch(input)
	if len(matches) > 2 {
		return CodeBlockResult{
			Text: matches[2],
			Type: matches[1],
		}
	}
	return CodeBlockResult{
		Text: input,
		Type: "",
	}
}

// CodeBlockState represents the state of code block extraction
type CodeBlockState int

const (
	SearchingOpening CodeBlockState = iota
	Open
	Closed
)

// ExtractCodeBlockStream extracts code blocks from a stream
func ExtractCodeBlockStream(inputStream <-chan string) <-chan CodeBlockResult {
	outputStream := make(chan CodeBlockResult)

	go func() {
		defer close(outputStream)

		// Collect all parts into a single string
		completeBuffer := strings.Builder{}
		for part := range inputStream {
			completeBuffer.WriteString(part)
		}

		// Use the same regex as ExtractCodeBlock
		completeStr := completeBuffer.String()
		re := regexp.MustCompile("```([a-zA-Z0-9.]*)(?:\n)?([\\s\\S]*?)(?:\n```|```)")
		matches := re.FindStringSubmatch(completeStr)

		if len(matches) > 2 {
			// Found a code block, return the content and type
			outputStream <- CodeBlockResult{
				Text: matches[2],
				Type: matches[1],
			}
		} else {
			// No code block found, return the original string
			outputStream <- CodeBlockResult{
				Text: completeStr,
				Type: "",
			}
		}
	}()

	return outputStream
}
