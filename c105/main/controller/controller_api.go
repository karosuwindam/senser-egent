package controller

import (
	"context"
	"senseregent/controller/senser"
)

type API struct {
}

func NewAPI() *API {
	return &API{}
}

func (a *API) ReadValue(ctx context.Context) (senser.SenserValue, error) {
	return senser.GetValue(ctx)
}

func (a *API) Resetsenser(ctx context.Context) error {
	return senser.Reset(ctx)
}
