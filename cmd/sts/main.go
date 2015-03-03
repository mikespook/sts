package main

import (
	"flag"
	"fmt"
	"io"
	"net/rpc"
	"os"
	"strings"

	"github.com/mikespook/golib/log"

	"github.com/peterh/liner"
)

const (
	historyFile = ".sts_history"
)

var (
	addr string
)

func init() {
	flag.StringVar(&addr, "addr", "", "PRC address of STS")
	flag.Parse()
}

func main() {
	line := liner.NewLiner()
	defer line.Close()

	line.SetCompleter(func(line string) (c []string) {
		for _, n := range cmdNames {
			if strings.HasPrefix(n, strings.ToLower(line)) {
				c = append(c, n)
			}
		}
		return
	})
	hf := os.Getenv("HOME") + "/" + historyFile
	if f, err := os.OpenFile(hf, os.O_RDONLY|os.O_CREATE, 0644); err != nil {
		log.Error(err)
		return
	} else {
		line.ReadHistory(f)
		f.Close()
	}
	defer func() {
		if f, err := os.OpenFile(hf, os.O_WRONLY, 0644); err != nil {
			log.Errorf("Error writing history file: %s", err)
		} else {
			line.WriteHistory(f)
			f.Close()
		}
	}()

	client, err := rpc.Dial("tcp", addr)
	if err != nil {
		log.Error(err)
		return
	}
	for {
		if cmd, err := line.Prompt(addr + "$ "); err != nil {
			if err != io.EOF {
				log.Errorf("Error reading line: %s", err)
			}
			fmt.Println("")
			return
		} else {
			line.AppendHistory(cmd)
			if cmd == "quit" {
				return
			}
			argv := strings.Split(cmd, " ")
			cmd = argv[0]
			if f, ok := cmds[cmd]; ok {
				if err := f(client, line, argv[1:]); err != nil {
					log.Error(err)
					return
				}
			} else if cmd != "" {
				fmt.Printf("%s: command not found\n", cmd)
			}
		}
	}
}
