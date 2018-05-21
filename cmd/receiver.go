package cmd

import (
	"log"
	"os"
	"os/signal"
	"path"

	"github.com/quickfixgo/quickfix"
	"github.com/rudeigerc/broker-gateway/receiver"
	"github.com/spf13/cobra"
)

var receiverCmd = &cobra.Command{
	Use:   "receiver",
	Short: "Run receiver",
	Long:  "Run receiver",
	Run: func(cmd *cobra.Command, args []string) {
		cfgFileName := path.Join("config", "receiver.cfg")

		cfg, err := os.Open(cfgFileName)
		if err != nil {
			log.Printf("Error opening %v, %v\n", cfgFileName, err)
			return
		}

		appSettings, err := quickfix.ParseSettings(cfg)
		if err != nil {
			log.Println("Error reading cfg,", err)
			return
		}

		logFactory := quickfix.NewScreenLogFactory()

		r := receiver.NewReceiver()

		acceptor, err := quickfix.NewAcceptor(r, quickfix.NewMemoryStoreFactory(), appSettings, logFactory)
		if err != nil {
			log.Printf("Unable to create Acceptor: %s\n", err)
			return
		}

		err = acceptor.Start()
		if err != nil {
			log.Printf("Unable to start Acceptor: %s\n", err)
			return
		}

		interrupt := make(chan os.Signal)
		signal.Notify(interrupt, os.Interrupt, os.Kill)
		<-interrupt

		acceptor.Stop()
		r.Stop()
	},
}
