# GoClaude

A powerful AI agent framework written in Go that integrates Claude AI with system tools. Built as part of the Codecrafters challenge, this agent can read files, write files, and execute bash commands on your system.

## Features

- **Claude AI Integration**: Uses Claude Haiku 4.5 or Kimi K2 (local mode) via OpenRouter API
- **Tool Use**: Agentic tool-calling capabilities with three built-in tools:
  - **Read Tool**: Read file contents
  - **Write Tool**: Create or overwrite files
  - **Bash Tool**: Execute shell commands
- **Agentic Loop**: Implements an agent loop that iteratively calls Claude and executes tools until completion
- **Flexible Configuration**: Support for custom API endpoints and local/remote model selection

## Prerequisites

- Go 1.25 or higher
- An OpenRouter API key (for Claude access)
- Environment variables configured

## Installation

1. Clone the repository:

```bash
git clone https://github.com/asutosh29/go-claude
cd go-claude
```

2. Install dependencies:

```bash
go mod download
```

3. Set up environment variables:
   Create a `.env` file in the root directory:

```
cp .env.sample .env
```

## Usage

Run the agent with a prompt:

```bash
go run ./app -p "Your prompt here"
```

Example:

```bash
go run ./app -p "Create a file named test.txt with the content 'Hello, World!'"
```

Or use the provided shell script:

```bash
./your_program.sh "Your prompt here"
```

## Project Structure

```
.
├── app/
│   ├── main.go          # Main agent loop and CLI entry point
│   ├── toolManager.go   # Tool registry and executor
│   ├── readTool.go      # File reading tool implementation
│   ├── writeTool.go     # File writing tool implementation
│   ├── bash.go          # Bash command execution tool
├── codecrafters.yml     # Codecrafters configuration
├── go.mod              # Go module definition
└── README.md           # This file
```

## How It Works

1. **Initialization**: The agent loads your API key and initializes Claude via OpenRouter
2. **Tool Registration**: Three tools (Read, Write, Bash) are registered with the agent
3. **Agent Loop**:
   - Sends your prompt and available tools to Claude
   - Claude decides which tools to use based on your request
   - The agent executes the requested tools
   - Results are sent back to Claude for further processing
   - Loop continues until Claude returns a final response (finish_reason = "stop")
4. **Output**: Final response is printed to stdout

## Dependencies

- `github.com/openai/openai-go/v3` - OpenRouter/OpenAI API client
- `github.com/joho/godotenv` - Environment variable loader

## Configuration

### API Selection

- **Remote (Default)**: Uses `anthropic/claude-haiku-4.5` via OpenRouter
- **Local**: Set `local=true` to use `moonshotai/kimi-k2-instruct-0905`

## Examples

### Create a file

```bash
go run ./app -p "Create a file called hello.go with a simple Hello World program"
```

### Read and analyze a file

```bash
go run ./app -p "Read main.go and explain what the agent loop does"
```

### Execute system commands

```bash
go run ./app -p "Execute the command 'ls -la' and tell me what files are in the current directory"
```
