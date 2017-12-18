package tagbox

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

// Teach teaches tagbox the image with a custom tag in the io.Reader.
// The tag is the string representation of the main thing on the image.
// The id should be a unique identifier for the image, usually the filename.
func (c *Client) Teach(image io.Reader, id, tag string) error {
	fn := id
	if fn == "" {
		fn = "image.dat"
	}
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	fw, err := w.CreateFormFile("file", fn)
	if err != nil {
		return err
	}
	_, err = io.Copy(fw, image)
	if err != nil {
		return err
	}
	if err := w.WriteField("tag", tag); err != nil {
		return err
	}
	if err := w.WriteField("id", id); err != nil {
		return err
	}
	if err = w.Close(); err != nil {
		return err
	}
	u, err := url.Parse(c.addr + "/tagbox/teach")
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
	return c.parseResponse(resp.Body)
}

// TeachURL teaches tagbox the image with a custom tag at the specified URL.
// See Teach for more information.
func (c *Client) TeachURL(imageURL *url.URL, id, tag string) error {
	u, err := url.Parse(c.addr + "/tagbox/teach")
	if err != nil {
		return err
	}
	if !u.IsAbs() {
		return errors.New("box address must be absolute")
	}
	if !imageURL.IsAbs() {
		return errors.New("url must be absolute")
	}
	form := url.Values{}
	form.Set("url", imageURL.String())
	form.Set("tag", tag)
	form.Set("id", id)
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
	err = c.parseResponse(resp.Body)
	if err != nil {
		return err
	}
	return nil
}

// TeachBase64 teaches tagbox the Base64 encoded image with a custom tag.
// See Teach for more information.
func (c *Client) TeachBase64(data, id, tag string) error {
	u, err := url.Parse(c.addr + "/tagbox/teach")
	if err != nil {
		return err
	}
	if !u.IsAbs() {
		return errors.New("box address must be absolute")
	}
	form := url.Values{}
	form.Set("base64", data)
	form.Set("tag", tag)
	form.Set("id", id)
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
	err = c.parseResponse(resp.Body)
	if err != nil {
		return err
	}
	return nil
}

// Remove makes tagbox to forget an image
func (c *Client) Remove(id string) error {
	if id == "" {
		return errors.New("id can not be empty")
	}
	u, err := url.Parse(c.addr + "/tagbox/teach")
	if err != nil {
		return err
	}
	if !u.IsAbs() {
		return errors.New("box address must be absolute")
	}

	q := u.Query()
	u.Path = u.Path + "/" + id
	u.RawQuery = q.Encode()
	req, err := http.NewRequest("DELETE", u.String(), nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json; charset=utf-8")
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	err = c.parseResponse(resp.Body)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) parseResponse(r io.Reader) error {
	var response struct {
		Success bool
		Error   string
	}
	if err := json.NewDecoder(r).Decode(&response); err != nil {
		return errors.Wrap(err, "decoding response")
	}
	if !response.Success {
		return ErrTagbox(response.Error)
	}
	return nil
}
