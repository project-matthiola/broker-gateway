package cmd

import (
	"os"
	"os/signal"

	"github.com/rudeigerc/broker-gateway/matcher"
	"github.com/spf13/cobra"
)

var matcherCmd = &cobra.Command{
	Use:   "matcher",
	Short: "Run matcher",
	Long:  "Run matcher",
	Run: func(cmd *cobra.Command, args []string) {
		m := matcher.NewMatcher()

		interrupt := make(chan os.Signal)
		signal.Notify(interrupt, os.Interrupt, os.Kill)
		<-interrupt

		m.Stop()
	},
}
