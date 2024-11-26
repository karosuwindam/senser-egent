package bme280

import (
	"context"
	"senseregent/controller/senser/i2c_senser/common"
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
	i2c = common.Init(BME280, I2C_BUS)
	return &API{}
}

func (api *API) Test(ctx context.Context) bool {
	return readICIDCheck(ctx)
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
		return BME280_SLEEP_MODE
	}
	api.calib = calibRead(ctx)
	return nil
}

func (api *API) ReadData(ctx context.Context) error {
	api.mu.Lock()
	defer api.mu.Unlock()
	if readStatus(ctx) == CtrMeasReg_Sleep {
		return BME280_SLEEP_MODE
	}
	if len(api.calib.hum) == 0 {
		api.calib = calibRead(ctx)
	}
	api.Press, api.Tmp, api.Hum = readSenserData(ctx, api.calib)
	return nil
}
