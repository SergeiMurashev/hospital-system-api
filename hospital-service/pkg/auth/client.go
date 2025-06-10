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
	url := fmt.Sprintf("%s/api/Authentication/Validate?accessToken=%s", c.baseURL, token)
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to validate token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("invalid token: %s", resp.Status)
	}

	return nil
}
