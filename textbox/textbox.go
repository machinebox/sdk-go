// Package textbox provides a client for accessing Textbox services.
package textbox

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/machinebox/sdk-go/boxutil"
	"github.com/pkg/errors"
)

// Analysis represents an analysis of text.
type Analysis struct {
	Sentences []Sentence `json:"sentences"`
	Keywords  []Keyword  `json:"keywords"`
}

// Sentence represents a single sentence of text.
type Sentence struct {
	// Text is the text of the sentence.
	Text string `json:"text"`
	// Start is the absolute start position of the sentence (in the original text).
	Start int `json:"start"`
	// Start is the absolute end position of the sentence (in the original text).
	End int `json:"end"`
	// Sentiment is a probability score (between 0 and 1) of the sentiment of the sentence;
	// higher is more positive, lower is more negative.
	Sentiment float64 `json:"sentiment"`
	// Entities represents entities discovered in the text.
	Entities []Entity `json:"entities"`
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

// Check passes the text from the Reader to Textbox for analysis.
func (c *Client) Check(r io.Reader) (*Analysis, error) {
	u, err := url.Parse(c.addr + "/textbox/check")
	if err != nil {
		return nil, err
	}
	if !u.IsAbs() {
		return nil, errors.New("box address must be absolute")
	}
	vals := url.Values{}
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	vals.Set("text", string(b))
	req, err := http.NewRequest("POST", u.String(), strings.NewReader(vals.Encode()))
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
		Success   bool
		Error     string
		Sentences []Sentence `json:"sentences"`
		Keywords  []Keyword  `json:"keywords"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, errors.Wrap(err, "decoding response")
	}
	if !response.Success {
		return nil, ErrTextbox(response.Error)
	}
	return &Analysis{
		Sentences: response.Sentences,
		Keywords:  response.Keywords,
	}, nil
}

// ErrTextbox represents an error from Textbox.
type ErrTextbox string

func (e ErrTextbox) Error() string {
	return "textbox: " + string(e)
}
