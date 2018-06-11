package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/quickfixgo/enum"
	"github.com/rudeigerc/broker-gateway/service"
	"github.com/rudeigerc/broker-gateway/tool"
)

func TradeHandler(c *gin.Context) {
	futuresID := c.Query("futures_id")
	traderName := c.Query("trader_name")
	page, err := strconv.Atoi(c.DefaultQuery("page", "0"))
	if err != nil {
		ErrorHandler(c, http.StatusBadRequest, "Invalid page number.")
		return
	}

	firm, err := validate(c)
	if err != nil {
		log.Print(err)
		ErrorHandler(c, http.StatusUnauthorized, "Invalid token.")
		return
	}

	total, trades := service.Trade{}.TradesWithCondition(firm.FirmID, futuresID, traderName, page)
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
			CreatedAt: trade.CreatedAt.String(),
		}
		data[index] = response
	}
	c.JSON(http.StatusOK, gin.H{
		"data":  data,
		"count": len(data),
		"total": total,
		"page":  page,
	})
}

func OrderHandler(c *gin.Context) {
	futuresID := c.Query("futures_id")
	traderName := c.Query("trader_name")
	status := c.Query("status")
	page, err := strconv.Atoi(c.DefaultQuery("page", "0"))
	if err != nil {
		ErrorHandler(c, http.StatusBadRequest, "Invalid page number.")
		return
	}

	firm, err := validate(c)
	if err != nil {
		log.Print(err)
		ErrorHandler(c, http.StatusUnauthorized, "Invalid token.")
		return
	}

	total, orders := service.Order{}.OrdersWithCondition(firm.FirmID, futuresID, traderName, status, page)
	data := make([]OrderResponse, len(orders))
	for index, order := range orders {
		response := OrderResponse{
			OrderID:      order.OrderID,
			OrderType:    tool.Convert(enum.OrdType(order.OrderType)),
			Side:         tool.Convert(enum.Side(order.Side)),
			FuturesID:    order.FuturesID,
			Firm:         service.Auth{}.FirmNameByID(order.FirmID),
			TraderName:   order.TraderName,
			Quantity:     order.Quantity.String(),
			OpenQuantity: order.OpenQuantity.String(),
			Price:        order.Price.String(),
			StopPrice:    order.StopPrice.String(),
			Status:       tool.Convert(enum.OrdStatus(order.Status)),
			CreatedAt:    order.CreatedAt.String(),
			UpdatedAt:    order.UpdatedAt.String(),
		}
		data[index] = response
	}
	c.JSON(http.StatusOK, gin.H{
		"data":  data,
		"count": len(data),
		"total": total,
		"page":  page,
	})
}
