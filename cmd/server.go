package cmd

import (
	"github.com/gin-gonic/gin"
	"github.com/rudeigerc/broker-gateway/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run server",
	Long:  "Run server",
	Run: func(cmd *cobra.Command, args []string) {
		router := gin.Default()

		v1 := router.Group("/api/v1")
		{
			v1.GET("/status", server.StatusHandler)
			v1.GET("/auth", server.ValidationHandler)
			v1.POST("/auth", server.AuthHandler)
		}

		router.Run(":" + viper.GetString("http.port"))
	},
}
