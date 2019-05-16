package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/vgxbj/microchain/core"
)

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

			_, ns := c.node.GetNodesOfRoutingTable()

			for _, n := range ns {
				c.terminal <- "---\n"
				c.terminal <- fmt.Sprintf("ID\t: %s\nAddress\t: %s\nLastseen: %d\n", n.PK(), n.Addr(), n.Lastseen)
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

					err := c.pingNode(addr)

					if err != nil {
						c.terminal <- err.Error() + "\n"
						continue
					}

					c.terminal <- fmt.Sprintf("%s responded and it's in good state ...\n", addr)
				}
			}()
		} else if joinNetworkOpt.MatchString(input) {
			// Join p2p network through node.
			// TODO: I don't want to implement this block now. Because I have very limited time to submit my final project.
		} else if queryTransactionsOpt.MatchString(input) {
			b, msg := checkQueryTransactionsCommand(input)
			if !b {
				c.terminal <- msg
				continue
			}

			n, ts := c.node.GetTransactionsOfPool()

			c.terminal <- fmt.Sprintf("Currently, there are %d transactions\n", n)

			for _, t := range ts {
				c.terminal <- "---\n"
				c.terminal <- fmt.Sprintf("ID\t : %s\nFrom\t : %s\nTo\t : %s\nData\t : %s\nTimestamp: %d\n", core.Base58Encode(t.ID()),
					core.Base58Encode(t.RequesterPK()),
					core.Base58Encode(t.RequesteePK()),
					string(t.Meta), t.Timestamp())
			}
		} else if sendTransactionOpt.MatchString(input) {
			// Send transaction to given node.
			b, msg, id, data := checkSendTransactionCommand(input)
			if !b {
				c.terminal <- msg
				continue
			}

			// TODO: c.sendTransaction(id, data)
			if c.node.PrevTransaction() == nil {
				c.terminal <- "Please generate genesis transaction first\n"
				continue
			}

			t := c.newPendingTransaction(id, data)

			c.sendPendingTransaction(t)
		} else if genesisOpt.MatchString(input) {
			// Generate genesis transaction.
			b, msg, data := checkGenesisCommand(input)
			if !b {
				c.terminal <- msg
				continue
			}

			if c.node.PrevTransaction() != nil {
				c.terminal <- "You cannot generate two genesis transaction\n"
				continue
			}

			c.broadcastGenesisTransaction(data)
		} else if queryPendingJobsOpt.MatchString(input) {
			// Query pending jobs.
			b, msg := checkQueryPendingCommand(input)
			if !b {
				c.terminal <- msg
				continue
			}

			// TODO:
			_, ts := c.node.GetPendingTransactions()

			for _, t := range ts {
				c.terminal <- "---\n"
				c.terminal <- fmt.Sprintf("ID %s\nFrom: %s\nTo: %s\nData: %s\n", core.Base58Encode(t.ID()),
					core.Base58Encode(t.RequesterPK()),
					core.Base58Encode(t.RequesteePK()),
					string(t.Meta))
			}

		} else if confirmReqOpt.MatchString(input) {
			b, msg, id := checkConfirmCommand(input)
			if !b {
				c.terminal <- msg
				continue
			}

			// TODO:
			c.confirmPendingTransaction(id)
		} else if input == "" {
			// Do nothing, intended leaving blank.
		} else {
			c.terminal <- fmt.Sprintf("Unknown command: %s\n", input)
		}
	}
}
