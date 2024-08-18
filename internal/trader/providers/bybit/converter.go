package bybit_provider

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/MaikovskiyS/TraderBot/internal/domain"
	trading "github.com/MaikovskiyS/TraderBot/internal/trader/services/traiding_service"
	"github.com/hirokisan/bybit/v2"
)

func convertCandleToDomain(kline bybit.V5GetKlineItem) (*domain.Candle, error) {
	time, err := strconv.Atoi(kline.StartTime)
	if err != nil {
		return nil, fmt.Errorf("convert start time: %w", err)
	}

	open, err := strconv.ParseFloat(kline.Open, 64)
	if err != nil {
		return nil, fmt.Errorf("convert open price: %w", err)
	}

	high, err := strconv.ParseFloat(kline.High, 64)
	if err != nil {
		return nil, fmt.Errorf("convert high price: %w", err)
	}

	low, err := strconv.ParseFloat(kline.Low, 64)
	if err != nil {
		return nil, fmt.Errorf("convert low price: %w", err)
	}

	close, err := strconv.ParseFloat(kline.Close, 64)
	if err != nil {
		return nil, fmt.Errorf("convert close price: %w", err)
	}

	volume, err := strconv.ParseFloat(kline.Volume, 64)
	if err != nil {
		return nil, fmt.Errorf("convert volume: %w", err)
	}

	return &domain.Candle{
		Time:   int64(time),
		Open:   open,
		High:   high,
		Low:    low,
		Close:  close,
		Volume: volume,
	}, nil
}

func convertTickerDomain(data bybit.V5GetTickersLinearInverseItem) *domain.Ticker {
	precision := 0
	parts := strings.Split(data.LastPrice, ".")
	if len(parts) == 2 {
		precision = len(parts[1])
	}

	return &domain.Ticker{
		Symbol:            string(data.Symbol),
		LastPrice:         data.LastPrice,
		IndexPrice:        data.IndexPrice,
		MarkPrice:         data.MarkPrice,
		PrevPrice24h:      data.PrevPrice24H,
		Price24hPcnt:      data.Price24HPcnt,
		HighPrice24h:      data.HighPrice24H,
		LowPrice24h:       data.LowPrice24H,
		PrevPrice1h:       data.PrevPrice1H,
		OpenInterest:      data.OpenInterest,
		OpenInterestValue: data.OpenInterestValue,
		Turnover24h:       data.Turnover24H,
		Volume24h:         data.Volume24H,
		Ask1Size:          data.Ask1Size,
		Bid1Price:         data.Bid1Price,
		Ask1Price:         data.Ask1Price,
		Bid1Size:          data.Bid1Size,
		Ask1Iv:            data.Ask1Price,
		Bid1Iv:            data.Bid1Price,
		MarkIv:            data.MarkPrice,
		Precision:         precision,
	}
}

func convertOrderRequest(req *trading.CreateOrderParams) bybit.V5CreateOrderParam {
	var tpSlMode bybit.TpSlMode = bybit.TpSlModePartial
	var orderType bybit.OrderType = bybit.OrderTypeLimit
	var positionIdx bybit.PositionIdx
	switch req.Side {
	case domain.BuySide:
		positionIdx = 1
	case domain.SellSide:
		positionIdx = 2
	}

	return bybit.V5CreateOrderParam{
		Category:              bybit.CategoryV5Linear,
		Symbol:                bybit.SymbolV5(req.Coin),
		Side:                  bybit.Side(req.Side),
		OrderType:             "Limit",
		Qty:                   req.Quantity,
		IsLeverage:            nil,
		Price:                 &req.OpenPrice,
		TriggerDirection:      nil,
		OrderFilter:           nil,
		TriggerPrice:          nil,
		TriggerBy:             nil,
		OrderIv:               new(string),
		TimeInForce:           nil,
		PositionIdx:           &positionIdx,
		OrderLinkID:           new(string),
		TakeProfit:            &req.TakeProfit,
		StopLoss:              &req.StopLoss,
		TpTriggerBy:           nil,
		SlTriggerBy:           nil,
		ReduceOnly:            new(bool),
		CloseOnTrigger:        new(bool),
		SmpType:               new(string),
		MarketMakerProtection: new(bool),
		TpSlMode:              &tpSlMode,
		TpLimitPrice:          &req.TakeProfit,
		SlLimitPrice:          &req.StopLoss,
		TpOrderType:           &orderType,
		SlOrderType:           &orderType,
		MarketUnit:            nil,
	}

}

func convertBidAskToDomain(data bybit.V5GetOrderbookResult) *domain.OrderBook {

	asks := make([]domain.OrderbookBidAsk, len(data.Asks))
	bids := make([]domain.OrderbookBidAsk, len(data.Bids))

	for i, ask := range data.Asks {
		asks[i] = convertBidAskValuesToDomain(ask)
	}

	for i, bid := range data.Bids {
		bids[i] = convertBidAskValuesToDomain(bid)
	}

	return &domain.OrderBook{
		Symbol: string(data.Symbol),
		Bids:   bids,
		Asks:   asks,
	}
}

func convertBidAskValuesToDomain(data bybit.V5GetOrderbookBidAsk) domain.OrderbookBidAsk {
	return domain.OrderbookBidAsk{
		Price:    convertStrToFloat(data.Price),
		Quantity: convertStrToFloat(data.Quantity),
	}
}
func convertStrToFloat(value string) float64 {
	float, _ := strconv.ParseFloat(value, 64)

	return float
}

