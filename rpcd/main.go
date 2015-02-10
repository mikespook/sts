package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/mikespook/golib/log"
	"github.com/mikespook/golib/signal"
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
	cfg, err := LoadConfig(config)
	if err != nil {
		log.Errorf("Start: %s", err)
		return
	}
	if err := log.Init(cfg.Log.File, log.StrToLevel(cfg.Log.Level), 0); err != nil {
		log.Errorf("Start: %s", err)
		return
	}
	log.Messagef("Starting: pid=%d, addr=%s, pwd=%s",
		os.Getpid(), cfg.Addr, cfg.Pwd)
	srv, err := NewRPC(cfg)
	if err != nil {
		log.Errorf("Start: %s", err)
		return
	}
	go func() {
		if err := srv.Serve(); err != nil {
			log.Errorf("Serve: %s", err)
		}
		signal.Send(os.Getpid(), os.Interrupt)
	}()
	signal.Bind(os.Interrupt, func() uint {
		if err := srv.Close(); err != nil {
			log.Errorf("Stopping: %s", err)
		}
		return signal.BreakExit
	})
	signal.Wait()
	log.Messagef("Stopped")
}
