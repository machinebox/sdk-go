package facebox_test

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/machinebox/sdk-go/facebox"
	"github.com/matryer/is"
)

func TestOpenState(t *testing.T) {
	is := is.New(t)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		is.Equal(r.Method, "GET")
		is.Equal(r.URL.Path, "/facebox/state")
		io.WriteString(w, `(pretend this is the state file)`)
	}))
	defer srv.Close()
	fb := facebox.New(srv.URL)
	f, err := fb.OpenState()
	is.NoErr(err)
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	is.NoErr(err)
	is.Equal(string(b), `(pretend this is the state file)`)
}

func TestPostState(t *testing.T) {
	is := is.New(t)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		is.Equal(r.Method, "POST")
		is.Equal(r.URL.Path, "/facebox/state")
		f, _, err := r.FormFile("file")
		is.NoErr(err)
		b, err := ioutil.ReadAll(f)
		is.NoErr(err)
		is.Equal(string(b), `(pretend this is the state file)`)
		io.WriteString(w, `{"success":true}`)
	}))
	defer srv.Close()
	fb := facebox.New(srv.URL)
	r := strings.NewReader(`(pretend this is the state file)`)
	err := fb.PostState(r)
	is.NoErr(err)
}

func TestPostStateError(t *testing.T) {
	is := is.New(t)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		is.Equal(r.Method, "POST")
		is.Equal(r.URL.Path, "/facebox/state")
		f, _, err := r.FormFile("file")
		is.NoErr(err)
		b, err := ioutil.ReadAll(f)
		is.NoErr(err)
		is.Equal(string(b), `(pretend this is the state file)`)
		io.WriteString(w, `{"success":false,"error":"something went wrong"}`)
	}))
	defer srv.Close()
	fb := facebox.New(srv.URL)
	r := strings.NewReader(`(pretend this is the state file)`)
	err := fb.PostState(r)
	is.True(err != nil)
	is.Equal(err.Error(), "facebox: something went wrong")
}

func TestPostStateURL(t *testing.T) {
	is := is.New(t)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		is.Equal(r.Method, "POST")
		is.Equal(r.URL.Path, "/facebox/state")
		is.Equal(r.FormValue("url"), "https://test.machinebox.io/test.facebox")
		io.WriteString(w, `{"success":true}`)
	}))
	defer srv.Close()
	fb := facebox.New(srv.URL)
	u, err := url.Parse("https://test.machinebox.io/test.facebox")
	is.NoErr(err)
	err = fb.PostStateURL(u)
	is.NoErr(err)
}

func TestPostStateURLError(t *testing.T) {
	is := is.New(t)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		is.Equal(r.Method, "POST")
		is.Equal(r.URL.Path, "/facebox/state")
		is.Equal(r.FormValue("url"), "https://test.machinebox.io/test.facebox")
		io.WriteString(w, `{"success":false,"error":"something went wrong"}`)
	}))
	defer srv.Close()
	fb := facebox.New(srv.URL)
	u, err := url.Parse("https://test.machinebox.io/test.facebox")
	is.NoErr(err)
	err = fb.PostStateURL(u)
	is.True(err != nil)
	is.Equal(err.Error(), "facebox: something went wrong")
}
