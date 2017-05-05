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

func TestTeachURL(t *testing.T) {
	is := is.New(t)
	imageURL, err := url.Parse("https://test.machinebox.io/image1.png")
	is.NoErr(err)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		is.Equal(r.URL.Path, "/facebox/teach")
		is.Equal(r.FormValue("url"), imageURL.String())
		is.Equal(r.FormValue("name"), "John Lennon")
		is.Equal(r.FormValue("id"), "john1.jpg")
		io.WriteString(w, `{
			"success": true
		}`)
	}))
	defer srv.Close()
	fb := facebox.New(srv.URL)
	err = fb.TeachURL(imageURL, "john1.jpg", "John Lennon")
	is.NoErr(err)
}

func TestTeachURLError(t *testing.T) {
	is := is.New(t)
	imageURL, err := url.Parse("https://test.machinebox.io/image1.png")
	is.NoErr(err)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		is.Equal(r.URL.Path, "/facebox/teach")
		is.Equal(r.FormValue("url"), imageURL.String())
		is.Equal(r.FormValue("name"), "John Lennon")
		is.Equal(r.FormValue("id"), "john1.jpg")
		io.WriteString(w, `{
			"success": false,
			"error": "something went wrong"
		}`)
	}))
	defer srv.Close()
	fb := facebox.New(srv.URL)
	err = fb.TeachURL(imageURL, "john1.jpg", "John Lennon")
	is.True(err != nil)
	is.Equal(err.Error(), "facebox: something went wrong")
}

func TestTeachImage(t *testing.T) {
	is := is.New(t)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		is.Equal(r.URL.Path, "/facebox/teach")
		is.Equal(r.FormValue("name"), "John Lennon")
		is.Equal(r.FormValue("id"), "john1.jpg")
		f, _, err := r.FormFile("file")
		is.NoErr(err)
		defer f.Close()
		b, err := ioutil.ReadAll(f)
		is.NoErr(err)
		is.Equal(string(b), `(pretend this is image data)`)
		io.WriteString(w, `{
			"success": true
		}`)
	}))
	defer srv.Close()

	fb := facebox.New(srv.URL)
	err := fb.Teach(strings.NewReader(`(pretend this is image data)`), "john1.jpg", "John Lennon")
	is.NoErr(err)

}

func TestTeachImageError(t *testing.T) {
	is := is.New(t)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		is.Equal(r.URL.Path, "/facebox/teach")
		is.Equal(r.FormValue("name"), "John Lennon")
		is.Equal(r.FormValue("id"), "john1.jpg")
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

	fb := facebox.New(srv.URL)
	err := fb.Teach(strings.NewReader(`(pretend this is image data)`), "john1.jpg", "John Lennon")
	is.True(err != nil)
	is.Equal(err.Error(), "facebox: something went wrong")

}
