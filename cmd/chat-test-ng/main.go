package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"

	"github.com/boletia/chat-test-ng/pkg/bot"
	log "github.com/sirupsen/logrus"
)

var wg sync.WaitGroup

func main() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

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

	flag.IntVar(&numBots, "bots", bot.DefaultNumBots, "-bots=<numbots>")
	flag.BoolVar(&cnf.SendMessages, "sendmessages", bot.DefaultSendMessages, "-sendmessages=<true|false>")
	flag.BoolVar(&cnf.WithGossiper, "gossiper", bot.DefaultWithGossiper, "-gossiper=<true|false>")
	flag.StringVar(&cnf.SudDomain, "subdomain", bot.DefaultSubdomain, "-subdomain=<subdomain>")
	flag.IntVar(&cnf.NumMessages, "messages", bot.DefaultNumMessages, "-messages=<num_messages>")
	flag.IntVar(&cnf.MinDelay, "mindelay", bot.DefaultMinDelay, "-mindelay=<delay_in_sec>")
	flag.IntVar(&cnf.MaxDelay, "maxdelay", bot.DefaultMaxDelay, "-maxdelay=<delay_in_sec>")
	flag.StringVar(&cnf.URL, "endpoint", bot.DefautlEndPoint, "-endpoint=<endpoint>")
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
		bot := bot.New(cnf, quit)
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
