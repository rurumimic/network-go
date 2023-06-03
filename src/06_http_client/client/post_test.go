// go test -v post_test.go

package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
)

type User struct {
	First string
	Last  string
}

func handlePostUser(t *testing.T) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func(r io.ReadCloser) {
			_, _ = io.Copy(io.Discard, r)
			_ = r.Close()
		}(r.Body)

		if r.Method != http.MethodPost {
			http.Error(w, "", http.StatusMethodNotAllowed)
			return
		}

		var u User
		err := json.NewDecoder(r.Body).Decode(&u)
		if err != nil {
			t.Error(err)
			http.Error(w, "Decode Failed", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusAccepted)
		t.Logf("User: %#v", u)
	}
}

func TestPostUser(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(handlePostUser(t)))
	defer ts.Close()

	resp, err := http.Get(ts.URL)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("Expected %d, got %d", http.StatusMethodNotAllowed, resp.StatusCode)
	}

	buf := new(bytes.Buffer)
	u := User{First: "Adam", Last: "Woodbeck"}
	err = json.NewEncoder(buf).Encode(u)
	if err != nil {
		t.Error(err)
	}

	resp, err = http.Post(ts.URL, "application/json", buf)
	if err != nil {
		t.Error(err)
	}
	t.Logf("StatusCode: %d", resp.StatusCode)
	if resp.StatusCode != http.StatusAccepted {
		t.Errorf("Expected %d, got %d", http.StatusAccepted, resp.StatusCode)
	}
	_ = resp.Body.Close()
}

/*

=== RUN   TestPostUser
    post_test.go:39: User: client.User{First:"Adam", Last:"Woodbeck"}
    post_test.go:66: StatusCode: 202
--- PASS: TestPostUser (0.00s)

*/

func TestMultipartPost(t *testing.T) {
	reqBody := new(bytes.Buffer)
	w := multipart.NewWriter(reqBody)

	for k, v := range map[string]string{
		"date":        time.Now().Format(time.RFC3339),
		"description": "From values with attached file",
	} {
		err := w.WriteField(k, v)
		if err != nil {
			t.Fatal(err)
		}
	}

	for i, file := range []string{
		"./files/hello.txt",
		"./files/goodbye.txt",
	} {
		t.Logf("Attaching file: %s", file)
		filePart, err := w.CreateFormFile(fmt.Sprintf("file%d", i+1), filepath.Base(file))
		if err != nil {
			t.Fatal(err)
		}

		f, err := os.Open(file)
		if err != nil {
			t.Fatal(err)
		}

		_, err = io.Copy(filePart, f)
		_ = f.Close()
		if err != nil {
			t.Fatal(err)
		}

	}

	err := w.Close()
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "http://httpbin.org/post", reqBody)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected %d, got %d", http.StatusOK, resp.StatusCode)
	}

	t.Logf("\n%s", b)
}

/*

=== RUN   TestMultipartPost
    post_test.go:106: Attaching file: ./files/hello.txt
    post_test.go:106: Attaching file: ./files/goodbye.txt
    post_test.go:155:
        {
          "args": {},
          "data": "",
          "files": {
            "file1": "Hello, world!",
            "file2": "Goodbye, world!"
          },
          "form": {
            "date": "2023-06-03T16:51:34+09:00",
            "description": "From values with attached file"
          },
          "headers": {
            "Accept-Encoding": "gzip",
            "Content-Length": "736",
            "Content-Type": "multipart/form-data; boundary=7ea725ca64f034e38e857c0a090f8d9abe87c1ccc9397ece5ecd374f838e",
            "Host": "httpbin.org",
            "User-Agent": "Go-http-client/1.1",
            "X-Amzn-Trace-Id": "Root=1-647af107-75050961471378f2540bf6e5"
          },
          "json": null,
          "origin": "121.170.132.15",
          "url": "http://httpbin.org/post"
        }
--- PASS: TestMultipartPost (2.93s)

*/
