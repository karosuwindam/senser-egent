package jsons

import (
	"net/http"
	"senseregent/webserver/common"
)

func Init(route string, mux *http.ServeMux) error {
	if route == "/" {
		route = ""
	}
	common.TraceHandleFunc(mux, "GET "+route, getJsons)
	return nil
}
