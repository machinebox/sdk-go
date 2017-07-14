package facebox_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/machinebox/sdk-go/facebox"
	"github.com/matryer/is"
)

func TestRename(t *testing.T) {
	is := is.New(t)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		is.Equal(r.URL.Path, "/facebox/teach/john1.jpg")
		is.Equal(r.Header.Get("Accept"), "application/json; charset=utf-8")
		is.Equal(r.Method, "PATCH")
		is.Equal(r.FormValue("name"), "Sir John L")
		io.WriteString(w, `{
			"success": true
		}`)
	}))
	defer srv.Close()

	fb := facebox.New(srv.URL)
	err := fb.Rename("john1.jpg", "Sir John L")
	is.NoErr(err)
}

func TestRenameAll(t *testing.T) {
	is := is.New(t)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		is.Equal(r.URL.Path, "/facebox/rename")
		is.Equal(r.Header.Get("Accept"), "application/json; charset=utf-8")
		is.Equal(r.Method, "POST")
		is.Equal(r.FormValue("from"), "John Lennon")
		is.Equal(r.FormValue("to"), "Sir John Lennon")
		io.WriteString(w, `{
			"success": true
		}`)
	}))
	defer srv.Close()

	fb := facebox.New(srv.URL)
	err := fb.RenameAll("John Lennon", "Sir John Lennon")
	is.NoErr(err)
}
