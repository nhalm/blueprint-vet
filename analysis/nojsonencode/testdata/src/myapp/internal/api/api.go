package api

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func BadHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{"err": "bad"}) // want `use chikit\.SetResponse / chikit\.SetError`
	_ = r
}

func GoodEncodeToBuffer(w http.ResponseWriter) {
	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(map[string]string{"ok": "yes"})
	_ = w
}

func GoodHandler(w http.ResponseWriter, r *http.Request) {
	_ = w
	_ = r
}
