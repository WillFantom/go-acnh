package goacnh

import (
	"os"

	"github.com/go-resty/resty/v2"
)

const (
	baseURL string = "https://acnhapi.com"
)

// Client facilitates interaction with the AC:NH API
type Client struct {
	restClient *resty.Client
}

// New creates a new instance of the AC:NH API client
func New() *Client {
	c := Client{
		restClient: resty.New(),
	}
	c.restClient.SetBaseURL(baseURL)
	return &c
}

func dirExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
