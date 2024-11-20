package config

import (
	"github.com/caarlos0/env/v6"
)

// Webサーバの設定
type WebConfig struct {
	Protocol   string `env:"WEB_PROTOCOL" envDefault:"tcp"`  //接続プロトコル
	Hostname   string `env:"WEB_HOST" envDefault:""`         //接続DNS名
	Port       string `env:"WEB_PORT" envDefault:"8080"`     //接続ポート
	StaticPage string `env:"WEB_FOLDER" envDefault:"./html"` //静的ページの参照先
}

type TracerData struct {
	GrpcOn bool `env:"TRACER_GRPC_ON" envDefault:"true"`
	// GrpcURL     string `env:"TRACER_GRPC_URL" envDefault:"localhost:4317"`
	GrpcURL     string `env:"TRACER_GRPC_URL" envDefault:"otel-grpc.bookserver.home:4317"`
	HttpURL     string `env:"TRACER_HTTP_URL" envDefault:"localhost:4318"`
	ServiceName string `env:"TRACER_SERVICE_NAME" envDefault:"senser-egent-test"`
	TracerUse   bool   `env:"TRACER_ON" envDefault:"true"`
}

var Web WebConfig
var TraData TracerData

// 環境設定
func Init() error {
	if err := env.Parse(&Web); err != nil {
		return err
	}
	if err := env.Parse(&TraData); err != nil {
		return err
	}
	return nil
}
