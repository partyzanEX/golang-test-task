package models

import (
	"encoding/json"
	"io"
	"golang.org/x/net/html"
)

// url list
type Urls []string

// parsing urls from request.Body
func (u *Urls) SetFromBody(body []byte) error {
	err := json.Unmarshal(body, &u)
	if err != nil {
		return err
	}

	u.unique()
	return nil
}

// delete duplicate URL
func (u *Urls) unique() {
	arr := *u
	a := make([]string, 0, len(arr))
	b := make(map[string]bool)

	for _, val := range arr {
		if _, ok := b[val]; !ok {
			b[val] = true
			a = append(a, val)
		}
	}

	*u = a
}

// structure for json-response
type UrlInfo struct {
	Url      string    `json:"url"`
	Meta     Meta      `json:"meta"`
	Elements []Element `json:"elements"`
	counter  map[string]uint
}

type Meta struct {
	Status        int    `json:"status"`
	ContentType   string `json:"content-type"`
	ContentLength int64  `json:"content-length"`
}

type Element struct {
	TagName string `json:"tag-name"`
	Count   uint   `json:"count"`
}

// parsing amd counting of elements
func (ui *UrlInfo) SetElements(body io.Reader) {
	tok := html.NewTokenizer(body)
	ui.counter = make(map[string]uint)

	for {
	repeat:
		switch tok.Next() {
		case html.ErrorToken:
			goto end
		case html.StartTagToken, html.SelfClosingTagToken:
			name, _ := tok.TagName()
			if name == nil {
				goto repeat
			}

			ui.AddTag(name)
		}
	}

end:
	for tagName, count := range ui.counter {
		ui.Elements = append(ui.Elements, Element{TagName: tagName, Count: count})
	}
}

// counting tags
func (ui *UrlInfo) AddTag(name []byte) {
	tagName := string(name)
	ui.counter[tagName]++
}
