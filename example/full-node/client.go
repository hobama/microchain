package main

import (
	"github.com/bosoncat/microchain/core"
	"net"
	"time"
)

type client struct {
	node     *core.Node
	terminal chan string
}

// Generate new client.
func newClient(ip string, port int) (*client, error) {
	// new client
	c := &client{
		node:     core.NewNode(ip, port),
		terminal: make(chan string),
	}

	// initialize network.
	go c.node.Run()

	// initialize print loop.
	go c.printLoop()

	// process incomming message.
	go c.processIncommingMessage()

	return c, nil
}

// Callback functions.
var resp = map[byte]func(core.IncommingMessage, *client){
	core.Ping: func(m core.IncommingMessage, c *client) {
		var p core.PingMsg

		// TODO: Error handling
		err := p.UnmarshalJson(m.Content.Data)
		if err != nil {
			return
		}

		if c.node.RoutingTable[string(p.PublicKey)] == nil {
			c.node.RoutingTable[string(p.PublicKey)] = &core.RemoteNode{
				PublicKey:  p.PublicKey,
				Address:    m.Conn.RemoteAddr().(*net.TCPAddr),
				Lastseen:   int(time.Now().Unix()),
				VerifiedBy: nil,
			}
		}
	},

	core.Join: func(m core.IncommingMessage, c *client) {
	},
}

// Process incomming message
func (c *client) processIncommingMessage() {
	for m := range c.node.MessageChannel {
		// If callback function exists, we should process this message.
		if resp[m.Content.Type] != nil {
			resp[m.Content.Type](m, c)
		}
	}
}
