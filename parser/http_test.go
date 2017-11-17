package parser

import (
	"encoding/json"
	"testing"
	"github.com/partyzanex/golang-test-task/log"
	"strings"
	"io/ioutil"
	"fmt"
	"github.com/partyzanex/golang-test-task/conf"
	"github.com/partyzanex/golang-test-task/models"
	"net/http/httptest"
	"net/http"
)

func TestHandler_ServeHTTP(t *testing.T) {
	html := "123<html lang=\"en\">X<head><meta charset=\"UTF-8\"><title></title></head>T<body><b></b><b></b></body></html>qwerty"
	expected := map[string]uint{
		"html": 1,
		"head": 1,
		"title": 1,
		"meta": 1,
		"body": 1,
		"b": 2,
	}

	t.Run("Test elements", func(t *testing.T) {
		// handler
		logger, _ := log.NewLogger("", true)
		config, _ := conf.NewConfig("../config.json")
		s := NewHttp(config, logger)
		s.Host = "localhost"
		s.Port = "3018"
		//go s.Serve()
		go func(h *Http) {
			s.Serve()
		}(s)

		go func(h *Http) {
			// test page
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("content-type", "application/json; charset=utf-8")
				w.Header().Set("X-test", "test-header")
				fmt.Fprint(w, html)
			}))

			data, err := json.Marshal([]string{ts.URL})
			resp, err := http.Post("http://" + s.GetAddr(), "application/json", strings.NewReader(string(data)))
			if err != nil {
				t.Fatal(err)
			}

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			var urls []models.UrlInfo
			err = json.Unmarshal(body, &urls)
			if err != nil {
				t.Fatal(err)
			}

			if len(urls) != 1 {
				t.Error("No lenght, expected 1")
			}

			for _, urlInfo := range urls {
				if urlInfo.Meta.Status != http.StatusOK {
					t.Error("Invalid http status, expected 200, host: " + urlInfo.Url)
				}
				if len(urlInfo.Elements) != 6 {
					t.Error("Invalid elements lenght, expected 6")
				}

				for _, elem := range urlInfo.Elements {
					if count, ok := expected[elem.TagName]; !ok {
						t.Error("Tag " + elem.TagName + " is not found")
					} else {
						if count != elem.Count {
							t.Error("Number of tags does not match, tags: " + string(elem.Count) + ", expected: " + string(count))
						}
					}
				}

				if urlInfo.Meta.ContentType != "text/html" {
					t.Error("Invalid content-type")
				}
			}
		}(s)
	})
}