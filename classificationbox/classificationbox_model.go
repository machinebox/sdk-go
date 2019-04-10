package classificationbox

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

// Model represents a single model inside Classificationbox.
type Model struct {
	// ID is the ID of the model.
	ID string `json:"id,omitempty"`
	// Name is the human readable name of the Model.
	Name string `json:"name"`
	// Options are optional Model settings to adjust the behaviour
	// of this Model within Classificationbox.
	Options *ModelOptions `json:"options,omitempty"`
	// Classes are the classes that this model can learn.
	Classes []string `json:"classes,omitempty"`
}

// NewModel creates a new Model.
func NewModel(id, name string, classes ...string) Model {
	return Model{
		ID:      id,
		Name:    name,
		Classes: classes,
	}
}

// ModelOptions describes the behaviours of a Model.
type ModelOptions struct {
	// Ngrams describes the n-grams for text analysis.
	Ngrams int `json:"ngrams,omitempty"`
	// Skipgrams describes the skip-grams for the text analysis.
	Skipgrams int `json:"skipgrams,omitempty"`
}

// Feature represents a single feature, to describe an input.
type Feature struct {
	// Key is the name of the Feature.
	Key string `json:"key"`
	// Value is the string value of this Feature.
	Value string `json:"value"`
	// Type is the type of the Feature.
	// Can be "number", "text", "keyword", "list", "image_url" or "image_base64"..
	Type string `json:"type"`
}

// CreateModel creates the Model in Classificationbox.
// If no ID is set, one will be assigned in the return Model.
func (c *Client) CreateModel(ctx context.Context, model Model) (Model, error) {
	u, err := url.Parse(c.addr + "/classificationbox/models")
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
	var outModel Model
	_, err = c.client.DoUnmarshal(req, &outModel)
	if err != nil {
		return outModel, err
	}
	return outModel, nil
}

// ListModels gets all models.
func (c *Client) ListModels(ctx context.Context) ([]Model, error) {
	u, err := url.Parse(c.addr + "/classificationbox/models")
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
	var response struct {
		Models []Model
	}
	_, err = c.client.DoUnmarshal(req, &response)
	if err != nil {
		return nil, err
	}
	return response.Models, nil
}

// GetModel gets a Model by its ID.
func (c *Client) GetModel(ctx context.Context, modelID string) (Model, error) {
	var model Model
	u, err := url.Parse(c.addr + "/" + path.Join("classificationbox", "models", modelID))
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
	_, err = c.client.DoUnmarshal(req, &model)
	if err != nil {
		return model, err
	}
	return model, nil
}

// DeleteModel deletes a Model by its ID.
func (c *Client) DeleteModel(ctx context.Context, modelID string) error {
	u, err := url.Parse(c.addr + "/" + path.Join("classificationbox", "models", modelID))
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
	_, err = c.client.DoUnmarshal(req, nil)
	if err != nil {
		return err
	}
	return nil
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

// FeatureImageBase64 makes a Feature that is a base64 encoded image.
func FeatureImageBase64(key string, data string) Feature {
	return Feature{
		Type:  "image_base64",
		Key:   key,
		Value: data,
	}
}

// ModelStats are the statistics for a Model.
type ModelStats struct {
	// Predictions is the number of predictions this model has made.
	Predictions int
	// Examples is the total number of examples this model has
	// been taught.
	Examples int
	// Classes is a list of statistics per class.
	Classes []ClassStats
}

// ClassStats contains per-class statistics.
type ClassStats struct {
	// Name is the name of the class.
	Name string
	// Examples is the number of examples of this class that
	// the model has been taught.
	Examples int
}

// GetModelStats gets the statistics for the specified model.
func (c *Client) GetModelStats(ctx context.Context, modelID string) (ModelStats, error) {
	var stats ModelStats
	u, err := url.Parse(c.addr + "/" + path.Join("classificationbox", "models", modelID, "stats"))
	if err != nil {
		return stats, err
	}
	if !u.IsAbs() {
		return stats, errors.New("box address must be absolute")
	}
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return stats, err
	}
	req = req.WithContext(ctx)
	req.Header.Set("Accept", "application/json; charset=utf-8")
	_, err = c.client.DoUnmarshal(req, &stats)
	if err != nil {
		return stats, err
	}
	return stats, nil
}
