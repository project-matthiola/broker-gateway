package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"path"

	"github.com/quickfixgo/quickfix"
	"github.com/rudeigerc/broker-gateway/receiver"
)

func main() {
	flag.Parse()

	cfgFileName := path.Join("config", "receiver.cfg")
	if flag.NArg() > 0 {
		cfgFileName = flag.Arg(0)
	}

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

	acceptor, err := quickfix.NewAcceptor(receiver.NewReceiver(), quickfix.NewMemoryStoreFactory(), appSettings, logFactory)
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
}
