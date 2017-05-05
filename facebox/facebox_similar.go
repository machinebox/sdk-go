package facebox

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/url"

	"github.com/pkg/errors"
)

// Similar checks the image in the io.Reader for similar faces.
func (c *Client) Similar(image io.Reader) ([]Similar, error) {
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
	u, err := url.Parse(c.addr + "/facebox/similar")
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
	return c.parseSimilarResponse(resp.Body)
}

// SimilarURL checks the image at the specified URL for similar faces.
func (c *Client) SimilarURL(imageURL *url.URL) ([]Similar, error) {
	u, err := url.Parse(c.addr + "/facebox/similar")
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
	return c.parseSimilarResponse(resp.Body)
}

func (c *Client) parseSimilarResponse(r io.Reader) ([]Similar, error) {
	var similarResponse struct {
		Success bool
		Error   string
		Similar []Similar
	}
	if err := json.NewDecoder(r).Decode(&similarResponse); err != nil {
		return nil, errors.Wrap(err, "decoding response")
	}
	if !similarResponse.Success {
		return nil, ErrFacebox(similarResponse.Error)
	}
	return similarResponse.Similar, nil
}
