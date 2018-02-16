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

func TestPredict(t *testing.T) {
	is := is.New(t)
	var apiCalls int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiCalls++
		is.Equal(r.Method, http.MethodPost)
		is.Equal(r.URL.Path, "/classificationbox/models/model1/predict")
		is.Equal(r.Header.Get("Accept"), "application/json; charset=utf-8")
		is.Equal(r.Header.Get("Content-Type"), "application/json; charset=utf-8")
		var req classificationbox.PredictRequest
		is.NoErr(json.NewDecoder(r.Body).Decode(&req))
		is.Equal(req.Inputs[0].Key, "title")
		resp := classificationbox.PredictResponse{
			Classes: []classificationbox.Class{
				{
					ID:    "choice1",
					Score: 0.7,
				},
				{
					ID:    "choice2",
					Score: 0.2,
				},
				{
					ID:    "choice3",
					Score: 0.1,
				},
			},
		}
		is.NoErr(json.NewEncoder(w).Encode(struct {
			classificationbox.PredictResponse
			Success bool `json:"success"`
		}{
			Success:         true,
			PredictResponse: resp,
		}))
	}))
	defer srv.Close()
	cb := classificationbox.New(srv.URL)
	predictReq := classificationbox.PredictRequest{
		Inputs: []classificationbox.Feature{
			{
				Key:   "title",
				Type:  "text",
				Value: "Machine Box releases new product",
			},
		},
	}
	predictResp, err := cb.Predict(context.Background(), "model1", predictReq)
	is.NoErr(err)
	is.Equal(apiCalls, 1)                 // apiCalls
	is.Equal(len(predictResp.Classes), 3) // len(predictResp.Choices)
}
