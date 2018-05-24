package cmd

import (
	"log"

	"github.com/micro/go-micro"
	"github.com/rudeigerc/broker-gateway/matcher"
	"github.com/spf13/cobra"
)

var matcherCmd = &cobra.Command{
	Use:   "matcher",
	Short: "Run matcher",
	Long:  "Run matcher",
	Run: func(cmd *cobra.Command, args []string) {
		m := matcher.NewMatcher()

		service := micro.NewService(
			micro.Name("github.com.rudeigerc.broker-gateway.matcher"),
			micro.BeforeStop(func() error {
				m.Stop()
				return nil
			}),
		)

		if err := service.Run(); err != nil {
			log.Fatal(err)
		}
	},
}
