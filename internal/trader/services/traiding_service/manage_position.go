package traidingservice

import (
	"context"
	"fmt"

	"github.com/MaikovskiyS/TraderBot/internal/domain"
)

func (s *tradingService) ManagePosition(ctx context.Context, p *domain.Position) error {
	s.Log.Debug().Any("current position", p).Send()
	// Str("ticker", p.Symbol).
	// Float64("mark price", p.MarkPrice).
	// Float64("avg price", p.AvgPrice).
	// Str("PNL", p.UnrealisedPnl).
	// Msg("current position")

	// Если позиция прошла 30% до takeProfit, фиксируем половину и обновляем stopLoss до точки входа
	ordersResp, err := s.GetOpenOrders(ctx)
	if err != nil {
		return fmt.Errorf("get open orders: %w", err)
	}

	var takeProfitOrder *domain.Order
	var stopOrdersLen int64
	for _, order := range ordersResp.Orders {
		// s.Log.Debug().Any("order", order).Send()
		if order.StopOrderType == domain.PartialTakeProfitOrderType {
			takeProfitOrder = order
		} else {
			stopOrdersLen += 1
		}
	}

	// если мы уже поставили новый stop, выходим из функции
	if stopOrdersLen == 2 {
		s.Log.Debug().Msg("stop loss already changed")
		return nil
	}

	movePercent := 0.30
	var targetPrice float64
	var isProfitable bool
	// Рассчитываем пороговое значение цены для частичного закрытия позиции
	if p.Side == domain.BuySide {
		targetPrice = p.AvgPrice + movePercent*(takeProfitOrder.TriggerPrice-p.AvgPrice)
		isProfitable = p.MarkPrice >= targetPrice
	} else if p.Side == domain.SellSide {
		targetPrice = p.AvgPrice - movePercent*(p.AvgPrice-takeProfitOrder.TriggerPrice)
		isProfitable = p.MarkPrice <= targetPrice
	}

	if isProfitable {
		s.Log.Debug().Msg("is profitable")
		if err := s.ClosePositionPartially(ctx, p); err != nil {
			return fmt.Errorf("failed to close partial position: %w", err)
		}
	}

	return nil
}

func (s *tradingService) ClosePositionPartially(ctx context.Context, position *domain.Position) error {

	err := s.Provider.SetStopLoss(ctx, &SetSlParams{
		Symbol:       position.Symbol,
		Side:         position.Side,
		StopLoss:     convertPriceToString(position.AvgPrice, position.Precision),
		Size:         fmt.Sprintf("%v", customRound(position.Size)),
		TpLimitPrice: convertPriceToString(position.TakeProfit, position.Precision),
	})
	if err != nil {
		return fmt.Errorf("set stop loss: %w", err)
	}
	s.Log.Debug().Msg("stop loss updated")

	// err = s.Provider.SetTakeProfit(ctx, &SetTpParams{
	// 	Symbol:     position.Symbol,
	// 	Side:       position.Side,
	// 	TakeProfit: convertPriceToString(position.MarkPrice, position.Precision),
	// 	Size:       fmt.Sprintf("%v", customRound(position.Size/2)),
	// })
	// if err != nil {
	// 	return fmt.Errorf("set take profit: %w", err)
	// }
	// s.Log.Debug().Float64("mark price", position.MarkPrice).Msg("take profit updated")

	// закрыть пол позиции
	// переставить stop-loss

	return nil
}

/*
	var stopOrderID string
	var stopOrderLinkID string
	var takeOrderID string
	var takeOrderLinkID string
	for _, order := range ordersResp.Orders {
		if order.Symbol == position.Symbol && order.StopOrderType == domain.StopLossOrderType {
			stopOrderID = order.OrderID
			stopOrderLinkID = order.OrderLinkID
		} else if order.Symbol == position.Symbol && order.StopOrderType == domain.TakeProfitOrderType {
			takeOrderID = order.OrderID
			takeOrderLinkID = order.OrderLinkID
		}
	}


	triggerStopPrice := convertPriceToString(position.AvgPrice, position.Precision)
	stopQty := fmt.Sprintf("%v", customRound(position.Size/2))

	err = s.Provider.AmendOrder(ctx, &AmendOrderParams{
		Symbol:       position.Symbol,
		OrderID:      &stopOrderID,
		OrderLinkID:  &stopOrderLinkID,
		Qty:          &stopQty,
		TriggerPrice: &triggerStopPrice,
	})
	if err != nil {
		return err
	}

	//TODO: дотестить

	//проблема в том, что TP должен быть больше mark price
	triggerTakePrice := convertPriceToString(position.MarkPrice, position.Precision)
	takeQty := fmt.Sprintf("%v", customRound(position.Size/2))

	s.Log.Debug().
		Str("takeQTY", takeQty).
		Msg("changing order's TP")

	err = s.Provider.AmendOrder(ctx, &AmendOrderParams{
		Symbol:       position.Symbol,
		OrderID:      &takeOrderID,
		OrderLinkID:  &takeOrderLinkID,
		Qty:          &takeQty,
		TriggerPrice: &triggerTakePrice,
	})
	if err != nil {
		return err
	}

	//set new takeProfit

*/
