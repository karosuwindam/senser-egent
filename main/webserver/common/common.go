package common

import (
	"net/http"
	"senseregent/config"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func TraceHandleFunc(mux *http.ServeMux, pattern string, handlerFunc func(http.ResponseWriter, *http.Request)) {
	if config.TraData.TracerUse {
		handler := otelhttp.NewHandler(http.HandlerFunc(handlerFunc), pattern)
		mux.Handle(pattern, handler)
		return
	}
	mux.HandleFunc(pattern, handlerFunc)
}
