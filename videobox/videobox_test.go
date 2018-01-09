package videobox_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/machinebox/sdk-go/videobox"
	"github.com/matryer/is"
)

func TestInfo(t *testing.T) {
	is := is.New(t)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		is.Equal(r.Method, "GET")
		is.Equal(r.URL.Path, "/info")
		is.Equal(r.Header.Get("Accept"), "application/json; charset=utf-8")
		_, err := io.WriteString(w, `{
			"name": "textbox",
			"version": 1,
			"build": "abcdefg",
			"status": "ready"
		}`)
		is.NoErr(err)
	}))
	defer srv.Close()
	b := videobox.New(srv.URL)
	info, err := b.Info()
	is.NoErr(err)
	is.Equal(info.Name, "textbox")
	is.Equal(info.Version, 1)
	is.Equal(info.Build, "abcdefg")
	is.Equal(info.Status, "ready")
}
