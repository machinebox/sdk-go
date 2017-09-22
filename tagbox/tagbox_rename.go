package tagbox

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

// Rename allows to change the custom tag for a given image by id
func (c *Client) Rename(id, tag string) error {
	if id == "" {
		return errors.New("id can not be empty")
	}
	if tag == "" {
		return errors.New("tag can not be empty")
	}
	u, err := url.Parse(c.addr + "/tagbox/teach")
	if err != nil {
		return err
	}
	if !u.IsAbs() {
		return errors.New("box address must be absolute")
	}
	form := url.Values{}
	form.Set("tag", tag)

	q := u.Query()
	u.Path = u.Path + "/" + id
	u.RawQuery = q.Encode()
	req, err := http.NewRequest("PATCH", u.String(), strings.NewReader(form.Encode()))
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

// RenameAll changes the tag for all the images
func (c *Client) RenameAll(oldTag, newTag string) error {
	if oldTag == "" {
		return errors.New("oldTag can not be empty")
	}
	if newTag == "" {
		return errors.New("newTag can not be empty")
	}
	u, err := url.Parse(c.addr + "/tagbox/rename")
	if err != nil {
		return err
	}
	if !u.IsAbs() {
		return errors.New("box address must be absolute")
	}
	form := url.Values{}
	form.Set("from", oldTag)
	form.Set("to", newTag)

	q := u.Query()
	u.RawQuery = q.Encode()
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
