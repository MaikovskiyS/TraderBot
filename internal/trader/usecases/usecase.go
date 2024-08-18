package usecases

import (
	"context"

	"github.com/MaikovskiyS/TraderBot/internal/domain"
	strategy "github.com/MaikovskiyS/TraderBot/internal/trader/services/strategy_service"
	trading "github.com/MaikovskiyS/TraderBot/internal/trader/services/traiding_service"
	"github.com/rs/zerolog"
)

type (
	TraidingService interface {
		Get20VolatilityTickers(ctx context.Context, interval string) (*trading.GetVolatilityTicketsResponse, error)
		GetVolatilityTickersWithOrderBooks(ctx context.Context, interval string) (*trading.GetTicketsWithOrderBookResponse, error)
		InitOrder(ctx context.Context, req *trading.CreateOrderRequest) error
		GetPositions(ctx context.Context) (*trading.GetPositionInfoResponse, error)
		ManagePosition(ctx context.Context, position *domain.Position) error
		IsTradeAvaliable(ctx context.Context, ticker string) (bool, error)
		ManageOrders(ctx context.Context) error
	}

	StrategyService interface {
		ApplySupportResistance(data []*domain.TickerCandels) (*strategy.StrategyResponse, error)
		ApplyBidsAsksStrategy(ctx context.Context, tickers []*domain.TickerOrderBook) (
			*strategy.StrategyResponse, error)
	}
)

type useCases struct {
	Traiding TraidingService
	Strategy StrategyService
	Log      zerolog.Logger
}

func New(traiding TraidingService, strategy StrategyService, log zerolog.Logger) *useCases {
	return &useCases{
		Traiding: traiding,
		Strategy: strategy,
		Log:      log,
	}
}
