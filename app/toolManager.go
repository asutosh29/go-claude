package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/openai/openai-go/v3"
)

type toolManager struct {
	tools []openai.ChatCompletionToolUnionParam
}

func NewToolManager() toolManager {
	return toolManager{
		tools: []openai.ChatCompletionToolUnionParam{},
	}
}

func (t *toolManager) AddTool(tool openai.ChatCompletionToolUnionParam) {
	t.tools = append(t.tools, tool)
}

func (t *toolManager) GetTools() []openai.ChatCompletionToolUnionParam {
	return t.tools
}

func (t *toolManager) ExecuteTool(toolCall openai.ChatCompletionMessageToolCallUnion) (string, error) {
	switch toolCall.Function.Name {
	case "Read":
		var toolCallArgs ReadToolArguments
		toolCallArgumentsJSON := toolCall.Function.Arguments
		if err := json.Unmarshal([]byte(toolCallArgumentsJSON), &toolCallArgs); err != nil {
			fmt.Fprintln(os.Stderr, err)
			panic("failed to unmarshall")

		}
		fileContent, err := ReadFile(toolCallArgs.FilePath)
		if err != nil {
			return "", err
		}
		return fileContent, nil
	case "Write":
		var toolCallArgs WriteToolArguments
		toolCallArgumentsJSON := toolCall.Function.Arguments
		if err := json.Unmarshal([]byte(toolCallArgumentsJSON), &toolCallArgs); err != nil {
			fmt.Fprintln(os.Stderr, err)
			panic("failed to unmarshall")
		}
		fileContent, err := WriteToFile(toolCallArgs.FilePath, toolCallArgs.Content)
		if err != nil {
			return "", err
		}
		return fileContent, nil

	default:
		return "", nil
	}
}
