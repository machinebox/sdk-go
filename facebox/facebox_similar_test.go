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

func TestSimilarURL(t *testing.T) {
	is := is.New(t)

	imageURL, err := url.Parse("https://test.machinebox.io/image1.png")
	is.NoErr(err)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		is.Equal(r.URL.Path, "/facebox/similar")
		is.Equal(r.Header.Get("Accept"), "application/json; charset=utf-8")
		is.Equal(r.FormValue("url"), imageURL.String())
		io.WriteString(w, `{
			"success": true,
			"similarCount": 3,
			"similar": [
				{
					"id": "file1.jpg",
					"name": "Ringo Starr"
				},
				{
					"id": "file2.jpg",
					"name": "Ringo Starr"
				},
				{
					"id": "file3.jpg",
					"name": "Ringo Starr"
				}
			]
		}`)
	}))
	defer srv.Close()

	fb := facebox.New(srv.URL)
	similar, err := fb.SimilarURL(imageURL)
	is.NoErr(err)

	is.Equal(len(similar), 3)
	is.Equal(similar[0].ID, "file1.jpg")
	is.Equal(similar[0].Name, "Ringo Starr")

	is.Equal(similar[1].ID, "file2.jpg")
	is.Equal(similar[1].Name, "Ringo Starr")

	is.Equal(similar[2].ID, "file3.jpg")
	is.Equal(similar[2].Name, "Ringo Starr")

}

func TestSimilarURLError(t *testing.T) {
	is := is.New(t)

	imageURL, err := url.Parse("https://test.machinebox.io/image1.png")
	is.NoErr(err)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		is.Equal(r.URL.Path, "/facebox/similar")
		is.Equal(r.Header.Get("Accept"), "application/json; charset=utf-8")
		is.Equal(r.FormValue("url"), imageURL.String())
		io.WriteString(w, `{
			"success": false,
			"error": "something went wrong"
		}`)
	}))
	defer srv.Close()

	fb := facebox.New(srv.URL)
	_, err = fb.SimilarURL(imageURL)
	is.True(err != nil)
	is.Equal(err.Error(), "facebox: something went wrong")

}

func TestSimilarImage(t *testing.T) {
	is := is.New(t)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		is.Equal(r.URL.Path, "/facebox/similar")
		is.Equal(r.Header.Get("Accept"), "application/json; charset=utf-8")
		f, _, err := r.FormFile("file")
		is.NoErr(err)
		defer f.Close()
		b, err := ioutil.ReadAll(f)
		is.NoErr(err)
		is.Equal(string(b), `(pretend this is image data)`)
		io.WriteString(w, `{
			"success": true,
			"similarCount": 3,
			"similar": [
				{
					"id": "file1.jpg",
					"name": "Ringo Starr"
				},
				{
					"id": "file2.jpg",
					"name": "Ringo Starr"
				},
				{
					"id": "file3.jpg",
					"name": "Ringo Starr"
				}
			]
		}`)
	}))
	defer srv.Close()

	fb := facebox.New(srv.URL)
	similar, err := fb.Similar(strings.NewReader(`(pretend this is image data)`))
	is.NoErr(err)

	is.Equal(len(similar), 3)
	is.Equal(similar[0].ID, "file1.jpg")
	is.Equal(similar[0].Name, "Ringo Starr")

	is.Equal(similar[1].ID, "file2.jpg")
	is.Equal(similar[1].Name, "Ringo Starr")

	is.Equal(similar[2].ID, "file3.jpg")
	is.Equal(similar[2].Name, "Ringo Starr")

}

func TestSimilarImageError(t *testing.T) {
	is := is.New(t)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		is.Equal(r.URL.Path, "/facebox/similar")
		is.Equal(r.Header.Get("Accept"), "application/json; charset=utf-8")
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
	_, err := fb.Similar(strings.NewReader(`(pretend this is image data)`))
	is.True(err != nil)
	is.Equal(err.Error(), "facebox: something went wrong")

}

func TestSimilarID(t *testing.T) {
	is := is.New(t)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		is.Equal(r.URL.Path, "/facebox/similar")
		is.Equal(r.Header.Get("Accept"), "application/json; charset=utf-8")
		is.Equal(r.FormValue("id"), "abc123")
		io.WriteString(w, `{
			"success": true,
			"similarCount": 3,
			"similar": [
				{
					"id": "file1.jpg",
					"name": "Ringo Starr"
				},
				{
					"id": "file2.jpg",
					"name": "Ringo Starr"
				},
				{
					"id": "file3.jpg",
					"name": "Ringo Starr"
				}
			]
		}`)
	}))
	defer srv.Close()

	fb := facebox.New(srv.URL)
	similar, err := fb.SimilarID("abc123")
	is.NoErr(err)

	is.Equal(len(similar), 3)
	is.Equal(similar[0].ID, "file1.jpg")
	is.Equal(similar[0].Name, "Ringo Starr")

	is.Equal(similar[1].ID, "file2.jpg")
	is.Equal(similar[1].Name, "Ringo Starr")

	is.Equal(similar[2].ID, "file3.jpg")
	is.Equal(similar[2].Name, "Ringo Starr")
}

func TestSimilarBase64(t *testing.T) {
	is := is.New(t)

	base64Str := "base64Str"

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		is.Equal(r.URL.Path, "/facebox/similar")
		is.Equal(r.Header.Get("Accept"), "application/json; charset=utf-8")
		is.Equal(r.FormValue("base64"), base64Str)
		io.WriteString(w, `{
			"success": true,
			"similarCount": 3,
			"similar": [
				{
					"id": "file1.jpg",
					"name": "Ringo Starr"
				},
				{
					"id": "file2.jpg",
					"name": "Ringo Starr"
				},
				{
					"id": "file3.jpg",
					"name": "Ringo Starr"
				}
			]
		}`)
	}))
	defer srv.Close()

	fb := facebox.New(srv.URL)
	similar, err := fb.SimilarBase64(base64Str)
	is.NoErr(err)

	is.Equal(len(similar), 3)

}
