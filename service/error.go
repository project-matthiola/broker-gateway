package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Error struct {
	Status  int    `json:"status"`
	Title   string `json:"title"`
	Message string `json:"message"`
}

func NewError(c *gin.Context, code int, message string) {
	response := &Error{
		Status:  code,
		Title:   http.StatusText(code),
		Message: message,
	}
	c.JSON(code, gin.H{"error": response})
}
