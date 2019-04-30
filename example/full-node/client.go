package main

import (
	"github.com/bosoncat/microchain/core"
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

	// initialize network
	go c.node.Run()

	// initialize print loop
	go c.printLoop()

	go func() {
		for p := range c.node.IncommingPacket {
			c.terminal <- string(p.Content)
		}
	}()

	return c, nil
}
