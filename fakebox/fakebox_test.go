package fakebox_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/machinebox/sdk-go/fakebox"
	"github.com/matryer/is"
)

func TestCheck(t *testing.T) {
	result := `
	{
	"success": true,
	"title": {
		"decision": "impartial",
		"score": 0.7757045030593872,
		"entities": [
			{
				"text": "China",
				"start": 0,
				"end": 4,
				"type": "place"
			}
		]
	},
	"content": {
		"decision": "bias",
		"score": 0.33232277631759644,
		"entities": [
			{
				"text": "Fuxing",
				"start": 21,
				"end": 26,
				"type": "place"
			},
			{
				"text": "300",
				"start": 74,
				"end": 76,
				"type": "cardinal"
			},
			{
				"text": "186mph",
				"start": 83,
				"end": 88,
				"type": "quantity"
			},
			{
				"text": "2011",
				"start": 94,
				"end": 97,
				"type": "date"
			},
			{
				"text": "two",
				"start": 109,
				"end": 111,
				"type": "cardinal"
			},
			{
				"text": "40",
				"start": 133,
				"end": 134,
				"type": "cardinal"
			},
			{
				"text": "next week",
				"start": 149,
				"end": 157,
				"type": "date"
			},
			{
				"text": "about 350",
				"start": 234,
				"end": 242,
				"type": "cardinal"
			}
		],
		"keywords": [
			{
				"keyword": "high speed"
			},
			{
				"keyword": "bullet train"
			},
			{
				"keyword": "speed"
			},
			{
				"keyword": "train"
			},
			{
				"keyword": "mph"
			},
			{
				"keyword": "km/h"
			},
			{
				"keyword": "rejuvenation"
			},
			{
				"keyword": "fuxing"
			},
			{
				"keyword": "crash"
			},
			{
				"keyword": "people"
			}
		]
	},
	"domain": {
		"domain": "bbc.co.uk",
		"category": "trusted"
	}
	}
	`
	is := is.New(t)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		is.Equal(r.Method, "POST")
		is.Equal(r.URL.Path, "/fakebox/check")
		is.Equal(r.Header.Get("Accept"), "application/json; charset=utf-8")
		io.WriteString(w, result)
	}))
	defer srv.Close()
	u, err := url.Parse("http://www.bbc.co.uk/news/technology-41011662")
	is.NoErr(err)

	tb := fakebox.New(srv.URL)
	res, err := tb.Check(`China relaunches world's fastest train`,
		`The top speed of the Fuxing or "rejuvenation" bullet trains was capped at 300km/h (186mph) in 2011 following two crashes that killed 40 people.
From next week, some of the trains will once again be allowed to run at a higher speed of about 350 km/h.`,
		u,
	)
	is.NoErr(err)

	// Title
	is.Equal(res.Title.Decision, "impartial")
	is.True(res.Title.Score > 0.6)
	is.Equal(len(res.Title.Entities), 1)
	is.Equal(res.Title.Entities[0].Start, 0)
	is.Equal(res.Title.Entities[0].End, 4)
	is.Equal(res.Title.Entities[0].Text, "China")
	is.Equal(res.Title.Entities[0].Type, "place")

	// Domain
	is.Equal(res.Domain.Category, "trusted")
	is.Equal(res.Domain.Domain, "bbc.co.uk")

	// Content
	is.Equal(res.Content.Decision, "bias")
	is.True(res.Content.Score < 0.4)
	is.Equal(len(res.Content.Keywords), 10)
	is.Equal(res.Content.Keywords[0].Keyword, "high speed")
	is.Equal(res.Content.Keywords[1].Keyword, "bullet train")
	is.Equal(len(res.Content.Entities), 8)
	is.Equal(res.Content.Entities[0].Start, 21)
	is.Equal(res.Content.Entities[0].End, 26)
	is.Equal(res.Content.Entities[0].Text, "Fuxing")
	is.Equal(res.Content.Entities[0].Type, "place")

}
