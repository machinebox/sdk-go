package classificationbox_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/machinebox/sdk-go/classificationbox"
	"github.com/matryer/is"
)

func TestCreateModel(t *testing.T) {
	is := is.New(t)
	var apiCalls int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiCalls++
		is.Equal(r.Method, http.MethodPost)
		is.Equal(r.URL.Path, "/classificationbox/models")
		is.Equal(r.Header.Get("Accept"), "application/json; charset=utf-8")
		is.Equal(r.Header.Get("Content-Type"), "application/json; charset=utf-8")
		var model classificationbox.Model
		is.NoErr(json.NewDecoder(r.Body).Decode(&model))
		is.Equal(model.Name, "My Model")
		model.ID = "1"
		is.NoErr(json.NewEncoder(w).Encode(struct {
			classificationbox.Model
			Success bool `json:"success"`
		}{
			Success: true,
			Model:   model,
		}))
	}))
	defer srv.Close()
	cb := classificationbox.New(srv.URL)
	inModel := classificationbox.NewModel("", "My Model", "class1", "class2", "class3")
	outModel, err := cb.CreateModel(context.Background(), inModel)
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
		is.Equal(r.URL.Path, "/classificationbox/models/model1")
		is.Equal(r.Header.Get("Accept"), "application/json; charset=utf-8")
		model := classificationbox.Model{
			ID: "model1",
		}
		is.NoErr(json.NewEncoder(w).Encode(struct {
			classificationbox.Model
			Success bool `json:"success"`
		}{
			Success: true,
			Model:   model,
		}))
	}))
	defer srv.Close()
	cb := classificationbox.New(srv.URL)
	outModel, err := cb.GetModel(context.Background(), "model1")
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
		is.Equal(r.URL.Path, "/classificationbox/models")
		is.Equal(r.Header.Get("Accept"), "application/json; charset=utf-8")
		model := classificationbox.Model{
			ID: "model1",
		}
		is.NoErr(json.NewEncoder(w).Encode(struct {
			Success bool `json:"success"`
			Models  []classificationbox.Model
		}{
			Success: true,
			Models:  []classificationbox.Model{model, model, model},
		}))
	}))
	defer srv.Close()
	cb := classificationbox.New(srv.URL)
	models, err := cb.ListModels(context.Background())
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
		is.Equal(r.URL.Path, "/classificationbox/models/model1")
		is.NoErr(json.NewEncoder(w).Encode(struct {
			Success bool `json:"success"`
		}{
			Success: true,
		}))
	}))
	defer srv.Close()
	cb := classificationbox.New(srv.URL)
	err := cb.DeleteModel(context.Background(), "model1")
	is.NoErr(err)
	is.Equal(apiCalls, 1) // apiCalls
}

func TestFeatureHelpers(t *testing.T) {
	is := is.New(t)

	var f classificationbox.Feature

	f = classificationbox.FeatureNumber("age", 20)
	is.Equal(f.Type, "number")
	is.Equal(f.Key, "age")
	is.Equal(f.Value, "20")

	f = classificationbox.FeatureText("title", "Machine box releases new box")
	is.Equal(f.Type, "text")
	is.Equal(f.Key, "title")
	is.Equal(f.Value, "Machine box releases new box")

	f = classificationbox.FeatureKeyword("city", "New York City")
	is.Equal(f.Type, "keyword")
	is.Equal(f.Key, "city")
	is.Equal(f.Value, "New York City")

	f = classificationbox.FeatureList("categories", "one", "two", "three")
	is.Equal(f.Type, "list")
	is.Equal(f.Key, "categories")
	is.Equal(f.Value, "one,two,three")

	f = classificationbox.FeatureImageURL("pic", "http://url.com/path/to/pic.jpg")
	is.Equal(f.Type, "image_url")
	is.Equal(f.Key, "pic")
	is.Equal(f.Value, "http://url.com/path/to/pic.jpg")

	f = classificationbox.FeatureImageBase64("pic", "pretendthisisimagedata")
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
		is.Equal(r.URL.Path, "/classificationbox/models/model1/stats")
		is.Equal(r.Header.Get("Accept"), "application/json; charset=utf-8")
		stats := classificationbox.ModelStats{
			Predictions: 50,
			Examples:    100,
			Classes: []classificationbox.ClassStats{
				{
					Name: "class1", Examples: 100,
				},
				{
					Name: "class2", Examples: 98,
				},
				{
					Name: "class3", Examples: 102,
				},
			},
		}
		is.NoErr(json.NewEncoder(w).Encode(struct {
			classificationbox.ModelStats
			Success bool `json:"success"`
		}{
			Success:    true,
			ModelStats: stats,
		}))
	}))
	defer srv.Close()
	cb := classificationbox.New(srv.URL)
	stats, err := cb.GetModelStats(context.Background(), "model1")
	is.NoErr(err)
	is.Equal(apiCalls, 1) // apiCalls
	is.Equal(stats.Predictions, 50)
	is.Equal(stats.Examples, 100)
	is.Equal(len(stats.Classes), 3)
}
