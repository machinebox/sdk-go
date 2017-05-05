package facebox

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/url"

	"github.com/pkg/errors"
)

// Check checks the image in the io.Reader for faces.
func (c *Client) Check(image io.Reader) ([]Face, error) {
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
	u, err := url.Parse(c.addr + "/facebox/check")
	if err != nil {
		return nil, err
	}
	if !u.IsAbs() {
		return nil, errors.New("box address must be absolute")
	}
	resp, err := c.HTTPClient.Post(u.String(), w.FormDataContentType(), &buf)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return c.parseCheckResponse(resp.Body)
}

// CheckURL checks the image at the specified URL for faces.
func (c *Client) CheckURL(imageURL *url.URL) ([]Face, error) {
	u, err := url.Parse(c.addr + "/facebox/check")
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
	resp, err := c.HTTPClient.PostForm(u.String(), form)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return c.parseCheckResponse(resp.Body)
}

func (c *Client) parseCheckResponse(r io.Reader) ([]Face, error) {
	var checkResponse struct {
		Success bool
		Error   string
		Faces   []Face
	}
	if err := json.NewDecoder(r).Decode(&checkResponse); err != nil {
		return nil, errors.Wrap(err, "decoding response")
	}
	if !checkResponse.Success {
		return nil, ErrFacebox(checkResponse.Error)
	}
	return checkResponse.Faces, nil
}
