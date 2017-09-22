// Package tagbox provides a client for accessing tagbox services.
package tagbox

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
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

// Check gets the tags for the image data provided.
func (c *Client) Check(image io.Reader) (CheckResponse, error) {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	fw, err := w.CreateFormFile("file", "image.dat")
	if err != nil {
		return CheckResponse{}, err
	}
	_, err = io.Copy(fw, image)
	if err != nil {
		return CheckResponse{}, err
	}
	if err = w.Close(); err != nil {
		return CheckResponse{}, err
	}
	u, err := url.Parse(c.addr + "/tagbox/check")
	if err != nil {
		return CheckResponse{}, err
	}
	if !u.IsAbs() {
		return CheckResponse{}, errors.New("box address must be absolute")
	}
	req, err := http.NewRequest("POST", u.String(), &buf)
	if err != nil {
		return CheckResponse{}, err
	}
	req.Header.Set("Accept", "application/json; charset=utf-8")
	req.Header.Set("Content-Type", w.FormDataContentType())
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return CheckResponse{}, err
	}
	return c.parseCheckResponse(resp.Body)
}

// CheckURL gets the tags for the image at the specified URL.
func (c *Client) CheckURL(imageURL *url.URL) (CheckResponse, error) {
	u, err := url.Parse(c.addr + "/tagbox/check")
	if err != nil {
		return CheckResponse{}, err
	}
	if !u.IsAbs() {
		return CheckResponse{}, errors.New("box address must be absolute")
	}
	if !imageURL.IsAbs() {
		return CheckResponse{}, errors.New("url must be absolute")
	}
	form := url.Values{}
	form.Set("url", imageURL.String())
	req, err := http.NewRequest("POST", u.String(), strings.NewReader(form.Encode()))
	if err != nil {
		return CheckResponse{}, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json; charset=utf-8")
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return CheckResponse{}, err
	}
	defer resp.Body.Close()
	return c.parseCheckResponse(resp.Body)
}

// parseCheckResponse parses the check response data.
func (c *Client) parseCheckResponse(r io.Reader) (CheckResponse, error) {
	var resp struct {
		Success bool
		Error   string
		CheckResponse
	}
	if err := json.NewDecoder(r).Decode(&resp); err != nil {
		return CheckResponse{}, errors.Wrap(err, "decoding response")
	}
	if !resp.Success {
		return CheckResponse{}, ErrTagbox(resp.Error)
	}
	return resp.CheckResponse, nil
}

// ErrTagbox represents an error from tagbox.
type ErrTagbox string

func (e ErrTagbox) Error() string {
	return "tagbox: " + string(e)
}
