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
