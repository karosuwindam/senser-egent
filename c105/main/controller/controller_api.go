package controller

import (
	"context"
	"senseregent/controller/sennser"
)

type API struct {
}

func NewAPI() *API {
	return &API{}
}

func (a *API) ReadValue(ctx context.Context) (sennser.SennserValue, error) {
	return sennser.GetValue(ctx)
}

func (a *API) ResetSennser(ctx context.Context) error {
	return sennser.Reset(ctx)
}
