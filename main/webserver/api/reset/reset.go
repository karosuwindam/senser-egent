package reset

import (
	"log/slog"
	"net/http"
	"senseregent/config"
)

func getReset(w http.ResponseWriter, r *http.Request) {
	ctx, span := config.TracerS(r.Context(), "getReset", "reset")
	defer span.End()
	slog.DebugContext(ctx, "getReset")
	//なんにもせずに、/metricsに移動を指示
}

func postReset(w http.ResponseWriter, r *http.Request) {
	ctx, span := config.TracerS(r.Context(), "postReset", "reset")
	defer span.End()
	slog.DebugContext(ctx, "postReset")
	//センサーの定期取得処理をリセットする
	//OKと結果を返す
	w.Write([]byte("OK"))
}
