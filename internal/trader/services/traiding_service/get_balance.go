package traidingservice

import "context"

func (a *tradingService) GetBalance(ctx context.Context) (float64, error) {
	return a.Provider.GetBalance(ctx)
}
