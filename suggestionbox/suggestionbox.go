package suggestionbox

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/machinebox/sdk-go/boxutil"
)

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

// ErrSuggestionbox represents an error from suggestionbox.
type ErrSuggestionbox string

func (e ErrSuggestionbox) Error() string {
	if string(e) == "" {
		return "suggestionbox: success != true (unknown error)"
	}
	return "suggestionbox: " + string(e)
}
