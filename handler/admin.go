package handler

import (
	"net/http"
	"time"

	"strings"

	"github.com/gin-gonic/gin"
	"github.com/quickfixgo/enum"
	"github.com/rudeigerc/broker-gateway/service"
	"github.com/rudeigerc/broker-gateway/tool"
	"github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
)

func OrderHandler(c *gin.Context) {
	orders := service.Order{}.Orders()
	c.JSON(http.StatusOK, gin.H{"data": orders})
}

type Trader struct {
	Firm   string `json:"firm"`
	Side   string `json:"side"`
	Trader string `json:"trader"`
}

type TradeResponse struct {
	TradeID    uuid.UUID       `json:"trade_id"`
	Quantity   decimal.Decimal `json:"quantity"`
	Price      decimal.Decimal `json:"price"`
	FuturesID  string          `json:"futures_id"`
	Initiator  Trader          `json:"initiator"`
	Completion Trader          `json:"completion"`
	CreatedAt  time.Time       `json:"trade_time"`
}

func TradeHandler(c *gin.Context) {
	tradeID := strings.Trim(c.Param("trade_id"), "/")
	if tradeID != "" {
		trade := service.Trade{}.TradeByID(tradeID)
		response := TradeResponse{
			TradeID:   trade.TradeID,
			Quantity:  trade.Quantity,
			Price:     trade.Price,
			FuturesID: trade.FuturesID,
			Initiator: Trader{
				Firm:   service.Auth{}.FirmNameByID(trade.InitiatorID),
				Side:   tool.Convert(enum.Side(trade.InitiatorSide)),
				Trader: trade.InitiatorName,
			},
			Completion: Trader{
				Firm:   service.Auth{}.FirmNameByID(trade.CompletionID),
				Side:   tool.Convert(enum.Side(trade.CompletionSide)),
				Trader: trade.CompletionName,
			},
			CreatedAt: trade.CreatedAt,
		}
		c.JSON(http.StatusOK, gin.H{"data": response})
	} else {
		trades := service.Trade{}.Trades()
		data := make([]TradeResponse, len(trades))
		for index, trade := range trades {
			response := TradeResponse{
				TradeID:   trade.TradeID,
				Quantity:  trade.Quantity,
				Price:     trade.Price,
				FuturesID: trade.FuturesID,
				Initiator: Trader{
					Firm:   service.Auth{}.FirmNameByID(trade.InitiatorID),
					Side:   tool.Convert(enum.Side(trade.InitiatorSide)),
					Trader: trade.InitiatorName,
				},
				Completion: Trader{
					Firm:   service.Auth{}.FirmNameByID(trade.CompletionID),
					Side:   tool.Convert(enum.Side(trade.CompletionSide)),
					Trader: trade.CompletionName,
				},
				CreatedAt: trade.CreatedAt,
			}
			data[index] = response
		}
		c.JSON(http.StatusOK, gin.H{"data": data})
	}
}
