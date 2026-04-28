package api

import (
	"net/http"

	"myapp/chikit"
)

type req struct{}

func BadDiscard(w http.ResponseWriter, r *http.Request) {
	chikit.JSON(r, &req{}) // want `chikit\.JSON result discarded`
	_ = w
}

func BadDiscardQuery(w http.ResponseWriter, r *http.Request) {
	chikit.Query(r, &req{}) // want `chikit\.Query result discarded`
	_ = w
}

func BadBlankAssign(w http.ResponseWriter, r *http.Request) {
	_ = chikit.JSON(r, &req{}) // want `chikit\.JSON result discarded via _`
	_ = w
}

func GoodGuard(w http.ResponseWriter, r *http.Request) {
	if !chikit.JSON(r, &req{}) {
		return
	}
	_ = w
}

func GoodInitGuard(w http.ResponseWriter, r *http.Request) {
	if ok := chikit.JSON(r, &req{}); !ok {
		return
	}
	_ = w
}
