package jsons

import (
	"log/slog"
	"net/http"
	"senseregent/controller"
)

func getJsons(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	slog.DebugContext(ctx, r.Method+":"+r.URL.Path, "Method", r.Method, "Path", r.URL.Path, "RemoteAddr", r.RemoteAddr)

	controllerAPI := controller.NewAPI()
	value, err := controllerAPI.ReadValue(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "ReadValue error", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//センサーデータをJSON形式に変換してデータを返す
	jsonData := value.ToJson()
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(jsonData))
}
