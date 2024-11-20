package reset

import (
	"net/http"
)

func Init(route string, mux *http.ServeMux) error {
	if route == "/" {
		route = ""
	}
	mux.HandleFunc("GET "+route, getReset)
	mux.HandleFunc("POST "+route, postReset)
	return nil
}
