package cmd

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/micro/go-web"
	"github.com/rudeigerc/broker-gateway/broadcaster"
	"github.com/rudeigerc/broker-gateway/mapper"
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

		mapper.NewDB()

		hub := broadcaster.NewHub()

		go hub.RunBroadcaster()
		go hub.RunOrderBookWatcher()
		go hub.RunTradeWatcher()

		defer func() {
			mapper.DB.Close()
			hub.EtcdClient.Close()
		}()

		router := gin.Default()

		router.GET("/broadcaster", func(c *gin.Context) {
			broadcaster.SocketHandler(hub, c)
		})

		service := web.NewService(
			web.Name("github.com.rudeigerc.broker-gateway.broadcaster"),
			web.Version("1.0.0"),
			web.Handler(router),
		)

		if err := service.Run(); err != nil {
			log.Fatalf("[cmd.broadcaster.broadcasterCmd] [FETAL] %s", err)
		}
	},
}
