package facebox

import (
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
)

// OpenState opens the state file for reading.
// Clients must call Close.
func (c *Client) OpenState() (io.ReadCloser, error) {
	u, err := url.Parse(c.addr + "/facebox/state")
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

// PostState uploads new state data.
func (c *Client) PostState(r io.Reader) error {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	fw, err := w.CreateFormFile("file", "image.dat")
	if err != nil {
		return err
	}
	_, err = io.Copy(fw, r)
	if err != nil {
		return err
	}
	if err = w.Close(); err != nil {
		return err
	}
	u, err := url.Parse(c.addr + "/facebox/state")
	if err != nil {
		return err
	}
	if !u.IsAbs() {
		return errors.New("box address must be absolute")
	}
	req, err := http.NewRequest("POST", u.String(), &buf)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json; charset=utf-8")
	req.Header.Set("Content-Type", w.FormDataContentType())
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.New(resp.Status)
	}
	return c.parseResponse(resp.Body)
}

// PostStateURL tells facebox to download the state file specified
// by the URL.
func (c *Client) PostStateURL(stateURL *url.URL) error {
	u, err := url.Parse(c.addr + "/facebox/state")
	if err != nil {
		return err
	}
	if !u.IsAbs() {
		return errors.New("box address must be absolute")
	}
	if !stateURL.IsAbs() {
		return errors.New("url must be absolute")
	}
	form := url.Values{}
	form.Set("url", stateURL.String())
	req, err := http.NewRequest("POST", u.String(), strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json; charset=utf-8")
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.New(resp.Status)
	}
	return c.parseResponse(resp.Body)
}
