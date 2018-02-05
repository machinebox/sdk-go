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
	// ID is the choice identifier being predicted.
	ID string `json:"id,omitempty"`
	// RewardID is the ID of the reward that should be sent if this
	// prediction was successful.
	RewardID string `json:"reward_id,omitempty"`
	// Score is a numerical value indicating how this prediction relates
	// to other predictions.
	Score float64 `json:"score,omitempty"`
}

// PredictRequest contains information about the prediction that Suggestionbox
// will make.
type PredictRequest struct {
	// Inputs is a list of Feature objects that will be used when
	// making the prediction.
	Inputs []Feature `json:"inputs,omitempty"`
}

// PredictResponse contains prediction choices.
type PredictResponse struct {
	// Choices contains the predictions.
	Choices []Prediction `json:"choices,omitempty"`
}

// Predict asks a Model to make a prediction.
func (c *Client) Predict(ctx context.Context, modelID string, request PredictRequest) (PredictResponse, error) {
	var response PredictResponse
	u, err := url.Parse(c.addr + "/" + path.Join("suggestionbox", "models", modelID, "predict"))
	if err != nil {
		return response, err
	}
	if !u.IsAbs() {
		return response, errors.New("box address must be absolute")
	}
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return response, errors.Wrap(err, "encoding request body")
	}
	req, err := http.NewRequest(http.MethodPost, u.String(), &buf)
	if err != nil {
		return response, err
	}
	req = req.WithContext(ctx)
	req.Header.Set("Accept", "application/json; charset=utf-8")
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	_, err = c.client.Do(req, &response)
	if err != nil {
		return response, err
	}
	return response, nil
}
