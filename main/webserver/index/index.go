package index

import (
	"fmt"
	"net/http"
)

func getIndex(w http.ResponseWriter, r *http.Request) {
	output := "<html><body><a href=\"/metrics\">metrics</a></body></html>"
	fmt.Fprintf(w, "%s", output)
}
