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

func TestSimilarURL(t *testing.T) {
	is := is.New(t)

	imageURL, err := url.Parse("https://test.machinebox.io/image1.png")
	is.NoErr(err)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		is.Equal(r.URL.Path, "/tagbox/similar")
		is.Equal(r.Header.Get("Accept"), "application/json; charset=utf-8")
		is.Equal(r.FormValue("url"), imageURL.String())
		io.WriteString(w, `{
			"success": true,
			"tagCount": 3,
			"similar": [
				{
					"tag": "monkeys",
					"confidence": 0.79,
					"id": "file1.jpg"
				},
				{
					"tag": "monkeys",
					"confidence": 0.60,
					"id": "file2.jpg"
				},
				{
					"tag": "monkeys",
					"confidence": 0.50,
					"id": "file3.jpg"
				}
			]
		}`)
	}))

	defer srv.Close()

	fb := tagbox.New(srv.URL)
	similar, err := fb.SimilarURL(imageURL)
	is.NoErr(err)

	is.Equal(len(similar), 3)
	is.Equal(similar[0].ID, "file1.jpg")
	is.Equal(similar[0].Tag, "monkeys")
	is.Equal(similar[0].Confidence, 0.79)

	is.Equal(similar[1].ID, "file2.jpg")
	is.Equal(similar[1].Tag, "monkeys")
	is.Equal(similar[1].Confidence, 0.6)

	is.Equal(similar[2].ID, "file3.jpg")
	is.Equal(similar[2].Tag, "monkeys")
	is.Equal(similar[2].Confidence, 0.5)

}

func TestSimilarBase64(t *testing.T) {
	is := is.New(t)

	base64Str := "base64Str"

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		is.Equal(r.URL.Path, "/tagbox/similar")
		is.Equal(r.Header.Get("Accept"), "application/json; charset=utf-8")
		is.Equal(r.FormValue("base64"), base64Str)
		io.WriteString(w, `{
			"success": true,
			"tagCount": 3,
			"similar": [
				{
					"tag": "monkeys",
					"confidence": 0.79,
					"id": "file1.jpg"
				},
				{
					"tag": "monkeys",
					"confidence": 0.60,
					"id": "file2.jpg"
				},
				{
					"tag": "monkeys",
					"confidence": 0.50,
					"id": "file3.jpg"
				}
			]
		}`)
	}))

	defer srv.Close()

	fb := tagbox.New(srv.URL)
	similar, err := fb.SimilarBase64(base64Str)
	is.NoErr(err)

	is.Equal(len(similar), 3)
	is.Equal(similar[0].ID, "file1.jpg")
	is.Equal(similar[0].Tag, "monkeys")
	is.Equal(similar[0].Confidence, 0.79)

	is.Equal(similar[1].ID, "file2.jpg")
	is.Equal(similar[1].Tag, "monkeys")
	is.Equal(similar[1].Confidence, 0.6)

	is.Equal(similar[2].ID, "file3.jpg")
	is.Equal(similar[2].Tag, "monkeys")
	is.Equal(similar[2].Confidence, 0.5)

}

func TestSimilarURLError(t *testing.T) {
	is := is.New(t)

	imageURL, err := url.Parse("https://test.machinebox.io/image1.png")
	is.NoErr(err)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		is.Equal(r.URL.Path, "/tagbox/similar")
		is.Equal(r.Header.Get("Accept"), "application/json; charset=utf-8")
		is.Equal(r.FormValue("url"), imageURL.String())
		io.WriteString(w, `{
			"success": false,
			"error": "something went wrong"
		}`)
	}))
	defer srv.Close()

	fb := tagbox.New(srv.URL)
	_, err = fb.SimilarURL(imageURL)
	is.True(err != nil)
	is.Equal(err.Error(), "tagbox: something went wrong")

}

func TestSimilarImage(t *testing.T) {
	is := is.New(t)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		is.Equal(r.URL.Path, "/tagbox/similar")
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
					"tag": "monkeys",
					"confidence": 0.79,
					"id": "file1.jpg"
				},
				{
					"tag": "monkeys",
					"confidence": 0.60,
					"id": "file2.jpg"
				},
				{
					"tag": "monkeys",
					"confidence": 0.50,
					"id": "file3.jpg"
				}
			]
		}`)
	}))
	defer srv.Close()

	fb := tagbox.New(srv.URL)
	similar, err := fb.Similar(strings.NewReader(`(pretend this is image data)`))
	is.NoErr(err)

	is.Equal(len(similar), 3)
	is.Equal(similar[0].ID, "file1.jpg")
	is.Equal(similar[0].Tag, "monkeys")
	is.Equal(similar[0].Confidence, 0.79)

	is.Equal(similar[1].ID, "file2.jpg")
	is.Equal(similar[1].Tag, "monkeys")
	is.Equal(similar[1].Confidence, 0.6)

	is.Equal(similar[2].ID, "file3.jpg")
	is.Equal(similar[2].Tag, "monkeys")
	is.Equal(similar[2].Confidence, 0.5)

}

func TestSimilarImageError(t *testing.T) {
	is := is.New(t)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		is.Equal(r.URL.Path, "/tagbox/similar")
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

	fb := tagbox.New(srv.URL)
	_, err := fb.Similar(strings.NewReader(`(pretend this is image data)`))
	is.True(err != nil)
	is.Equal(err.Error(), "tagbox: something went wrong")

}

func TestSimilarID(t *testing.T) {
	is := is.New(t)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		is.Equal(r.URL.Path, "/tagbox/similar")
		is.Equal(r.Header.Get("Accept"), "application/json; charset=utf-8")
		is.Equal(r.FormValue("id"), "abc123")
		io.WriteString(w, `{
			"success": true,
			"similarCount": 3,
			"similar": [
				{
					"tag": "monkeys",
					"confidence": 0.79,
					"id": "file1.jpg"
				},
				{
					"tag": "monkeys",
					"confidence": 0.60,
					"id": "file2.jpg"
				},
				{
					"tag": "monkeys",
					"confidence": 0.50,
					"id": "file3.jpg"
				}
			]
		}`)
	}))
	defer srv.Close()

	fb := tagbox.New(srv.URL)
	similar, err := fb.SimilarID("abc123")
	is.NoErr(err)

	is.Equal(len(similar), 3)
	is.Equal(similar[0].ID, "file1.jpg")
	is.Equal(similar[0].Tag, "monkeys")
	is.Equal(similar[0].Confidence, 0.79)

	is.Equal(similar[1].ID, "file2.jpg")
	is.Equal(similar[1].Tag, "monkeys")
	is.Equal(similar[1].Confidence, 0.6)

	is.Equal(similar[2].ID, "file3.jpg")
	is.Equal(similar[2].Tag, "monkeys")
	is.Equal(similar[2].Confidence, 0.5)
}
