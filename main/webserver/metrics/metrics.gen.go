package metrics

import (
	"net/http"
	"senseregent/webserver/common"
)

func Init(route string, mux *http.ServeMux) error {
	if route == "/" {
		route = ""
	}
	common.TraceHandleFunc(mux, "GET "+route, getMetrics)

	return nil
}
