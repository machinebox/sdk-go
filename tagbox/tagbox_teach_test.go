package tagbox_test

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/machinebox/sdk-go/tagbox"
	"github.com/matryer/is"
)

func TestTeachURL(t *testing.T) {
	is := is.New(t)
	imageURL, err := url.Parse("https://test.machinebox.io/image1.png")
	is.NoErr(err)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		is.Equal(r.URL.Path, "/tagbox/teach")
		is.Equal(r.Header.Get("Accept"), "application/json; charset=utf-8")
		is.Equal(r.FormValue("url"), imageURL.String())
		is.Equal(r.FormValue("tag"), "monkeys")
		is.Equal(r.FormValue("id"), "image1.jpg")
		io.WriteString(w, `{
			"success": true
		}`)
	}))
	defer srv.Close()
	fb := tagbox.New(srv.URL)
	err = fb.TeachURL(imageURL, "image1.jpg", "monkeys")
	is.NoErr(err)
}

func TestTeachURLError(t *testing.T) {
	is := is.New(t)
	imageURL, err := url.Parse("https://test.machinebox.io/image1.png")
	is.NoErr(err)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		is.Equal(r.URL.Path, "/tagbox/teach")
		is.Equal(r.Header.Get("Accept"), "application/json; charset=utf-8")
		is.Equal(r.FormValue("url"), imageURL.String())
		is.Equal(r.FormValue("tag"), "monkeys")
		is.Equal(r.FormValue("id"), "image1.jpg")
		io.WriteString(w, `{
			"success": false,
			"error": "something went wrong"
		}`)
	}))
	defer srv.Close()
	fb := tagbox.New(srv.URL)
	err = fb.TeachURL(imageURL, "image1.jpg", "monkeys")
	is.True(err != nil)
	is.Equal(err.Error(), "tagbox: something went wrong")
}

func TestTeachBase64(t *testing.T) {
	is := is.New(t)
	base64Str := "base64Str"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		is.Equal(r.URL.Path, "/tagbox/teach")
		is.Equal(r.Header.Get("Accept"), "application/json; charset=utf-8")
		is.Equal(r.FormValue("base64"), base64Str)
		is.Equal(r.FormValue("tag"), "monkeys")
		is.Equal(r.FormValue("id"), "image1.jpg")
		io.WriteString(w, `{
			"success": true
		}`)
	}))
	defer srv.Close()
	fb := tagbox.New(srv.URL)
	err := fb.TeachBase64(base64Str, "image1.jpg", "monkeys")
	is.NoErr(err)
}

func TestTeachImage(t *testing.T) {
	is := is.New(t)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		is.Equal(r.URL.Path, "/tagbox/teach")
		is.Equal(r.Header.Get("Accept"), "application/json; charset=utf-8")
		is.Equal(r.FormValue("tag"), "monkeys")
		is.Equal(r.FormValue("id"), "image1.jpg")
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

	fb := tagbox.New(srv.URL)
	err := fb.Teach(strings.NewReader(`(pretend this is image data)`), "image1.jpg", "monkeys")
	is.NoErr(err)

}

func TestTeachImageError(t *testing.T) {
	is := is.New(t)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		is.Equal(r.URL.Path, "/tagbox/teach")
		is.Equal(r.Header.Get("Accept"), "application/json; charset=utf-8")
		is.Equal(r.FormValue("tag"), "monkeys")
		is.Equal(r.FormValue("id"), "image1.jpg")
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

	fb := tagbox.New(srv.URL)
	err := fb.Teach(strings.NewReader(`(pretend this is image data)`), "image1.jpg", "monkeys")
	is.True(err != nil)
	is.Equal(err.Error(), "tagbox: something went wrong")

}

func TestRemove(t *testing.T) {
	is := is.New(t)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		is.Equal(r.URL.Path, "/tagbox/teach/image1.jpg")
		is.Equal(r.Header.Get("Accept"), "application/json; charset=utf-8")
		is.Equal(r.Method, "DELETE")
		io.WriteString(w, `{
			"success": true
		}`)
	}))
	defer srv.Close()

	fb := tagbox.New(srv.URL)
	err := fb.Remove("image1.jpg")
	is.NoErr(err)
}
