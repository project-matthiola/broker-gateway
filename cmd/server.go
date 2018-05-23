package cmd

import (
	"github.com/gin-gonic/gin"
	"github.com/rudeigerc/broker-gateway/handler"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run server",
	Long:  "Run server",
	Run: func(cmd *cobra.Command, args []string) {
		if !viper.GetBool("gin.debug") {
			gin.SetMode(gin.ReleaseMode)
		}

		router := gin.Default()

		v1 := router.Group("/api/v1")
		{
			v1.GET("/status", handler.StatusHandler)

			auth := v1.Group("/auth")
			{
				auth.GET("", handler.ValidationHandler)
				auth.POST("", handler.AuthHandler)
			}
		}

		router.Run(":" + viper.GetString("gin.port"))
	},
}

func init() {
	serverCmd.PersistentFlags().IntP("gin.port", "p", 8080, "port of server")

	viper.BindPFlags(serverCmd.PersistentFlags())
}
