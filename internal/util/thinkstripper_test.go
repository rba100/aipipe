package util

import (
	"reflect"
	"testing"
	"time"
)

func TestStripThinkTagsStream(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "simple text without think tags",
			input:    []string{"Hello"},
			expected: []string{"Hello"},
		},
		{
			name:     "complete think block in chunks",
			input:    []string{"<think>", "some", "thinking", "</think>", "Hello"},
			expected: []string{"Hello"},
		},
		{
			name:     "incomplete think block",
			input:    []string{"<think>", "some", " thinking"},
			expected: []string{"<think>some thinking"},
		},
		{
			name:     "prefix text means no think block",
			input:    []string{"Hello", "<thin", "k>", "</think>World"},
			expected: []string{"Hello<thin", "k>", "</think>World"},
		},
		{
			name:     "complete think block in one chunk",
			input:    []string{"<think>hidden thought</think>visible", " output"},
			expected: []string{"visible", " output"},
		},
		{
			name:     "nested-looking think tags are not considered think tags",
			input:    []string{"<think>", "outer", "<think>", "inner", "</think>", "start here", "</think>", "more"},
			expected: []string{"start here", "</think>", "more"},
		},
		{
			name:     "think tag with empty content",
			input:    []string{"<think>", "</think>", "content"},
			expected: []string{"content"},
		},
		{
			name:     "split closing tag",
			input:    []string{"<think>", "some text", "</t", "hink>after"},
			expected: []string{"after"},
		},
		{
			name:     "split opening tag",
			input:    []string{"<th", "ink>some text", "</t", "hink>after"},
			expected: []string{"after"},
		},
		{
			name:     "first seven chars are buffered and not emitted immediately",
			input:    []string{"Hello", "World", "This", "Is", "A", "Test"},
			expected: []string{"HelloWorld", "This", "Is", "A", "Test"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputStream := make(chan string)
			resultStream := StripThinkTagsStream(inputStream)

			go func() {
				defer close(inputStream)
				for _, s := range tt.input {
					inputStream <- s
				}
			}()

			var result []string
			for s := range resultStream {
				result = append(result, s)
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("got %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestStripThinkTagsStream_emission_order(t *testing.T) {
	t.Run("single non-thinking part", func(t *testing.T) {

		inputStream := make(chan string)
		resultStream := StripThinkTagsStream(inputStream)

		inputStream <- "Hi!"
		close(inputStream)

		output := <-resultStream
		if output != "Hi!" {
			t.Errorf("expected Hi!, got %q", output)
		}
	})
	t.Run("multiple non-thinking parts", func(t *testing.T) {

		inputStream := make(chan string)
		resultStream := StripThinkTagsStream(inputStream)

		inputStream <- "Hello, World!"

		select {
		case output1 := <-resultStream:
			if output1 != "Hello, World!" {
				t.Errorf("expected Hello, World!, got %q", output1)
			}
		case <-time.After(500 * time.Millisecond):
			t.Fatal("timeout waiting for output1")
		}

		inputStream <- " My name is John Doe"

		select {
		case output2 := <-resultStream:
			if output2 != " My name is John Doe" {
				t.Errorf("expected  My name is John Doe., got %q", output2)
			}
		case <-time.After(500 * time.Millisecond):
			t.Fatal("timeout waiting for output2")
		}

		inputStream <- "!"

		select {
		case output3 := <-resultStream:
			if output3 != "!" {
				t.Errorf("expected !, got %q", output3)
			}
		case <-time.After(500 * time.Millisecond):
			t.Fatal("timeout waiting for output3")
		}

		close(inputStream)
	})
}
