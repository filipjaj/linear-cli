package main

import (
	"context"
	"fmt"

	"github.com/filipjaj/linear-cli/internal/credentials"
	"google.golang.org/genai"
)

func createGeminiClient() (*genai.Client, error) {
	apiKey, err := credentials.GetGoogleAPIKey()
	if err != nil {
		return nil, fmt.Errorf("failed to get Google API key from keyring: %w", err)
	}
	client, err := genai.NewClient(context.Background(), &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, err
	}
	return client, nil
}

func createGeminiChat(client *genai.Client) (*genai.Chat, error) {
	var model string = "gemini-2.5-flash-lite"

	config := &genai.GenerateContentConfig{
		ResponseMIMEType: "application/json",
		ResponseSchema: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"title":       {Type: genai.TypeString},
				"description": {Type: genai.TypeString},
			},
			PropertyOrdering: []string{"title", "description"},
		},
	}

	chat, err := client.Chats.Create(context.Background(), model, config, nil)
	if err != nil {
		return nil, err
	}

	return chat, nil
}

func SendMessage(ctx context.Context, chat *genai.Chat, message string) (*genai.GenerateContentResponse, error) {
	text, err := chat.SendMessage(ctx, genai.Part{
		Text: master_prompt,
	}, genai.Part{
		Text: message,
	})

	if err != nil {
		return nil, err
	}
	return text, nil
}
