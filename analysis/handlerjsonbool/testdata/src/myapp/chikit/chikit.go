// Package chikit is a test stub of the real github.com/nhalm/chikit
// surface used by the analyzer. analysistest's GOPATH-style testdata
// can't reach external modules, so we stand the relevant shape up here.
package chikit

import "net/http"

func JSON(r *http.Request, v any) bool  { return true }
func Query(r *http.Request, v any) bool { return true }
