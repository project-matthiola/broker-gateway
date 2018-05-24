package cmd

import (
	"log"
	"os"
	"path"

	"github.com/micro/go-micro"
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

		service := micro.NewService(
			micro.Name("github.com.rudeigerc.broker-gateway.receiver"),
			micro.BeforeStop(func() error {
				acceptor.Stop()
				r.Stop()
				return nil
			}),
		)

		if err := service.Run(); err != nil {
			log.Fatal(err)
		}

	},
}
