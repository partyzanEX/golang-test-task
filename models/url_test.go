package models

import (
	"testing"
	"net/http/httptest"
	"net/http"
	"strings"
	"encoding/json"
	"io/ioutil"
)

var req = []string{"https://google.com", "https://yandex.ru", "https://yandex.ru", "https://yandex.ru/qwerty"}

func TestUrls_SetFromRequest(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}
		defer r.Body.Close()

		urls := Urls{}
		urls.SetFromBody(body)
		if len(urls) != 3 {
			t.Error("Invalid urls lenght, expected 4")
		}
	}))

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatal(err)
	}

	_, err = http.Post(ts.URL, "application/json", strings.NewReader(string(data)))
	if err != nil {
		t.Fatal(err)
	}
}
