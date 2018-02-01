package suggestionbox

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// Model represents a single model inside Suggestionbox.
type Model struct {
	// ID is the ID of the model.
	ID string `json:"id,omitempty"`
	// Name is the human readable name of the Model.
	Name string `json:"name,omitempty"`
	// Options are optional Model settings to adjust the behaviour
	// of this Model within Suggestionbox.
	Options *ModelOptions `json:"options,omitempty"`
	// Choices are the options this Model will select from.
	Choices []Choice `json:"choices,omitempty"`
}

// Feature represents a single feature, to describe an input or a choice
// for example age:28 or location:"London".
type Feature struct {
	// Key is the name of the Feature.
	Key string `json:"key,omitempty"`
	// Value is the string value of this Feature.
	Value string `json:"value,omitempty"`
	// Type is the type of the Feature.
	// Can be "number", "text", "keyword", "list", "image_url" or "image_base64"..
	Type string `json:"type,omitempty"`
}

// Choice is an option with features.
type Choice struct {
	// ID is a unique ID for this choice.
	ID string `json:"id,omitempty"`
	// Features holds all the Feature objects that describe
	// this choice.
	Features []Feature `json:"features,omitempty"`
}

// ModelOptions describes the behaviours of a Model.
type ModelOptions struct {
	// Expiration is the time to wait for the reward before it expires.
	Expiration time.Duration `json:"expiration,omitempty"`

	// Epsilon enables proportionate exploiting vs exploring ratio.
	Epsilon float64 `json:"epsilon,omitempty"`

	// SoftmaxLambda enables adaptive exploiting vs exploring ratio.
	SoftmaxLambda float64 `json:"softmax_lambda,omitempty"`

	// Ngrams describes the n-grams for text analysis.
	Ngrams int `json:"ngrams,omitempty"`
	// Skipgrams describes the skip-grams for the text analysis.
	Skipgrams int `json:"skipgrams,omitempty"`
}

// CreateModel creates the Model in Suggestionbox.
// If no ID is set, one will be assigned in the return Model.
func (c *Client) CreateModel(ctx context.Context, model Model) (Model, error) {
	u, err := url.Parse(c.addr + "/suggestionbox/models")
	if err != nil {
		return model, err
	}
	if !u.IsAbs() {
		return model, errors.New("box address must be absolute")
	}
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(model); err != nil {
		return model, errors.Wrap(err, "encoding request body")
	}
	req, err := http.NewRequest(http.MethodPost, u.String(), &buf)
	if err != nil {
		return model, err
	}
	req = req.WithContext(ctx)
	req.Header.Set("Accept", "application/json; charset=utf-8")
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return model, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return model, errors.New(resp.Status)
	}
	var response struct {
		Success bool
		Error   string
		Model
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return model, errors.Wrap(err, "decoding response")
	}
	if !response.Success {
		return model, ErrSuggestionbox(response.Error)
	}
	return response.Model, nil
}

// FeatureNumber makes a numerical Feature.
func FeatureNumber(key string, value float64) Feature {
	return Feature{
		Type:  "number",
		Key:   key,
		Value: fmt.Sprintf("%v", value),
	}
}

// FeatureText makes a textual Feature that will be tokenized.
// Use FeatureKeyword for values that should not be tokenized.
func FeatureText(key string, text string) Feature {
	return Feature{
		Type:  "text",
		Key:   key,
		Value: text,
	}
}

// FeatureKeyword makes a textual Feature that will not be tokenized.
// Use FeatureList to provide multiple keywords in a single Feature.
// Use Text for bodies of text that should be tokenized.
func FeatureKeyword(key string, keyword string) Feature {
	return Feature{
		Type:  "keyword",
		Key:   key,
		Value: keyword,
	}
}

// FeatureList makes a Feature made up of multiple keywords.
func FeatureList(key string, keywords ...string) Feature {
	return Feature{
		Type:  "list",
		Key:   key,
		Value: strings.Join(keywords, ","),
	}
}

// FeatureImageURL makes a Feature that points to a hosted image.
func FeatureImageURL(key string, url string) Feature {
	return Feature{
		Type:  "image_url",
		Key:   key,
		Value: url,
	}
}

// FeatureImageBase64 makes a Feature that is base 64 encoded.
func FeatureImageBase64(key string, data string) Feature {
	return Feature{
		Type:  "image_base64",
		Key:   key,
		Value: data,
	}
}
