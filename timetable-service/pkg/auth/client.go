package auth

import (
	"fmt"
	"net/http"
	"os"
)

type Client interface {
	ValidateToken(token string) error
}

type client struct {
	baseURL string
}

func NewClient() Client {
	return &client{
		baseURL: os.Getenv("ACCOUNT_SERVICE_URL"),
	}
}

func (c *client) ValidateToken(token string) error {
	url := fmt.Sprintf("%s/api/v1/auth/validate", c.baseURL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("invalid token")
	}

	return nil
}
