package traidingservice

import (
	"context"
	"fmt"

	"github.com/MaikovskiyS/TraderBot/internal/domain"
)

const candelsLimit = 60

func (t *tradingService) Get20VolatilityTickers(ctx context.Context, interval string) (*GetVolatilityTicketsResponse, error) {

	tickers, err := t.Provider.GetTickers(ctx)
	if err != nil {
		return nil, fmt.Errorf("get tickers: %w", err)
	}

	topTickers := tickers.Get10MostVolume24h()
	tickerCandels := make([]*domain.TickerCandels, 0)
	for _, ticker := range topTickers {
		if ticker.Symbol == "BTCUSDT" ||
			ticker.Symbol == "ETHUSDT" ||
			ticker.Symbol == "BNBUSDT" ||
			ticker.Symbol == "BTCPERP" ||
			ticker.Symbol == "LTCUSDT" {
			continue
		}

		candels, err := t.Provider.GetCandels(ctx, &GetCandelsRequest{
			Symbol:   ticker.Symbol,
			Interval: interval,
			Limit:    candelsLimit,
		})
		if err != nil {
			return nil, fmt.Errorf("get candels by ticker %s: %w", ticker.Symbol, err)
		}

		tickerCandels = append(tickerCandels, &domain.TickerCandels{
			Symbol: domain.TickerPrecision{
				Symbol:    ticker.Symbol,
				Precision: ticker.Precision,
			},
			Candels: candels,
		})
	}

	return &GetVolatilityTicketsResponse{
		TickersWithCandels: tickerCandels,
	}, nil

}

func (t *tradingService) GetVolatilityTickersWithOrderBooks(ctx context.Context, interval string) (*GetTicketsWithOrderBookResponse, error) {

	tickers, err := t.Provider.GetTickers(ctx)
	if err != nil {
		return nil, fmt.Errorf("get tickers: %w", err)
	}

	topTickers := tickers.Get10MostVolume24h()
	tickerOrderBook := make([]*domain.TickerOrderBook, 0)
	for _, ticker := range topTickers {

		orderBook, err := t.Provider.GetOrderBook(ctx, ticker.Symbol)
		if err != nil {
			return nil, fmt.Errorf("get candels by ticker %s: %w", ticker.Symbol, err)
		}

		tickerOrderBook = append(tickerOrderBook, &domain.TickerOrderBook{
			Symbol: domain.TickerPrecision{
				Symbol:    ticker.Symbol,
				Precision: ticker.Precision,
			},
			OrderBook: orderBook,
		})
	}

	return &GetTicketsWithOrderBookResponse{
		TickersWithOrderBook: tickerOrderBook,
	}, nil

}
