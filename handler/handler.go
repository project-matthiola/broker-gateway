package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/rudeigerc/broker-gateway/model"
	"github.com/rudeigerc/broker-gateway/service"
	"github.com/spf13/viper"
)

type Error struct {
	Status  int    `json:"status"`
	Title   string `json:"title"`
	Message string `json:"message"`
}

func StatusHandler(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Normal"})
}

func ErrorHandler(c *gin.Context, code int, message string) {
	response := &Error{
		Status:  code,
		Title:   http.StatusText(code),
		Message: message,
	}
	c.JSON(code, gin.H{"error": response})
}

func validate(c *gin.Context) (model.Firm, error) {
	tokenString := c.Request.Header.Get("Authorization")

	if tokenString != "" {
		token, _ := jwt.ParseWithClaims(strings.Split(tokenString, " ")[1], &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(viper.GetString("auth.secret")), nil
		})

		if claims, ok := token.Claims.(*jwt.StandardClaims); ok && token.Valid {
			firmID, _ := strconv.Atoi(claims.Subject)
			firm := service.Auth{}.Validate(firmID)
			return firm, nil
		}
	}
	return model.Firm{}, fmt.Errorf("[hanlder.auth.validate] [ERROR] Invalid token")
}

func adminValidate(c *gin.Context) error {
	tokenString := c.Request.Header.Get("Authorization")

	if tokenString != "" {
		token, _ := jwt.ParseWithClaims(strings.Split(tokenString, " ")[1], &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(viper.GetString("auth.admin.secret")), nil
		})

		if _, ok := token.Claims.(*jwt.StandardClaims); ok && token.Valid {
			return nil
		}
	}
	return fmt.Errorf("[hanlder.auth.adminValidate] [ERROR] Invalid token")
}
