package mbhttp_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/machinebox/sdk-go/internal/mbhttp"
	"github.com/matryer/is"
)

func TestDo(t *testing.T) {
	is := is.New(t)
	type obj struct {
		Field1 string `json:"field1"`
		Field2 int    `json:"field2"`
		Field3 bool   `json:"field3"`
	}
	in := obj{Field1: "in", Field2: 123, Field3: true}
	out := obj{Field1: "in", Field2: 123, Field3: true}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var requestObj obj
		is.Equal(r.Method, http.MethodPost)
		is.Equal(r.URL.Path, "/something")
		is.NoErr(json.NewDecoder(r.Body).Decode(&requestObj))
		is.Equal(requestObj.Field1, in.Field1)
		is.Equal(requestObj.Field2, in.Field2)
		is.Equal(requestObj.Field3, in.Field3)
		is.NoErr(json.NewEncoder(w).Encode(struct {
			Success bool `json:"success"`
			obj
		}{
			Success: true,
			obj:     out,
		}))
	}))
	defer srv.Close()
	var buf bytes.Buffer
	is.NoErr(json.NewEncoder(&buf).Encode(in))
	req, err := http.NewRequest(http.MethodPost, srv.URL+"/something", &buf)
	c := mbhttp.New("testbox", http.DefaultClient)
	var actualOut obj
	resp, err := c.Do(req, &actualOut)
	is.NoErr(err)
	defer resp.Body.Close()
	is.Equal(actualOut.Field1, out.Field1)
	is.Equal(actualOut.Field2, out.Field2)
	is.Equal(actualOut.Field3, out.Field3)
}

func TestDoBoxError(t *testing.T) {
	is := is.New(t)
	type obj struct {
		Field1 string `json:"field1"`
		Field2 int    `json:"field2"`
		Field3 bool   `json:"field3"`
	}
	in := obj{Field1: "in", Field2: 123, Field3: true}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var requestObj obj
		is.Equal(r.Method, http.MethodPost)
		is.Equal(r.URL.Path, "/something")
		is.NoErr(json.NewDecoder(r.Body).Decode(&requestObj))
		is.Equal(requestObj.Field1, in.Field1)
		is.Equal(requestObj.Field2, in.Field2)
		is.Equal(requestObj.Field3, in.Field3)
		is.NoErr(json.NewEncoder(w).Encode(struct {
			Success bool   `json:"success"`
			Error   string `json:"error"`
		}{
			Success: false,
			Error:   "something went wrong",
		}))
	}))
	defer srv.Close()
	var buf bytes.Buffer
	is.NoErr(json.NewEncoder(&buf).Encode(in))
	req, err := http.NewRequest(http.MethodPost, srv.URL+"/something", &buf)
	c := mbhttp.New("testbox", http.DefaultClient)
	var actualOut obj
	_, err = c.Do(req, &actualOut)
	is.True(err != nil)
	is.Equal(err.Error(), "testbox: something went wrong")
}

func TestDoBoxMissingError(t *testing.T) {
	is := is.New(t)
	type obj struct {
		Field1 string `json:"field1"`
		Field2 int    `json:"field2"`
		Field3 bool   `json:"field3"`
	}
	in := obj{Field1: "in", Field2: 123, Field3: true}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var requestObj obj
		is.Equal(r.Method, http.MethodPost)
		is.Equal(r.URL.Path, "/something")
		is.NoErr(json.NewDecoder(r.Body).Decode(&requestObj))
		is.Equal(requestObj.Field1, in.Field1)
		is.Equal(requestObj.Field2, in.Field2)
		is.Equal(requestObj.Field3, in.Field3)
		is.NoErr(json.NewEncoder(w).Encode(struct {
			Success bool   `json:"success"`
			Error   string `json:"error"`
		}{
			Success: false,
		}))
	}))
	defer srv.Close()
	var buf bytes.Buffer
	is.NoErr(json.NewEncoder(&buf).Encode(in))
	req, err := http.NewRequest(http.MethodPost, srv.URL+"/something", &buf)
	c := mbhttp.New("testbox", http.DefaultClient)
	var actualOut obj
	_, err = c.Do(req, &actualOut)
	is.True(err != nil)
	is.Equal(err.Error(), "testbox: an unknown error occurred in the box")
}

func TestDoHTTPError(t *testing.T) {
	is := is.New(t)
	type obj struct {
		Field1 string `json:"field1"`
		Field2 int    `json:"field2"`
		Field3 bool   `json:"field3"`
	}
	in := obj{Field1: "in", Field2: 123, Field3: true}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var requestObj obj
		is.Equal(r.Method, http.MethodPost)
		is.Equal(r.URL.Path, "/something")
		is.NoErr(json.NewDecoder(r.Body).Decode(&requestObj))
		is.Equal(requestObj.Field1, in.Field1)
		is.Equal(requestObj.Field2, in.Field2)
		is.Equal(requestObj.Field3, in.Field3)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
	}))
	defer srv.Close()
	var buf bytes.Buffer
	is.NoErr(json.NewEncoder(&buf).Encode(in))
	req, err := http.NewRequest(http.MethodPost, srv.URL+"/something", &buf)
	c := mbhttp.New("testbox", http.DefaultClient)
	var actualOut obj
	_, err = c.Do(req, &actualOut)
	is.True(err != nil)
	is.Equal(err.Error(), "testbox: 500 Internal Server Error")
}
