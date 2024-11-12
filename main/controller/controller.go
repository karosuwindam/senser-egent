package controller

import (
	"context"
	"fmt"
	"time"
)

var shutdown chan struct{}
var done chan struct{}

func Init() error {

	shutdown = make(chan struct{}, 1)
	done = make(chan struct{}, 1)
	return nil
}

func Run(ctx context.Context) error {
	go func(ctx context.Context) {
		for {
		}
	}(ctx)
	select {
	case <-shutdown:
		done <- struct{}{}
	}
	return nil
}

func Stop(ctx context.Context) error {
	shutdown <- struct{}{}
	select {
	case <-done:
		break
	case <-time.After(5 * time.Second):
		return fmt.Errorf("time over 5 sec")
	}
	return nil
}

type API struct {
}
