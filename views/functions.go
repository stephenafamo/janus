package views

import (
	"encoding/json"
	"html/template"
	"net/url"
	"strings"
	"time"
)

var funcMap = template.FuncMap{
	"ToUpper": strings.ToUpper,
	"ToLower": strings.ToLower,
	"ToTitle": strings.Title,
	"Year": func() int {
		year, _, _ := time.Now().Date()
		return year
	},
	"divFloat": func(a, b float64) float64 {
		return a / b
	},
	"divInt": func(a, b int) float64 {
		return float64(a) / float64(b)
	},
	"divUint": func(a, b uint) float64 {
		return float64(a) / float64(b)
	},
	"Concatenate": func(ss ...string) string {
		return strings.Join(ss, "")
	},
	"ToHTML": func(src string) template.HTML {
		return template.HTML(src)
	},
	"ToJSON": func(src interface{}) (string, error) {
		bytes, err := json.Marshal(src)
		return string(bytes), err
	},
	"Now": func() time.Time {
		return time.Now()
	},
	"Iterate": func(count uint) []uint {
		var i uint
		var Items []uint
		for i = 0; i < (count); i++ {
			Items = append(Items, i)
		}
		return Items
	},
	"nextURL": func(url url.URL, after string) (string, error) {
		aURL := url // make a copy
		vals := aURL.Query()
		vals.Set("after", after)
		vals.Del("before")

		aURL.RawQuery = vals.Encode()
		return aURL.String(), nil
	},
	"previousURL": func(url url.URL, before string) (string, error) {
		aURL := url // make a copy
		vals := aURL.Query()
		vals.Set("before", before)
		vals.Del("after")

		aURL.RawQuery = vals.Encode()
		return aURL.String(), nil
	},
}
