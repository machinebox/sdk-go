package videobox

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

// VideoAnalysis describes the results of a video processing operation.
type VideoAnalysis struct {
	// Ready indicates whether the results are ready or not.
	Ready   bool     `json:"ready"`
	Facebox *Facebox `json:"facebox,omitempty"`
	Tagbox  *Tagbox  `json:"tagbox,omitempty"`
	Nudebox *Nudebox `json:"nudebox,omitempty"`
}

// Facebox holds box specific results.
type Facebox struct {
	Faces      []Item `json:"faces"`
	ErrorCount int    `json:"errorsCount"`
	LastErr    string `json:"lastError,omitempty"`
}

// Tagbox holds box specific results.
type Tagbox struct {
	Tags       []Item `json:"tags"`
	ErrorCount int    `json:"errorsCount"`
	LastErr    string `json:"lastError,omitempty"`
}

// Nudebox holds box specific results.
type Nudebox struct {
	Nudity     []Item `json:"nudity"`
	ErrorCount int    `json:"errorsCount"`
	LastErr    string `json:"lastError,omitempty"`
}

// Item describes a single entity that was discovered at
// one or many instances throughout the video.
type Item struct {
	// Key is a string describing the item.
	// For Facebox it will be the name.
	// For Tagbox it will be the tag.
	// For Nudebox it will be a description of the nudity detected.
	Key string `json:"key"`
	// Instances holds the time and frame ranges where this Item
	// appears in the video.
	Instances []Range `json:"instances"`
}

// Range describes a period of time within the video.
type Range struct {
	// Start is the start frame.
	Start int `json:"start"`
	// End is the end frame.
	End int `json:"end"`
	// StartMS is the start time in milliseconds.
	StartMS int `json:"start_ms"`
	// EndMS is the end time in milliseconds.
	EndMS int `json:"end_ms"`
	// Confidence is the maximum confidence of any instance in this
	// range.
	Confidence float64 `json:"confidence,omitempty"`
}

// Results gets the results of a video processing operation.
// This should be called after the Video.Status is StatusCompleted.
func (c *Client) Results(id string) (*VideoAnalysis, error) {
	u, err := url.Parse(c.addr + "/videobox/results/" + id)
	if err != nil {
		return nil, err
	}
	if !u.IsAbs() {
		return nil, errors.New("box address must be absolute")
	}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json; charset=utf-8")
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, errors.New(resp.Status)
	}
	var results VideoAnalysis
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, err
	}
	return &results, nil
}
