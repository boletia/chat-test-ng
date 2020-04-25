package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"

	"github.com/boletia/chat-test-ng/pkg/bot"
	"github.com/boletia/chat-test-ng/pkg/wsocket"
	log "github.com/sirupsen/logrus"
)

var wg sync.WaitGroup

func main() {
	bots, conf := readConfig()
	quit := launch(bots, conf)

	go stop(bots, quit)

	log.WithFields(log.Fields{"bots": bots}).Info("waiting for bots")
	wg.Wait()

	log.Info("end")
}

func readConfig() (int, bot.Conf) {
	cnf := bot.Conf{}
	numBots := 0

	flag.IntVar(&numBots, "bots", 0, "-bots=<numbots>")
	flag.BoolVar(&cnf.SendMessages, "sendmessages", false, "-sendmessages=<true|false>")
	flag.BoolVar(&cnf.WithGossiper, "gossiper", false, "-gossiper=<true|false>")
	flag.StringVar(&cnf.SudDomain, "subdomain", "el-show-de-producto-online", "-subdomain=<subdomain>")
	flag.IntVar(&cnf.NumMessages, "messages", 0, "-messages=<num_messages>")
	flag.Int64Var(&cnf.MinDelay, "mindelay", 10, "-mindelay=<delay_in_sec>")
	flag.Int64Var(&cnf.MaxDelay, "maxdelay", 30, "-maxdelay=<delay_in_sec>")
	flag.StringVar(&cnf.URL, "endpoint", "wss://6vfdhz6o24.execute-api.us-east-1.amazonaws.com/beta", "-endpoint=<endpoint>")
	flag.Parse()

	log.WithFields(log.Fields{
		"bots":         numBots,
		"sendmessages": cnf.SendMessages,
		"gossiper":     cnf.WithGossiper,
		"subdomain":    cnf.SudDomain,
		"messages":     cnf.NumMessages,
		"mindelay":     cnf.MinDelay,
		"maxdelay":     cnf.MaxDelay,
		"endpoint":     cnf.URL,
	}).Info("read params")

	return numBots, cnf
}

func launch(bots int, cnf bot.Conf) chan bool {
	quit := make(chan bool)

	for i := 0; i < bots; i++ {
		cnf.NickName = fmt.Sprintf("bot-%d", i)
		sock, err := wsocket.New(cnf.URL)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
				"bot":   cnf.NickName,
			}).Error("unable to connect")
			continue
		}

		bot := bot.New(cnf, sock, quit)
		wg.Add(1)
		go bot.Start(&wg)
	}

	return quit
}

func stop(bots int, quit chan bool) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	select {
	case <-interrupt:
		log.Infof("sending quit message to channels")
		for i := 0; i < bots; i++ {
			go func() {
				quit <- true
			}()
		}
	}
}
