package jsons

import (
	"net/http"
)

func Init(route string, mux *http.ServeMux) error {
	if route == "/" {
		route = ""
	}
	mux.HandleFunc("GET "+route, getJsons)
	return nil
}
