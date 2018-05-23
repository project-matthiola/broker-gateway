package server

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
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
	FirmName  string    `json:"firm_name" binding:"required"`
	ExpiresAt time.Time `json:"expires_at" binding:"required" time_format:"2006-01-02 15:04:05"`
}

func AuthHandler(c *gin.Context) {
	var json Auth
	if err := c.ShouldBindJSON(&json); err == nil {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
			Audience:  json.FirmName,
			ExpiresAt: json.ExpiresAt.Unix(),
			Id:        uuid.NewV1().String(),
			IssuedAt:  time.Now().Unix(),
		})
		tokenString, _ := token.SignedString([]byte(viper.GetString("auth.secret")))
		c.JSON(http.StatusOK, gin.H{"token": tokenString})
	}
	return
}

func ValidationHandler(c *gin.Context) {
	tokenString := c.Request.Header.Get("Authorization")

	if tokenString != "" {
		token, _ := jwt.ParseWithClaims(strings.Split(tokenString, " ")[1], &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(viper.GetString("auth.secret")), nil
		})

		if claims, ok := token.Claims.(*jwt.StandardClaims); ok && token.Valid {
			log.Println(claims.Audience)
			return
		}
	}

	response := ErrorResponse{
		http.StatusUnauthorized,
		http.StatusText(http.StatusUnauthorized),
		"Invalid token.",
	}
	c.JSON(http.StatusUnauthorized, response)
}
