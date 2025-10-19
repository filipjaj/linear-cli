/*
Linear CLI is a CLI tool for creating Linear issues using AI. Simply describe your task in natural language, and let Gemini AI generate a properly formatted issue.

Usage:

	linear-cli <description...>
	linear-cli ai <description...>
	linear-cli create <title>
	linear-cli --version
	linear-cli --help
*/
package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/filipjaj/linear-cli/internal/credentials"
	"github.com/filipjaj/linear-cli/internal/linear"
	"github.com/urfave/cli/v3"
	"google.golang.org/genai"
)

var (
	// Set these at build time: -ldflags "-X main.version=1.0.0 -X main.commit=abc123"
	version = "dev"
	commit  = "unknown"
)

const master_prompt = `Du er en AI-assistent som hjelper med å opprette oppgaver i Linear.
Basert på en gitt oppgavebeskrivelse skal du lage en kort, presis og meningsfull oppgave i JSON-formatet:

{"title": "title", "description": "description"}

Retningslinjer:

Tittel

Skal være en kort, handlingsorientert setning som oppsummerer hva som må gjøres.

Bruk maks 10–12 ord.

Fjern fyllord og irrelevante detaljer.

Beskrivelse

Skal inneholde hele eller deler av den originale inputen, spesielt hvis den inneholder kontekst eller tekniske detaljer.

Oppsummer eller rens bort overflødig tekst, men behold nødvendig informasjon for at en utvikler skal forstå oppgaven uten mer kontekst.

Hvis inputen allerede er godt formulert, bruk den nesten verbatim.

Generelle regler

Skriv på samme språk som inputen (norsk eller engelsk).

Ikke legg til informasjon som ikke finnes i inputen.

Ikke bruk markdown, kun gyldig JSON.

Output skal alltid ha nøyaktig følgende nøkler: title og description.

Eksempel

Input:

Feil ved oppdatering av deals — knappen «oppdater» gjør ingenting etter at man endrer tid.

Output:

{
  "title": "Fiks oppdateringsknapp for endring av deal-tid",
  "description": "Knappen «oppdater» fungerer ikke etter at man endrer tidspunktet på en eksisterende deal. Må undersøkes og rettes."
}`

// createIssueWithAI processes a message through AI and creates a Linear issue
func createIssueWithAI(ctx context.Context, chat *genai.Chat, linearClient *linear.Client, userID string, teamID string, message string) error {
	res, err := SendMessage(ctx, chat, message)
	if err != nil {
		return fmt.Errorf("AI request failed: %v", err)
	}

	text := res.Text()
	var issue struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}
	err = json.Unmarshal([]byte(text), &issue)
	if err != nil {
		return fmt.Errorf("failed to parse AI response: %v", err)
	}

	fmt.Printf("AI generated: %s\n", issue.Title)

	success, err := linearClient.CreateIssue(issue.Title, issue.Description, userID, teamID)
	if err != nil {
		return fmt.Errorf("failed to create issue: %v", err)
	}

	if success {
		fmt.Println("✓ Issue created successfully")
	}

	return nil
}

func promptForCredential(reader *bufio.Reader, name string, required bool) (string, error) {
	fmt.Printf("Enter %s%s: ", name, func() string {
		if required {
			return ""
		}
		return " (optional, press Enter to skip)"
	}())
	value, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	value = strings.TrimSpace(value)
	if required && value == "" {
		return "", fmt.Errorf("%s is required", name)
	}
	return value, nil
}

