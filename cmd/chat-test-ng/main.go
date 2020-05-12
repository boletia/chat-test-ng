package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/boletia/chat-test-ng/pkg/bot"
	"github.com/boletia/chat-test-ng/pkg/dynamo"
	log "github.com/sirupsen/logrus"
)

var wg sync.WaitGroup
var totalCalls []int

func main() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	bots, conf := readConfig()
	quit := launch(bots, conf)

	go stop(bots, quit)

	log.WithFields(log.Fields{"bots": bots}).Info("waiting for bots")
	wg.Wait()

	greatTotal := 0
	for _, call := range totalCalls {
		greatTotal += call
	}

	log.WithFields(log.Fields{
		"socket-operations": greatTotal,
	}).Info("end")
}

func readConfig() (int, bot.Conf) {
	cnf := bot.Conf{}
	numBots := 0

	flag.IntVar(&numBots, "bots", bot.DefaultNumBots, "-bots=<numbots>")
	flag.BoolVar(&cnf.SendMessages, "sendmessages", bot.DefaultSendMessages, "-sendmessages=<true|false>")
	flag.BoolVar(&cnf.WithGossiper, "gossiper", bot.DefaultWithGossiper, "-gossiper=<true|false>")
	flag.StringVar(&cnf.SudDomain, "subdomain", bot.DefaultSubdomain, "-subdomain=<subdomain>")
	flag.IntVar(&cnf.NumMessages, "messages", bot.DefaultNumMessages, "-messages=<num_messages>")
	flag.Int64Var(&cnf.MinDelay, "mindelay", bot.DefaultMinDelay, "-mindelay=<delay_in_msec>")
	flag.Int64Var(&cnf.MaxDelay, "maxdelay", bot.DefaultMaxDelay, "-maxdelay=<delay_in_msec>")
	flag.StringVar(&cnf.URL, "endpoint", bot.DefautlEndPoint, "-endpoint=<endpoint>")
	flag.IntVar(&cnf.Ramping, "ramping", bot.DefaultRamping, "-ramping=<bots/sec>")
	flag.BoolVar(&cnf.OnlyError, "onlyerrors", false, "-onlyerrors=<true|false>")
	flag.BoolVar(&cnf.Sent2Dynamo, "send2dynamo", false, "-send2dynamo=<true|false>")
	flag.Parse()

	if (cnf.MaxDelay - cnf.MinDelay) <= 0 {
		log.WithFields(log.Fields{
			"min": cnf.MinDelay,
			"max": cnf.MaxDelay,
		}).Fatal("bad delay numbers")
	}

	log.WithFields(log.Fields{
		"bots":         numBots,
		"sendmessages": cnf.SendMessages,
		"gossiper":     cnf.WithGossiper,
		"subdomain":    cnf.SudDomain,
		"messages":     cnf.NumMessages,
		"mindelay":     cnf.MinDelay,
		"maxdelay":     cnf.MaxDelay,
		"endpoint":     cnf.URL,
		"onlyerrors":   cnf.OnlyError,
		"send2dynamo":  cnf.Sent2Dynamo,
	}).Info("read params")

	return numBots, cnf
}

func launch(bots int, cnf bot.Conf) chan bool {
	var dyDB dynamo.DyDB

	quit := make(chan bool)
	totalCalls = make([]int, bots)

	if cnf.Sent2Dynamo {
		dyDB = dynamo.New()
	}

	for i := 0; i < bots; i++ {
		cnf.NickName = fmt.Sprintf("bot-%d", i)
		bot := bot.New(cnf, quit)
		bot.AddDynamo(dyDB)
		wg.Add(1)
		go bot.Start(&wg, &totalCalls[i])
		time.Sleep((time.Second / time.Duration(cnf.Ramping)))
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
