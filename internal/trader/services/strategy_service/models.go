package strategyservice

import (
	"github.com/MaikovskiyS/TraderBot/internal/domain"
)

type StrategyResponse struct {
	Trade *Trade
}

type Trade struct {
	Symbol        domain.TickerPrecision
	Side          domain.Side
	OpenPrice     float64
	ClosePrice    *float64
	TargetCandle  *domain.Candle
	ProfitPercent float64
}
