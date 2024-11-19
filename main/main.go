package main

import (
	"context"
	"os"
	"os/signal"
	"senseregent/config"
	"senseregent/controller"
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
	if err := controller.Init(); err != nil {
		panic(err)
	}

}

func main() {

	ctx := context.Background()
	ctxmain, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	config.TracerStart(config.TraData.GrpcURL, config.TraData.ServiceName, ctxmain)
	defer config.TracerStop(ctxmain)
	ctxController, canncelController := context.WithCancel(ctx)
	go func(ctx context.Context) {
		if err := controller.Run(ctx); err != nil {
			panic(err)
		}
	}(ctxController)
	ctxWeb, canncelweb := context.WithCancel(ctx)
	go func(ctx context.Context) {
		if err := webserver.Start(ctx); err != nil {
			panic(err)
		}
	}(ctxWeb)

	<-ctxmain.Done()

	webserver.Stop(context.Background())
	controller.Stop(context.Background())
	canncelController()

	canncelweb()

}
