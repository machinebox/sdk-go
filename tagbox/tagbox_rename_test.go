package tagbox_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/machinebox/sdk-go/tagbox"
	"github.com/matryer/is"
)

func TestRename(t *testing.T) {
	is := is.New(t)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		is.Equal(r.URL.Path, "/tagbox/teach/image1.jpg")
		is.Equal(r.Header.Get("Accept"), "application/json; charset=utf-8")
		is.Equal(r.Method, "PATCH")
		is.Equal(r.FormValue("tag"), "monkeys")
		io.WriteString(w, `{
			"success": true
		}`)
	}))
	defer srv.Close()

	fb := tagbox.New(srv.URL)
	err := fb.Rename("image1.jpg", "monkeys")
	is.NoErr(err)
}

func TestRenameAll(t *testing.T) {
	is := is.New(t)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		is.Equal(r.URL.Path, "/tagbox/rename")
		is.Equal(r.Header.Get("Accept"), "application/json; charset=utf-8")
		is.Equal(r.Method, "POST")
		is.Equal(r.FormValue("from"), "monkeys")
		is.Equal(r.FormValue("to"), "apes")
		io.WriteString(w, `{
			"success": true
		}`)
	}))
	defer srv.Close()

	fb := tagbox.New(srv.URL)
	err := fb.RenameAll("monkeys", "apes")
	is.NoErr(err)
}
