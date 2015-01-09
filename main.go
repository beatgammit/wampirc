package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	irc "github.com/fluffle/goirc/client"
	"gopkg.in/jcelliott/turnpike.v2"
)

var (
	ircAddr, wampAddr string
	ircPass           string
	channel           string
	wampRealm         string
	nick              string
)

func init() {
	flag.StringVar(&ircAddr, "ircaddr", "localhost:6667", "address of IRC server")
	flag.StringVar(&wampAddr, "wampaddr", "localhost:8000", "address of WAMP server")
	flag.StringVar(&wampRealm, "wamprealm", "turnpike.examples", "IRC realm to connect to")
	flag.StringVar(&nick, "nick", "nick", "nickname of the irc bot")
	flag.StringVar(&ircPass, "pass", "", "password to use when connecting to the irc server")
	flag.StringVar(&channel, "channel", "#channel", "irc channel to connect to")
	flag.Parse()
}

type client struct {
	quit chan bool
	tp   *turnpike.Client
	conn *irc.Conn
}

func New() (*client, error) {
	tp, err := turnpike.NewWebsocketClient(turnpike.JSON, "ws://"+wampAddr)
	if err != nil {
		return nil, err
	} else if _, err = tp.JoinRealm("turnpike.examples", turnpike.ALLROLES, nil); err != nil {
		return nil, err
	}

	return &client{quit: make(chan bool), tp: tp}, nil
}

func (c *client) handleWampMsg(args []interface{}, kwargs map[string]interface{}) {
	if len(args) == 2 {
		if from, ok := args[0].(string); !ok {
			log.Println("First argument not a string:", args[0])
		} else if msg, ok := args[1].(string); !ok {
			log.Println("Second argument not a string:", args[1])
		} else {
			log.Printf("%s: %s", from, msg)
			c.conn.Privmsg(channel, from+" -> "+msg)
		}
	}
}

func (c *client) handleIrcMsg(conn *irc.Conn, line *irc.Line) {
	c.tp.Publish("chat", []interface{}{line.Nick, strings.Join(line.Args, " ")}, nil)
}

func (c *client) connected(conn *irc.Conn, line *irc.Line) {
	fmt.Println("Joining channel")
	c.tp.Subscribe("chat", c.handleWampMsg)

	c.conn = conn

	conn.Join(channel)
	conn.Privmsg(channel, "I'm a hippopotamus!")
}

func (c *client) disconnected(conn *irc.Conn, line *irc.Line) {
	c.quit <- true
}

func (c *client) run() {
	<-c.quit
}

func main() {
	// Creating a simple IRC client is simple.
	cl, err := New()
	if err != nil {
		log.Fatalln("Error connecting to WAMP:", err)
	}
	c := irc.SimpleClient(nick)

	// Add handlers to do things here!
	// e.g. join a channel on connect.
	c.HandleFunc("connected", cl.connected)
	// And a signal on disconnect
	c.HandleFunc("disconnected", cl.disconnected)
	c.HandleFunc("privmsg", cl.handleIrcMsg)

	if err := c.ConnectTo(ircAddr, ircPass); err != nil {
		log.Fatalln("Connection error:", err)
	}

	cl.run()
}
