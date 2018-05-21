package server

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type ErrorResponse struct {
	Status  int    `json:"status"`
	Title   string `json:"title"`
	Message string `json:"message"`
}

func StatusHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Normal",
	})
}

type Auth struct {
	FirmName  string `json:"firm_name" binding:"required"`
	ValidTime int    `json:"valid_time" binding:"required"`
}

func AuthHandler(c *gin.Context) {
	var json Auth
	if err := c.ShouldBindJSON(&json); err == nil {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"firm_name": json.FirmName,
		})
		tokenString, _ := token.SignedString([]byte(viper.GetString("auth.secret")))
		c.JSON(http.StatusOK, gin.H{"token": tokenString})
	}
	return
}

func ValidationHandler(c *gin.Context) {
	tokenString := c.Request.Header.Get("Authorization")
	if tokenString == "" {
		response := ErrorResponse{
			http.StatusUnauthorized,
			http.StatusText(http.StatusUnauthorized),
			"Invalid token.",
		}
		c.JSON(http.StatusUnauthorized, response)
	} else {
		token, err := jwt.Parse(strings.Split(tokenString, " ")[1], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(viper.GetString("auth.secret")), nil
		})

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			log.Println(claims["firm_name"])
		} else {
			log.Println(err)
		}
	}
}
