package jsons

import (
	"log/slog"
	"net/http"
	"senseregent/config"
)

func getJsons(w http.ResponseWriter, r *http.Request) {
	ctx, span := config.TracerS(r.Context(), "getJsons", "jsons")
	defer span.End()
	slog.WarnContext(ctx, "getJsons")
}
