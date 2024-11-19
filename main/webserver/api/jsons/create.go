package jsons

import (
	"log/slog"
	"net/http"
	"senseregent/config"
	"senseregent/controller"
)

func getJsons(w http.ResponseWriter, r *http.Request) {
	ctx, span := config.TracerS(r.Context(), "getJsons", "jsons")
	defer span.End()
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
