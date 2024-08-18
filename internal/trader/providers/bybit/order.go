package bybit_provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/MaikovskiyS/TraderBot/internal/domain"
	trading "github.com/MaikovskiyS/TraderBot/internal/trader/services/traiding_service"
	"github.com/hirokisan/bybit/v2"
)

func (p *BybitProvider) CreateOrder(ctx context.Context, req *trading.CreateOrderParams) error {
	_, err := p.BybitClient.V5().Order().CreateOrder(convertOrderRequest(req))
	if err != nil {
		return fmt.Errorf("place order: %w", err)
	}
	return nil
}

func (p *BybitProvider) SetLeverage(ctx context.Context, ticker string) error {
	_, err := p.BybitClient.V5().Position().SetLeverage(bybit.V5SetLeverageParam{
		Category:     bybit.CategoryV5Linear,
		Symbol:       bybit.SymbolV5(ticker),
		BuyLeverage:  "1",
		SellLeverage: "1",
	})
	if err != nil {
		var bybitErr *bybit.ErrorResponse
		if errors.As(err, &bybitErr) {
			if bybitErr.RetCode == 110043 {
				return nil
			}
		}

		return fmt.Errorf("set leverage: %w", err)
	}

	return nil
}

func (p *BybitProvider) GetOpenClosedOrdersByTicker(ctx context.Context, ticker string) (*trading.GetOpenOrdersResponse, error) {
	settle := bybit.CoinUSDT
	openOnly := 0
	//symbol := bybit.SymbolV5(ticker)
	resp, err := p.BybitClient.V5().Order().GetOpenOrders(bybit.V5GetOpenOrdersParam{
		Category:   bybit.CategoryV5Linear,
		Symbol:     nil,
		BaseCoin:   nil,
		SettleCoin: &settle,
		OpenOnly:   &openOnly,
	})
	if err != nil {
		return nil, fmt.Errorf("get open orders: %w", err)
	}

	if resp.Result.List == nil {
		return nil, ErrEmptyResponse
	}

	orders := make([]*domain.Order, len(resp.Result.List))
	for i, order := range resp.Result.List {
		orders[i] = convertOrderToDomain(order)
	}

	return &trading.GetOpenOrdersResponse{
		Orders: orders,
	}, nil
}

func (p *BybitProvider) AmendOrder(ctx context.Context, params *trading.AmendOrderParams) error {
	_, err := p.BybitClient.V5().Order().AmendOrder(convertAmendOrderParams(params))
	if err != nil {
		return fmt.Errorf("amend order: %w", err)
	}

	return nil
}

func (p *BybitProvider) SetStopLoss(ctx context.Context, params *trading.SetSlParams) error {
	_, err := p.BybitClient.V5().Position().SetTradingStop(convertSlParams(params))
	if err != nil {
		return err
	}

	return nil
}

func (p *BybitProvider) SetTakeProfit(ctx context.Context, params *trading.SetTpParams) error {
	_, err := p.BybitClient.V5().Position().SetTradingStop(convertTpParams(params))
	if err != nil {
		return err
	}

	return nil
}

func (p *BybitProvider) CancelOrder(ctx context.Context, params *trading.CancelOrderParams) error {
	_, err := p.BybitClient.V5().Order().CancelOrder(bybit.V5CancelOrderParam{
		Category:    bybit.CategoryV5Linear,
		Symbol:      bybit.SymbolV5(params.Symbol),
		OrderID:     params.OrderID,
		OrderLinkID: params.OrderLinkID,
	})
	if err != nil {
		return fmt.Errorf("cancel order: %w", err)
	}

	return nil
}
