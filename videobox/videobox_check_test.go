package videobox_test

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/machinebox/sdk-go/videobox"
	"github.com/matryer/is"
)

func TestCheckURL(t *testing.T) {
	is := is.New(t)

	imageURL, err := url.Parse("https://test.machinebox.io/image1.png")
	is.NoErr(err)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		is.Equal(r.URL.Path, "/videobox/check")
		is.Equal(r.Header.Get("Accept"), "application/json; charset=utf-8")
		is.Equal(r.FormValue("url"), imageURL.String())

		is.Equal(r.FormValue("resultsDuration"), "1h0m0s")
		is.Equal(r.FormValue("skipframes"), "2")
		is.Equal(r.FormValue("skipseconds"), "3")
		is.Equal(r.FormValue("frameWidth"), "100")
		is.Equal(r.FormValue("frameHeight"), "120")
		is.Equal(r.FormValue("faceboxThreshold"), "0.75")
		is.Equal(r.FormValue("tagboxInclude"), "custom")
		is.Equal(r.FormValue("tagboxThreshold"), "0.7")
		is.Equal(r.FormValue("nudeboxThreshold"), "0.2")

		_, err := io.WriteString(w, `{
			"success": true,
			"id": "5a50b8067eced76bad103c53dd0f5226",
			"status": "pending",
			"framesComplete": 0,
			"millisecondsComplete": 0
		}`)
		is.NoErr(err)
	}))
	defer srv.Close()

	vb := videobox.New(srv.URL)
	options := videobox.NewCheckOptions()
	options.ResultsDuration(1 * time.Hour)
	options.SkipFrames(2)
	options.SkipSeconds(3)
	options.FrameWidth(100)
	options.FrameHeight(120)
	options.FaceboxThreshold(0.75)
	options.TagboxIncludeCustom()
	options.TagboxThreshold(0.7)
	options.NudeboxThreshold(0.2)
	video, err := vb.CheckURL(imageURL, options)
	is.NoErr(err)
	is.Equal(video.ID, "5a50b8067eced76bad103c53dd0f5226")

}

func TestCheckURLError(t *testing.T) {
	is := is.New(t)

	imageURL, err := url.Parse("https://test.machinebox.io/image1.png")
	is.NoErr(err)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		is.Equal(r.URL.Path, "/videobox/check")
		is.Equal(r.Header.Get("Accept"), "application/json; charset=utf-8")
		is.Equal(r.FormValue("url"), imageURL.String())
		_, err := io.WriteString(w, `{
			"success": false,
			"error": "something went wrong"
		}`)
		is.NoErr(err)
	}))
	defer srv.Close()

	vb := videobox.New(srv.URL)
	_, err = vb.CheckURL(imageURL, nil)
	is.True(err != nil)
	is.Equal(err.Error(), "videobox: something went wrong")

}

func TestCheckImage(t *testing.T) {
	is := is.New(t)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		is.Equal(r.URL.Path, "/videobox/check")
		is.Equal(r.Header.Get("Accept"), "application/json; charset=utf-8")
		f, _, err := r.FormFile("file")
		is.NoErr(err)
		defer f.Close()
		b, err := ioutil.ReadAll(f)
		is.NoErr(err)
		is.Equal(string(b), `(pretend this is image data)`)
		is.Equal(r.FormValue("resultsDuration"), "1h0m0s")
		is.Equal(r.FormValue("skipframes"), "2")
		is.Equal(r.FormValue("frameWidth"), "100")
		is.Equal(r.FormValue("frameHeight"), "120")
		is.Equal(r.FormValue("faceboxThreshold"), "0.75")
		is.Equal(r.FormValue("tagboxInclude"), "custom")
		is.Equal(r.FormValue("tagboxThreshold"), "0.7")
		is.Equal(r.FormValue("nudeboxThreshold"), "0.2")
		_, err = io.WriteString(w, `{
			"success": true,
			"id": "5a50b8067eced76bad103c53dd0f5226",
			"status": "pending",
			"framesComplete": 0,
			"millisecondsComplete": 0
		}`)
		is.NoErr(err)
	}))
	defer srv.Close()

	vb := videobox.New(srv.URL)
	options := videobox.NewCheckOptions()
	options.ResultsDuration(1 * time.Hour)
	options.SkipFrames(2)
	options.SkipSeconds(3)
	options.FrameWidth(100)
	options.FrameHeight(120)
	options.FaceboxThreshold(0.75)
	options.TagboxIncludeCustom()
	options.TagboxThreshold(0.7)
	options.NudeboxThreshold(0.2)
	r := strings.NewReader(`(pretend this is image data)`)
	video, err := vb.Check(r, options)
	is.NoErr(err)
	is.Equal(video.ID, "5a50b8067eced76bad103c53dd0f5226")

}

func TestCheckImageError(t *testing.T) {
	is := is.New(t)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		is.Equal(r.URL.Path, "/videobox/check")
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

	vb := videobox.New(srv.URL)
	_, err := vb.Check(strings.NewReader(`(pretend this is image data)`), nil)
	is.True(err != nil)
	is.Equal(err.Error(), "videobox: something went wrong")

}

func TestCheckBase64(t *testing.T) {
	is := is.New(t)

	base64Str := `base64Str`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		is.Equal(r.URL.Path, "/videobox/check")
		is.Equal(r.Header.Get("Accept"), "application/json; charset=utf-8")
		is.Equal(r.FormValue("base64"), base64Str)
		_, err := io.WriteString(w, `{
			"success": true,
			"id": "5a50b8067eced76bad103c53dd0f5226",
			"status": "pending",
			"framesComplete": 0,
			"millisecondsComplete": 0
		}`)
		is.NoErr(err)
	}))
	defer srv.Close()

	vb := videobox.New(srv.URL)
	options := videobox.NewCheckOptions()
	options.ResultsDuration(1 * time.Hour)
	options.SkipFrames(2)
	options.SkipSeconds(3)
	options.FrameWidth(100)
	options.FrameHeight(120)
	options.FaceboxThreshold(0.75)
	options.TagboxIncludeCustom()
	options.TagboxThreshold(0.7)
	options.NudeboxThreshold(0.2)
	video, err := vb.CheckBase64(base64Str, options)
	is.NoErr(err)

	is.Equal(video.ID, "5a50b8067eced76bad103c53dd0f5226")

}
