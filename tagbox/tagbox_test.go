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

func TestInfo(t *testing.T) {
	is := is.New(t)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		is.Equal(r.Method, "GET")
		is.Equal(r.URL.Path, "/info")
		io.WriteString(w, `{
			"name": "tagbox",
			"version": 1,
			"build": "abcdefg",
			"status": "ready"
		}`)
	}))
	defer srv.Close()
	tb := tagbox.New(srv.URL)
	info, err := tb.Info()
	is.NoErr(err)
	is.Equal(info.Name, "tagbox")
	is.Equal(info.Version, 1)
	is.Equal(info.Build, "abcdefg")
	is.Equal(info.Status, "ready")
}

func TestCheckURL(t *testing.T) {
	is := is.New(t)

	imageURL, err := url.Parse("https://test.machinebox.io/image1.png")
	is.NoErr(err)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		is.Equal(r.URL.Path, "/tagbox/check")
		is.Equal(r.FormValue("url"), imageURL.String())
		io.WriteString(w, `{
			"success": true,
			"tags": [
				{"tag":"one", "confidence":0.9},
				{"tag":"two", "confidence":0.8},
				{"tag":"three", "confidence":0.7}
			]
		}`)
	}))
	defer srv.Close()

	tb := tagbox.New(srv.URL)
	tags, err := tb.CheckURL(imageURL)
	is.NoErr(err)

	is.Equal(len(tags), 3)
	is.Equal(tags[0].Tag, "one")
	is.Equal(tags[0].Confidence, 0.9)
	is.Equal(tags[1].Tag, "two")
	is.Equal(tags[1].Confidence, 0.8)
	is.Equal(tags[2].Tag, "three")
	is.Equal(tags[2].Confidence, 0.7)

}

func TestCheckURLError(t *testing.T) {
	is := is.New(t)

	imageURL, err := url.Parse("https://test.machinebox.io/image1.png")
	is.NoErr(err)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		is.Equal(r.URL.Path, "/tagbox/check")
		is.Equal(r.FormValue("url"), imageURL.String())
		io.WriteString(w, `{
			"success": false,
			"error": "something went wrong"
		}`)
	}))
	defer srv.Close()

	tb := tagbox.New(srv.URL)
	_, err = tb.CheckURL(imageURL)
	is.True(err != nil)
	is.Equal(err.Error(), "tagbox: something went wrong")

}

func TestCheckImage(t *testing.T) {
	is := is.New(t)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		is.Equal(r.URL.Path, "/tagbox/check")
		f, _, err := r.FormFile("file")
		is.NoErr(err)
		defer f.Close()
		b, err := ioutil.ReadAll(f)
		is.NoErr(err)
		is.Equal(string(b), `(pretend this is image data)`)
		io.WriteString(w, `{
			"success": true,
			"tags": [
				{"tag":"one", "confidence":0.9},
				{"tag":"two", "confidence":0.8},
				{"tag":"three", "confidence":0.7}
			]
		}`)
	}))
	defer srv.Close()

	tb := tagbox.New(srv.URL)
	tags, err := tb.Check(strings.NewReader(`(pretend this is image data)`))
	is.NoErr(err)

	is.Equal(len(tags), 3)
	is.Equal(tags[0].Tag, "one")
	is.Equal(tags[0].Confidence, 0.9)
	is.Equal(tags[1].Tag, "two")
	is.Equal(tags[1].Confidence, 0.8)
	is.Equal(tags[2].Tag, "three")
	is.Equal(tags[2].Confidence, 0.7)

}

func TestCheckImageError(t *testing.T) {
	is := is.New(t)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		is.Equal(r.URL.Path, "/tagbox/check")
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

	tb := tagbox.New(srv.URL)
	_, err := tb.Check(strings.NewReader(`(pretend this is image data)`))
	is.True(err != nil)
	is.Equal(err.Error(), "tagbox: something went wrong")

}
