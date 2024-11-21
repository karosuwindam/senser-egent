package common

import "github.com/pkg/errors"

var (
	CONSTEX_DATA_NOTFOND = errors.New("error context data not fond I2C Data")
	I2C_WRITE_TIMEOUT    = errors.New("error I2c Write timeout")
	I2C_READ_TIMEOUT     = errors.New("error I2c Read timeout")
	CONTEXT_DONE         = errors.New("error context Done")
)
