package main

import (
	"log"
	"os"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/shared"
)

var WriteToolParams = map[string]any{
	"type": "object",
	"properties": map[string]any{
		"file_path": map[string]any{
			"type":        "string",
			"description": "The path of the file to write to",
		},
		"content": map[string]any{
			"type":        "string",
			"description": "The content to write to the file",
		},
	},
	"required": []string{"file_path", "content"},
}

type WriteToolArguments struct {
	FilePath string `json:"file_path"`
	Content  string `json:"content"`
}

var WriteFileTool = openai.ChatCompletionFunctionTool(shared.FunctionDefinitionParam{
	Name:        "Write",
	Description: openai.String("Write content to a file"),
	Parameters:  WriteToolParams,
})

func WriteToFile(filePath string, content string) (string, error) {
	data := []byte(content)
	// The third argument, 0644, sets the file permissions (read/write for user, read-only for others).
	err := os.WriteFile(filePath, data, 0644)
	if err != nil {
		log.Fatal(err)
	}
	return content, nil
}
