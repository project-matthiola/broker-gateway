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

	firmID, err := strconv.Atoi(c.DefaultQuery("firm_id", "-1"))
	if err != nil {
		ErrorHandler(c, http.StatusBadRequest, "Invalid firm ID.")
		return
	}

	futuresID := c.Query("futures_id")
	traderName := c.Query("trader_name")
	status := c.Query("status")

	page, err := strconv.Atoi(c.DefaultQuery("page", "0"))
	if err != nil {
		ErrorHandler(c, http.StatusBadRequest, "Invalid page number.")
		return
	}

	total, orders := service.Order{}.OrdersWithCondition(firmID, futuresID, traderName, status, page)
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
	CreatedAt  string          `json:"trade_time"`
}

func AdminTradeHandler(c *gin.Context) {
	err := adminValidate(c)
	if err != nil {
		log.Print(err)
		ErrorHandler(c, http.StatusUnauthorized, "Invalid token.")
		return
	}

	firmID, err := strconv.Atoi(c.DefaultQuery("firm_id", "-1"))
	if err != nil {
		ErrorHandler(c, http.StatusBadRequest, "Invalid firm ID.")
		return
	}

	futuresID := c.Query("futures_id")
	traderName := c.Query("trader_name")

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
			CreatedAt: trade.CreatedAt.String(),
		}
		c.JSON(http.StatusOK, gin.H{"data": response})
	} else {
		total, trades := service.Trade{}.TradesWithCondition(firmID, futuresID, traderName, page)
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
}

func AdminFirmHandler(c *gin.Context) {
	err := adminValidate(c)
	if err != nil {
		log.Print(err)
		ErrorHandler(c, http.StatusUnauthorized, "Invalid token.")
		return
	}
	firms := service.Firm{}.Firms()
	c.JSON(http.StatusOK, gin.H{"data": firms})
}

type Children struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

type FuturesResponse struct {
	Value    string     `json:"value"`
	Label    string     `json:"label"`
	Children []Children `json:"children"`
}

func AdminFuturesHandler(c *gin.Context) {
	err := adminValidate(c)
	if err != nil {
		log.Print(err)
		ErrorHandler(c, http.StatusUnauthorized, "Invalid token.")
		return
	}
	futures := service.Futures{}.Futures()

	var data []FuturesResponse
	for k, v := range futures {
		var childrens []Children
		for _, f := range v {
			childrens = append(childrens, Children{
				Value: f.FuturesID,
				Label: f.FuturesID,
			})
		}
		response := FuturesResponse{
			Value:    k,
			Label:    k,
			Children: childrens,
		}
		data = append(data, response)
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}
