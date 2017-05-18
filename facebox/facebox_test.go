package facebox_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/machinebox/sdk-go/facebox"
	"github.com/matryer/is"
)

func TestInfo(t *testing.T) {
	is := is.New(t)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		is.Equal(r.Method, "GET")
		is.Equal(r.URL.Path, "/info")
		is.Equal(r.Header.Get("Accept"), "application/json; charset=utf-8")
		io.WriteString(w, `{
			"name": "facebox",
			"version": 1,
			"build": "abcdefg",
			"status": "ready"
		}`)
	}))
	defer srv.Close()
	fb := facebox.New(srv.URL)
	info, err := fb.Info()
	is.NoErr(err)
	is.Equal(info.Name, "facebox")
	is.Equal(info.Version, 1)
	is.Equal(info.Build, "abcdefg")
	is.Equal(info.Status, "ready")
}
