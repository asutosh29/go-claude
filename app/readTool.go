package main

import (
	"fmt"
	"os"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/shared"
)

var ReadToolParams = map[string]any{
	"type": "object",
	"properties": map[string]any{
		"file_path": map[string]any{
			"type":        "string",
			"description": "The path to the file to read",
		},
	},
	"required": []string{"file_path"},
}

type ReadToolArguments struct {
	FilePath string `json:"file_path"`
}

var ReadFileTool = openai.ChatCompletionFunctionTool(shared.FunctionDefinitionParam{
	Name:        "Read",
	Description: openai.String("Read and return the contents of a file"),
	Parameters:  ReadToolParams,
})

func ReadFile(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("File reading error", err)
		return "", err
	}
	fileContent := string(content)
	// fmt.Println(fileContent)
	return fileContent, nil
}
