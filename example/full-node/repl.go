package main

import (
	"bufio"
	"fmt"
	"github.com/bosoncat/microchain/core"
	"os"
	"strings"
)

func (c *client) repl() {
	reader := bufio.NewReader(os.Stdin)

	for {
		c.terminal <- "> "
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "nodes":
			for _, n := range c.node.RoutingTable {
				base58PK := core.Base58Encode(n.PublicKey[:32])
				c.terminal <- fmt.Sprintf("# %s Lastseen: %d Public key: %s\n", n.Address.String(), n.Lastseen, base58PK)
			}
		}
	}
}

func (c *client) printLoop() {
	for s := range c.terminal {
		fmt.Printf("%s", s)
	}
}
