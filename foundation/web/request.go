package web

import (
	"encoding/json"
	"github.com/dimfeld/httptreemux/v5"
	"net/http"
)

// Param returns param from the request based on route's wildcards.
func Param(r *http.Request, key string) string {
	 pm := httptreemux.ContextParams(r.Context())
	 return pm[key]
}

// Decode reads the body of an HTTP request for a JSON document.
// If val is a struct - then it is checked for validation tags.
func Decode(r *http.Request, val interface{}) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(val); err != nil {
		return err
	}

	return nil
}
