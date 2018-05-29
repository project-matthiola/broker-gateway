package cmd

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/micro/go-web"
	"github.com/rudeigerc/broker-gateway/broadcaster"
	"github.com/rudeigerc/broker-gateway/handler"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var broadcasterCmd = &cobra.Command{
	Use:   "broadcaster",
	Short: "Run broadcaster",
	Long:  "Run broadcaster",
	Run: func(cmd *cobra.Command, args []string) {
		if !viper.GetBool("gin.debug") {
			gin.SetMode(gin.ReleaseMode)
		}

		go broadcaster.HandleBroadcast()

		router := gin.Default()

		router.GET("/ping", handler.PingHandler)

		service := web.NewService(
			web.Name("github.com.rudeigerc.broker-gateway.broadcaster"),
			web.Version("1.0.0"),
			web.Address(":"+viper.GetString("websocket.port")),
			web.Handler(router),
		)

		if err := service.Run(); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	serverCmd.PersistentFlags().Int("websocket.port", 8000, "port of WebSocket server")

	viper.BindPFlags(serverCmd.PersistentFlags())
}
