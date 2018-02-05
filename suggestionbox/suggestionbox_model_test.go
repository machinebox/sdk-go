package suggestionbox_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/machinebox/sdk-go/suggestionbox"
	"github.com/matryer/is"
)

func TestCreateModel(t *testing.T) {
	is := is.New(t)
	var apiCalls int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiCalls++
		is.Equal(r.Method, http.MethodPost)
		is.Equal(r.URL.Path, "/suggestionbox/models")
		is.Equal(r.Header.Get("Accept"), "application/json; charset=utf-8")
		is.Equal(r.Header.Get("Content-Type"), "application/json; charset=utf-8")
		var model suggestionbox.Model
		is.NoErr(json.NewDecoder(r.Body).Decode(&model))
		is.Equal(model.Name, "My Model")
		model.ID = "1"
		is.NoErr(json.NewEncoder(w).Encode(struct {
			suggestionbox.Model
			Success bool `json:"success"`
		}{
			Success: true,
			Model:   model,
		}))
	}))
	defer srv.Close()
	sb := suggestionbox.New(srv.URL)
	inModel := suggestionbox.NewModel("", "My Model",
		suggestionbox.NewChoice("choice1",
			suggestionbox.FeatureText("title", "Machine Box releases new product"),
		),
	)
	outModel, err := sb.CreateModel(context.Background(), inModel)
	is.NoErr(err)
	is.Equal(apiCalls, 1)      // apiCalls
	is.Equal(outModel.ID, "1") // outModel.ID
}

func TestGetModel(t *testing.T) {
	is := is.New(t)
	var apiCalls int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiCalls++
		is.Equal(r.Method, http.MethodGet)
		is.Equal(r.URL.Path, "/suggestionbox/models/model1")
		is.Equal(r.Header.Get("Accept"), "application/json; charset=utf-8")
		model := suggestionbox.Model{
			ID: "model1",
		}
		is.NoErr(json.NewEncoder(w).Encode(struct {
			suggestionbox.Model
			Success bool `json:"success"`
		}{
			Success: true,
			Model:   model,
		}))
	}))
	defer srv.Close()
	sb := suggestionbox.New(srv.URL)
	outModel, err := sb.GetModel(context.Background(), "model1")
	is.NoErr(err)
	is.Equal(apiCalls, 1)           // apiCalls
	is.Equal(outModel.ID, "model1") // outModel.ID
}

func TestListModels(t *testing.T) {
	is := is.New(t)
	var apiCalls int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiCalls++
		is.Equal(r.Method, http.MethodGet)
		is.Equal(r.URL.Path, "/suggestionbox/models")
		is.Equal(r.Header.Get("Accept"), "application/json; charset=utf-8")
		model := suggestionbox.Model{
			ID: "model1",
		}
		is.NoErr(json.NewEncoder(w).Encode(struct {
			Success bool `json:"success"`
			Models  []suggestionbox.Model
		}{
			Success: true,
			Models:  []suggestionbox.Model{model, model, model},
		}))
	}))
	defer srv.Close()
	sb := suggestionbox.New(srv.URL)
	models, err := sb.ListModels(context.Background())
	is.NoErr(err)
	is.Equal(apiCalls, 1)    // apiCalls
	is.Equal(len(models), 3) // len(models)
}

func TestDeleteModel(t *testing.T) {
	is := is.New(t)
	var apiCalls int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiCalls++
		is.Equal(r.Method, http.MethodDelete)
		is.Equal(r.URL.Path, "/suggestionbox/models/model1")
		is.NoErr(json.NewEncoder(w).Encode(struct {
			Success bool `json:"success"`
		}{
			Success: true,
		}))
	}))
	defer srv.Close()
	sb := suggestionbox.New(srv.URL)
	err := sb.DeleteModel(context.Background(), "model1")
	is.NoErr(err)
	is.Equal(apiCalls, 1) // apiCalls
}

func TestNewModel(t *testing.T) {
	is := is.New(t)

	m := suggestionbox.NewModel("model1", "My Model",
		suggestionbox.NewChoice("choice1",
			suggestionbox.FeatureKeyword("city", "New York City"),
		),
	)
	is.Equal(m.ID, "model1")
	is.Equal(m.Name, "My Model")
	is.Equal(len(m.Choices), 1)
	is.Equal(m.Choices[0].ID, "choice1")
}

func TestNewChoice(t *testing.T) {
	is := is.New(t)

	c := suggestionbox.NewChoice("choice1",
		suggestionbox.FeatureKeyword("city", "New York City"),
	)
	is.Equal(c.ID, "choice1")
	is.Equal(len(c.Features), 1)
	is.Equal(c.Features[0].Key, "city")
}

func TestFeatureHelpers(t *testing.T) {
	is := is.New(t)

	var f suggestionbox.Feature

	f = suggestionbox.FeatureNumber("age", 20)
	is.Equal(f.Type, "number")
	is.Equal(f.Key, "age")
	is.Equal(f.Value, "20")

	f = suggestionbox.FeatureText("title", "Machine box releases new box")
	is.Equal(f.Type, "text")
	is.Equal(f.Key, "title")
	is.Equal(f.Value, "Machine box releases new box")

	f = suggestionbox.FeatureKeyword("city", "New York City")
	is.Equal(f.Type, "keyword")
	is.Equal(f.Key, "city")
	is.Equal(f.Value, "New York City")

	f = suggestionbox.FeatureList("categories", "one", "two", "three")
	is.Equal(f.Type, "list")
	is.Equal(f.Key, "categories")
	is.Equal(f.Value, "one,two,three")

	f = suggestionbox.FeatureImageURL("pic", "http://url.com/path/to/pic.jpg")
	is.Equal(f.Type, "image_url")
	is.Equal(f.Key, "pic")
	is.Equal(f.Value, "http://url.com/path/to/pic.jpg")

	f = suggestionbox.FeatureImageBase64("pic", "pretendthisisimagedata")
	is.Equal(f.Type, "image_base64")
	is.Equal(f.Key, "pic")
	is.Equal(f.Value, "pretendthisisimagedata")

}

func TestGetModelStats(t *testing.T) {
	is := is.New(t)
	var apiCalls int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiCalls++
		is.Equal(r.Method, http.MethodGet)
		is.Equal(r.URL.Path, "/suggestionbox/models/model1/stats")
		is.Equal(r.Header.Get("Accept"), "application/json; charset=utf-8")
		stats := suggestionbox.ModelStats{
			Predictions:  1,
			Rewards:      2,
			RewardRatio:  3.3,
			Explores:     4,
			Exploits:     5,
			ExploreRatio: 6.6,
		}
		is.NoErr(json.NewEncoder(w).Encode(struct {
			suggestionbox.ModelStats
			Success bool `json:"success"`
		}{
			Success:    true,
			ModelStats: stats,
		}))
	}))
	defer srv.Close()
	sb := suggestionbox.New(srv.URL)
	stats, err := sb.GetModelStats(context.Background(), "model1")
	is.NoErr(err)
	is.Equal(apiCalls, 1) // apiCalls
	is.Equal(stats.Predictions, 1)
	is.Equal(stats.Rewards, 2)
	is.Equal(stats.RewardRatio, 3.3)
	is.Equal(stats.Explores, 4)
	is.Equal(stats.Exploits, 5)
	is.Equal(stats.ExploreRatio, 6.6)
}
