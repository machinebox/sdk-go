package suggestionbox

import (
	"bytes"
	"context"
	"encoding/json"
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
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, errors.New(resp.Status)
	}
	return resp.Body, nil
}

// PostState uploads new state data and returns the Model that was contained
// in the state file.
func (c *Client) PostState(ctx context.Context, r io.Reader) (Model, error) {
	var response struct {
		Success bool
		Error   string
		Model
	}
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	fw, err := w.CreateFormFile("file", "image.dat")
	if err != nil {
		return response.Model, err
	}
	_, err = io.Copy(fw, r)
	if err != nil {
		return response.Model, err
	}
	if err = w.Close(); err != nil {
		return response.Model, err
	}
	u, err := url.Parse(c.addr + "/suggestionbox/state")
	if err != nil {
		return response.Model, err
	}
	if !u.IsAbs() {
		return response.Model, errors.New("box address must be absolute")
	}
	req, err := http.NewRequest("POST", u.String(), &buf)
	if err != nil {
		return response.Model, err
	}
	req.Header.Set("Accept", "application/json; charset=utf-8")
	req.Header.Set("Content-Type", w.FormDataContentType())
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return response.Model, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return response.Model, errors.New(resp.Status)
	}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return response.Model, errors.Wrap(err, "decoding response")
	}
	if !response.Success {
		return response.Model, ErrSuggestionbox(response.Error)
	}
	return response.Model, nil
}

// PostStateURL tells Suggestionbox to download the state file specified
// by the URL and returns the Model that was contained in the state file.
func (c *Client) PostStateURL(ctx context.Context, stateURL *url.URL) (Model, error) {
	var response struct {
		Success bool
		Error   string
		Model
	}
	u, err := url.Parse(c.addr + "/suggestionbox/state")
	if err != nil {
		return response.Model, err
	}
	if !u.IsAbs() {
		return response.Model, errors.New("box address must be absolute")
	}
	if !stateURL.IsAbs() {
		return response.Model, errors.New("url must be absolute")
	}
	form := url.Values{}
	form.Set("url", stateURL.String())
	req, err := http.NewRequest("POST", u.String(), strings.NewReader(form.Encode()))
	if err != nil {
		return response.Model, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json; charset=utf-8")
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return response.Model, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return response.Model, errors.New(resp.Status)
	}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return response.Model, errors.Wrap(err, "decoding response")
	}
	if !response.Success {
		return response.Model, ErrSuggestionbox(response.Error)
	}
	return response.Model, nil
}
