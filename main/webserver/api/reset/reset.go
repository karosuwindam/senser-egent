package reset

import (
	"log/slog"
	"net/http"
	"senseregent/config"
	"senseregent/controller"
)

func getReset(w http.ResponseWriter, r *http.Request) {
	ctx, span := config.TracerS(r.Context(), "getReset", "reset")
	defer span.End()
	slog.DebugContext(ctx, r.Method+":"+r.URL.Path, "Method", r.Method, "Path", r.URL.Path, "RemoteAddr", r.RemoteAddr)
	//なんにもせずに、/metricsに移動を指示
	http.Redirect(w, r, "/metrics", http.StatusFound)
}

func postReset(w http.ResponseWriter, r *http.Request) {
	ctx, span := config.TracerS(r.Context(), "postReset", "reset")
	defer span.End()
	controllerAPI := controller.NewAPI()

	slog.DebugContext(ctx, r.Method+":"+r.URL.Path, "Method", r.Method, "Path", r.URL.Path, "RemoteAddr", r.RemoteAddr)
	//センサーの定期取得処理をリセットする
	if err := controllerAPI.ResetSenser(ctx); err != nil {
		slog.ErrorContext(ctx, "ResetSenser error", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//OKと結果を返す
	w.Write([]byte("OK"))
}
