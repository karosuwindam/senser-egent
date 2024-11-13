package bme280

import (
	"context"
	"fmt"
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
	api.calib = bme280_cal{}
}

func (api *API) CalibRead(ctx context.Context) error {
	api.mu.Lock()
	defer api.mu.Unlock()

	if readStatus(ctx) == CtrMeasReg_Sleep {
		return fmt.Errorf("BME280 Sleep Mode")
	}
	api.calib = calibRead(ctx)
	return nil
}

func (api *API) ReadData(ctx context.Context) error {
	api.mu.Lock()
	defer api.mu.Unlock()
	if readStatus(ctx) == CtrMeasReg_Sleep {
		return fmt.Errorf("BME280 Sleep Mode")
	}
	if len(api.calib.hum) == 0 {
		api.CalibRead(ctx)
	}
	api.Press, api.Tmp, api.Hum = readSenserData(ctx, api.calib)
	return nil
}
