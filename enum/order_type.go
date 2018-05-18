package enum

type OrderType int

const (
	OrderType_MARKET = iota
	OrderType_LIMIT
	OrderType_STOP
	OrderType_CANCEL
)
