// Package mbhttp provides helpers common across all SDKs.
package mbhttp

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

// Client makes requests and handles common Machine Box error cases.
type Client struct {
	boxname string

	// HTTPClient is the underlying http.Client that will be
	// used to make requests.
	HTTPClient *http.Client
}

// New makes a new Client.
func New(boxname string, client *http.Client) *Client {
	return &Client{
		boxname:    boxname,
		HTTPClient: client,
	}
}

// Do makes the request and unmarshals the response into v.
func (c *Client) Do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "read response data")
	}
	var o struct {
		Success bool
		Error   string
	}
	if err := json.Unmarshal(b, &o); err != nil {
		if resp.StatusCode < 200 || resp.StatusCode >= 400 {
			return nil, errors.Errorf("%s: %s", c.boxname, resp.Status)
		}
		return nil, errors.Wrap(err, "decode common response data")
	}
	if !o.Success {
		if o.Error == "" {
			o.Error = "an unknown error occurred in the box"
		}
		return nil, errors.Errorf("%s: %s", c.boxname, o.Error)
	}
	if err := json.Unmarshal(b, &v); err != nil {
		return nil, errors.Wrap(err, "decode response data")
	}
	return resp, nil
}
