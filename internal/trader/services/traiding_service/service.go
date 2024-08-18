package traidingservice

import (
	"context"
	"errors"
	"fmt"

	"github.com/MaikovskiyS/TraderBot/internal/domain"
	"github.com/rs/zerolog"
)

type ExchangeProvider interface {
	GetCandels(ctx context.Context, req *GetCandelsRequest) ([]*domain.Candle, error)
	CreateOrder(ctx context.Context, req *CreateOrderParams) error
	GetTickers(ctx context.Context) (domain.Tickers, error)
	GetOrderBook(ctx context.Context, symbol string) (*domain.OrderBook, error)
	GetBalance(ctx context.Context) (float64, error)
	GetPositionInfo(ctx context.Context) (*GetPositionInfoResponse, error)
	SetLeverage(ctx context.Context, ticker string) error
	GetOpenClosedOrdersByTicker(ctx context.Context, ticker string) (*GetOpenOrdersResponse, error)
	AmendOrder(ctx context.Context, params *AmendOrderParams) error
	SetStopLoss(ctx context.Context, params *SetSlParams) error
	SetTakeProfit(ctx context.Context, params *SetTpParams) error
	CancelOrder(ctx context.Context, params *CancelOrderParams) error
	GetTickerInfo(ctx context.Context, symbol string) (*domain.Ticker, error)
}

var (
	ErrNoPositions          = errors.New("doesn't have positions")
	ErrPositionAlreadyExist = errors.New("position already exist")
	ErrOrderAlreadyExist    = errors.New("order already exist")

	ErrNoOrders = errors.New("doesn't have orders")
)

type setting struct {
	riskPercentage float64
	candelsLimit   int
}
type tradingService struct {
	setting
	Log      zerolog.Logger
	Provider ExchangeProvider
}

func New(provider ExchangeProvider, log zerolog.Logger) *tradingService {
	return &tradingService{
		setting: setting{
			riskPercentage: 2.00,
			candelsLimit:   60,
		},
		Provider: provider,
		Log:      log,
	}
}

func convertPriceToString(price float64, precision int) string {
	return fmt.Sprintf("%.*f", precision, price)
}
