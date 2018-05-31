package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
