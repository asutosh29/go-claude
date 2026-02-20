package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/joho/godotenv"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

// Tool Definitions
func ReadFile(filePath string) {
	f, err := os.Open(filePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer f.Close()

	// Create a byte slice (buffer) of a specific chunk size (e.g., 5 bytes).
	buffer := make([]byte, 5)
	fmt.Fprintln(os.Stderr)
	for {
		// Read up to len(buffer) bytes from the file.
		n, err := f.Read(buffer)

		// Process the bytes that were read (from index 0 to n-1).
		if n > 0 {
			fmt.Printf("%s", string(buffer[:n]))
			// fmt.Printf("%d bytes: %s\n", n, string(buffer[:n]))
		}

		// Check for EOF (End of File) or any other errors.
		if err == io.EOF {
			break // Exit the loop when the end of the file is reached
		} else if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}

func main() {
	godotenv.Load()

	var prompt string
	flag.StringVar(&prompt, "p", "", "Prompt to send to LLM")
	flag.Parse()

	if prompt == "" {
		panic("Prompt must not be empty")
	}
	apiKey := os.Getenv("OPENROUTER_API_KEY")
	baseUrl := os.Getenv("OPENROUTER_BASE_URL")
	if baseUrl == "" {
		baseUrl = "https://openrouter.ai/api/v1"
	}

	if apiKey == "" {
		panic("Env variable OPENROUTER_API_KEY not found")
	}

	var modelName = ""
	if os.Getenv("local") == "true" {
		modelName = "moonshotai/kimi-k2-instruct-0905"
	} else {
		modelName = "anthropic/claude-haiku-4.5"
	}

	// Tool setup
	Tools := []openai.ChatCompletionToolUnionParam{
		ReadFileTool,
	}

	client := openai.NewClient(option.WithAPIKey(apiKey), option.WithBaseURL(baseUrl))
	resp, err := client.Chat.Completions.New(context.Background(),
		openai.ChatCompletionNewParams{
			Model: modelName,
			Messages: []openai.ChatCompletionMessageParamUnion{
				{
					OfUser: &openai.ChatCompletionUserMessageParam{
						Content: openai.ChatCompletionUserMessageParamContentUnion{
							OfString: openai.String(prompt),
						},
					},
				},
			},
			Tools: Tools,
		},
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	if len(resp.Choices) == 0 {
		panic("No choices in response")
	}

	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Fprintln(os.Stderr, "Logs from your program will appear here!")

	// TODO: Uncomment the line below to pass the first stage
	fmt.Print(resp.Choices[0].Message.Content)
	for _, Choice := range resp.Choices {
		toolCalls := Choice.Message.ToolCalls
		if len(toolCalls) != 0 {
			for _, toolCall := range toolCalls {

				switch toolCall.Function.Name {
				case "Read":
					var toolCallArgs ReadToolArguments
					toolCallArgumentsJSON := toolCall.Function.Arguments
					if err := json.Unmarshal([]byte(toolCallArgumentsJSON), &toolCallArgs); err != nil {
						fmt.Fprintln(os.Stderr, err)
						panic("failed to unmarshall")

					}
					ReadFile(toolCallArgs.FilePath)
				}

			}
		}
	}
}
