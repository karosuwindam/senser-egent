package config

import (
	"log/slog"

	"github.com/caarlos0/env/v6"
)

// Webサーバの設定
type WebConfig struct {
	Protocol string `env:"WEB_PROTOCOL" envDefault:"tcp"` //接続プロトコル
	Hostname string `env:"WEB_HOST" envDefault:""`        //接続DNS名
	Port     string `env:"WEB_PORT" envDefault:"8080"`    //接続ポート
}

var Web WebConfig

// 環境設定
func Init() error {
	if err := env.Parse(&Web); err != nil {
		return err
	}
	logger := LoggerConfig(slog.LevelInfo)
	slog.SetDefault(logger)
	return nil
}
