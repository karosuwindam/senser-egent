package index

import (
	"fmt"
	"log/slog"
	"net/http"
)

func getIndex(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	slog.DebugContext(ctx, r.Method+":"+r.URL.Path, "Method", r.Method, "Path", r.URL.Path, "RemoteAddr", r.RemoteAddr)

	output := "<html><body><a href=\"/metrics\">metrics</a></body></html>"
	fmt.Fprintf(w, "%s", output)
}
