package domain

type Side string

const (
	BuySide  Side = "Buy"
	SellSide Side = "Sell"
)

const (
	StopLossOrderType          = "StopLoss"
	TakeProfitOrderType        = "TakeProfit"
	PartialStopLossOrderType   = "PartialStopLoss"
	PartialTakeProfitOrderType = "PartialTakeProfit"
)

type OrderStatus string

const (

	// OrderStatusCreated :
	OrderStatusCreated = OrderStatus("Created")
	// OrderStatusRejected :
	OrderStatusRejected = OrderStatus("Rejected")
	// OrderStatusNew :
	OrderStatusNew = OrderStatus("New")
	// OrderStatusPartiallyFilled :
	OrderStatusPartiallyFilled = OrderStatus("PartiallyFilled")
	// OrderStatusFilled :
	OrderStatusFilled = OrderStatus("Filled")
	// OrderStatusCancelled :
	OrderStatusCancelled = OrderStatus("Cancelled")
	// OrderStatusPendingCancel :
	OrderStatusPendingCancel = OrderStatus("PendingCancel")

	// OrderStatusUntriggered : Only for conditional orders
	OrderStatusUntriggered = OrderStatus("Untriggered")
	// OrderStatusDeactivated : Only for conditional orders
	OrderStatusDeactivated = OrderStatus("Deactivated")
	// OrderStatusTriggered : Only for conditional orders
	OrderStatusTriggered = OrderStatus("Triggered")
	// OrderStatusActive : Only for conditional orders
	OrderStatusActive = OrderStatus("Active")
)

type Order struct {
	Symbol             string
	OrderLinkID        string
	OrderID            string
	CancelType         string
	AvgPrice           float64
	StopOrderType      string
	LastPriceOnCreated string
	OrderStatus        OrderStatus
	TakeProfit         float64
	Price              float64
	OrderIv            string
	CreatedTime        string
	Side               Side
	TriggerPrice       float64
	Qty                float64
	StopLoss           float64
	Precision          int
}
