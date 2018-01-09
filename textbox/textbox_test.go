package textbox_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/machinebox/sdk-go/textbox"
	"github.com/matryer/is"
)

func TestInfo(t *testing.T) {
	is := is.New(t)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		is.Equal(r.Method, "GET")
		is.Equal(r.URL.Path, "/info")
		is.Equal(r.Header.Get("Accept"), "application/json; charset=utf-8")
		io.WriteString(w, `{
			"name": "textbox",
			"version": 1,
			"build": "abcdefg",
			"status": "ready"
		}`)
	}))
	defer srv.Close()
	b := textbox.New(srv.URL)
	info, err := b.Info()
	is.NoErr(err)
	is.Equal(info.Name, "textbox")
	is.Equal(info.Version, 1)
	is.Equal(info.Build, "abcdefg")
	is.Equal(info.Status, "ready")
}

func TestCheck(t *testing.T) {
	result := `{
	"success": true,
	"sentences": [
		{
			"text": "I really love Machina, who is the MachineBox mascot designed by Ashley McNamara.",
			"start": 0,
			"end": 80,
			"sentiment": 0.7128883004188538,
			"entities": [
				{
					"text": "Machina",
					"start": 14,
					"end": 21,
					"type": "person"
				},
				{
					"text": "MachineBox",
					"start": 34,
					"end": 44,
					"type": "organization"
				},
				{
					"text": "Ashley McNamara",
					"start": 64,
					"end": 79,
					"type": "person"
				}
			]
		}
	],
	"keywords": [
		{
			"keyword": "ashley mcnamara"
		},
		{
			"keyword": "machinebox mascot"
		},
		{
			"keyword": "machina"
		}
	]
}`
	is := is.New(t)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		is.Equal(r.Method, "POST")
		is.Equal(r.URL.Path, "/textbox/check")
		is.Equal(r.Header.Get("Accept"), "application/json; charset=utf-8")
		io.WriteString(w, result)
	}))
	defer srv.Close()
	tb := textbox.New(srv.URL)
	src := strings.NewReader(`I really love Machina, who is the MachineBox mascot designed by Ashley McNamara.`)
	res, err := tb.Check(src)
	is.NoErr(err)
	is.Equal(len(res.Keywords), 3)
	is.Equal(res.Keywords[0].Keyword, "ashley mcnamara")
	is.Equal(res.Keywords[1].Keyword, "machinebox mascot")
	is.Equal(res.Keywords[2].Keyword, "machina")
	is.Equal(len(res.Sentences), 1)
	is.Equal(res.Sentences[0].Text, "I really love Machina, who is the MachineBox mascot designed by Ashley McNamara.")
	is.Equal(res.Sentences[0].Start, 0)
	is.Equal(res.Sentences[0].End, 80)
	is.Equal(res.Sentences[0].Sentiment, 0.7128883004188538)
	is.Equal(len(res.Sentences[0].Entities), 3)
	is.Equal(res.Sentences[0].Entities[0].Text, "Machina")
	is.Equal(res.Sentences[0].Entities[0].Type, "person")
	is.Equal(res.Sentences[0].Entities[0].Start, 14)
	is.Equal(res.Sentences[0].Entities[0].End, 21)
	is.Equal(res.Sentences[0].Entities[1].Text, "MachineBox")
	is.Equal(res.Sentences[0].Entities[1].Type, "organization")
	is.Equal(res.Sentences[0].Entities[1].Start, 34)
	is.Equal(res.Sentences[0].Entities[1].End, 44)
	is.Equal(res.Sentences[0].Entities[2].Text, "Ashley McNamara")
	is.Equal(res.Sentences[0].Entities[2].Type, "person")
	is.Equal(res.Sentences[0].Entities[2].Start, 64)
	is.Equal(res.Sentences[0].Entities[2].End, 79)
}
