package traidingservice

import (
	"context"
	"fmt"
	"math"

	"github.com/MaikovskiyS/TraderBot/internal/domain"
)

const riskPercentage = 1.00

func (s *tradingService) InitOrder(ctx context.Context, req *CreateOrderRequest) error {
	balance, err := s.Provider.GetBalance(ctx)
	if err != nil {
		return fmt.Errorf("get balance: %w", err)
	}
	balance = balance - (balance * 0.1)

	err = s.Provider.SetLeverage(ctx, req.Ticker.Symbol)
	if err != nil {
		return err
	}

	params := createOrderParams(req, balance)

	s.Log.Debug().Any("order params", params).Send()

	return s.Provider.CreateOrder(ctx, params)
}

func (s *tradingService) IsTradeAvaliable(ctx context.Context, ticker string) (bool, error) {
	info, err := s.Provider.GetPositionInfo(ctx)
	if err != nil {
		return false, fmt.Errorf("get positions: %w", err)
	}

	if len(info.Positions) >= 1 {
		return false, ErrPositionAlreadyExist
	}

	for _, info := range info.Positions {
		if string(info.Symbol) == ticker {
			return false, ErrPositionAlreadyExist
		}
	}

	respOrders, err := s.Provider.GetOpenClosedOrdersByTicker(ctx, ticker)
	if err != nil {
		return false, fmt.Errorf("get orders: %w", err)
	}

	if len(respOrders.Orders) >= 1 {
		return false, ErrOrderAlreadyExist
	}

	for _, order := range respOrders.Orders {
		if order.Symbol == ticker {
			return false, ErrOrderAlreadyExist
		}
	}

	return true, nil
}
func createOrderParams(req *CreateOrderRequest, balance float64) *CreateOrderParams {
	var takeProfit string

	stopLoss := calculateOrderStopLoss(req.TargetCandle, req.Side)
	openPrice := convertPriceToString(req.OpenPrice, req.Ticker.Precision)
	sl := convertPriceToString(stopLoss, req.Ticker.Precision)
	qty := calculateOrderQuantity(
		balance,
		req.OpenPrice,
		stopLoss)

	if req.ClosePrice != nil {
		tP := calculateTakeProfit(req.OpenPrice, stopLoss, req.Side)
		takeProfit = convertPriceToString(tP, req.Ticker.Precision)
	}
	return &CreateOrderParams{
		Coin:       req.Ticker.Symbol,
		OpenPrice:  openPrice,
		Quantity:   qty,
		StopLoss:   sl,
		TakeProfit: takeProfit,
		Side:       req.Side,
	}
}

func calculateOrderQuantity(balance, entryPrice, stopLossPrice float64) string {
	// Вычисляем максимальный допустимый убыток
	maxLoss := balance * (riskPercentage / 100)

	// Вычисляем возможные потери на единицу актива
	lossPerUnit := math.Abs(entryPrice - stopLossPrice)

	// Рассчитываем количество активов, чтобы не превышать допустимый убыток
	quantity := maxLoss / lossPerUnit

	// Проверяем, чтобы количество активов не превышало доступное количество по балансу
	maxQuantityByBalance := balance / entryPrice
	if quantity > maxQuantityByBalance {
		quantity = maxQuantityByBalance
	}

	// Используем customRound для корректного округления значения
	roundedQty := customRound(quantity)

	// Форматируем значение как строку
	return fmt.Sprintf("%v", roundedQty)
}

func calculateTakeProfit(entryPrice, stopLossPrice float64, side domain.Side) float64 {
	// Рассчитываем риск
	risk := 0.0
	if side == domain.BuySide {
		risk = entryPrice - stopLossPrice
	} else if side == domain.SellSide {
		risk = stopLossPrice - entryPrice
	}

	// Рассчитываем цель (take profit) с соотношением риск/прибыль 1:3
	takeProfit := 0.0
	if side == domain.BuySide {
		takeProfit = entryPrice + 3*risk
	} else if side == domain.SellSide {
		takeProfit = entryPrice - 3*risk
	}

	return takeProfit
}

func calculateOrderStopLoss(candle *domain.Candle, side domain.Side) float64 {
	percentage := candle.GetMovePercent()

	movePercentage := (percentage / 100) / 2 // 50% от движения цены в свече
	var stopLossPrice float64

	if side == domain.BuySide {
		// Для покупок отнимаем процент от минимума свечи
		stopLossPrice = candle.Low - (candle.Low * movePercentage)
	} else {
		// Для продаж добавляем процент к максимуму свечи
		stopLossPrice = candle.High + (candle.High * movePercentage)
	}

	return stopLossPrice
}

func customRound(value float64) float64 {
	// Округляем до целого, если значение больше 1
	if value >= 1 {
		return math.Floor(value)
	}

	// Если значение меньше 1, то сокращаем до определенного количества знаков после запятой
	if value < 1 {
		factor := 1.0

		// Определяем количество знаков после запятой
		for value*factor < 1 {
			factor *= 10
		}

		// Округляем до двух значимых цифр после первой ненулевой
		value = math.Floor(value*factor*10+0.5) / (factor * 10)
	}

	return value
}
