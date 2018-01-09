package videobox

import (
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

// Delete removes the results for a video.
func (c *Client) Delete(id string) error {
	u, err := url.Parse(c.addr + "/videobox/results/" + id)
	if err != nil {
		return err
	}
	if !u.IsAbs() {
		return errors.New("box address must be absolute")
	}
	req, err := http.NewRequest(http.MethodDelete, u.String(), nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json; charset=utf-8")
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.New(resp.Status)
	}
	return nil
}
