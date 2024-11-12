package bme280

import (
	"context"
	"sync"
)

type API struct {
	calib bme280_cal
	Press float64 //気圧
	Tmp   float64 //温度
	Hum   float64 //気圧
	mu    sync.Mutex
}

func APIInit() *API {
	return &API{}
}

func (api *API) Up(ctx context.Context) {
	up(ctx)
}

func (api *API) Down(ctx context.Context) {
	down(ctx)
}

func (api *API) CalibRead(ctx context.Context) {
	api.mu.Lock()
	defer api.mu.Unlock()
	api.calib = calibRead(ctx)
}

func (api *API) ReadData(ctx context.Context) {
	api.mu.Lock()
	defer api.mu.Unlock()
	api.Press, api.Tmp, api.Hum = readSenserData(ctx, api.calib)
}
