package httputil

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/mdzio/go-lib/logging"
)

func init() {
	var l logging.LogLevel
	err := l.Set(os.Getenv("LOG_LEVEL"))
	if err == nil {
		logging.SetLevel(l)
	}
}

func TestSingleAuthHandler(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("test"))
	})
	auth := &SingleAuthHandler{
		Handler:  h,
		User:     "My User",
		Password: "My Password",
		Realm:    "The Realm",
	}
	srv := httptest.NewServer(auth)
	defer srv.Close()

	resp, err := http.Get(srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusUnauthorized {
		t.Error(resp.StatusCode)
	}
	if resp.Header.Get("WWW-Authenticate") != "Basic realm=\"The Realm\", charset=\"UTF-8\"" {
		t.Error(resp.Header.Get("WWW-Authenticate"))
	}

	req, _ := http.NewRequest(http.MethodGet, srv.URL, nil)
	req.SetBasicAuth("My User", "My Password")
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Error(resp.StatusCode)
	}
	b, _ := ioutil.ReadAll(resp.Body)
	if string(b) != "test" {
		t.Error(string(b))
	}

	req, _ = http.NewRequest(http.MethodGet, srv.URL, nil)
	req.SetBasicAuth("My User2", "My Password")
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusUnauthorized {
		t.Error(resp.StatusCode)
	}
	if resp.Header.Get("WWW-Authenticate") != "Basic realm=\"The Realm\", charset=\"UTF-8\"" {
		t.Error(resp.Header.Get("WWW-Authenticate"))
	}
}
