package reset

import (
	"log/slog"
	"net/http"
	"senseregent/controller"
)

func getReset(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	slog.DebugContext(ctx, r.Method+":"+r.URL.Path, "Method", r.Method, "Path", r.URL.Path, "RemoteAddr", r.RemoteAddr)

	//なんにもせずに、/metricsに移動を指示
	http.Redirect(w, r, "/metrics", http.StatusFound)
}

func postReset(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	slog.DebugContext(ctx, r.Method+":"+r.URL.Path, "Method", r.Method, "Path", r.URL.Path, "RemoteAddr", r.RemoteAddr)

	controllerAPI := controller.NewAPI()

	slog.DebugContext(ctx, "postReset")
	//センサーの定期取得処理をリセットする
	if err := controllerAPI.Resetsenser(ctx); err != nil {
		slog.ErrorContext(ctx, "Resetsenser error", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//OKと結果を返す
	w.Write([]byte("OK"))
}
