package main

import (
	"net"
	"time"

	"github.com/bosoncat/microchain/core"
)

type client struct {
	node     *core.Node
	terminal chan string
}

// Generate new client.
func newClient(ip string, port int, l *core.Logger) (*client, error) {
	// new client
	n, err := core.NewNode(ip, port)
	if err != nil {
		return nil, err
	}

	c := &client{
		node:     n,
		terminal: make(chan string),
	}

	// initialize network.
	err = c.node.Run()
	if err != nil {
		return nil, err
	}

	// initialize print loop.
	go c.printLoop()

	// process incomming message.
	go c.processIncommingMessage(l)

	return c, nil
}

// Callback functions.
var resp = map[byte]func(core.IncommingMessage, *client, *core.Logger){
	// Response for Ping
	core.Ping: func(m core.IncommingMessage, c *client, l *core.Logger) {
		var p core.PingMsg

		defer m.Conn.Close()

		err := p.UnmarshalJson(m.Content.Data)
		if err != nil {
			l.Error.Println(err)
			return
		}

		// Test if node is existed in routing table, if not, add it into routing table.
		if c.node.RoutingTable[string(p.PublicKey)] == nil {
			c.node.RoutingTable[string(p.PublicKey)] = &core.RemoteNode{
				PublicKey:  p.PublicKey,
				Address:    m.Conn.RemoteAddr().(*net.TCPAddr),
				Lastseen:   int(time.Now().Unix()),
				VerifiedBy: nil,
			}
		}

		m.Conn.Write([]byte("pong"))
	},

	// Response for Join.
	core.Join: func(m core.IncommingMessage, c *client, l *core.Logger) {
	},
}

// Process incomming message
func (c *client) processIncommingMessage(l *core.Logger) {
	for m := range c.node.MessageChannel {
		// If callback function exists, we should process this message.
		if resp[m.Content.Type] != nil {
			resp[m.Content.Type](m, c, l)
		}
	}
}
