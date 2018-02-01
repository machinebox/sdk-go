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
	inModel := suggestionbox.Model{
		Name: "My Model",
		Choices: []suggestionbox.Choice{
			suggestionbox.NewChoice("choice1", suggestionbox.FeatureText("title", "Machine Box releases new product")),
		},
	}
	outModel, err := sb.CreateModel(context.Background(), inModel)
	is.NoErr(err)
	is.Equal(apiCalls, 1)      // apiCalls
	is.Equal(outModel.ID, "1") // outModel.ID
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
