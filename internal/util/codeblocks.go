package util

import (
	"regexp"
	"strings"
)

// ExtractCodeBlock extracts a code block from a string
func ExtractCodeBlock(input string) string {
	// Use a regex pattern that can handle empty code blocks
	re := regexp.MustCompile("```[a-zA-Z0-9.]*\n([\\s\\S]*?)(?:\n```|```)")
	matches := re.FindStringSubmatch(input)
	if len(matches) > 1 {
		return matches[1]
	}
	return input
}

// CodeBlockState represents the state of code block extraction
type CodeBlockState int

const (
	SearchingOpening CodeBlockState = iota
	Open
	Closed
)

// ExtractCodeBlockStream extracts code blocks from a stream
func ExtractCodeBlockStream(inputStream <-chan string) <-chan string {
	outputStream := make(chan string)

	go func() {
		defer close(outputStream)

		// Collect all parts into a single string
		completeBuffer := strings.Builder{}
		for part := range inputStream {
			completeBuffer.WriteString(part)
		}

		// Use the same regex as ExtractCodeBlock
		completeStr := completeBuffer.String()
		re := regexp.MustCompile("```[a-zA-Z0-9.]*\n([\\s\\S]*?)(?:\n```|```)")
		matches := re.FindStringSubmatch(completeStr)

		if len(matches) > 1 {
			// Found a code block, return the content
			outputStream <- matches[1]
		} else {
			// No code block found, return the original string
			outputStream <- completeStr
		}
	}()

	return outputStream
}
