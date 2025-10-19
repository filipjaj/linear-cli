package linear

import (
	"context"
	"fmt"
	"net/http"

	"github.com/filipjaj/linear-cli/internal/credentials"
	"github.com/hasura/go-graphql-client"
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
			apiKey, err := credentials.GetLinearAPIKey()
			if err != nil {
				// Log error but don't panic - let API call fail with auth error
				fmt.Printf("Warning: failed to get Linear API key from keyring: %v\n", err)
				return
			}
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

func (c *Client) GetCycle() (*Issues, error) {
	client := c.graphqlClient
	var resp struct {
		Viewer struct {
			Cycle struct {
				Issues Issues `json:"issues"`
			} `json:"cycle"`
		} `json:"viewer"`
	}
	
	err := client.Query(context.Background(), &resp, nil)
	if err != nil {
		return nil, err
	}
	
	return &resp.Viewer.Cycle.Issues, nil
}
