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

// Prediction is a predicted choice.
type Prediction struct {
	ID       string  `json:"id,omitempty"`
	RewardID string  `json:"reward_id,omitempty"`
	Score    float64 `json:"score,omitempty"`
}

// PredictRequest contains information about the prediction that Suggestionbox
// will make.
type PredictRequest struct {
	Inputs []Feature `json:"inputs,omitempty"`
}

// PredictResponse contains prediction choices.
type PredictResponse struct {
	Choices []Prediction `json:"choices,omitempty"`
}

// Predict asks a Model to make a prediction.
func (c *Client) Predict(ctx context.Context, modelID string, request PredictRequest) (PredictResponse, error) {
	var response struct {
		Success bool
		Error   string
		PredictResponse
	}
	u, err := url.Parse(c.addr + "/" + path.Join("suggestionbox", "models", modelID, "predict"))
	if err != nil {
		return response.PredictResponse, err
	}
	if !u.IsAbs() {
		return response.PredictResponse, errors.New("box address must be absolute")
	}
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return response.PredictResponse, errors.Wrap(err, "encoding request body")
	}
	req, err := http.NewRequest(http.MethodPost, u.String(), &buf)
	if err != nil {
		return response.PredictResponse, err
	}
	req = req.WithContext(ctx)
	req.Header.Set("Accept", "application/json; charset=utf-8")
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return response.PredictResponse, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return response.PredictResponse, errors.New(resp.Status)
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return response.PredictResponse, errors.Wrap(err, "decoding response")
	}
	if !response.Success {
		return response.PredictResponse, ErrSuggestionbox(response.Error)
	}
	return response.PredictResponse, nil
}
