package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/mikespook/golib/log"
	"github.com/mikespook/golib/signal"
	"github.com/mikespook/sts"
)

var config string

func init() {
	flag.Usage = func() {
		fmt.Println(intro)
	}
	flag.StringVar(&config, "config", "", "Path to the configuration file")
	flag.Parse()
}

func main() {
	if config == "" {
		flag.PrintDefaults()
		return
	}
	log.Messagef("STS: Secure Tunnel Server: Starting")
	cfg, err := sts.LoadConfig(config)
	if err != nil {
		log.Errorf("Load config: %s", err)
		return
	}
	if err := log.Init(cfg.Log.File, log.StrToLevel(cfg.Log.Level), 0); err != nil {
		log.Errorf("Init log: %s", err)
		return
	}
	tunnel := sts.New(cfg)
	go func() {
		if err := tunnel.Serve(); err != nil {
			log.Errorf("Serve: %s", err)
		}
		signal.Send(os.Getpid(), os.Interrupt)
	}()
	signal.Bind(os.Interrupt, func() uint {
		tunnel.Close()
		return signal.BreakExit
	})
	signal.Wait()
	log.Messagef("STS: Secure Tunnel Server: Shutdown")
}
