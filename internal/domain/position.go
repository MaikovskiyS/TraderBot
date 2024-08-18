package domain

import "math"

type Position struct {
	Symbol          string
	AvgPrice        float64
	TakeProfit      float64
	PositionValue   string
	TpSlMode        string
	TrailingStop    string
	UnrealisedPnl   string
	MarkPrice       float64
	PositionBalance string
	Side            Side
	Size            float64
	PositionStatus  string
	StopLoss        float64
	CreatedTime     string
	Precision       int
}

func (p Position) MovePricePercent() float64 {
	if p.MarkPrice == 0 {
		return 0
	}

	move := ((p.MarkPrice - p.AvgPrice) / p.AvgPrice) * 100

	return math.Abs(move)
}
