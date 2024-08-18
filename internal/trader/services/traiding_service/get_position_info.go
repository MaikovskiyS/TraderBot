package traidingservice

import (
	"context"
	"fmt"
)

func (a *tradingService) GetPositions(ctx context.Context) (*GetPositionInfoResponse, error) {
	resp, err := a.Provider.GetPositionInfo(ctx)
	if err != nil {
		return nil, fmt.Errorf("get positions: %w", err)
	}

	if resp.Positions == nil || len(resp.Positions) == 0 {
		return nil, ErrNoPositions
	}

	return resp, nil

}

func (a *tradingService) GetOpenOrders(ctx context.Context) (*GetOpenOrdersResponse, error) {
	return a.Provider.GetOpenClosedOrdersByTicker(ctx, "")
}
