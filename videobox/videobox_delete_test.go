package videobox_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/machinebox/sdk-go/videobox"
	"github.com/matryer/is"
)

func TestDelete(t *testing.T) {
	is := is.New(t)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		is.Equal(r.Method, "DELETE")
		is.Equal(r.URL.Path, "/videobox/results/5a50b8067eced76bad103c53dd0f5226")
		is.Equal(r.Header.Get("Accept"), "application/json; charset=utf-8")
	}))
	defer srv.Close()
	vb := videobox.New(srv.URL)
	err := vb.Delete("5a50b8067eced76bad103c53dd0f5226")
	is.NoErr(err)
}
