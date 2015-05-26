package main

import (
	"fmt"
	"net/rpc"
	"time"

	"github.com/mikespook/sts/model"
	"gopkg.in/mgo.v2/bson"

	"github.com/peterh/liner"
)

type cmdFunc func(*rpc.Client, *liner.State, []string) error

var (
	cmds = map[string]cmdFunc{
		"sessions": sessions,
		"agents":   agents,
		"stat":     stat,
		"agent":    agent,
		"session":  session,
		"restart":  restart,
		"cutoff":   cutoff,
		"kickoff":  kickoff,
		"?":        help,
		"help":     help,
	}
	cmdNames = []string{"session", "sessions", "stat", "agent", "agents", "restart", "cutoff", "kickoff", "help/?"}
)

func help(client *rpc.Client, line *liner.State, argv []string) error {
	if len(argv) == 0 {
		fmt.Println("Following commands can be used:")
		for _, v := range cmdNames {
			fmt.Printf("\t%s\n", v)
		}
		return nil
	}
	switch argv[0] {
	case "session":
		fmt.Println("Display the detail information of a session.\nUsage: session [session-id]")
	case "agent":
		fmt.Println("Display the detail information of a agent.\nUsage: agent [agent-id]")
	case "cutoff":
		fmt.Println("Cut off an agent connection.\nUsage: cutoff [agent-id]")
	case "kickoff":
		fmt.Println("Kick of a session.\nUsage: kickoff [session-id]")
	default:
		fmt.Printf("Usage: %s\n", argv[0])
	}
	return nil
}

func sessions(client *rpc.Client, line *liner.State, argv []string) error {
	var m model.Sessions
	if err := client.Call("Stat.Sessions", "", &m); err != nil {
		return err
	}
	fmt.Printf("%24s\t%25s\t%20s\t%21s\t%s\n", "[Session ID]", "[Estblished Time]", "[User]", "[Address]", "[Num of Agents]")
	for k, v := range m.M {
		fmt.Printf("%24s\t%25s\t%20s\t%21s\t%d\n", k.Hex(), v.ETime.Format(time.RFC3339), v.User, v.RemoteAddr, len(v.Agents))
	}
	return nil
}

func agents(client *rpc.Client, line *liner.State, argv []string) error {
	var m model.Agents
	if err := client.Call("Stat.Agents", "", &m); err != nil {
		return err
	}
	fmt.Printf("%24s\t%24s\t%25s\t%20s\t%21s\n", "[Agent ID]", "[Session ID]", "[Estblished Time]", "[User]", "[Target]")
	for k, v := range m.M {
		fmt.Printf("%24s\t%24s\t%25s\t%20s\t%21s\n", k.Hex(), v.SessionId.Hex(), v.ETime.Format(time.RFC3339), v.User, v.RemoteAddr)
	}
	return nil
}

func stat(client *rpc.Client, line *liner.State, argv []string) error {
	var m model.Stat
	if err := client.Call("Stat.Stat", struct{}{}, &m); err != nil {
		return err
	}
	fmt.Printf("%25s\t%17s\t%17s\n", "[Estblished Time]", "[Num of Sessions]", "[Num of Agents]")
	fmt.Printf("%25s\t%17d\t%17d\n", m.ETime.Format(time.RFC3339), m.Sessions, m.Agents)
	return nil
}

func session(client *rpc.Client, line *liner.State, argv []string) error {
	if argv == nil || len(argv) < 1 {
		fmt.Println("session [id]")
		return nil
	}
	if !bson.IsObjectIdHex(argv[0]) {
		fmt.Println("session [id]")
		return nil
	}
	id := bson.ObjectIdHex(argv[0])
	var m model.Session
	if err := client.Call("Stat.Session", id, &m); err != nil {
		return err
	}
	fmt.Printf("ID: %s\t\n", m.Id.Hex())
	fmt.Printf("Estblished Time: %s\t\n", m.ETime.Format(time.RFC3339))
	fmt.Printf("User: %s\t\n", m.User)
	fmt.Printf("Address: %s\t\n", m.RemoteAddr)
	fmt.Printf("Client Version: %s\t\n", m.ClientVersion)
	fmt.Println("Agents")
	fmt.Printf("\t%24s\t%25s\t%21s\n", "[ID]", "[Estblished Time]", "[Target]")
	for k, v := range m.Agents {
		fmt.Printf("\t%24s\t%25s\t%21s\n", k.Hex(), v.ETime.Format(time.RFC3339), v.RemoteAddr)

	}
	return nil
}

func agent(client *rpc.Client, line *liner.State, argv []string) error {
	if argv == nil || len(argv) < 1 {
		fmt.Println("agent [id]")
		return nil
	}
	if !bson.IsObjectIdHex(argv[0]) {
		fmt.Println("agent [id]")
		return nil
	}
	id := bson.ObjectIdHex(argv[0])
	var m model.Agent
	if err := client.Call("Stat.Agent", id, &m); err != nil {
		return err
	}
	fmt.Printf("%24s\t%24s\t%25s\t%20s\t%21s\n", "[Agent ID]", "[Session ID]", "[Estblished Time]", "[User]", "[Target]")
	fmt.Printf("%24s\t%24s\t%25s\t%20s\t%21s\n", m.Id.Hex(), m.SessionId.Hex(), m.ETime.Format(time.RFC3339), m.User, m.RemoteAddr)
	return nil
}

func restart(client *rpc.Client, line *liner.State, argv []string) error {
	if err := client.Call("Ctrl.Restart", struct{}{}, &struct{}{}); err != nil {
		return err
	}
	return stat(client, line, argv)
}

func cutoff(client *rpc.Client, line *liner.State, argv []string) error {
	if !bson.IsObjectIdHex(argv[0]) {
		fmt.Println("cutoff [id]")
		return nil
	}
	id := bson.ObjectIdHex(argv[0])
	if err := client.Call("Ctrl.Cutoff", id, &struct{}{}); err != nil {
		return err
	}
	return agents(client, line, argv)
}

func kickoff(client *rpc.Client, line *liner.State, argv []string) error {
	if !bson.IsObjectIdHex(argv[0]) {
		fmt.Println("kickoff [id]")
		return nil
	}
	id := bson.ObjectIdHex(argv[0])
	if err := client.Call("Ctrl.Kickoff", id, &struct{}{}); err != nil {
		return err
	}
	return sessions(client, line, argv)
}
