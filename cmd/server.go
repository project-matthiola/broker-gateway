package cmd

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/micro/go-web"
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

		service := web.NewService(
			web.Name("github.com.rudeigerc.broker-gateway.server"),
			web.Version("1.0.0"),
			web.Address(":"+viper.GetString("gin.port")),
			web.Handler(router),
		)

		if err := service.Run(); err != nil {
			log.Fatalf("[cmd.server.serverCmd] [FETAL] %s", err)
		}
	},
}

func init() {
	serverCmd.PersistentFlags().IntP("gin.port", "p", 8080, "port of HTTP server")

	viper.BindPFlags(serverCmd.PersistentFlags())
}
