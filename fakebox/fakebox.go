// Package fakebox provides a client for accessing Fakebox services.
package fakebox

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/machinebox/sdk-go/boxutil"
	"github.com/pkg/errors"
)

// Analysis represents an analysis of title, content and domain.
type Analysis struct {
	// Title is the response object for the title analysis.
	Title Title `json:"title"`
	// Content is the response object for the content analysis.
	Content Content `json:"content"`
	// Domain is the response object for the domain analysis.
	Domain Domain `json:"domain"`
}

// Title is the response object for the title analysis.
type Title struct {
	// Decision is the string representing the decision could be bias/unsure/impartial.
	Decision string `json:"decision,omitempty"`
	// Score is the numeric score of the decision is between 0.00 (bias) and 1.00 (impartial).
	Score float64 `json:"score,omitempty"`
	// Entities represents entities discovered in the text.
	Entities []Entity `json:"entities,omitempty"`
}

// Content is the response object for the content analysis.
type Content struct {
	// Decision is the string representing the decision could be bias/unsure/impartial.
	Decision string `json:"decision,omitempty"`
	// Score is the numeric score of the decision is between 0.00 (bias) and 1.00 (impartial).
	Score float64 `json:"score,omitempty"`
	// Entities represents entities discovered in the text.
	Entities []Entity `json:"entities,omitempty"`
	// Keywords are the most relevant keywords extracted from the text.
	Keywords []Keyword `json:"keywords"`
}

// Domain is the response object for the domain analysis.
type Domain struct {
	// Domain is the domain extracted from the URL.
	Domain string `json:"domain,omitempty"`
	// Category is one of the listed on the API docs.
	Category string `json:"category,omitempty"`
}

// Entity represents an entity discovered in the text.
type Entity struct {
	// Type is a string describing the kind of entity.
	Type string `json:"type"`
	// Text is the text of the entity.
	Text string `json:"text"`
	// Start is the absolute start position of the entity (in the original text).
	Start int `json:"start"`
	// Start is the absolute end position of the entity (in the original text).
	End int `json:"end"`
}

// Keyword represents a key word.
type Keyword struct {
	Keyword string `json:"keyword"`
}

// Client is an HTTP client that can make requests to the box.
type Client struct {
	addr string

	// HTTPClient is the http.Client that will be used to
	// make requests.
	HTTPClient *http.Client
}

// New makes a new Client.
func New(addr string) *Client {
	return &Client{
		addr: addr,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
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
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, errors.New(resp.Status)
	}
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, err
	}
	return &info, nil
}

// Check passes the text from the Reader to fakebox for analysis.
func (c *Client) Check(title string, content string, u *url.URL) (*Analysis, error) {
	uu, err := url.Parse(c.addr + "/fakebox/check")
	if err != nil {
		return nil, err
	}
	if !u.IsAbs() {
		return nil, errors.New("box address must be absolute")
	}
	vals := url.Values{}
	vals.Set("title", title)
	vals.Set("content", content)
	vals.Set("url", u.String())

	req, err := http.NewRequest("POST", uu.String(), strings.NewReader(vals.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json; charset=utf-8")
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, errors.New(resp.Status)
	}
	var response struct {
		Success bool
		Error   string

		Title   Title   `json:"title"`
		Content Content `json:"content"`
		Domain  Domain  `json:"domain"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, errors.Wrap(err, "decoding response")
	}
	if !response.Success {
		return nil, ErrFakebox(response.Error)
	}
	return &Analysis{
		Title:   response.Title,
		Content: response.Content,
		Domain:  response.Domain,
	}, nil
}

// ErrFakebox represents an error from Fakebox.
type ErrFakebox string

func (e ErrFakebox) Error() string {
	return "fakebox: " + string(e)
}
