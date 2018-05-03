package classificationbox_test

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/machinebox/sdk-go/classificationbox"
	"github.com/matryer/is"
)

func TestOpenState(t *testing.T) {
	is := is.New(t)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		is.Equal(r.Method, "GET")
		is.Equal(r.URL.Path, "/classificationbox/state/model1")
		io.WriteString(w, `(pretend this is the state file)`)
	}))
	defer srv.Close()
	cb := classificationbox.New(srv.URL)
	f, err := cb.OpenState(context.Background(), "model1")
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
		is.Equal(r.URL.Path, "/classificationbox/state")
		is.Equal(r.Header.Get("Accept"), "application/json; charset=utf-8")
		f, _, err := r.FormFile("file")
		is.NoErr(err)
		b, err := ioutil.ReadAll(f)
		is.NoErr(err)
		is.Equal(string(b), `(pretend this is the state file)`)
		io.WriteString(w, `{"success":true,"id":"model1"}`)
	}))
	defer srv.Close()
	cb := classificationbox.New(srv.URL)
	r := strings.NewReader(`(pretend this is the state file)`)
	model, err := cb.PostState(context.Background(), r, false)
	is.NoErr(err)
	is.Equal(model.ID, "model1")
}

func TestPostStateError(t *testing.T) {
	is := is.New(t)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		is.Equal(r.Method, "POST")
		is.Equal(r.URL.Path, "/classificationbox/state")
		is.Equal(r.Header.Get("Accept"), "application/json; charset=utf-8")
		f, _, err := r.FormFile("file")
		is.NoErr(err)
		b, err := ioutil.ReadAll(f)
		is.NoErr(err)
		is.Equal(string(b), `(pretend this is the state file)`)
		io.WriteString(w, `{"success":false,"error":"something went wrong"}`)
	}))
	defer srv.Close()
	cb := classificationbox.New(srv.URL)
	r := strings.NewReader(`(pretend this is the state file)`)
	_, err := cb.PostState(context.Background(), r, false)
	is.True(err != nil)
	is.Equal(err.Error(), "classificationbox: something went wrong")
}

func TestPostStateURL(t *testing.T) {
	is := is.New(t)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		is.Equal(r.Method, "POST")
		is.Equal(r.URL.Path, "/classificationbox/state")
		is.Equal(r.Header.Get("Accept"), "application/json; charset=utf-8")
		is.Equal(r.FormValue("url"), "https://test.machinebox.io/test.classificationbox")
		is.Equal(r.FormValue("predict_only"), "true")
		io.WriteString(w, `{"success":true,"id":"model1"}`)
	}))
	defer srv.Close()
	cb := classificationbox.New(srv.URL)
	u, err := url.Parse("https://test.machinebox.io/test.classificationbox")
	is.NoErr(err)
	model, err := cb.PostStateURL(context.Background(), u, true)
	is.NoErr(err)
	is.Equal(model.ID, "model1")
}

func TestPostStateURLError(t *testing.T) {
	is := is.New(t)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		is.Equal(r.Method, "POST")
		is.Equal(r.URL.Path, "/classificationbox/state")
		is.Equal(r.Header.Get("Accept"), "application/json; charset=utf-8")
		is.Equal(r.FormValue("url"), "https://test.machinebox.io/test.classificationbox")
		io.WriteString(w, `{"success":false,"error":"something went wrong"}`)
	}))
	defer srv.Close()
	cb := classificationbox.New(srv.URL)
	u, err := url.Parse("https://test.machinebox.io/test.classificationbox")
	is.NoErr(err)
	_, err = cb.PostStateURL(context.Background(), u, true)
	is.True(err != nil)
	is.Equal(err.Error(), "classificationbox: something went wrong")
}
