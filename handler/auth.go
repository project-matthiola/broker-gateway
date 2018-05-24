package handler

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/rudeigerc/broker-gateway/service"
	"github.com/satori/go.uuid"
	"github.com/spf13/viper"
)

type authRequest struct {
	FirmName  string    `json:"firm_name" binding:"required"`
	ExpiresAt time.Time `json:"expires_at" binding:"required" time_format:"2006-01-02 15:04:05"`
}

type authResponse struct {
	FirmID   int    `json:"firm_id"`
	FirmName string `json:"firm_name"`
	Credit   int    `json:"credit"`
	Token    string `json:"token"`
}

func AuthHandler(c *gin.Context) {
	var req authRequest

	if err := c.ShouldBindJSON(&req); err == nil {
		firm := service.Auth{}.Sign(req.FirmName)

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
			Audience:  req.FirmName,
			ExpiresAt: req.ExpiresAt.Unix(),
			Id:        uuid.NewV1().String(),
			IssuedAt:  time.Now().Unix(),
			Subject:   strconv.Itoa(firm.FirmID),
		})
		tokenString, _ := token.SignedString([]byte(viper.GetString("auth.secret")))

		res := authResponse{
			FirmID:   firm.FirmID,
			FirmName: firm.FirmName,
			Credit:   firm.Credit,
			Token:    tokenString,
		}
		c.JSON(http.StatusOK, gin.H{"data": res})
	} else {
		ErrorHandler(c, http.StatusBadRequest, "Invalid request message framing.")
	}
}

func ValidationHandler(c *gin.Context) {
	tokenString := c.Request.Header.Get("Authorization")

	if tokenString != "" {
		token, _ := jwt.ParseWithClaims(strings.Split(tokenString, " ")[1], &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(viper.GetString("auth.secret")), nil
		})

		if claims, ok := token.Claims.(*jwt.StandardClaims); ok && token.Valid {
			firmID, _ := strconv.Atoi(claims.Subject)
			firm := service.Auth{}.Validate(firmID)
			res := authResponse{
				FirmID:   firm.FirmID,
				FirmName: firm.FirmName,
				Credit:   firm.Credit,
				Token:    strings.Split(tokenString, " ")[1],
			}
			c.JSON(http.StatusOK, gin.H{"data": res})
			return
		}
	}

	ErrorHandler(c, http.StatusUnauthorized, "Invalid token.")
}
