package cmd

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/micro/go-web"
	"github.com/rudeigerc/broker-gateway/handler"
	"github.com/rudeigerc/broker-gateway/mapper"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run HTTP server",
	Long:  "Run HTTP server",
	Run: func(cmd *cobra.Command, args []string) {
		if !viper.GetBool("gin.debug") {
			gin.SetMode(gin.ReleaseMode)
		}

		router := gin.Default()

		v1 := router.Group("/server/api/v1")
		{
			v1.GET("/status", handler.StatusHandler)

			auth := v1.Group("/auth")
			{
				auth.GET("", handler.ValidationHandler)
				auth.POST("", handler.AuthHandler)
			}

			admin := v1.Group("/admin")
			{
				admin.POST("/auth", handler.AdminAuthHandler)
				admin.GET("/orders", handler.AdminOrderHandler)
				admin.GET("/trades/*trade_id", handler.AdminTradeHandler)
			}

			v1.GET("/orders", handler.OrderHandler)
			v1.GET("/trades", handler.TradeHandler)
		}

		mapper.NewDB()
		defer mapper.DB.Close()

		service := web.NewService(
			web.Name("github.com.rudeigerc.broker-gateway.server"),
			web.Version("1.0.0"),
			web.Handler(router),
		)

		if err := service.Run(); err != nil {
			log.Fatalf("[cmd.server.serverCmd] [FETAL] %s", err)
		}
	},
}
