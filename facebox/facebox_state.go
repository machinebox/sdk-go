package facebox

import (
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"net/url"
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
	resp, err := c.HTTPClient.Get(u.String())
	if err != nil {
		return nil, err
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
	resp, err := c.HTTPClient.Post(u.String(), w.FormDataContentType(), &buf)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
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
	resp, err := c.HTTPClient.PostForm(u.String(), form)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return c.parseResponse(resp.Body)
}
