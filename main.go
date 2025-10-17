package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/filipjaj/linear-ai/internal/linear"
	"github.com/urfave/cli/v3"
)

var (
	// Set these at build time: -ldflags "-X main.version=1.0.0 -X main.commit=abc123"
	version = "dev"
	commit  = "unknown"
)


const master_prompt = `Du er en AI-assistent som hjelper med å opprette oppgaver i Linear. Basert på beskrivelsen av oppgaven skal du lage en oppgave i Linear, med en passende tittel, beskrivelse målet er å være så kort og konsis som mulig for å beholde nødvendige informasjon, outputen skal være json i formatet {"title": "title", "description": "description"}`

func main() {
	cli.VersionPrinter = func(cmd *cli.Command) {
		fmt.Printf("%s %s (commit %s)\n", cmd.Root().Name, version, commit)
	}
	linearClient := linear.NewClient()
	geminiClient, err := createGeminiClient()
	if err != nil {
		fmt.Println("failed to create gemini client", err)
		return
	}
	chat, err := createGeminiChat(geminiClient)
	if err != nil {
		fmt.Println("failed to create gemini chat", err)
		return
	}

	user, err := linearClient.GetUser()
	if err != nil {
		fmt.Println(err)
		return
	}

	app := &cli.Command{
		Name:  "linear",
		Usage: "a cli for linear, help you create tasks with ai",
		Version:              version,
		EnableShellCompletion: true,
		Commands: []*cli.Command{
			{
				Name:  "create",
				Usage: "create an issue",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					fmt.Println("created issue")
					var title string
					if cmd.NArg() > 0 {
						title = cmd.Args().First()
					}
					if title == "" {
						return fmt.Errorf("title is required")
					}
					i, err := linearClient.CreateIssue(title,"", user.ID, "05071e04-d370-43c6-97cb-2a83b3214b78")
					if err != nil {
						fmt.Println(err)
						return err
					}
					fmt.Println(i)
					return nil
				},
			},
			{
				Name:  "ai",
				Usage: "create an issue with ai",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					fmt.Println("created issue with ai")
					var title string
					if cmd.NArg() > 0 {
						title = cmd.Args().First()
					}
					if title == "" {
						return fmt.Errorf("title is required")
					}
					res, err := SendMessage(chat, title)
					if err != nil {
						fmt.Println(err)
						return err
					}

					text := res.Text()
					var issue struct {
						Title string `json:"title"`
						Description string `json:"description"`
					}
					err = json.Unmarshal([]byte(text), &issue)
					if err != nil {
						fmt.Println(err)
						return err
					}
					fmt.Println(issue)
					i, err := linearClient.CreateIssue(issue.Title,  issue.Description,user.ID, "05071e04-d370-43c6-97cb-2a83b3214b78")
					if err != nil {
						fmt.Println(err)
						return err
					}
					fmt.Println(i)
					return nil
				},
			},
		},
	}
	app.Run(context.Background(), os.Args)
}
