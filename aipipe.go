package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"

	openai "github.com/sashabaranov/go-openai"
)

var (
	GROQ_API_KEY  = os.Getenv("GROQ_API_KEY")
	GROQ_ENDPOINT = os.Getenv("GROQ_ENDPOINT")
	GROQ_MODEL    = os.Getenv("GROQ_MODEL")

	systemMessage = "You are a helpful assistant. If the user has asked for something written, put it in a code block (```), otherwise just provide the answer. If you do use a codeblock, all other text is ignored."
)

func getGroqCompletion(client *openai.Client, userMessage, groqModel string) (string, error) {
	req := openai.ChatCompletionRequest{
		Model: groqModel,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: systemMessage,
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: userMessage,
			},
		},
	}
	resp, err := client.CreateChatCompletion(context.Background(), req)
	if err != nil {
		return "", err
	}
	return resp.Choices[0].Message.Content, nil
}

func extractCodeBlock(completion string) string {
	println("Extracting code block")
	codeBlockPattern := "```[a-zA-Z0-9.]*\\n([\\s\\S]+?)\\n```"
	re := regexp.MustCompile(codeBlockPattern)
	match := re.FindStringSubmatch(completion)
	if match != nil {
		return match[1]
	}
	return completion
}

func isPipedInput() bool {
	fi, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}
	return fi.Mode()&os.ModeNamedPipe != 0
}

func main() {
	var codeBlock = flag.Bool("cb", false, "Return only the code block in the completion.")
	flag.Parse()

	flag.Usage = func() {
		fmt.Println("Usage: go run aipipe.go \"query\" > output.txt")
	}

	var prompt string
	if isPipedInput() {
		input, err := io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading stdin:", err)
			os.Exit(1)
		}
		prompt = string(input)
	} else {
		args := flag.Args()
		if len(args) < 1 {
			flag.Usage()
			return
		}
		prompt = args[0]
	}

	var client *openai.Client

	config := openai.DefaultConfig(GROQ_API_KEY)
	config.BaseURL = GROQ_ENDPOINT
	client = openai.NewClientWithConfig(config)

	var completion string
	var err error
	completion, err = getGroqCompletion(client, prompt, GROQ_MODEL)

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	var output string
	if *codeBlock {
		output = extractCodeBlock(completion)
	} else {
		output = completion
	}

	fmt.Print(output)
}
