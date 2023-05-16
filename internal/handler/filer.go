package handler

import (
	"net/http"
)

func memberTest(r *http.Request, param string) []string {
	vals := []string{}
	for _, stringVal := range r.URL.Query()[param] {
		if stringVal != "" {
			vals = append(vals, stringVal)
		}
	}
	if len(vals) == 0 {
		vals = []string{"news", "sport", "music", "kids", "hobbies", "programming", "art", "cooking", "other"}
	}
	return vals
}
