package facebox

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
	req, err := http.NewRequest("POST", u.String(), &buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json; charset=utf-8")
	req.Header.Set("Content-Type", w.FormDataContentType())
	var checkResponse struct {
		Faces []Face
	}
	_, err = mbhttp.New("facebox", c.HTTPClient).DoUnmarshal(req, &checkResponse)
	if err != nil {
		return nil, err
	}
	return checkResponse.Faces, nil
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
	req, err := http.NewRequest("POST", u.String(), strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json; charset=utf-8")
	var checkResponse struct {
		Faces []Face
	}
	_, err = mbhttp.New("facebox", c.HTTPClient).DoUnmarshal(req, &checkResponse)
	if err != nil {
		return nil, err
	}
	return checkResponse.Faces, nil
}

// CheckBase64 checks the Base64 encoded image for faces.
func (c *Client) CheckBase64(data string) ([]Face, error) {
	return c.checkBase64WithOptions(data, nil)
}

// CheckBase64WithFaceprint checks the Base64 encoded image for faces and the object returned including the faceprints
func (c *Client) CheckBase64WithFaceprint(data string) ([]Face, error) {
	return c.checkBase64WithOptions(data, map[string]string{"faceprint": "true"})
}

func (c *Client) checkBase64WithOptions(data string, options map[string]string) ([]Face, error) {
	u, err := url.Parse(c.addr + "/facebox/check")
	if err != nil {
		return nil, err
	}
	if !u.IsAbs() {
		return nil, errors.New("box address must be absolute")
	}
	form := url.Values{}
	form.Set("base64", data)
	for k, v := range options {
		form.Set(k, v)
	}
	req, err := http.NewRequest("POST", u.String(), strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json; charset=utf-8")
	var checkResponse struct {
		Faces []Face
	}
	_, err = mbhttp.New("facebox", c.HTTPClient).DoUnmarshal(req, &checkResponse)
	if err != nil {
		return nil, err
	}
	return checkResponse.Faces, nil
}
