package bme280

import (
	"context"
	"log/slog"
	"senseregent/controller/senser/i2c_senser/common"
)

const (
	SENSER_NAME string = "BME280"
)

var (
	BME280 = uint8(0x76)
	// BME280 = uint8(0x77)
	I2C_BUS = 1
	i2c     *common.I2CBUS
)

const (
	BME280_ID         byte = 0xD0
	BME280_CTRL_HUM   byte = 0xF2
	BME280_STATUS     byte = 0xF3
	BME280_CTRL_MEAS  byte = 0xF4
	BME280_CONFIG     byte = 0xF5
	BME280_PRESS_MSB  byte = 0xF7
	BME280_PRESS_LSB  byte = 0xF8
	BME280_PRESS_XLSB byte = 0xF9
	BME280_TEMP_MSB   byte = 0xFA
	BME280_TEMP_LSB   byte = 0xFB
	BME280_TEMP_XLSB  byte = 0xFC
	BME280_HUM_MSB    byte = 0xFD
	BME280_HUM_LSB    byte = 0xFE

	BME280_CALIB00 byte = 0x88
	BME280_CALIB25 byte = 0xA1
	BME280_CALIB26 byte = 0xE1
	BME280_CALIB41 byte = 0xF0
)

const (
	BME280_ID_DATA byte = 0x60
)

type osrs_t int //Temperature oversampling
type osrs_p int //Pressure oversampling
type mode int   //register settings mode

const (
	CtrMeasReg_t_SKIP            osrs_t = 0
	CtrMeasReg_t_Oversampling_1  osrs_t = 1
	CtrMeasReg_t_Oversampling_2  osrs_t = 2
	CtrMeasReg_t_Oversampling_4  osrs_t = 3
	CtrMeasReg_t_Oversampling_8  osrs_t = 4
	CtrMeasReg_t_Oversampling_16 osrs_t = 5
	CtrMeasReg_p_SKIP            osrs_p = 0
	CtrMeasReg_p_Oversampling_1  osrs_p = 1
	CtrMeasReg_p_Oversampling_2  osrs_p = 2
	CtrMeasReg_p_Oversampling_4  osrs_p = 3
	CtrMeasReg_p_Oversampling_8  osrs_p = 4
	CtrMeasReg_p_Oversampling_16 osrs_p = 5
	CtrMeasReg_Sleep             mode   = 0
	CtrMeasReg_Forced            mode   = 1
	CtrMeasReg_Normal            mode   = 3
)

type t_sb int     //Controls inactive duration stanbay[ms]
type filter int   //] Filter coefficient
type spi3w_en int //Enables 3-wire SPI interface when set to ‘1’.

const (
	CONFIG_t_sb_0_5   t_sb     = 0 //Forced to Sleep 0.5ms
	CONFIG_t_sb_62_5  t_sb     = 1 //Forced to Sleep 62.5ms
	CONFIG_t_sb_125   t_sb     = 2 //Forced to Sleep 125ms
	CONFIG_t_sb_250   t_sb     = 3 //Forced to Sleep 250ms
	CONFIG_t_sb_500   t_sb     = 4 //Forced to Sleep 500ms
	CONFIG_t_sb_1000  t_sb     = 5 //Forced to Sleep 1000ms
	CONFIG_t_sb_10    t_sb     = 6 //Forced to Sleep 10ms
	CONFIG_t_sb_20    t_sb     = 7 //Forced to Sleep 20ms
	CONFIG_FILTER_OFF filter   = 0
	CONFIG_FILTER_2   filter   = 1
	CONFIG_FILTER_4   filter   = 2
	CONFIG_FILTER_8   filter   = 3
	CONFIG_FILTER_16  filter   = 4
	CONFIG_SPI_3w_OFF spi3w_en = 0
	CONFIG_SPI_3w_ON  spi3w_en = 1
)

type osrs_h int //Humidity oversampling

