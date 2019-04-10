package objectbox

import (
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"

	"github.com/machinebox/sdk-go/internal/mbhttp"
)

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
	u, err := url.Parse(c.addr + "/objectbox/state")
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
	_, err = mbhttp.New("objectbox", c.HTTPClient).DoUnmarshal(req, nil)
	if err != nil {
		return err
	}
	return nil
}

// PostStateURL tells objectbox to download the state file specified
// by the URL.
func (c *Client) PostStateURL(stateURL *url.URL) error {
	u, err := url.Parse(c.addr + "/objectbox/state")
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
	client := mbhttp.New("objectbox", c.HTTPClient)
	_, err = client.DoUnmarshal(req, nil)
	if err != nil {
		return err
	}
	return nil
}
