package suggestionbox

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/pkg/errors"
)

// Model represents a single model inside Suggestionbox.
type Model struct {
	// ID is the ID of the model.
	ID string `json:"id,omitempty"`
	// Name is the human readable name of the Model.
	Name string `json:"name"`
	// Options are optional Model settings to adjust the behaviour
	// of this Model within Suggestionbox.
	Options *ModelOptions `json:"options,omitempty"`
	// Choices are the options this Model will select from.
	Choices []Choice `json:"choices,omitempty"`
}

// NewModel makes a new Model.
func NewModel(id, name string, choices ...Choice) Model {
	return Model{
		ID:      id,
		Name:    name,
		Choices: choices,
	}
}

// Feature represents a single feature, to describe an input or a choice
// for example age:28 or location:"London".
type Feature struct {
	// Key is the name of the Feature.
	Key string `json:"key"`
	// Value is the string value of this Feature.
	Value string `json:"value"`
	// Type is the type of the Feature.
	// Can be "number", "text", "keyword", "list", "image_url" or "image_base64"..
	Type string `json:"type"`
}

// Choice is an option with features.
type Choice struct {
	// ID is a unique ID for this choice.
	ID string `json:"id"`
	// Features holds all the Feature objects that describe
	// this choice.
	Features []Feature `json:"features,omitempty"`
}

// NewChoice creates a new Choice.
func NewChoice(id string, features ...Feature) Choice {
	return Choice{
		ID:       id,
		Features: features,
	}
}

// ModelOptions describes the behaviours of a Model.
type ModelOptions struct {
	// RewardExpirationSeconds is the number of seconds to wait for
	// the reward before it expires.
	RewardExpirationSeconds int `json:"reward_expiration_seconds,omitempty"`

	// Epsilon enables proportionate exploiting vs exploring ratio.
	Epsilon float64 `json:"epsilon,omitempty"`

	// SoftmaxLambda enables adaptive exploiting vs exploring ratio.
	SoftmaxLambda float64 `json:"softmax_lambda,omitempty"`

	// Ngrams describes the n-grams for text analysis.
	Ngrams int `json:"ngrams,omitempty"`
	// Skipgrams describes the skip-grams for the text analysis.
	Skipgrams int `json:"skipgrams,omitempty"`
}

// ListModels gets a Model by its ID.
func (c *Client) ListModels(ctx context.Context) ([]Model, error) {
	u, err := url.Parse(c.addr + "/suggestionbox/models")
	if err != nil {
		return nil, err
	}
	if !u.IsAbs() {
		return nil, errors.New("box address must be absolute")
	}
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	req.Header.Set("Accept", "application/json; charset=utf-8")
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, errors.New(resp.Status)
	}
	var response struct {
		Success bool
		Error   string
		Models  []Model
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, errors.Wrap(err, "decoding response")
	}
	if !response.Success {
		return nil, ErrSuggestionbox(response.Error)
	}
	return response.Models, nil
}

// GetModel gets a Model by its ID.
func (c *Client) GetModel(ctx context.Context, modelID string) (Model, error) {
	var model Model
	u, err := url.Parse(c.addr + "/" + path.Join("suggestionbox", "models", modelID))
	if err != nil {
		return model, err
	}
	if !u.IsAbs() {
		return model, errors.New("box address must be absolute")
	}
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return model, err
	}
	req = req.WithContext(ctx)
	req.Header.Set("Accept", "application/json; charset=utf-8")
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

// DeleteModel gets a Model by its ID.
func (c *Client) DeleteModel(ctx context.Context, modelID string) error {
	u, err := url.Parse(c.addr + "/" + path.Join("suggestionbox", "models", modelID))
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
	req = req.WithContext(ctx)
	req.Header.Set("Accept", "application/json; charset=utf-8")
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.New(resp.Status)
	}
	var response struct {
		Success bool
		Error   string
		Model
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return errors.Wrap(err, "decoding response")
	}
	if !response.Success {
		return ErrSuggestionbox(response.Error)
	}
	return nil
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