const (
	CTRL_HUM_SKIP            osrs_h = 0
	CTRL_HUM_Oversampling_1  osrs_h = 1
	CTRL_HUM_Oversampling_2  osrs_h = 2
	CTRL_HUM_Oversampling_4  osrs_h = 3
	CTRL_HUM_Oversampling_8  osrs_h = 4
	CTRL_HUM_Oversampling_16 osrs_h = 5
)

type Bme280_Vaule struct {
	Press float64
	Temp  float64
	Hum   float64
}

type bme280_cal struct {
	press []int
	temp  []int
	hum   []int
}

func ctrlMeasReg_set(v ...interface{}) byte {
	var regdata []int = make([]int, 3)
	for _, value := range v {
		switch value.(type) {
		case osrs_t:
			tmp, _ := value.(osrs_t)
			regdata[0] = int(tmp)
		case osrs_p:
			tmp, _ := value.(osrs_p)
			regdata[1] = int(tmp)
		case mode:
			tmp, _ := value.(mode)
			regdata[2] = int(tmp)
		}
	}
	return byte((regdata[0] << 5) | (regdata[1] << 2) | regdata[2])
}

func config_reg_set(v ...interface{}) byte {
	var regdata []int = make([]int, 3)
	for _, value := range v {
		switch value.(type) {
		case t_sb:
			tmp, _ := value.(t_sb)
			regdata[0] = int(tmp)
		case filter:
			tmp, _ := value.(filter)
			regdata[1] = int(tmp)
		case spi3w_en:
			tmp, _ := value.(spi3w_en)
			regdata[2] = int(tmp)
		}
	}
	return byte((regdata[0] << 5) | (regdata[1] << 2) | regdata[2])
}

func ctrl_hum_reg_set(v osrs_h) byte {
	return byte(v)
}

func up(ctx context.Context) {
	slog.DebugContext(ctx, "BME280 up")
	ctrl_meas_reg := ctrlMeasReg_set(
		CtrMeasReg_p_Oversampling_1, CtrMeasReg_t_Oversampling_1, CtrMeasReg_Normal,
	)
	config_reg := config_reg_set(
		CONFIG_FILTER_OFF, CONFIG_t_sb_1000, CONFIG_SPI_3w_OFF,
	)
	ctrl_hum_reg := ctrl_hum_reg_set(
		CTRL_HUM_Oversampling_1,
	)

	commands := []byte{
		BME280_CTRL_HUM,
		BME280_CTRL_MEAS,
		BME280_CONFIG,
	}
	datas := []byte{
		ctrl_hum_reg,
		ctrl_meas_reg,
		config_reg,
	}
	for i, command := range commands {
		ctx = common.WriteI2cContext(ctx, common.I2c{
			Command: command,
			Data:    datas[i],
		})
		if err := i2c.WriteByte(ctx); err != nil {
			slog.WarnContext(ctx, "BME280 writeI2c", "err", err)

		}
	}
}

func down(ctx context.Context) {
	slog.DebugContext(ctx, "BME280 down")

	ctrl_meas_reg := ctrlMeasReg_set(
		CtrMeasReg_p_SKIP, CtrMeasReg_t_SKIP, CtrMeasReg_Sleep,
	)
	config_reg := config_reg_set(
		CONFIG_FILTER_OFF, CONFIG_t_sb_0_5, CONFIG_SPI_3w_OFF,
	)
	ctrl_hum_reg := ctrl_hum_reg_set(
		CTRL_HUM_SKIP,
	)
	commands := []byte{
		BME280_CTRL_HUM,
		BME280_CTRL_MEAS,
		BME280_CONFIG,
	}
	datas := []byte{
		ctrl_hum_reg,
		ctrl_meas_reg,
		config_reg,
	}
	for i, command := range commands {
		ctx = common.WriteI2cContext(ctx, common.I2c{
			Command: command,
			Data:    datas[i],
		})
		if err := i2c.WriteByte(ctx); err != nil {
			slog.WarnContext(ctx, "BME280 writeI2c", "err", err)

		}
	}
}

