package main

import (
	"os/exec"
	"strings"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/shared"
)

var BashToolParams = map[string]any{
	"type": "object",
	"properties": map[string]any{
		"command": map[string]any{
			"type":        "string",
			"description": "The command to execute",
		},
	},
	"required": []string{"command"},
}

type BashToolArguments struct {
	Command string `json:"command"`
}

var BashFileTool = openai.ChatCompletionFunctionTool(shared.FunctionDefinitionParam{
	Name:        "Bash",
	Description: openai.String("Execute a shell command"),
	Parameters:  BashToolParams,
})

func RunBashCommand(command string) (string, error) {
	parts := strings.Split(command, " ")
	var program string
	var args []string
	var cmd *exec.Cmd
	program = parts[0]
	if len(parts) == 1 {
		cmd = exec.Command(program)
	} else {
		args = parts[1:]
		cmd = exec.Command(program, args...)
	}
	// fmt.Println("Running: ", program, args)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}
