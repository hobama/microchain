package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/bosoncat/microchain/core"
)

// Options
var queryNodesOpt = regexp.MustCompile(`nodes`)
var pingNodeOpt = regexp.MustCompile(`ping`)
var joinNetworkOpt = regexp.MustCompile(`join`)
var sendTransactionOpt = regexp.MustCompile(`tran`)

// Common used regex
var addrRegex = regexp.MustCompile(`(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}:\d+`)

func (c *client) repl() {
	reader := bufio.NewReader(os.Stdin)

	for {
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if queryNodesOpt.MatchString(input) {
			// Query nodes in routing table.
			if b, msg := checkQueryNodesCommand(input); !b {
				c.terminal <- msg
				continue
			}

			c.terminal <- fmt.Sprintf("Currently, %d nodes in routing table.\n", len(c.node.RoutingTable))

			for k, n := range c.node.RoutingTable {
				c.terminal <- fmt.Sprintf("# %s Lastseen: %d Public key: %s\n", n.Address, n.Lastseen, k)
			}
		} else if pingNodeOpt.MatchString(input) {
			// Ping node.
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
						c.terminal <- err.Error() + "\n"
						continue
					}

					c.terminal <- fmt.Sprintf("%s responded and it's in good state ...\n", addr)
				}
			}()
		} else if joinNetworkOpt.MatchString(input) {
			// Join p2p network through node.
		} else if sendTransactionOpt.MatchString(input) {
			// Send transaction to given node.
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
	p := core.NewPingMessage(c.node.Keypair.Public, c.node.IP+":"+strconv.Itoa(c.node.Port))
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