func readICIDCheck(ctx context.Context) bool {
	slog.DebugContext(ctx, "BME280 readICIDCheck")

	ctx = common.WriteI2cContext(ctx, common.I2c{
		Command:  BME280_ID,
		ReadSize: 1,
	})
	if buf, err := i2c.ReadByte(ctx); err != nil {
		slog.ErrorContext(ctx, "BM280 Test Read Error", "err", err)
		return false
	} else if buf[0] != BME280_ID_DATA {
		slog.WarnContext(ctx, "BM280 Test header error data", "ID", buf[0])
		return false
	} else {
		slog.DebugContext(ctx, "BME280 Test Read Bme280 ID", "ID", buf[0])
	}
	return true
}

func readStatus(ctx context.Context) mode {
	slog.DebugContext(ctx, "BME280 readStatus")

	ctx = common.WriteI2cContext(ctx, common.I2c{
		Command:  BME280_CTRL_MEAS,
		ReadSize: 1,
	})
	if buf, err := i2c.ReadByte(ctx); err != nil {
		slog.ErrorContext(ctx, "BM280 readStatus Read Error", "err", err)
		return CtrMeasReg_Sleep
	} else {
		switch buf[0] & 0x03 {
		case 0x0:
			return CtrMeasReg_Sleep
		case 0x1:
			return CtrMeasReg_Forced
		case 0x2:
			return CtrMeasReg_Forced
		case 0x3:
			return CtrMeasReg_Normal
		}
	}
	return CtrMeasReg_Sleep

}

func rawRead(ctx context.Context) (int, int, int) {
	slog.DebugContext(ctx, "BME280 rawRead")

	ctx = common.WriteI2cContext(ctx, common.I2c{
		Command:  BME280_PRESS_MSB,
		ReadSize: 8,
	})
	buf, err := i2c.ReadByte(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "error BME280 read data by i2c", "err", err)
		return -1, -1, -1
	}
	tmp := 0
	for _, bt := range buf {
		tmp += int(bt)
	}
	if tmp == 0 {
		slog.ErrorContext(ctx, "error BME280 read data byte")
		return -1, -1, -1
	}
	press := int(buf[0])<<12 | int(buf[1])<<4 | int(buf[2])>>4
	temp := int(buf[3])<<12 | int(buf[4])<<4 | int(buf[5])>>4
	hum := int(buf[6])<<8 | int(buf[7])
	slog.DebugContext(ctx, "BME280 read raw data", "press", press, "temp", temp, "hum", hum)
	return press, temp, hum
}

func calibRead(ctx context.Context) bme280_cal {
	slog.DebugContext(ctx, "BME280 calibRead")

	var out bme280_cal

	var calib []int
	commands := []byte{
		BME280_CALIB00,
		BME280_CALIB26,
	}
	sizes := []int{
		int(BME280_CALIB25 - BME280_CALIB00),
		int(BME280_CALIB41 - BME280_CALIB26 + 1),
	}
	for i, command := range commands {
		ctx = common.WriteI2cContext(ctx, common.I2c{
			Command:  command,
			ReadSize: sizes[i],
		})
		buf, err := i2c.ReadByte(ctx)
		if err != nil {
			slog.WarnContext(ctx, "BME280 writeI2c", "err", err)
		}
		for _, b := range buf {
			calib = append(calib, int(b))
		}
	}
	for i := 0; i < 3; i++ { //0-5
		num := (calib[1+i*2] << 8) | calib[0+i*2]
		if i != 0 {
			if (num & 0x8000) != 0 {
				num = (-num ^ 0xffff) + 1
			}
		}
		out.temp = append(out.temp, num)
	}
	for i := 0; i < 9; i++ { //6-23
		num := (calib[7+i*2] << 8) | calib[6+i*2]
		if i != 0 {
			if (num & 0x8000) != 0 {
				num = (-num ^ 0xffff) + 1
			}
		}
		out.press = append(out.press, num)
	}
	//24-31
	out.hum = append(out.hum, calib[24])
	out.hum = append(out.hum, (calib[26]<<8)|calib[25])
	out.hum = append(out.hum, calib[27])
	out.hum = append(out.hum, (calib[28]<<4)|(0x0F&calib[29]))
	out.hum = append(out.hum, (calib[30]<<4)|((calib[29]>>4)&0x0F))
	out.hum = append(out.hum, calib[31])
	for i := 0; i < len(out.hum); i++ {
		if (i != 0) && (i != 2) {
			if (out.hum[i] & 0x8000) != 0 {
				out.hum[i] = (-out.hum[i] ^ 0xffff) + 1
			}
		}
	}
	slog.DebugContext(ctx, "BME280 read calib data", "press", out.press, "temp", out.temp, "hum", out.hum)

	return out
}

