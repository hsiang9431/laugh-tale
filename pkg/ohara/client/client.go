package client

import (
	"net/http"

	"laugh-tale/pkg/ohara/types"
)

type Client struct {
	URL        string
	HTTPClient *http.Client
}

// called by poneglyph to discover service
func (c *Client) Discover() error {

	return nil
}

// called by poneglyph to retrieve privilege escalation key and payload decryption key
func (c *Client) GetKey() (types.EncKey, error) {

	return types.EncKey{}, nil
}
