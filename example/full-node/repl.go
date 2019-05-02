package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/bosoncat/microchain/core"
)

// Options
var queryNodesOpt = regexp.MustCompile(`nodes`)
var pingNodeOpt = regexp.MustCompile(`ping`)

// Common used regex
var addrRegex = regexp.MustCompile(`(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}:\d+`)

func (c *client) repl() {
	reader := bufio.NewReader(os.Stdin)

	for {
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if queryNodesOpt.MatchString(input) {
			if b, msg := checkQueryNodesCommand(input); !b {
				c.terminal <- msg
				continue
			}

			for _, n := range c.node.RoutingTable {
				base58PK := core.Base58Encode(n.PublicKey[:32])
				c.terminal <- fmt.Sprintf("# %s Lastseen: %d Public key: %s\n", n.Address.String(), n.Lastseen, base58PK)
			}
		} else if pingNodeOpt.MatchString(input) {
			b, msg, addrs := checkPingNodeCommand(input)
			if !b {
				c.terminal <- msg
				continue
			}

			go func() {
				for _, addr := range addrs {
					c.terminal <- fmt.Sprintf("ping node %s ...\n", addr)

					err := c.pingNode(addr, c.pingNodeCallBack)

					if err != nil {
						c.terminal <- err.Error()
						continue
					}

					c.terminal <- fmt.Sprintf("%s responded and it's in good state ...\n", addr)
				}
			}()
		} else if input == "" {
			// Do nothing
		} else {
			c.terminal <- fmt.Sprintf("Unknown command: %s\n", input)
		}
	}
}

func checkQueryNodesCommand(s string) (bool, string) {
	if s != "nodes" {
		return false, fmt.Sprintf("Unknown command: %s, do you mean: nodes ?\n", s)
	}

	return true, ""
}

func checkPingNodeCommand(s string) (bool, string, []string) {
	if !strings.HasPrefix(s, "ping") {
		return false, fmt.Sprintf("Unknown command: %s, do you mean: ping ?\n", s), []string{}
	}

	// Remove `ping`
	s = strings.TrimSpace(s[4:])

	addrs := addrRegex.FindAllString(s, -1)

	if len(addrs) == 0 {
		return false, "Please check the address that you want to ping\n", []string{}
	}

	return true, "", addrs
}

func (c *client) pingNode(addr string, pingNodeCallback func([]byte) error) error {

	p := core.NewPingMessage(*c.node)
	pjson, err := p.MarshalJson()
	if err != nil {
		return err
	}

	err = c.node.Send(addr, pjson, pingNodeCallback)

	if err != nil {
		return err
	}

	return nil
}

func (c *client) pingNodeCallBack(data []byte) error {
	if string(data) != "pong" {
		return fmt.Errorf("Invalid response for ping")
	}

	return nil
}

func (c *client) printLoop() {
	for s := range c.terminal {
		fmt.Printf("%s", s)
	}
}
