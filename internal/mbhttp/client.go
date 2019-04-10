// Package mbhttp provides helpers common across all SDKs.
package mbhttp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

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

// DoUnmarshal makes the request and unmarshals the response into v.
// The Body in the Response will be closed after calling this method.
func (c *Client) DoUnmarshal(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "read response data")
	}
	if len(b) == 0 {
		if resp.StatusCode < 200 || resp.StatusCode >= 400 {
			return nil, errors.Errorf("%s: %s", c.boxname, resp.Status)
		}
		return resp, nil
	}
	var o struct {
		Success bool
		Error   string
	}
	if err := json.Unmarshal(b, &o); err != nil {
		if resp.StatusCode < 200 || resp.StatusCode >= 400 {
			return nil, errors.Errorf("%s: %d: %s", c.boxname, resp.StatusCode, strings.TrimSpace(string(b)))
		}
		return nil, errors.Wrap(err, "decode common response data")
	}
	if !o.Success {
		if o.Error == "" {
			o.Error = fmt.Sprintf("%d: %s", resp.StatusCode, strings.TrimSpace(string(b)))
		}
		return nil, errors.Errorf("%s: %s", c.boxname, o.Error)
	}
	if err := json.Unmarshal(b, &v); err != nil {
		return nil, errors.Wrap(err, "decode response data")
	}
	return resp, nil
}
