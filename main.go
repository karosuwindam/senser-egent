package main

import (
	"context"
	"os"
	"os/signal"
	"senseregent/config"
	"senseregent/webserver"
	"syscall"
)

func init() {
	if err := config.Init(); err != nil {
		panic(err)
	}
	if err := webserver.Init(); err != nil {
		panic(err)
	}

}

func main() {

	ctx := context.Background()
	ctxmain, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	config.TracerStart(config.TraData.GrpcURL, config.TraData.ServiceName, ctxmain)
	defer config.TracerStop(ctxmain)
	ctxWeb, canncelweb := context.WithCancel(ctx)
	go func(ctx context.Context) {
		if err := webserver.Start(ctx); err != nil {
			panic(err)
		}
	}(ctxWeb)

	<-ctxmain.Done()

	webserver.Stop(ctxWeb)

	canncelweb()

}
