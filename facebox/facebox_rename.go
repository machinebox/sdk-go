package facebox

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/machinebox/sdk-go/internal/mbhttp"
	"github.com/pkg/errors"
)

// Rename allows to change the name for a given face
func (c *Client) Rename(id, name string) error {
	if id == "" {
		return errors.New("id can not be empty")
	}
	if name == "" {
		return errors.New("name can not be empty")
	}
	u, err := url.Parse(c.addr + "/facebox/teach")
	if err != nil {
		return err
	}
	if !u.IsAbs() {
		return errors.New("box address must be absolute")
	}
	form := url.Values{}
	form.Set("name", name)

	q := u.Query()
	u.Path = u.Path + "/" + id
	u.RawQuery = q.Encode()
	req, err := http.NewRequest("PATCH", u.String(), strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json; charset=utf-8")
	_, err = mbhttp.New("facebox", c.HTTPClient).DoUnmarshal(req, nil)
	if err != nil {
		return err
	}
	return nil
}

// RenameAll changes the name for all the faces that match a given name
func (c *Client) RenameAll(oldName, newName string) error {
	if oldName == "" {
		return errors.New("oldName can not be empty")
	}
	if newName == "" {
		return errors.New("newName can not be empty")
	}
	u, err := url.Parse(c.addr + "/facebox/rename")
	if err != nil {
		return err
	}
	if !u.IsAbs() {
		return errors.New("box address must be absolute")
	}
	form := url.Values{}
	form.Set("from", oldName)
	form.Set("to", newName)

	q := u.Query()
	u.RawQuery = q.Encode()
	req, err := http.NewRequest("POST", u.String(), strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json; charset=utf-8")
	_, err = mbhttp.New("facebox", c.HTTPClient).DoUnmarshal(req, nil)
	if err != nil {
		return err
	}
	return nil
}
