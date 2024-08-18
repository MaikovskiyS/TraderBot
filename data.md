	var symbol bybit.SymbolV5 = "BNBUSDT"
	// price := "525.70"
	// err = createOrder(symbol, bybitProvider, price)
	// if err != nil {
	// 	return err
	// }

	// triggerPrice := "510.00"
	// var triggerBy bybit.TriggerBy = bybit.TriggerByLastPrice
	// var triggerDirection bybit.TriggerDirection = bybit.TriggerDirectionFall

	takeProfit := "540"
	takeProfitSize := "0.01"
	stopLoss := "500"
	stopLossSize := "0.02"
	var idx bybit.PositionIdx = bybit.PositionIdxHedgeBuy
	var tpSlMode bybit.TpSlMode = bybit.TpSlModePartial
	var orderType bybit.OrderType = bybit.OrderTypeLimit
	

	_, err = bybitProvider.BybitClient.V5().Position().SetTradingStop(bybit.V5SetTradingStopParam{
		Category:     bybit.CategoryV5Linear,
		Symbol:       symbol,
		PositionIdx:  idx,
		TakeProfit:   &takeProfit,
		StopLoss:     nil,
		TrailingStop: new(string),
		TpTriggerBy:  nil,
		SlTriggerBy:  nil,
		ActivePrice:  nil,
		TpSize:       &takeProfitSize,
		SlSize:       nil,
		TpslMode:     &tpSlMode,
		TpLimitPrice: &takeProfit,
		SlLimitPrice: nil,
		TpOrderType:  &orderType,
		SlOrderType:  nil,
	})
	if err != nil {
		return err
	}

	resp, err := bybitProvider.BybitClient.V5().Order().GetOpenOrders(bybit.V5GetOpenOrdersParam{
		Category: bybit.CategoryV5Linear,
		Symbol:   &symbol,
	})
	if err != nil {
		return err
	}

	var orderID string
	for _, order := range resp.Result.List {

		fmt.Printf("order: %v\n", order)
		fmt.Printf("orderStopType: %v\n", order.StopOrderType)
		if order.StopOrderType == "PartialStopLoss" {
			orderID = order.OrderID
		}

	}

	partialQTY := "0.01"
	//stopLoss := "515.00"
	tgPrice := "520.00"
	var trigBy bybit.TriggerBy = bybit.TriggerByMarkPrice
	_, err = bybitProvider.BybitClient.V5().Order().AmendOrder(bybit.V5AmendOrderParam{
		Category:     bybit.CategoryV5Linear,
		Symbol:       symbol,
		OrderID:      &orderID,
		TriggerPrice: &tgPrice,
		Qty:          &partialQTY,
		Price:        &stopLoss,
		TakeProfit:   new(string),
		StopLoss:     nil,
		TpTriggerBy:  nil,
		SlTriggerBy:  &trigBy,
		TriggerBy:    &trigBy,
	})
	if err != nil {
		return err
	}
	fmt.Printf("\"end\": %v\n", "end")





/*
resp, err := trading.Provider.GetOpenClosedOrdersByTicker(context.Background(), "CRVUSDT")
if err != nil {
	return err
}

for _, order := range resp.Orders {
	fmt.Printf("order: %v\n", order)

	if order.StopOrderType == domain.TakeProfitOrderType {
		tp := "0.3155"
		qty := "40"
		err := trading.Provider.AmendOrder(context.Background(), &traidingservice.AmendOrderParams{
			Symbol:       "CRVUSDT",
			OrderID:      &order.OrderID,
			OrderLinkID:  &order.OrderLinkID,
			OrderIv:      nil,
			TriggerPrice: &tp,
			Qty:          &qty,
			Price:        nil,
			TakeProfit:   nil,
			StopLoss:     new(string),
		})
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
}
*/

func createOrder(symbol bybit.SymbolV5, pr *bybit_provider.BybitProvider, price string) error {

	// closeOnTrigger := true
	reduceOnly := false

	// triggerPrice := "510.00"
	// var triggerBy bybit.TriggerBy = bybit.TriggerByLastPrice
	// var triggerDirection bybit.TriggerDirection = bybit.TriggerDirectionFall

	takeProfit := "540"
	//takeProfitLimit := "535"
	stopLoss := "500"
	//stopLossLimit := "505"
	var idx bybit.PositionIdx = bybit.PositionIdxHedgeBuy
	var tpSlMode bybit.TpSlMode = bybit.TpSlModePartial
	var orderType bybit.OrderType = bybit.OrderTypeLimit

	_, err := pr.BybitClient.V5().Order().CreateOrder(bybit.V5CreateOrderParam{
		Category:     bybit.CategoryV5Linear,
		Symbol:       symbol,
		Side:         bybit.SideBuy,
		OrderType:    bybit.OrderTypeLimit,
		Qty:          "0.02",
		Price:        &price,
		PositionIdx:  &idx,
		TakeProfit:   &takeProfit,
		StopLoss:     &stopLoss,
		TpTriggerBy:  nil, //by last price
		SlTriggerBy:  nil, //by last price
		ReduceOnly:   &reduceOnly,
		TpSlMode:     &tpSlMode,
		TpLimitPrice: nil,
		SlLimitPrice: nil,
		TpOrderType:  &orderType,
		SlOrderType:  &orderType,
		MarketUnit:   nil,
	})
	if err != nil {
		return err
	}
	return nil
}
