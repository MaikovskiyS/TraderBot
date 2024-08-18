package bybit_provider

import (
	"context"
	"fmt"

	"github.com/MaikovskiyS/TraderBot/internal/domain"
	trading "github.com/MaikovskiyS/TraderBot/internal/trader/services/traiding_service"
	"github.com/hirokisan/bybit/v2"
)

func (p *BybitProvider) GetCandels(ctx context.Context, req *trading.GetCandelsRequest) ([]*domain.Candle, error) {
	resp, err := p.BybitClient.V5().Market().GetKline(bybit.V5GetKlineParam{
		Category: bybit.CategoryV5Linear,
		Symbol:   bybit.SymbolV5(req.Symbol),
		Interval: bybit.Interval(req.Interval),
		Start:    nil,
		End:      nil,
		Limit:    &req.Limit,
	})
	if err != nil {
		return nil, fmt.Errorf("klines service do: %w", err)
	}

	if resp.Result.List == nil {
		return nil, ErrEmptyResponse
	}

	candels := make([]*domain.Candle, len(resp.Result.List))
	for i, kline := range resp.Result.List {
		candle, err := convertCandleToDomain(kline)
		if err != nil {
			return nil, err
		}
		candels[i] = candle
	}

	return candels, nil
}

func (p *BybitProvider) GetTickers(ctx context.Context) (domain.Tickers, error) {
	resp, err := p.BybitClient.V5().Market().GetTickers(bybit.V5GetTickersParam{
		Category: bybit.CategoryV5Linear,
		Symbol:   nil,
		BaseCoin: nil,
		ExpDate:  nil,
	})

	if err != nil {
		return nil, fmt.Errorf("tickers service do: %w", err)
	}

	if resp.Result.LinearInverse.List == nil {
		return nil, ErrEmptyResponse
	}

	tickers := make([]*domain.Ticker, len(resp.Result.LinearInverse.List))
	for i, ticker := range resp.Result.LinearInverse.List {
		tickers[i] = convertTickerDomain(ticker)
	}

	return tickers, nil
}

func (p *BybitProvider) GetTickerInfo(ctx context.Context, symbol string) (*domain.Ticker, error) {
	s := bybit.SymbolV5(symbol)
	resp, err := p.BybitClient.V5().Market().GetTickers(bybit.V5GetTickersParam{
		Category: bybit.CategoryV5Linear,
		Symbol:   &s,
		BaseCoin: nil,
		ExpDate:  nil,
	})

	if err != nil {
		return nil, fmt.Errorf("tickers service do: %w", err)
	}

	if resp.Result.LinearInverse.List == nil || len(resp.Result.LinearInverse.List) != 1 {
		return nil, ErrEmptyResponse
	}

	return convertTickerDomain(resp.Result.LinearInverse.List[0]), nil
}

func (p *BybitProvider) GetOrderBook(ctx context.Context, symbol string) (*domain.OrderBook, error) {
	limit := 100
	resp, err := p.BybitClient.V5().Market().GetOrderbook(bybit.V5GetOrderbookParam{
		Category: bybit.CategoryV5Linear,
		Symbol:   bybit.SymbolV5(symbol),
		Limit:    &limit,
	})
	if err != nil {
		return nil, fmt.Errorf("get order book: %w", err)
	}

	return convertBidAskToDomain(resp.Result), nil
}
