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

func TestCheckURL(t *testing.T) {
	is := is.New(t)

	imageURL, err := url.Parse("https://test.machinebox.io/image1.png")
	is.NoErr(err)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		is.Equal(r.URL.Path, "/facebox/check")
		is.Equal(r.FormValue("url"), imageURL.String())
		io.WriteString(w, `{
			"success": true,
			"facesCount": 3,
			"faces": [
				{
					"rect": { "top": 0, "left": 0, "width": 120, "height": 120 },
					"id": "file1.jpg",
					"name": "John Lennon",
					"matched": true
				},
				{
					"rect": { "top": 200, "left": 200, "width": 100, "height": 100 },
					"id": "file6.jpg",
					"name": "Ringo Starr",
					"matched": true
				},
				{
					"rect": { "top": 50, "left": 50, "width": 100, "height": 100 },
					"matched": false
				}
			]
		}`)
	}))
	defer srv.Close()

	fb := facebox.New(srv.URL)
	faces, err := fb.CheckURL(imageURL)
	is.NoErr(err)

	is.Equal(len(faces), 3)
	is.Equal(faces[0].Rect.Top, 0)
	is.Equal(faces[0].Rect.Left, 0)
	is.Equal(faces[0].Rect.Width, 120)
	is.Equal(faces[0].Rect.Height, 120)
	is.Equal(faces[0].ID, "file1.jpg")
	is.Equal(faces[0].Name, "John Lennon")
	is.Equal(faces[0].Matched, true)

	is.Equal(faces[1].Rect.Top, 200)
	is.Equal(faces[1].Rect.Left, 200)
	is.Equal(faces[1].Rect.Width, 100)
	is.Equal(faces[1].Rect.Height, 100)
	is.Equal(faces[1].ID, "file6.jpg")
	is.Equal(faces[1].Name, "Ringo Starr")
	is.Equal(faces[1].Matched, true)

	is.Equal(faces[2].Rect.Top, 50)
	is.Equal(faces[2].Rect.Left, 50)
	is.Equal(faces[2].Rect.Width, 100)
	is.Equal(faces[2].Rect.Height, 100)
	is.Equal(faces[2].Matched, false)
	is.Equal(faces[2].ID, "")
	is.Equal(faces[2].Name, "")

}

func TestCheckURLError(t *testing.T) {
	is := is.New(t)

	imageURL, err := url.Parse("https://test.machinebox.io/image1.png")
	is.NoErr(err)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		is.Equal(r.URL.Path, "/facebox/check")
		is.Equal(r.FormValue("url"), imageURL.String())
		io.WriteString(w, `{
			"success": false,
			"error": "something went wrong"
		}`)
	}))
	defer srv.Close()

	fb := facebox.New(srv.URL)
	_, err = fb.CheckURL(imageURL)
	is.True(err != nil)
	is.Equal(err.Error(), "facebox: something went wrong")

}

func TestCheckImage(t *testing.T) {
	is := is.New(t)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		is.Equal(r.URL.Path, "/facebox/check")
		f, _, err := r.FormFile("file")
		is.NoErr(err)
		defer f.Close()
		b, err := ioutil.ReadAll(f)
		is.NoErr(err)
		is.Equal(string(b), `(pretend this is image data)`)
		io.WriteString(w, `{
			"success": true,
			"facesCount": 3,
			"faces": [
				{
					"rect": { "top": 0, "left": 0, "width": 120, "height": 120 },
					"id": "file1.jpg",
					"name": "John Lennon",
					"matched": true
				},
				{
					"rect": { "top": 200, "left": 200, "width": 100, "height": 100 },
					"id": "file6.jpg",
					"name": "Ringo Starr",
					"matched": true
				},
				{
					"rect": { "top": 50, "left": 50, "width": 100, "height": 100 },
					"matched": false
				}
			]
		}`)
	}))
	defer srv.Close()

	fb := facebox.New(srv.URL)
	faces, err := fb.Check(strings.NewReader(`(pretend this is image data)`))
	is.NoErr(err)

	is.Equal(len(faces), 3)
	is.Equal(faces[0].Rect.Top, 0)
	is.Equal(faces[0].Rect.Left, 0)
	is.Equal(faces[0].Rect.Width, 120)
	is.Equal(faces[0].Rect.Height, 120)
	is.Equal(faces[0].ID, "file1.jpg")
	is.Equal(faces[0].Name, "John Lennon")
	is.Equal(faces[0].Matched, true)

	is.Equal(faces[1].Rect.Top, 200)
	is.Equal(faces[1].Rect.Left, 200)
	is.Equal(faces[1].Rect.Width, 100)
	is.Equal(faces[1].Rect.Height, 100)
	is.Equal(faces[1].ID, "file6.jpg")
	is.Equal(faces[1].Name, "Ringo Starr")
	is.Equal(faces[1].Matched, true)

	is.Equal(faces[2].Rect.Top, 50)
	is.Equal(faces[2].Rect.Left, 50)
	is.Equal(faces[2].Rect.Width, 100)
	is.Equal(faces[2].Rect.Height, 100)
	is.Equal(faces[2].Matched, false)
	is.Equal(faces[2].ID, "")
	is.Equal(faces[2].Name, "")

}

func TestCheckImageError(t *testing.T) {
	is := is.New(t)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		is.Equal(r.URL.Path, "/facebox/check")
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
	_, err := fb.Check(strings.NewReader(`(pretend this is image data)`))
	is.True(err != nil)
	is.Equal(err.Error(), "facebox: something went wrong")

}