func main() {
	cli.VersionPrinter = func(cmd *cli.Command) {
		fmt.Printf("%s %s (commit %s)\n", cmd.Root().Name, version, commit)
	}

	app := &cli.Command{
		Name:                  "linear",
		Usage:                 "a cli for linear, help you create tasks with ai",
		Version:               version,
		EnableShellCompletion: true,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			// Default action: process all args through AI pipeline
			if cmd.NArg() == 0 {
				return fmt.Errorf("please provide a task description")
			}

			// Get credentials and initialize clients
			teamID, err := credentials.GetLinearTeamID()
			if err != nil {
				return fmt.Errorf("failed to get Linear team ID: %w. Run 'linear setup' first", err)
			}

			linearClient := linear.NewClient()
			geminiClient, err := createGeminiClient()
			if err != nil {
				return fmt.Errorf("failed to create gemini client: %w", err)
			}
			chat, err := createGeminiChat(geminiClient)
			if err != nil {
				return fmt.Errorf("failed to create gemini chat: %w", err)
			}

			user, err := linearClient.GetUser()
			if err != nil {
				return err
			}

			// Join all arguments into a single message
			message := strings.Join(cmd.Args().Slice(), " ")
			return createIssueWithAI(ctx, chat, &linearClient, user.ID, teamID, message)
		},
		Commands: []*cli.Command{
			{
				Name:  "setup",
				Usage: "setup credentials in system keyring",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					reader := bufio.NewReader(os.Stdin)

					fmt.Println("Linear CLI Setup")
					fmt.Println("----------------")
					fmt.Println("This will store your credentials securely in your system keyring.\n")

					linearAPIKey, err := promptForCredential(reader, "Linear API Key", true)
					if err != nil {
						return err
					}

					googleAPIKey, err := promptForCredential(reader, "Google API Key", true)
					if err != nil {
						return err
					}

					linearTeamID, err := promptForCredential(reader, "Linear Team ID", true)
					if err != nil {
						return err
					}

					// Save to keyring
					if err := credentials.SetLinearAPIKey(linearAPIKey); err != nil {
						return fmt.Errorf("failed to save Linear API key: %w", err)
					}

					if err := credentials.SetGoogleAPIKey(googleAPIKey); err != nil {
						return fmt.Errorf("failed to save Google API key: %w", err)
					}

					if err := credentials.SetLinearTeamID(linearTeamID); err != nil {
						return fmt.Errorf("failed to save Linear team ID: %w", err)
					}

					fmt.Println("\n✓ Credentials saved successfully to system keyring!")
					return nil
				},
			},
			{
				Name:  "create",
				Usage: "create an issue",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					var title string
					if cmd.NArg() > 0 {
						title = cmd.Args().First()
					}
					if title == "" {
						return fmt.Errorf("title is required")
					}

					teamID, err := credentials.GetLinearTeamID()
					if err != nil {
						return fmt.Errorf("failed to get Linear team ID: %w. Run 'linear setup' first", err)
					}

					linearClient := linear.NewClient()
					user, err := linearClient.GetUser()
					if err != nil {
						return err
					}

					success, err := linearClient.CreateIssue(title, "", user.ID, teamID)
					if err != nil {
						return err
					}

					if success {
						fmt.Println("✓ Issue created successfully")
					}
					return nil
				},
			},
			{
				Name:  "ai",
				Usage: "create an issue with ai (explicit subcommand, same as default behavior)",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					if cmd.NArg() == 0 {
						return fmt.Errorf("please provide a task description")
					}

					teamID, err := credentials.GetLinearTeamID()
					if err != nil {
						return fmt.Errorf("failed to get Linear team ID: %w. Run 'linear setup' first", err)
					}

					linearClient := linear.NewClient()
					geminiClient, err := createGeminiClient()
					if err != nil {
						return fmt.Errorf("failed to create gemini client: %w", err)
					}
					chat, err := createGeminiChat(geminiClient)
					if err != nil {
						return fmt.Errorf("failed to create gemini chat: %w", err)
					}

					user, err := linearClient.GetUser()
					if err != nil {
						return err
					}

					// Join all arguments into a single message
					message := strings.Join(cmd.Args().Slice(), " ")
					return createIssueWithAI(ctx, chat, &linearClient, user.ID, teamID, message)
				},
			},
		},
	}
	app.Run(context.Background(), os.Args)
}
