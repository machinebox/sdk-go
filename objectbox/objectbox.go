// Package objectbox provides a client for accessing Objectbox services.
package objectbox

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"github.com/machinebox/sdk-go/boxutil"
	"github.com/pkg/errors"
)

type Object struct {
	Rect  Rect    `json:"rect"`
	Score float64 `json:"score"`
}

type Rect struct {
	Top    int `json:"top"`
	Left   int `json:"left"`
	Width  int `json:"width"`
	Height int `json:"height"`
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

// ErrObjectbox represents an error from Objectbox.
type ErrObjectbox string

func (e ErrObjectbox) Error() string {
	return "objectbox: " + string(e)
}

// CheckResponse is all the data from /check request to objectbox
type CheckResponse struct {
	Detectors []CheckDetectorResponse `json:"detectors"`
}

type CheckDetectorResponse struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Objects []Object `json:"objects"`
}
