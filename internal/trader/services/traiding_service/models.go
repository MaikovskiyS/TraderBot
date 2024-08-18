package traidingservice

import "github.com/MaikovskiyS/TraderBot/internal/domain"

type GetVolatilityTicketsResponse struct {
	TickersWithCandels []*domain.TickerCandels
}

type GetTicketsWithOrderBookResponse struct {
	TickersWithOrderBook []*domain.TickerOrderBook
}

type CreateOrderRequest struct {
	Ticker       domain.TickerPrecision
	OpenPrice    float64
	ClosePrice   *float64
	Side         domain.Side
	TargetCandle *domain.Candle
}

type GetCandelsRequest struct {
	Symbol   string
	Interval string
	Limit    int
}

type CreateOrderParams struct {
	Coin       string
	OpenPrice  string
	Quantity   string
	StopLoss   string
	TakeProfit string
	Side       domain.Side
}

type AmendOrderParams struct {
	Symbol string

	OrderID      *string
	OrderLinkID  *string
	OrderIv      *string
	TriggerPrice *string
	Qty          *string
	Price        *string
	TakeProfit   *string
	StopLoss     *string
}

type SetSlParams struct {
	Symbol       string
	Side         domain.Side
	StopLoss     string
	Size         string
	TpLimitPrice string
}

type SetTpParams struct {
	Symbol     string
	Side       domain.Side
	TakeProfit string
	Size       string
}

type GetPositionInfoResponse struct {
	Positions []*domain.Position
}

type GetOpenOrdersResponse struct {
	Orders []*domain.Order
}

type CancelOrderParams struct {
	Symbol      string  `json:"symbol"`
	OrderID     *string `json:"orderId,omitempty"`
	OrderLinkID *string `json:"orderLinkId,omitempty"`
}
