package reset

import (
	"net/http"
	"senseregent/webserver/common"
)

func Init(route string, mux *http.ServeMux) error {
	if route == "/" {
		route = ""
	}
	common.TraceHandleFunc(mux, "GET "+route, getReset)
	common.TraceHandleFunc(mux, "POST "+route, postReset)
	return nil
}
