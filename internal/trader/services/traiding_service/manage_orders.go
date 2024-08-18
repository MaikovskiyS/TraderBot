package traidingservice

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"time"
)

func parseUnixTimestamp(timestamp string) (int64, error) {
	return strconv.ParseInt(timestamp, 10, 64)
}

func (s *tradingService) ManageOrders(ctx context.Context) error {
	ordersResp, err := s.GetOpenOrders(ctx)
	if err != nil {
		return fmt.Errorf("get open orders: %w", err)
	}

	if ordersResp.Orders == nil || len(ordersResp.Orders) == 0 {
		return nil
	}

	currentTime := time.Now().UnixMilli()

	for _, order := range ordersResp.Orders {
		// Парсим время создания ордера
		orderCreatedTime, err := parseUnixTimestamp(order.CreatedTime)
		if err != nil {
			fmt.Println("Error parsing time:", err)
			return err
		}

		ticker, err := s.Provider.GetTickerInfo(ctx, order.Symbol)
		if err != nil {
			return fmt.Errorf("get ticker info: %w", err)
		}

		// Считаем возраст ордера в миллисекундах
		orderAgeMs := currentTime - orderCreatedTime
		orderAgeMinutes := orderAgeMs / (1000 * 60)

		// Считаем разницу в цене
		priceDifference := math.Abs((order.Price-ticker.LastPriceFloat64())/ticker.LastPriceFloat64()) * 100

		s.Log.Debug().
			Float64("priceDiff", priceDifference).
			Int64("orderAgeMinutes", orderAgeMinutes).
			Msg("Order details")

		// Проверяем условия отмены ордера
		if orderAgeMs > int64(10*time.Minute.Milliseconds()) || priceDifference > 0.6 {
			if err := s.Provider.CancelOrder(ctx, &CancelOrderParams{
				Symbol:      order.Symbol,
				OrderID:     &order.OrderID,
				OrderLinkID: &order.OrderLinkID,
			}); err != nil {
				return fmt.Errorf("cancel order: %w", err)
			}
		}

		s.Log.Debug().Str("ticker", order.Symbol).Msg("order deleted")
	}

	return nil
}
