package metrics

import (
	"log/slog"
	"net/http"
	"senseregent/controller"
)

func getMetrics(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	controllerAPI := controller.NewAPI()
	value, err := controllerAPI.ReadValue(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "ReadValue error", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//PromQL形式に変換してデータを返す
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(value.ToPromQL()))

}