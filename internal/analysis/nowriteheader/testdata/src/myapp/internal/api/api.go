package api

import "net/http"

func BadHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest) // want `use chikit\.SetResponse or chikit\.SetError`
	_ = r
}

func AlsoBad(w http.ResponseWriter) {
	w.WriteHeader(500) // want `use chikit\.SetResponse or chikit\.SetError`
}

func GoodHandler(w http.ResponseWriter, r *http.Request) {
	_ = w
	_ = r
}
