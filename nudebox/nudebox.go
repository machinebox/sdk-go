// Package nudebox provides a client for accessing nudebox services.
package nudebox

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

// Check gets the nudity probability for the image data provided.
func (c *Client) Check(image io.Reader) (float64, error) {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	fw, err := w.CreateFormFile("file", "image.dat")
	if err != nil {
		return 0, err
	}
	_, err = io.Copy(fw, image)
	if err != nil {
		return 0, err
	}
	if err = w.Close(); err != nil {
		return 0, err
	}
	u, err := url.Parse(c.addr + "/nudebox/check")
	if err != nil {
		return 0, err
	}
	if !u.IsAbs() {
		return 0, errors.New("box address must be absolute")
	}
	resp, err := c.HTTPClient.Post(u.String(), w.FormDataContentType(), &buf)
	if err != nil {
		return 0, err
	}
	return c.parseCheckResponse(resp.Body)
}

// CheckURL gets the nudity probability for the image at the specified URL.
func (c *Client) CheckURL(imageURL *url.URL) (float64, error) {
	u, err := url.Parse(c.addr + "/nudebox/check")
	if err != nil {
		return 0, err
	}
	if !u.IsAbs() {
		return 0, errors.New("box address must be absolute")
	}
	if !imageURL.IsAbs() {
		return 0, errors.New("url must be absolute")
	}
	form := url.Values{}
	form.Set("url", imageURL.String())
	resp, err := c.HTTPClient.PostForm(u.String(), form)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	return c.parseCheckResponse(resp.Body)
}

// parseCheckResponse parses the check response data.
func (c *Client) parseCheckResponse(r io.Reader) (float64, error) {
	var checkResponse struct {
		Success bool
		Error   string
		Nude    float64
	}
	if err := json.NewDecoder(r).Decode(&checkResponse); err != nil {
		return 0, errors.Wrap(err, "decoding response")
	}
	if !checkResponse.Success {
		return 0, ErrNudebox(checkResponse.Error)
	}
	return checkResponse.Nude, nil
}

// ErrNudebox represents an error from nudebox.
type ErrNudebox string

func (e ErrNudebox) Error() string {
	return "nudebox: " + string(e)
}
