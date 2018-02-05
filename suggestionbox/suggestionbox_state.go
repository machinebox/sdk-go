package suggestionbox

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/pkg/errors"
)

// OpenState opens the state file for the specified model for reading.
// Clients must call Close.
func (c *Client) OpenState(ctx context.Context, modelID string) (io.ReadCloser, error) {
	u, err := url.Parse(c.addr + "/" + path.Join("suggestionbox", "state", modelID))
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
	req = req.WithContext(ctx)
	resp, err := c.client.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

// PostState uploads new state data and returns the Model that was contained
// in the state file.
func (c *Client) PostState(ctx context.Context, r io.Reader) (Model, error) {
	var model Model
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	fw, err := w.CreateFormFile("file", "image.dat")
	if err != nil {
		return model, err
	}
	_, err = io.Copy(fw, r)
	if err != nil {
		return model, err
	}
	if err = w.Close(); err != nil {
		return model, err
	}
	u, err := url.Parse(c.addr + "/suggestionbox/state")
	if err != nil {
		return model, err
	}
	if !u.IsAbs() {
		return model, errors.New("box address must be absolute")
	}
	req, err := http.NewRequest("POST", u.String(), &buf)
	if err != nil {
		return model, err
	}
	req = req.WithContext(ctx)
	req.Header.Set("Accept", "application/json; charset=utf-8")
	req.Header.Set("Content-Type", w.FormDataContentType())
	_, err = c.client.Do(req, &model)
	if err != nil {
		return model, err
	}
	return model, nil
}

// PostStateURL tells Suggestionbox to download the state file specified
// by the URL and returns the Model that was contained in the state file.
func (c *Client) PostStateURL(ctx context.Context, stateURL *url.URL) (Model, error) {
	var model Model
	u, err := url.Parse(c.addr + "/suggestionbox/state")
	if err != nil {
		return model, err
	}
	if !u.IsAbs() {
		return model, errors.New("box address must be absolute")
	}
	if !stateURL.IsAbs() {
		return model, errors.New("url must be absolute")
	}
	form := url.Values{}
	form.Set("url", stateURL.String())
	req, err := http.NewRequest("POST", u.String(), strings.NewReader(form.Encode()))
	if err != nil {
		return model, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json; charset=utf-8")
	req = req.WithContext(ctx)
	_, err = c.client.Do(req, &model)
	if err != nil {
		return model, err
	}
	return model, nil
}
