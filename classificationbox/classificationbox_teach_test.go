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

func TestTeach(t *testing.T) {
	is := is.New(t)
	var apiCalls int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiCalls++
		is.Equal(r.Method, http.MethodPost)
		is.Equal(r.URL.Path, "/classificationbox/models/model1/teach")
		is.Equal(r.Header.Get("Accept"), "application/json; charset=utf-8")
		is.Equal(r.Header.Get("Content-Type"), "application/json; charset=utf-8")
		var req classificationbox.Example
		is.NoErr(json.NewDecoder(r.Body).Decode(&req))
		is.Equal(req.Inputs[0].Key, "title")
		is.NoErr(json.NewEncoder(w).Encode(struct {
			Success bool `json:"success"`
		}{
			Success: true,
		}))
	}))
	defer srv.Close()
	cb := classificationbox.New(srv.URL)
	example := classificationbox.Example{
		Inputs: []classificationbox.Feature{
			{
				Key:   "title",
				Type:  "text",
				Value: "Machine Box releases new product",
			},
		},
	}
	err := cb.Teach(context.Background(), "model1", example)
	is.NoErr(err)
	is.Equal(apiCalls, 1) // apiCalls
}
