package objectbox

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

// Check gets the objects for the image data provided.
func (c *Client) Check(image io.Reader) (CheckResponse, error) {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	fw, err := w.CreateFormFile("file", "image.dat")
	if err != nil {
		return CheckResponse{}, err
	}
	_, err = io.Copy(fw, image)
	if err != nil {
		return CheckResponse{}, err
	}
	if err = w.Close(); err != nil {
		return CheckResponse{}, err
	}
	u, err := url.Parse(c.addr + "/objectbox/check")
	if err != nil {
		return CheckResponse{}, err
	}
	if !u.IsAbs() {
		return CheckResponse{}, errors.New("box address must be absolute")
	}
	req, err := http.NewRequest("POST", u.String(), &buf)
	if err != nil {
		return CheckResponse{}, err
	}
	req.Header.Set("Accept", "application/json; charset=utf-8")
	req.Header.Set("Content-Type", w.FormDataContentType())
	var response CheckResponse
	_, err = mbhttp.New("objectbox", c.HTTPClient).DoUnmarshal(req, &response)
	if err != nil {
		return CheckResponse{}, err
	}
	return response, nil
}

// CheckURL gets the tags for the image at the specified URL.
func (c *Client) CheckURL(imageURL *url.URL) (CheckResponse, error) {
	u, err := url.Parse(c.addr + "/objectbox/check")
	if err != nil {
		return CheckResponse{}, err
	}
	if !u.IsAbs() {
		return CheckResponse{}, errors.New("box address must be absolute")
	}
	if !imageURL.IsAbs() {
		return CheckResponse{}, errors.New("url must be absolute")
	}
	form := url.Values{}
	form.Set("url", imageURL.String())
	req, err := http.NewRequest("POST", u.String(), strings.NewReader(form.Encode()))
	if err != nil {
		return CheckResponse{}, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json; charset=utf-8")
	client := mbhttp.New("objectbox", c.HTTPClient)
	var response CheckResponse
	_, err = client.DoUnmarshal(req, &response)
	if err != nil {
		return CheckResponse{}, err
	}
	return response, nil
}

// CheckBase64 gets the tags for the image in the encoded Base64 data string.
func (c *Client) CheckBase64(data string) (CheckResponse, error) {
	u, err := url.Parse(c.addr + "/objectbox/check")
	if err != nil {
		return CheckResponse{}, err
	}
	if !u.IsAbs() {
		return CheckResponse{}, errors.New("box address must be absolute")
	}
	form := url.Values{}
	form.Set("base64", data)
	req, err := http.NewRequest("POST", u.String(), strings.NewReader(form.Encode()))
	if err != nil {
		return CheckResponse{}, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json; charset=utf-8")
	client := mbhttp.New("objectbox", c.HTTPClient)
	var response CheckResponse
	_, err = client.DoUnmarshal(req, &response)
	if err != nil {
		return CheckResponse{}, err
	}
	return response, nil
}
