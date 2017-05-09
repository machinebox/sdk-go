// Package tagbox provides a client for accessing tagbox services.
package tagbox

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"time"

	"github.com/machinebox/sdk-go/x/boxutil"
	"github.com/pkg/errors"
)

// Tag represents a single tag that describes an image.
type Tag struct {
	// Tag is the tag string.
	Tag string
	// Confidence is a probability number between 0 and 1.
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
	resp, err := c.HTTPClient.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, err
	}
	return &info, nil
}

// Check gets the tags for the image data provided.
func (c *Client) Check(image io.Reader) ([]Tag, error) {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	fw, err := w.CreateFormFile("file", "image.dat")
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(fw, image)
	if err != nil {
		return nil, err
	}
	if err = w.Close(); err != nil {
		return nil, err
	}
	u, err := url.Parse(c.addr + "/tagbox/check")
	if err != nil {
		return nil, err
	}
	if !u.IsAbs() {
		return nil, errors.New("box address must be absolute")
	}
	resp, err := c.HTTPClient.Post(u.String(), w.FormDataContentType(), &buf)
	if err != nil {
		return nil, err
	}
	return c.parseCheckResponse(resp.Body)
}

// CheckURL gets the tags for the image at the specified URL.
func (c *Client) CheckURL(imageURL *url.URL) ([]Tag, error) {
	u, err := url.Parse(c.addr + "/tagbox/check")
	if err != nil {
		return nil, err
	}
	if !u.IsAbs() {
		return nil, errors.New("box address must be absolute")
	}
	if !imageURL.IsAbs() {
		return nil, errors.New("url must be absolute")
	}
	form := url.Values{}
	form.Set("url", imageURL.String())
	resp, err := c.HTTPClient.PostForm(u.String(), form)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return c.parseCheckResponse(resp.Body)
}

// parseCheckResponse parses the check response data.
func (c *Client) parseCheckResponse(r io.Reader) ([]Tag, error) {
	var checkResponse struct {
		Success bool
		Error   string
		Tags    []Tag
	}
	if err := json.NewDecoder(r).Decode(&checkResponse); err != nil {
		return nil, errors.Wrap(err, "decoding response")
	}
	if !checkResponse.Success {
		return nil, ErrTagbox(checkResponse.Error)
	}
	return checkResponse.Tags, nil
}

// ErrTagbox represents an error from tagbox.
type ErrTagbox string

func (e ErrTagbox) Error() string {
	return "tagbox: " + string(e)
}