func (t *bme280_cal) CalibTemp(b_temp int) (float64, float64) {
	tmp := float64(b_temp)
	var calib []float64
	for _, flt := range t.temp {
		calib = append(calib, float64(flt))
	}
	v1 := (tmp/16384.0 - calib[0]/1024.0) * calib[1]
	v2 := (tmp/131072.0 - calib[0]/8192.0) * (tmp/131072.0 - calib[0]/8192.0) * calib[2]
	t_fine := v1 + v2
	temperature := t_fine / 5120.0
	return temperature, t_fine

}

func (t *bme280_cal) CalibPress(b_press int, t_fine float64) float64 {
	tmp := float64(b_press)
	var calib []float64
	for _, flt := range t.press {
		calib = append(calib, float64(flt))
	}

	v1 := (t_fine / 2.0) - 64000.0
	v2 := (((v1 / 4.0) * (v1 / 4.0)) / 2048) * calib[5]
	v2 = v2 + ((v1 * calib[4]) * 2.0)
	v2 = (v2 / 4.0) + (calib[3] * 65536.0)
	v1 = (((calib[2] * (((v1 / 4.0) * (v1 / 4.0)) / 8192)) / 8) + ((calib[1] * v1) / 2.0)) / 262144
	v1 = ((32768 + v1) * calib[0]) / 32768
	if v1 != 0 {
		pressure := ((1048576 - tmp) - (v2 / 4096)) * 3125
		if pressure < 0 {
			pressure = (pressure * 2.0) / v1
		} else {
			pressure = (pressure / v1) * 2
		}
		v1 = (calib[8] * (((pressure / 8.0) * (pressure / 8.0)) / 8192.0)) / 4096
		v2 = ((pressure / 4.0) * calib[7]) / 8192.0
		return pressure / 100
	}
	return 0
}

func (t *bme280_cal) CalibHum(b_hum int, t_fine float64) float64 {
	tmp := float64(b_hum)
	var calib []float64
	for _, flt := range t.hum {
		calib = append(calib, float64(flt))
	}
	var_h := t_fine - 76800.0
	if var_h != 0 {
		var_h = (tmp - (calib[3]*64.0 + calib[4]/16384.0*var_h)) * (calib[1] / 65536.0 * (1.0 + calib[5]/67108864.0*var_h*(1.0+calib[2]/67108864.0*var_h)))
	} else {
		return 0
	}
	var_h = var_h * (1.0 - calib[0]*var_h/524288.0)
	if var_h > 100 {
		var_h = 100
	} else if var_h < 0 {
		var_h = 0
	}
	return var_h
}

// readSenserData(calib bme280_cal) (float64, float64, float64)
//
// c_press, c_tmp, c_hum
func readSenserData(ctx context.Context, calib bme280_cal) (float64, float64, float64) {
	slog.DebugContext(ctx, "BME280 readSenserData")

	press, temp, hum := rawRead(ctx)
	if press == -1 && temp == -1 && hum == -1 {
		return -1, -1, -1
	}
	c_tmp, t_fine := calib.CalibTemp(temp)
	c_press := calib.CalibPress(press, t_fine)
	c_hum := calib.CalibHum(hum, t_fine)

	if c_hum <= 0 || c_tmp < -40 || c_tmp > 85 {
		return -1, -1, -1
	}
	slog.DebugContext(ctx, "BME280 readSenserData", "press", c_press, "temp", c_tmp, "hum", c_hum)
	return c_press, c_tmp, c_hum

}
