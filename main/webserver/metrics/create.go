package metrics

import (
	"log/slog"
	"net/http"
	"senseregent/config"
)

func getMetrics(w http.ResponseWriter, r *http.Request) {
	ctx, span := config.TracerS(r.Context(), "getMetrics", "metrics")
	defer span.End()
	slog.DebugContext(ctx, "getMetrics")
	//Prometheusが回収できる形式に出力する
}