func convertPositionInfoToDomain(data bybit.V5GetPositionInfoItem) *domain.Position {
	precision := 0
	parts := strings.Split(data.AvgPrice, ".")
	if len(parts) == 2 {
		precision = len(parts[1])
	}
	return &domain.Position{
		Symbol:          string(data.Symbol),
		AvgPrice:        convertStrToFloat(data.AvgPrice),
		TakeProfit:      convertStrToFloat(data.TakeProfit),
		PositionValue:   data.PositionValue,
		TpSlMode:        string(data.TpSlMode),
		TrailingStop:    data.TrailingStop,
		UnrealisedPnl:   data.UnrealisedPnl,
		MarkPrice:       convertStrToFloat(data.MarkPrice),
		PositionBalance: data.PositionBalance,
		Side:            domain.Side(data.Side),
		Size:            convertStrToFloat(data.Size),
		PositionStatus:  data.PositionStatus,
		StopLoss:        convertStrToFloat(data.StopLoss),
		CreatedTime:     data.CreatedTime,
		Precision:       precision,
	}
}

func convertOrderToDomain(data bybit.V5GetOrder) *domain.Order {
	precision := 0
	parts := strings.Split(data.AvgPrice, ".")
	if len(parts) == 2 {
		precision = len(parts[1])
	}

	return &domain.Order{
		Symbol:             string(data.Symbol),
		OrderLinkID:        data.OrderLinkID,
		OrderID:            data.OrderID,
		CancelType:         data.CancelType,
		AvgPrice:           convertStrToFloat(data.AvgPrice),
		StopOrderType:      data.StopOrderType,
		LastPriceOnCreated: data.LastPriceOnCreated,
		OrderStatus:        domain.OrderStatus(data.OrderStatus),
		TakeProfit:         convertStrToFloat(data.TakeProfit),
		Price:              convertStrToFloat(data.Price),
		OrderIv:            data.OrderIv,
		CreatedTime:        data.CreatedTime,
		Side:               domain.Side(data.Side),
		TriggerPrice:       convertStrToFloat(data.TriggerPrice),
		Qty:                convertStrToFloat(data.Qty),
		StopLoss:           convertStrToFloat(data.StopLoss),
		Precision:          precision,
	}
}

func convertAmendOrderParams(params *trading.AmendOrderParams) bybit.V5AmendOrderParam {
	return bybit.V5AmendOrderParam{
		Category:     bybit.CategoryV5Linear,
		Symbol:       bybit.SymbolV5(params.Symbol),
		OrderID:      params.OrderID,
		OrderLinkID:  params.OrderLinkID,
		OrderIv:      params.OrderIv,
		TriggerPrice: params.TriggerPrice,
		Qty:          params.Qty,
		Price:        params.Price,
		TakeProfit:   params.TakeProfit,
		StopLoss:     params.StopLoss,
	}
}

func convertTpParams(data *trading.SetTpParams) bybit.V5SetTradingStopParam {
	var tpSlMode bybit.TpSlMode = bybit.TpSlModePartial
	var orderType bybit.OrderType = bybit.OrderTypeLimit

	var idx bybit.PositionIdx
	switch data.Side {
	case domain.BuySide:
		idx = bybit.PositionIdxHedgeBuy
	case domain.SellSide:
		idx = bybit.PositionIdxHedgeSell
	}

	return bybit.V5SetTradingStopParam{
		Category:     bybit.CategoryV5Linear,
		Symbol:       bybit.SymbolV5(data.Symbol),
		PositionIdx:  idx,
		TakeProfit:   &data.TakeProfit,
		StopLoss:     nil,
		TrailingStop: nil,
		TpTriggerBy:  nil,
		SlTriggerBy:  nil,
		ActivePrice:  nil,
		TpSize:       &data.Size,
		SlSize:       nil,
		TpslMode:     &tpSlMode,
		TpLimitPrice: &data.TakeProfit,
		SlLimitPrice: nil,
		TpOrderType:  &orderType,
		SlOrderType:  nil,
	}
}

func convertSlParams(data *trading.SetSlParams) bybit.V5SetTradingStopParam {
	var tpSlMode bybit.TpSlMode = bybit.TpSlModePartial
	var orderType bybit.OrderType = bybit.OrderTypeLimit

	var idx bybit.PositionIdx
	switch data.Side {
	case domain.BuySide:
		idx = bybit.PositionIdxHedgeBuy
	case domain.SellSide:
		idx = bybit.PositionIdxHedgeSell
	}

	return bybit.V5SetTradingStopParam{
		Category:     bybit.CategoryV5Linear,
		Symbol:       bybit.SymbolV5(data.Symbol),
		PositionIdx:  idx,
		TakeProfit:   nil,
		StopLoss:     &data.StopLoss,
		TrailingStop: nil,
		TpTriggerBy:  nil,
		SlTriggerBy:  nil,
		ActivePrice:  nil,
		TpSize:       nil,
		SlSize:       &data.Size,
		TpslMode:     &tpSlMode,
		TpLimitPrice: nil,
		SlLimitPrice: &data.StopLoss,
		TpOrderType:  nil,
		SlOrderType:  &orderType,
	}
}
