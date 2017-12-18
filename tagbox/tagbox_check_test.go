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

func TestCheckURL(t *testing.T) {
	is := is.New(t)

	imageURL, err := url.Parse("https://test.machinebox.io/image1.png")
	is.NoErr(err)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		is.Equal(r.URL.Path, "/tagbox/check")
		is.Equal(r.Header.Get("Accept"), "application/json; charset=utf-8")
		is.Equal(r.FormValue("url"), imageURL.String())
		io.WriteString(w, `{
			"success": true,
			"tags": [
				{"tag":"one", "confidence":0.9},
				{"tag":"two", "confidence":0.8},
				{"tag":"three", "confidence":0.7}
			],
			"custom_tags": [
				{"tag": "monkeys","confidence": 0.58,"id": "monkeys2.jpg"},
				{"tag": "bonobos","confidence": 0.4,"id": "monkeys3.jpg"}
			]
		}`)
	}))
	defer srv.Close()

	tb := tagbox.New(srv.URL)
	chk, err := tb.CheckURL(imageURL)
	is.NoErr(err)

	tags := chk.Tags
	is.Equal(len(tags), 3)
	is.Equal(tags[0].Tag, "one")
	is.Equal(tags[0].Confidence, 0.9)
	is.Equal(tags[1].Tag, "two")
	is.Equal(tags[1].Confidence, 0.8)
	is.Equal(tags[2].Tag, "three")
	is.Equal(tags[2].Confidence, 0.7)

	ctags := chk.CustomTags
	is.Equal(len(ctags), 2)
	is.Equal(ctags[0].Tag, "monkeys")
	is.Equal(ctags[0].Confidence, 0.58)
	is.Equal(ctags[0].ID, "monkeys2.jpg")
	is.Equal(ctags[1].Tag, "bonobos")
	is.Equal(ctags[1].Confidence, 0.4)
	is.Equal(ctags[1].ID, "monkeys3.jpg")

}

func TestCheckBase64(t *testing.T) {
	is := is.New(t)

	base64Str := "base64Str"

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		is.Equal(r.URL.Path, "/tagbox/check")
		is.Equal(r.Header.Get("Accept"), "application/json; charset=utf-8")
		is.Equal(r.FormValue("base64"), base64Str)
		io.WriteString(w, `{
			"success": true,
			"tags": [
				{"tag":"one", "confidence":0.9},
				{"tag":"two", "confidence":0.8},
				{"tag":"three", "confidence":0.7}
			],
			"custom_tags": [
				{"tag": "monkeys","confidence": 0.58,"id": "monkeys2.jpg"},
				{"tag": "bonobos","confidence": 0.4,"id": "monkeys3.jpg"}
			]
		}`)
	}))
	defer srv.Close()

	tb := tagbox.New(srv.URL)
	chk, err := tb.CheckBase64(base64Str)
	is.NoErr(err)

	tags := chk.Tags
	is.Equal(len(tags), 3)
	is.Equal(tags[0].Tag, "one")
	is.Equal(tags[0].Confidence, 0.9)
	is.Equal(tags[1].Tag, "two")
	is.Equal(tags[1].Confidence, 0.8)
	is.Equal(tags[2].Tag, "three")
	is.Equal(tags[2].Confidence, 0.7)

	ctags := chk.CustomTags
	is.Equal(len(ctags), 2)
	is.Equal(ctags[0].Tag, "monkeys")
	is.Equal(ctags[0].Confidence, 0.58)
	is.Equal(ctags[0].ID, "monkeys2.jpg")
	is.Equal(ctags[1].Tag, "bonobos")
	is.Equal(ctags[1].Confidence, 0.4)
	is.Equal(ctags[1].ID, "monkeys3.jpg")

}

func TestCheckURLError(t *testing.T) {
	is := is.New(t)

	imageURL, err := url.Parse("https://test.machinebox.io/image1.png")
	is.NoErr(err)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		is.Equal(r.URL.Path, "/tagbox/check")
		is.Equal(r.Header.Get("Accept"), "application/json; charset=utf-8")
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
		is.Equal(r.Header.Get("Accept"), "application/json; charset=utf-8")
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
	chk, err := tb.Check(strings.NewReader(`(pretend this is image data)`))
	is.NoErr(err)

	tags := chk.Tags

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

	tb := tagbox.New(srv.URL)
	_, err := tb.Check(strings.NewReader(`(pretend this is image data)`))
	is.True(err != nil)
	is.Equal(err.Error(), "tagbox: something went wrong")

}
