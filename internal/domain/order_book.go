package domain

type OrderbookBidAsk struct {
	Price    float64
	Quantity float64
}

type OrderBook struct {
	Symbol string
	Bids   []OrderbookBidAsk
	Asks   []OrderbookBidAsk
}

/*
func CalculateStopLoss(orderBook OrderBook, pivotPrice float64, isLong bool, rangeLimit float64) (float64, error) {

}
*/
