package classificationbox

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"path"

	"github.com/pkg/errors"
)

// Example is a set of Feature properties with their associated Class
// which is used to teach Classificationbox models.
type Example struct {
	Class  string    `json:"class"`
	Inputs []Feature `json:"inputs"`
}

// Teach gives an Example to a model for it to learn from.
func (c *Client) Teach(ctx context.Context, modelID string, example Example) error {
	u, err := url.Parse(c.addr + "/" + path.Join("classificationbox", "models", modelID, "teach"))
	if err != nil {
		return err
	}
	if !u.IsAbs() {
		return errors.New("box address must be absolute")
	}
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(example); err != nil {
		return errors.Wrap(err, "encoding request body")
	}
	req, err := http.NewRequest(http.MethodPost, u.String(), &buf)
	if err != nil {
		return err
	}
	req = req.WithContext(ctx)
	req.Header.Set("Accept", "application/json; charset=utf-8")
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	_, err = c.client.Do(req, nil)
	if err != nil {
		return err
	}
	return nil
}
