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
			{
				ID: "choice1",
				Features: []suggestionbox.Feature{
					{
						Key:   "title",
						Type:  "text",
						Value: "Machine Box releases new product",
					},
				},
			},
		},
	}
	outModel, err := sb.CreateModel(context.Background(), inModel)
	is.NoErr(err)
	is.Equal(apiCalls, 1)      // apiCalls
	is.Equal(outModel.ID, "1") // outModel.ID
}
