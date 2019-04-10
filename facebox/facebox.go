// Package facebox provides a client for accessing Facebox services.
package facebox

import (
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/machinebox/sdk-go/internal/mbhttp"

	"github.com/machinebox/sdk-go/boxutil"
)

// Face represents a face in an image.
type Face struct {
	Rect       Rect
	ID         string
	Name       string
	Matched    bool
	Confidence float64
	Faceprint  string
}

// Rect represents the coordinates of a face within an image.
type Rect struct {
	Top, Left     int
	Width, Height int
}

// Similar represents a similar face.
type Similar struct {
	ID         string
	Name       string
	Confidence float64
}

// Client is an HTTP client that can make requests to the box.
type Client struct {
	addr string

	// HTTPClient is the http.Client that will be used to
	// make requests.
	HTTPClient *http.Client
}

// make sure the Client implements boxutil.Box
var _ boxutil.Box = (*Client)(nil)

// New creates a new Client.
func New(addr string) *Client {
	return &Client{
		addr: addr,
		HTTPClient: &http.Client{
			Timeout: 1 * time.Minute,
		},
	}
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
	_, err = mbhttp.New("facebox", c.HTTPClient).DoUnmarshal(req, &info)
	if err != nil {
		return nil, err
	}
	return &info, nil
}
