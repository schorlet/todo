package todo

import (
	"encoding/json"
	"html/template"
	"net/http"
)

// writeJSON encodes the specified value in the given Response.
func writeJSON(w http.ResponseWriter, v interface{}, status int) error {
	// set application/json header
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	// json encoder
	var encoder = json.NewEncoder(w)

	// encodes interface to json
	return encoder.Encode(v)
}

// readJSON decodes the specified value from the given Request.
func readJSON(r *http.Request, v interface{}) error {
	// json decoder
	var decoder = json.NewDecoder(r.Body)

	// decodes json body to interface
	return decoder.Decode(v)
}

func templateJSRaw(v interface{}) template.JS {
	var a, _ = json.Marshal(v)
	return template.JS(a)
}

func templateJSStr(v interface{}) template.JSStr {
	var a, _ = json.Marshal(v)
	return template.JSStr(a)
}
