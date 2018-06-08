package handler

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/quickfixgo/enum"
	"github.com/rudeigerc/broker-gateway/service"
	"github.com/rudeigerc/broker-gateway/tool"
	"github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
	"github.com/spf13/viper"
)

type adminAuthRequest struct {
	Name      string    `json:"name" binding:"required"`
	ExpiresAt time.Time `json:"expires_at" binding:"required" time_format:"2006-01-02 15:04:05"`
}

type adminAuthResponse struct {
	Name  string `json:"name"`
	Token string `json:"token"`
}

func AdminAuthHandler(c *gin.Context) {
	var req adminAuthRequest

	if err := c.ShouldBindJSON(&req); err == nil {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
			Audience:  req.Name,
			ExpiresAt: req.ExpiresAt.Unix(),
			Id:        uuid.NewV1().String(),
			IssuedAt:  time.Now().Unix(),
		})
		tokenString, _ := token.SignedString([]byte(viper.GetString("auth.admin.secret")))

		res := adminAuthResponse{
			Name:  req.Name,
			Token: tokenString,
		}
		c.JSON(http.StatusOK, gin.H{"data": res})
	} else {
		ErrorHandler(c, http.StatusBadRequest, "Invalid request message framing.")
	}
}

type OrderResponse struct {
	OrderID      uuid.UUID `json:"order_id" gorm:"primary_key"`
	OrderType    string    `json:"order_type"`
	Side         string    `json:"side"`
	FuturesID    string    `json:"futures_id"`
	Firm         string    `json:"firm"`
	TraderName   string    `json:"trader_name"`
	Quantity     string    `json:"quantity"`
	OpenQuantity string    `json:"open_quantity"`
	Price        string    `json:"price"`
	StopPrice    string    `json:"stop_price"`
	Status       string    `json:"status"`
	CreatedAt    string    `json:"created_at"`
	UpdatedAt    string    `json:"updated_at"`
}

func AdminOrderHandler(c *gin.Context) {
	err := adminValidate(c)
	if err != nil {
		log.Print(err)
		ErrorHandler(c, http.StatusUnauthorized, "Invalid token.")
		return
	}

	page, err := strconv.Atoi(c.DefaultQuery("page", "0"))
	if err != nil {
		ErrorHandler(c, http.StatusBadRequest, "Invalid page number.")
		return
	}

	orders := service.Order{}.OrdersWithPage(page)
	data := make([]OrderResponse, len(orders))
	for index, order := range orders {
		response := OrderResponse{
			OrderID:      order.OrderID,
			OrderType:    tool.Convert(order.OrderType),
			Side:         tool.Convert(order.Side),
			FuturesID:    order.FuturesID,
			Firm:         service.Auth{}.FirmNameByID(order.FirmID),
			TraderName:   order.TraderName,
			Quantity:     order.Quantity.String(),
			OpenQuantity: order.OpenQuantity.String(),
			Price:        order.Price.String(),
			StopPrice:    order.StopPrice.String(),
			Status:       tool.Convert(order.Status),
			CreatedAt:    order.CreatedAt.String(),
			UpdatedAt:    order.UpdatedAt.String(),
		}
		data[index] = response
	}
	c.JSON(http.StatusOK, gin.H{
		"data":  data,
		"count": len(data),
		"page":  page,
	})
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

func AdminTradeHandler(c *gin.Context) {
	err := adminValidate(c)
	if err != nil {
		log.Print(err)
		ErrorHandler(c, http.StatusUnauthorized, "Invalid token.")
		return
	}

	page, err := strconv.Atoi(c.DefaultQuery("page", "0"))
	if err != nil {
		ErrorHandler(c, http.StatusBadRequest, "Invalid page number.")
		return
	}

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
		trades := service.Trade{}.TradesWithPage(page)
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
		c.JSON(http.StatusOK, gin.H{
			"data":  data,
			"count": len(data),
			"page":  page,
		})
	}
}
