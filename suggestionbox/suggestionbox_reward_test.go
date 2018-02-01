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

func TestReward(t *testing.T) {
	is := is.New(t)
	var apiCalls int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiCalls++
		is.Equal(r.Method, http.MethodPost)
		is.Equal(r.URL.Path, "/suggestionbox/models/model1/rewards")
		is.Equal(r.Header.Get("Accept"), "application/json; charset=utf-8")
		is.Equal(r.Header.Get("Content-Type"), "application/json; charset=utf-8")
		var reward suggestionbox.Reward
		is.NoErr(json.NewDecoder(r.Body).Decode(&reward))
		is.Equal(reward.RewardID, "reward1")
		is.Equal(reward.Value, float64(1))
		is.NoErr(json.NewEncoder(w).Encode(struct {
			Success bool `json:"success"`
		}{
			Success: true,
		}))
	}))
	defer srv.Close()
	sb := suggestionbox.New(srv.URL)
	reward := suggestionbox.Reward{
		RewardID: "reward1",
		Value:    1,
	}
	err := sb.Reward(context.Background(), "model1", reward)
	is.NoErr(err)
	is.Equal(apiCalls, 1) // apiCalls
}
