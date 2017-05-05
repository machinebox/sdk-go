package nudebox_test

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/machinebox/sdk-go/nudebox"
	"github.com/matryer/is"
)

func TestInfo(t *testing.T) {
	is := is.New(t)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		is.Equal(r.Method, "GET")
		is.Equal(r.URL.Path, "/info")
		io.WriteString(w, `{
			"name": "nudebox",
			"version": 1,
			"build": "abcdefg",
			"status": "ready"
		}`)
	}))
	defer srv.Close()
	nb := nudebox.New(srv.URL)
	info, err := nb.Info()
	is.NoErr(err)
	is.Equal(info.Name, "nudebox")
	is.Equal(info.Version, 1)
	is.Equal(info.Build, "abcdefg")
	is.Equal(info.Status, "ready")
}

func TestCheckURL(t *testing.T) {
	is := is.New(t)

	imageURL, err := url.Parse("https://test.machinebox.io/image1.png")
	is.NoErr(err)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		is.Equal(r.URL.Path, "/nudebox/check")
		is.Equal(r.FormValue("url"), imageURL.String())
		io.WriteString(w, `{
			"success": true,
			"nude": 0.8
		}`)
	}))
	defer srv.Close()

	nb := nudebox.New(srv.URL)
	nude, err := nb.CheckURL(imageURL)
	is.NoErr(err)

	is.Equal(nude, 0.8)

}

func TestCheckURLError(t *testing.T) {
	is := is.New(t)

	imageURL, err := url.Parse("https://test.machinebox.io/image1.png")
	is.NoErr(err)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		is.Equal(r.URL.Path, "/nudebox/check")
		is.Equal(r.FormValue("url"), imageURL.String())
		io.WriteString(w, `{
			"success": false,
			"error": "something went wrong"
		}`)
	}))
	defer srv.Close()

	nb := nudebox.New(srv.URL)
	_, err = nb.CheckURL(imageURL)
	is.True(err != nil)
	is.Equal(err.Error(), "nudebox: something went wrong")

}

func TestCheckImage(t *testing.T) {
	is := is.New(t)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		is.Equal(r.URL.Path, "/nudebox/check")
		f, _, err := r.FormFile("file")
		is.NoErr(err)
		defer f.Close()
		b, err := ioutil.ReadAll(f)
		is.NoErr(err)
		is.Equal(string(b), `(pretend this is image data)`)
		io.WriteString(w, `{
			"success": true,
			"nude": 0.23
		}`)
	}))
	defer srv.Close()

	nb := nudebox.New(srv.URL)
	nude, err := nb.Check(strings.NewReader(`(pretend this is image data)`))
	is.NoErr(err)
	is.Equal(nude, 0.23)

}

func TestCheckImageError(t *testing.T) {
	is := is.New(t)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		is.Equal(r.URL.Path, "/nudebox/check")
		f, _, err := r.FormFile("file")
		is.NoErr(err)
		defer f.Close()
		b, err := ioutil.ReadAll(f)
		is.NoErr(err)
		is.Equal(string(b), `(pretend this is image data)`)
		io.WriteString(w, `{
			"success": false,
			"error": "something went wrong"
		}`)
	}))
	defer srv.Close()

	nb := nudebox.New(srv.URL)
	_, err := nb.Check(strings.NewReader(`(pretend this is image data)`))
	is.True(err != nil)
	is.Equal(err.Error(), "nudebox: something went wrong")

}
