# linear-cli

A smart CLI tool for creating Linear issues using AI. Simply describe your task in natural language, and let Gemini AI generate a properly formatted issue.

## Features

- 🤖 **AI-powered**: Uses Google Gemini to generate well-structured issue titles and descriptions
- ⚡ **Quick creation**: Default action makes it fast to create issues without subcommands
- 🇳🇴 **Norwegian support**: Works seamlessly with Norwegian language input
- 🔧 **Manual mode**: Option to create issues directly without AI

## Prerequisites

- Go 1.25.2 or later
- A [Linear API key](https://linear.app/settings/api)
- A [Google API key](https://ai.google.dev/gemini-api/docs/api-key) for Gemini
- Your Linear Team ID

## Installation

```bash
go install github.com/filipjaj/linear-cli
```

## Configuration

Create a `.env` file in your working directory or set environment variables:

```bash
LINEAR_API_KEY="lin_api_..."
GOOGLE_API_KEY="your-google-api-key"
LINEAR_TEAM_ID="your-team-id"
```

You can copy `.env.example` to get started:

```bash
cp .env.example .env
```

## Usage

### Default AI Mode (Recommended)

Simply type `linear-cli` followed by your task description:

```bash
linear-cli kjøpe kake til kontoret
linear-cli husk å kalle inn kandidater til intervju
linear-cli fix the bug in the login flow
```

All arguments are combined and processed through AI to generate a proper Linear issue.

### Explicit AI Command

Same as default, but using the `ai` subcommand:

```bash
linear-cli ai legge til dark mode support
```

### Manual Creation

Create an issue directly without AI processing:

```bash
linear-cli create "Issue Title"
```

## Commands

| Command | Description |
|---------|-------------|
| `linear-cli <description...>` | **Default**: Create issue using AI (all args combined) |
| `linear-cli ai <description...>` | Create issue using AI (explicit subcommand) |
| `linear-cli create <title>` | Create issue manually without AI |
| `linear-cli --version` | Show version information |
| `linear-cli --help` | Show help |

## How It Works

1. You provide a task description in natural language
2. Gemini AI processes it and generates:
   - A concise, actionable title
   - A clear description with necessary details
3. The issue is automatically created in your Linear team
4. You get instant confirmation ✓

## Development

```bash
# Clone the repository
git clone https://github.com/filipjaj/linear-cli.git
cd linear-cli

# Install dependencies
go mod download

# Run locally
go run .

# Build
go build -o linear-cli

# Install locally
go install
```

## License

MIT
