package linear

import (
	// this will automatically load your .env file:

	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/hasura/go-graphql-client"

	_ "github.com/joho/godotenv/autoload"
)

// Client -
type Client struct {
	graphqlClient *graphql.Client
	
}

// NewClient -
func NewClient() Client {
	
	
	return Client{
		graphqlClient: graphql.NewClient("https://api.linear.app/graphql", http.DefaultClient).
		WithRequestModifier(func(r *http.Request) {
			apiKey := os.Getenv("LINEAR_API_KEY")
			r.Header.Set("Authorization", apiKey)
		}),
	}
}


type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
}

func (c *Client) GetUser() (*User, error) {
	client := c.graphqlClient

	var resp struct {
		Viewer User `json:"viewer"`
	}

	err := client.Query(context.Background(), &resp, nil)
	if err != nil {
		return nil, err
	}
	
	
	return &resp.Viewer, nil
}

type Issues struct {
	Nodes []struct {
		ID   string `json:"id"`
		Title string `json:"title"`
	} `json:"nodes"`
}

type IssueInput struct {
	Title string `json:"title"`
	TeamId string `json:"teamId"`
	Description string `json:"description"`
	AssigneeId string `json:"assigneeId"`
}

func (IssueInput) GetGraphQLType() string { return "IssueCreateInput" }


func (c *Client) CreateIssue(issueTitle string,description string, assigneeID string,  teamID string) (bool, error) {
	client := c.graphqlClient

	var resp struct {
		CreateIssue struct {
			LastSyncId int `graphql:"lastSyncId"`
			Success bool `graphql:"success"`
		} `graphql:"issueCreate(input: $issue)"`
	}
	issue := IssueInput{
		Title: issueTitle,	
		TeamId: teamID,
		Description: description,
		AssigneeId: assigneeID,
	}

	vars := map[string]interface{}{
		"issue": issue,
	}

	err := client.Mutate(context.Background(), &resp, vars)
	if err != nil {
		return false, fmt.Errorf("error creating issue: %v", err)
	}
	
	
	return resp.CreateIssue.Success, nil
}
