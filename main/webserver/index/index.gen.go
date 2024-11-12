package index

import (
	"net/http"
	"senseregent/webserver/common"
)

func Init(mux *http.ServeMux) error {
	common.TraceHandleFunc(mux, "/", getIndex)
	return nil
}
