package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/alexander256/ascendex/internal/exchange"
	"github.com/sirupsen/logrus"
)

func main() {
	os.Exit(run())
}

func run() int {
	stopAppCh := make(chan os.Signal)
	signal.Notify(stopAppCh, syscall.SIGTERM, syscall.SIGINT)

	readChan := make(chan exchange.BestOrderBook)
	l := logrus.New()

	a := exchange.NewAscendex(l)
	if err := a.Connection(); err != nil {
		l.Error(err)
		return 1
	}
	defer a.Disconnect()

	go a.ReadMessagesFromChannel(readChan)
	if err := a.SubscribeToChannel("BTC/USDT"); err != nil {
		l.Error(err)
		return 1
	}
	var exitSignal os.Signal
LOOP:
	for {
		select {
		case exitSignal = <-stopAppCh:
			break LOOP
		case msg := <-readChan:
			fmt.Println(msg)
		}

	}

	if exitSignal == syscall.SIGINT || exitSignal == syscall.SIGTERM {
		return 0
	}
	return 1
}
