package videobox_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/machinebox/sdk-go/videobox"
	"github.com/matryer/is"
)

func TestStatus(t *testing.T) {
	is := is.New(t)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		is.Equal(r.Method, "GET")
		is.Equal(r.URL.Path, "/videobox/status/5a50b8067eced76bad103c53dd0f5226")
		is.Equal(r.Header.Get("Accept"), "application/json; charset=utf-8")
		_, err := io.WriteString(w, `{
			"success": true,
			"id": "5a50b8067eced76bad103c53dd0f5226",
			"status": "processing"
		}`)
		is.NoErr(err)
	}))
	defer srv.Close()
	vb := videobox.New(srv.URL)
	video, err := vb.Status("5a50b8067eced76bad103c53dd0f5226")
	is.NoErr(err)
	is.Equal(video.ID, "5a50b8067eced76bad103c53dd0f5226")
	is.Equal(video.Status, videobox.StatusProcessing)
}
