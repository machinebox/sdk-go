// Package suggestionbox provides a client for accessing Suggestionbox services.
package suggestionbox

import (
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/machinebox/sdk-go/boxutil"
	"github.com/machinebox/sdk-go/internal/mbhttp"
)

// Client is an HTTP client that can make requests to the box.
type Client struct {
	addr   string
	client *mbhttp.Client
}

// make sure the Client implements boxutil.Box
var _ boxutil.Box = (*Client)(nil)

// New makes a new Client for the box at the specified address.
func New(addr string) *Client {
	c := &Client{
		addr: addr,
	}
	c.SetClient(&http.Client{Timeout: 1 * time.Minute})
	return c
}

// SetClient sets the http.Client to use when making requests.
func (c *Client) SetClient(client *http.Client) {
	c.client = mbhttp.New("suggestionbox", client)
}

// Info gets the details about the box.
func (c *Client) Info() (*boxutil.Info, error) {
	var info boxutil.Info
	u, err := url.Parse(c.addr + "/info")
	if err != nil {
		return nil, err
	}
	if !u.IsAbs() {
		return nil, errors.New("box address must be absolute")
	}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json; charset=utf-8")
	_, err = c.client.Do(req, &info)
	if err != nil {
		return nil, err
	}
	return &info, nil
}
