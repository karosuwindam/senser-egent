package index

import (
	"net/http"
)

func Init(mux *http.ServeMux) error {
	mux.HandleFunc("/", getIndex)
	return nil
}
