package facebox_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/machinebox/sdk-go/facebox"
	"github.com/matryer/is"
)

func TestCompareFaceprints(t *testing.T) {
	is := is.New(t)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		is.Equal(r.URL.Path, "/facebox/faceprint/compare")
		is.Equal(r.Header.Get("Accept"), "application/json; charset=utf-8")
		io.WriteString(w, `{
			"success": true,			
			"confidences": [ 0.1, 0.2, 0.3]
		}`)
	}))
	defer srv.Close()

	fb := facebox.New(srv.URL)
	con, err := fb.CompareFaceprints("target", []string{"candidate1", "candidate2", "candidate3"})
	is.NoErr(err)

	is.Equal(len(con), 3)
	is.Equal(con[0], 0.1)
	is.Equal(con[1], 0.2)
	is.Equal(con[2], 0.3)
}
