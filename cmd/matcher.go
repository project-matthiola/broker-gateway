package cmd

import (
	"log"
	"time"

	"github.com/micro/go-micro"
	"github.com/rudeigerc/broker-gateway/mapper"
	"github.com/rudeigerc/broker-gateway/matcher"
	"github.com/spf13/cobra"
)

var matcherCmd = &cobra.Command{
	Use:   "matcher",
	Short: "Run matcher",
	Long:  "Run matcher",
	Run: func(cmd *cobra.Command, args []string) {
		m := matcher.NewMatcher()
		mapper.NewDB()

		defer func() {
			m.Stop()
			mapper.DB.Close()
		}()

		service := micro.NewService(
			micro.Name("github.com.rudeigerc.broker-gateway.matcher"),
			micro.RegisterTTL(time.Minute),
			micro.RegisterInterval(time.Second*30),
		)

		if err := service.Run(); err != nil {
			log.Fatalf("[cmd.matcher.matcherCmd] [FETAL] %s", err)
		}
	},
}
