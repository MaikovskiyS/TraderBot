package usecases

import (
	"context"
	"fmt"

	traidingservice "github.com/MaikovskiyS/TraderBot/internal/trader/services/traiding_service"
)

func (u *useCases) RunSupportResistance(ctx context.Context) error {
	interval := "5"

	resp, err := u.Traiding.Get20VolatilityTickers(ctx, interval)
	if err != nil {
		return fmt.Errorf("get volatility tickers: %w", err)
	}

	// проверить ордера на время создания и на сколько текущая цена отличается от от цены открытия ордера
	// получить ордера, если по этому символу нет позиции, и цена отличается от ордера больше чем на 1 процент, закрыть stopLoss takeProfit
	err = u.Traiding.ManageOrders(ctx)
	if err != nil {
		return err
	}

	orderData, err := u.Strategy.ApplySupportResistance(resp.TickersWithCandels)
	if err != nil {
		return fmt.Errorf("apply support-resistence strategy: %w", err)
	}

	if orderData == nil {
		return nil
	}

	trade := orderData.Trade

	u.Log.Debug().Any("trade", trade)

	isAvaliable, err := u.Traiding.IsTradeAvaliable(ctx, trade.Symbol.Symbol)
	if err != nil {
		return fmt.Errorf("is trade avaliable: %w", err)
	}

	if !isAvaliable {
		u.Log.Debug().
			Str("ticker", trade.Symbol.Symbol).
			Msg("trade is not avaliable")

		return nil
	}

	u.Log.Debug().
		Any("request", trade.Symbol).
		Str("Symbol", trade.Symbol.Symbol).
		Str("Profit", fmt.Sprintf("%.2f", trade.ProfitPercent)).
		Send()

	err = u.Traiding.InitOrder(ctx, &traidingservice.CreateOrderRequest{
		Ticker:       orderData.Trade.Symbol,
		OpenPrice:    trade.OpenPrice,
		ClosePrice:   trade.ClosePrice,
		Side:         trade.Side,
		TargetCandle: trade.TargetCandle,
	})
	if err != nil {
		u.Log.Error().Err(err).Send()

		return fmt.Errorf("init order: %w", err)
	}

	return nil
}
