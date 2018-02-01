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

func TestPredict(t *testing.T) {
	is := is.New(t)
	var apiCalls int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiCalls++
		is.Equal(r.Method, http.MethodPost)
		is.Equal(r.URL.Path, "/suggestionbox/models/model1/predict")
		is.Equal(r.Header.Get("Accept"), "application/json; charset=utf-8")
		is.Equal(r.Header.Get("Content-Type"), "application/json; charset=utf-8")
		var req suggestionbox.PredictRequest
		is.NoErr(json.NewDecoder(r.Body).Decode(&req))
		is.Equal(req.Inputs[0].Key, "title")
		resp := suggestionbox.PredictResponse{
			Choices: []suggestionbox.Prediction{
				{
					ID:       "choice1",
					RewardID: "reward1",
					Score:    0.7,
				},
				{
					ID:       "choice2",
					RewardID: "reward2",
					Score:    0.2,
				},
				{
					ID:       "choice3",
					RewardID: "reward3",
					Score:    0.1,
				},
			},
		}
		is.NoErr(json.NewEncoder(w).Encode(struct {
			suggestionbox.PredictResponse
			Success bool `json:"success"`
		}{
			Success:         true,
			PredictResponse: resp,
		}))
	}))
	defer srv.Close()
	sb := suggestionbox.New(srv.URL)
	predictReq := suggestionbox.PredictRequest{
		Inputs: []suggestionbox.Feature{
			{
				Key:   "title",
				Type:  "text",
				Value: "Machine Box releases new product",
			},
		},
	}
	predictResp, err := sb.Predict(context.Background(), "model1", predictReq)
	is.NoErr(err)
	is.Equal(apiCalls, 1)                 // apiCalls
	is.Equal(len(predictResp.Choices), 3) // len(predictResp.Choices)
}
