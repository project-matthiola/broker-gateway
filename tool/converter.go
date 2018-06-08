package tool

import "github.com/quickfixgo/enum"

func Convert(origin interface{}) string {
	switch origin.(type) {
	case enum.OrdType:
		switch origin {
		case enum.OrdType_MARKET:
			return "MARKET"
		case enum.OrdType_LIMIT:
			return "LIMIT"
		case enum.OrdType_STOP:
			return "STOP"
		case enum.OrdType_STOP_LIMIT:
			return "STOP_LIMIT"
		}
	case enum.Side:
		switch origin {
		case enum.Side_BUY:
			return "BUY"
		case enum.Side_SELL:
			return "SELL"
		}
	case enum.OrdStatus:
		switch origin {
		case enum.OrdStatus_NEW:
			return "NEW"
		case enum.OrdStatus_PARTIALLY_FILLED:
			return "PARTIALLY_FILLED"
		case enum.OrdStatus_FILLED:
			return "FILLED"
		case enum.OrdStatus_CANCELED:
			return "CANCELED"
		case enum.OrdStatus_REJECTED:
			return "REJECTED"
		case enum.OrdStatus_PENDING_NEW:
			return "PENDING_NEW"
		}
	}
	return ""
}
