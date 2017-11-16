package parser

// url list
type Urls []string

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
	ContentLength int    `json:"content-length"`
}

type Element struct {
	TagName string `json:"tag-name"`
	Count   uint   `json:"count"`
}
