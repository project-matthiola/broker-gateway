package cmd

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/micro/go-web"
	"github.com/rudeigerc/broker-gateway/broadcaster"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var broadcasterCmd = &cobra.Command{
	Use:   "broadcaster",
	Short: "Run WebSocket server",
	Long:  "Run WebSocket server",
	Run: func(cmd *cobra.Command, args []string) {
		if !viper.GetBool("gin.debug") {
			gin.SetMode(gin.ReleaseMode)
		}

		hub := broadcaster.NewHub()
		go hub.RunBroadcaster()
		go hub.RunWatcher()

		router := gin.Default()

		router.GET("/futures", func(c *gin.Context) {
			broadcaster.FuturesSocketHandler(hub, c)
		})

		service := web.NewService(
			web.Name("github.com.rudeigerc.broker-gateway.broadcaster"),
			web.Version("1.0.0"),
			web.Address(":"+viper.GetString("websocket.port")),
			web.Handler(router),
		)

		if err := service.Run(); err != nil {
			log.Fatalf("[cmd.broadcaster.broadcasterCmd] [FETAL] %s", err)
		}
	},
}

func init() {
	serverCmd.PersistentFlags().Int("websocket.port", 8000, "port of WebSocket server")

	viper.BindPFlags(serverCmd.PersistentFlags())
}
