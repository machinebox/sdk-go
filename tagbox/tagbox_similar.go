package tagbox

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"

	"github.com/machinebox/sdk-go/internal/mbhttp"
	"github.com/pkg/errors"
)

// Similar checks the image in the io.Reader for similar
// images based on tags previously taught.
func (c *Client) Similar(image io.Reader) ([]Tag, error) {
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
	u, err := url.Parse(c.addr + "/tagbox/similar")
	if err != nil {
		return nil, err
	}
	if !u.IsAbs() {
		return nil, errors.New("box address must be absolute")
	}
	req, err := http.NewRequest("POST", u.String(), &buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json; charset=utf-8")
	req.Header.Set("Content-Type", w.FormDataContentType())
	var similarResponse struct {
		Similar []Tag
	}
	_, err = mbhttp.New("tagbox", c.HTTPClient).DoUnmarshal(req, &similarResponse)
	if err != nil {
		return nil, err
	}
	return similarResponse.Similar, nil
}

// SimilarURL checks the image at the specified URL for similar
// images based on tags previously taught.
func (c *Client) SimilarURL(imageURL *url.URL) ([]Tag, error) {
	u, err := url.Parse(c.addr + "/tagbox/similar")
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
	req, err := http.NewRequest("POST", u.String(), strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json; charset=utf-8")
	var similarResponse struct {
		Similar []Tag
	}
	_, err = mbhttp.New("tagbox", c.HTTPClient).DoUnmarshal(req, &similarResponse)
	if err != nil {
		return nil, err
	}
	return similarResponse.Similar, nil
}

// SimilarBase64 checks the image at the specified URL for similar
// images based on tags previously taught.
func (c *Client) SimilarBase64(data string) ([]Tag, error) {
	u, err := url.Parse(c.addr + "/tagbox/similar")
	if err != nil {
		return nil, err
	}
	if !u.IsAbs() {
		return nil, errors.New("box address must be absolute")
	}
	form := url.Values{}
	form.Set("base64", data)
	req, err := http.NewRequest("POST", u.String(), strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json; charset=utf-8")
	var similarResponse struct {
		Similar []Tag
	}
	_, err = mbhttp.New("tagbox", c.HTTPClient).DoUnmarshal(req, &similarResponse)
	if err != nil {
		return nil, err
	}
	return similarResponse.Similar, nil
}

// SimilarID returns similar images based on the ID provided
func (c *Client) SimilarID(id string) ([]Tag, error) {
	u, err := url.Parse(c.addr + "/tagbox/similar")
	if err != nil {
		return nil, err
	}
	if !u.IsAbs() {
		return nil, errors.New("box address must be absolute")
	}
	if id == "" {
		return nil, errors.New("id can not be empty")
	}
	q := u.Query()
	q.Set("id", id)
	u.RawQuery = q.Encode()
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json; charset=utf-8")
	var similarResponse struct {
		Similar []Tag
	}
	_, err = mbhttp.New("tagbox", c.HTTPClient).DoUnmarshal(req, &similarResponse)
	if err != nil {
		return nil, err
	}
	return similarResponse.Similar, nil
}
