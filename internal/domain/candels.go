package domain

type Candle struct {
	Time   int64
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume float64
}

func (c Candle) GetMiddlePrice() float64 {
	return (c.High + c.Low) / 2
}

func (c Candle) GetMovePercent() float64 {
	if c.Low == 0 {
		return 0
	}

	move := ((c.High - c.Low) / c.Low) * 100

	return move
}

func (c Candle) GetMoveAbsolute() float64 {
	move := c.High - c.Low
	if move < 0 {
		return -move
	}
	return move
}

type TickerCandels struct {
	Symbol  TickerPrecision
	Candels []*Candle
}

type TickerOrderBook struct {
	Symbol    TickerPrecision
	OrderBook *OrderBook
}
