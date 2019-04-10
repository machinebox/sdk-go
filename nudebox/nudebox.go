// Package nudebox provides a client for accessing Nudebox services.
package nudebox

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/machinebox/sdk-go/boxutil"
	"github.com/machinebox/sdk-go/internal/mbhttp"
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
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json; charset=utf-8")
	_, err = mbhttp.New("nudebox", c.HTTPClient).DoUnmarshal(req, &info)
	if err != nil {
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
	req, err := http.NewRequest("POST", u.String(), &buf)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Accept", "application/json; charset=utf-8")
	req.Header.Set("Content-Type", w.FormDataContentType())
	var checkResponse struct {
		Nude float64
	}
	_, err = mbhttp.New("nudebox", c.HTTPClient).DoUnmarshal(req, &checkResponse)
	if err != nil {
		return 0, err
	}
	return checkResponse.Nude, nil
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
	req, err := http.NewRequest("POST", u.String(), strings.NewReader(form.Encode()))
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json; charset=utf-8")
	var checkResponse struct {
		Nude float64
	}
	_, err = mbhttp.New("nudebox", c.HTTPClient).DoUnmarshal(req, &checkResponse)
	if err != nil {
		return 0, err
	}
	return checkResponse.Nude, nil
}

// CheckBase64 gets the nudity probability for the Base64 encoded image.
func (c *Client) CheckBase64(data string) (float64, error) {
	u, err := url.Parse(c.addr + "/nudebox/check")
	if err != nil {
		return 0, err
	}
	if !u.IsAbs() {
		return 0, errors.New("box address must be absolute")
	}
	form := url.Values{}
	form.Set("base64", data)
	req, err := http.NewRequest("POST", u.String(), strings.NewReader(form.Encode()))
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json; charset=utf-8")
	var checkResponse struct {
		Nude float64
	}
	_, err = mbhttp.New("nudebox", c.HTTPClient).DoUnmarshal(req, &checkResponse)
	if err != nil {
		return 0, err
	}
	return checkResponse.Nude, nil
}

// ErrNudebox represents an error from Nudebox.
type ErrNudebox string

func (e ErrNudebox) Error() string {
	return "nudebox: " + string(e)
}
