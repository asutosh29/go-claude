package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

// Tool Definitions

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

	toolManager := NewToolManager()
	toolManager.AddTool(ReadFileTool)
	toolManager.AddTool(WriteFileTool)
	toolManager.AddTool(BashFileTool)

	// Messages Setup
	var messages []openai.ChatCompletionMessageParamUnion
	var systemPrompt = "You are an helpfull assistant. Strictly adhere to the instructions provided by the user. Be precise."
	messages = append(messages, openai.SystemMessage(systemPrompt))
	messages = append(messages, openai.UserMessage(prompt))

	// Tool setup
	Tools := toolManager.GetTools()

	aiCtx := context.Background()
	client := openai.NewClient(option.WithAPIKey(apiKey), option.WithBaseURL(baseUrl))
	// Make the Initial Request

	// Agent Loop
	for {
		// fmt.Println("=== LOGS ===")
		// for _, msg := range messages {
		// 	fmt.Printf("%+v\n", msg)
		// }
		resp, err := client.Chat.Completions.New(aiCtx,
			openai.ChatCompletionNewParams{
				Model:    modelName,
				Messages: messages,
				Tools:    Tools,
			},
		)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		if len(resp.Choices) == 0 {
			panic("No choices in response")
		}
		if resp.Choices[0].FinishReason == "stop" {
			fmt.Print(resp.Choices[0].Message.Content)
			break
		}

		// Manage Message history
		messages = append(messages, resp.Choices[0].Message.ToParam())

		// !!!! Do not show the thinking process for while submitting!!!
		// fmt.Fprintln(os.Stdout, resp.Choices[0].Message.Content)
		// fmt.Fprintln(os.Stdout)
		for _, Choice := range resp.Choices {
			toolCalls := Choice.Message.ToolCalls
			if len(toolCalls) != 0 {
				for _, toolCall := range toolCalls {
					toolCallId := toolCall.ID
					result, err := toolManager.ExecuteTool(toolCall)
					if err != nil {
						fmt.Fprintf(os.Stderr, "error: %v\n", err)
						os.Exit(1)
					}
					messages = append(messages, openai.ToolMessage(result, toolCallId))
				}
			}
		}
	}

}
