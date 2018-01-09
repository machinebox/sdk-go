// Package tagbox provides a client for accessing Tagbox services.
package tagbox

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"github.com/machinebox/sdk-go/boxutil"
	"github.com/pkg/errors"
)

// Tag represents a single tag that describes an image.
type Tag struct {
	// Tag is the tag string.
	Tag string
	// Confidence is a probability number between 0 and 1.
	Confidence float64
	// ID is unique identifier of the image, previosly teach
	ID string
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

// New makes a new Client for the box at the specified address.
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
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, err
	}
	return &info, nil
}

// CheckResponse is all the data from /check request to tagbox
type CheckResponse struct {
	// Tags are the standard tags returned
	Tags []Tag `json:"tags"`
	// CustomTags are the custom tags (previously teach) that match
	CustomTags []Tag `json:"custom_tags"`
}

// ErrTagbox represents an error from Tagbox.
type ErrTagbox string

func (e ErrTagbox) Error() string {
	return "tagbox: " + string(e)
}
