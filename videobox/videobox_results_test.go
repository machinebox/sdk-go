package videobox_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/machinebox/sdk-go/videobox"
	"github.com/matryer/is"
)

func TestResults(t *testing.T) {
	is := is.New(t)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		is.Equal(r.Method, "GET")
		is.Equal(r.URL.Path, "/videobox/results/5a50b8067eced76bad103c53dd0f5226")
		is.Equal(r.Header.Get("Accept"), "application/json; charset=utf-8")
		_, err := io.WriteString(w, resultsPayload)
		is.NoErr(err)
	}))
	defer srv.Close()
	vb := videobox.New(srv.URL)
	results, err := vb.Results("5a50b8067eced76bad103c53dd0f5226")
	is.NoErr(err)
	is.Equal(results.Facebox.FacesCount, 1)               // results.Facebox.FacesCount
	is.Equal(results.Tagbox.TagsCount, 3)                 // results.Tagbox.TagsCount
	is.Equal(len(results.Nudebox.Nudity[0].Instances), 5) // results.Nudebox.Nudity[0].InstancesCount
}

const resultsPayload = `{
	"success": true,
	"ready": true,
	"facebox": {
		"facesCount": 1,
		"faces": [
			{
				"key": "Unknown faces",
				"instancesCount": 3,
				"instances": [
					{
						"start": 24,
						"end": 144,
						"start_ms": 1000,
						"end_ms": 6006
					},
					{
						"start": 336,
						"end": 528,
						"start_ms": 14013,
						"end_ms": 22022
					},
					{
						"start": 720,
						"end": 720,
						"start_ms": 30029,
						"end_ms": 30029
					}
				]
			}
		],
		"errorsCount": 0
	},
	"tagbox": {
		"tagsCount": 3,
		"tags": [
			{
				"key": "candle",
				"instancesCount": 3,
				"instances": [
					{
						"start": 168,
						"end": 168,
						"start_ms": 7006,
						"end_ms": 7006
					},
					{
						"start": 216,
						"end": 216,
						"start_ms": 9009,
						"end_ms": 9009
					},
					{
						"start": 312,
						"end": 312,
						"start_ms": 13012,
						"end_ms": 13012
					}
				]
			},
			{
				"key": "crutch",
				"instancesCount": 1,
				"instances": [
					{
						"start": 504,
						"end": 504,
						"start_ms": 21021,
						"end_ms": 21021
					}
				]
			},
			{
				"key": "miniskirt",
				"instancesCount": 1,
				"instances": [
					{
						"start": 72,
						"end": 72,
						"start_ms": 3003,
						"end_ms": 3003
					}
				]
			}
		],
		"errorsCount": 0
	},
	"nudebox": {
		"nudity": [
			{
				"key": "greater than 0.5 chance of nuditiy",
				"instancesCount": 5,
				"instances": [
					{
						"start": 264,
						"end": 312,
						"start_ms": 11011,
						"end_ms": 13012
					},
					{
						"start": 360,
						"end": 360,
						"start_ms": 15014,
						"end_ms": 15014
					},
					{
						"start": 408,
						"end": 408,
						"start_ms": 17017,
						"end_ms": 17017
					},
					{
						"start": 456,
						"end": 552,
						"start_ms": 19019,
						"end_ms": 23023
					},
					{
						"start": 720,
						"end": 720,
						"start_ms": 30029,
						"end_ms": 30029
					}
				]
			}
		],
		"errorsCount": 0
	}
}`
