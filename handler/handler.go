package handler

import "github.com/gin-gonic/gin"

func StatusHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Normal",
	})
}
