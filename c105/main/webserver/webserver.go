package webserver

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"senseregent/config"
	"senseregent/webserver/api"
	"senseregent/webserver/index"
	"senseregent/webserver/metrics"
	"time"
)

// SetupServer
// サーバ動作の設定
type SetupServer struct {
	protocol string // Webサーバーのプロトコル
	hostname string //Webサーバのホスト名
	port     string //Webサーバの解放ポート

	mux *http.ServeMux //webサーバのmux
}

// Server
// Webサーバの管理情報
type Server struct {
	// Webサーバの管理関数
	srv *http.Server
	// 解放の管理関数
	l net.Listener
}

var srv *http.Server // Webサーバの管理関数

var cfg SetupServer

var shutdown chan bool
var done chan bool

func Init() error {
	shutdown = make(chan bool, 1)
	done = make(chan bool, 1)
	cfg = SetupServer{
		protocol: config.Web.Protocol,
		hostname: config.Web.Hostname,
		port:     config.Web.Port,
		mux:      http.NewServeMux(),
	}
	if err := api.Init(cfg.mux); err != nil {
		return err
	}
	metrics.Init("/metrics", cfg.mux)
	index.Init(cfg.mux)
	return nil
}

func Start(ctx context.Context) error {
	var err error = nil
	srv = &http.Server{
		Addr:         cfg.hostname + ":" + cfg.port,
		BaseContext:  func(_ net.Listener) context.Context { return ctx },
		Handler:      cfg.mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	l, err := net.Listen(cfg.protocol, srv.Addr)
	if err != nil {
		return err
	}
	slog.InfoContext(ctx, "Start Server", "IP", cfg.hostname, "Port", cfg.port)
	go func() {
		if err = srv.Serve(l); err != nil && err != http.ErrServerClosed {
			panic(err)
		} else {
			err = nil
		}
	}()
	select {
	case <-shutdown:
		done <- true
		break
	case <-ctx.Done():
		return nil
	}
	return err
}

func Stop(ctx context.Context) error {
	if srv == nil {
		return nil
	}
	ctx, _ = context.WithTimeout(ctx, time.Second)
	if err := srv.Shutdown(ctx); err != nil {
		return err
	}
	shutdown <- true

	select {
	case <-done:
		slog.InfoContext(ctx, "Server Shutdown")
		break
	case <-ctx.Done():
		slog.ErrorContext(ctx, "Server Shutdown", "Error", ctx.Err())
		break
	case <-time.After(time.Microsecond * 500):
		slog.ErrorContext(ctx, "Server Shutdown time out over 500 ms")
		break
	}
	return nil
}
