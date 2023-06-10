// go test -v pitfall_test.go

package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandlerWriteHeader(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("Bad request"))
		w.WriteHeader(http.StatusBadRequest)
	}
	r := httptest.NewRequest(http.MethodGet, "http://test", nil)
	w := httptest.NewRecorder()
	handler(w, r)
	t.Logf("Response Status: %q", w.Result().Status)

	handler = func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Bad request"))
	}
	r = httptest.NewRequest(http.MethodGet, "http://test", nil)
	w = httptest.NewRecorder()
	handler(w, r)
	t.Logf("Response Status: %q", w.Result().Status)

	handler = func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Bad request", http.StatusBadRequest)
	}
	r = httptest.NewRequest(http.MethodGet, "http://test", nil)
	w = httptest.NewRecorder()
	handler(w, r)
	t.Logf("Response Status: %q", w.Result().Status)
}

/*

=== RUN   TestHandlerWriteHeader
    pitfall_test.go:19: Response Status: "200 OK"
    pitfall_test.go:28: Response Status: "400 Bad Request"
    pitfall_test.go:36: Response Status: "400 Bad Request"
--- PASS: TestHandlerWriteHeader (0.00s)
PASS

*/
