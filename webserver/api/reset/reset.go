package reset

import (
	"log/slog"
	"net/http"
	"senseregent/config"
)

func getReset(w http.ResponseWriter, r *http.Request) {
	ctx, span := config.TracerS(r.Context(), "getReset", "reset")
	defer span.End()
	slog.WarnContext(ctx, "getReset")
}

func postReset(w http.ResponseWriter, r *http.Request) {
	ctx, span := config.TracerS(r.Context(), "postReset", "reset")
	defer span.End()
	slog.WarnContext(ctx, "postReset")
}
