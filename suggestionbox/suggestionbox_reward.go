package suggestionbox

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"path"

	"github.com/pkg/errors"
)

// Reward is used to inform Suggestionbox of a successful prediction.
type Reward struct {
	// RewardID is the ID of the reward being reported.
	RewardID string `json:"reward_id"`
	// Value is the weight of the reward.
	// Usually 1.
	Value float64 `json:"value,omitempty"`
}

// Reward tells Suggestionbox about a successful prediction.
func (c *Client) Reward(ctx context.Context, modelID string, reward Reward) error {
	u, err := url.Parse(c.addr + "/" + path.Join("suggestionbox", "models", modelID, "rewards"))
	if err != nil {
		return err
	}
	if !u.IsAbs() {
		return errors.New("box address must be absolute")
	}
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(reward); err != nil {
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
