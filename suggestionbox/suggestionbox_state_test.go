package suggestionbox_test

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/machinebox/sdk-go/suggestionbox"
	"github.com/matryer/is"
)

func TestOpenState(t *testing.T) {
	is := is.New(t)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		is.Equal(r.Method, "GET")
		is.Equal(r.URL.Path, "/suggestionbox/state/model1")
		io.WriteString(w, `(pretend this is the state file)`)
	}))
	defer srv.Close()
	sb := suggestionbox.New(srv.URL)
	f, err := sb.OpenState(context.Background(), "model1")
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
		is.Equal(r.URL.Path, "/suggestionbox/state")
		is.Equal(r.Header.Get("Accept"), "application/json; charset=utf-8")
		f, _, err := r.FormFile("file")
		is.NoErr(err)
		b, err := ioutil.ReadAll(f)
		is.NoErr(err)
		is.Equal(string(b), `(pretend this is the state file)`)
		io.WriteString(w, `{"success":true,"id":"model1"}`)
	}))
	defer srv.Close()
	sb := suggestionbox.New(srv.URL)
	r := strings.NewReader(`(pretend this is the state file)`)
	model, err := sb.PostState(context.Background(), r)
	is.NoErr(err)
	is.Equal(model.ID, "model1")
}

func TestPostStateError(t *testing.T) {
	is := is.New(t)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		is.Equal(r.Method, "POST")
		is.Equal(r.URL.Path, "/suggestionbox/state")
		is.Equal(r.Header.Get("Accept"), "application/json; charset=utf-8")
		f, _, err := r.FormFile("file")
		is.NoErr(err)
		b, err := ioutil.ReadAll(f)
		is.NoErr(err)
		is.Equal(string(b), `(pretend this is the state file)`)
		io.WriteString(w, `{"success":false,"error":"something went wrong"}`)
	}))
	defer srv.Close()
	sb := suggestionbox.New(srv.URL)
	r := strings.NewReader(`(pretend this is the state file)`)
	_, err := sb.PostState(context.Background(), r)
	is.True(err != nil)
	is.Equal(err.Error(), "suggestionbox: something went wrong")
}

func TestPostStateURL(t *testing.T) {
	is := is.New(t)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		is.Equal(r.Method, "POST")
		is.Equal(r.URL.Path, "/suggestionbox/state")
		is.Equal(r.Header.Get("Accept"), "application/json; charset=utf-8")
		is.Equal(r.FormValue("url"), "https://test.machinebox.io/test.suggestionbox")
		io.WriteString(w, `{"success":true,"id":"model1"}`)
	}))
	defer srv.Close()
	sb := suggestionbox.New(srv.URL)
	u, err := url.Parse("https://test.machinebox.io/test.suggestionbox")
	is.NoErr(err)
	model, err := sb.PostStateURL(context.Background(), u)
	is.NoErr(err)
	is.Equal(model.ID, "model1")
}

func TestPostStateURLError(t *testing.T) {
	is := is.New(t)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		is.Equal(r.Method, "POST")
		is.Equal(r.URL.Path, "/suggestionbox/state")
		is.Equal(r.Header.Get("Accept"), "application/json; charset=utf-8")
		is.Equal(r.FormValue("url"), "https://test.machinebox.io/test.suggestionbox")
		io.WriteString(w, `{"success":false,"error":"something went wrong"}`)
	}))
	defer srv.Close()
	sb := suggestionbox.New(srv.URL)
	u, err := url.Parse("https://test.machinebox.io/test.suggestionbox")
	is.NoErr(err)
	_, err = sb.PostStateURL(context.Background(), u)
	is.True(err != nil)
	is.Equal(err.Error(), "suggestionbox: something went wrong")
}
