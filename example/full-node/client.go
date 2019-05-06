package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/bosoncat/microchain/core"
)

type client struct {
	node     *core.Node
	terminal chan string
	logger   *core.Logger
}

// Callback functions.
var resp = map[byte]func(core.IncommingMessage, *client){
	core.Ping: pingResp,

	core.SyncNodes: syncNodesResp,

	core.SyncTransactions: syncTransactionsResp,

	core.SendTransaction: sendTransactionResp,
}

// Generate new client.
func newClient(ip string, nodePort, webPort int, l *core.Logger) (*client, error) {
	// new client
	n, err := core.NewNode(ip, nodePort)
	if err != nil {
		return nil, err
	}

	c := &client{
		node:     n,
		terminal: make(chan string),
		logger:   l,
	}

	// initialize network.
	err = c.node.Run()
	if err != nil {
		return nil, err
	}

	// initialize print loop.
	go c.printLoop()

	// process incomming message.
	go c.processIncommingMessage()

	// initialize web server.
	go c.runWebServer(webPort)

	// maintain routing table.
	go c.maintainRoutingTable()

	// broadcast nodes.
	go c.broadcastKnownNodes()

	// broadcast transactions
	go c.broadcastTransactionsPool()

	return c, nil
}

// Process incomming message
func (c *client) processIncommingMessage() {
	for m := range c.node.MessageChannel {
		// If callback function exists, we should process this message.
		if resp[m.Content.Type] != nil {
			resp[m.Content.Type](m, c)
		}

		m.Conn.Close()
	}
}

// Callback function for ping request.
func pingResp(m core.IncommingMessage, c *client) {
	var p core.PingData

	err := p.UnmarshalJson(m.Content.Data)
	if err != nil {
		c.logger.Error.Println(err)
		return
	}

	// Test if node is existed in routing table, if not, add it into routing table.
	if !c.node.IsInRoutingTable(p.PublicKey) {
		c.node.CheckAndAddNodeToRoutingTable(core.RemoteNode{
			PublicKey:  p.PublicKey,
			Address:    p.Address,
			Lastseen:   int(time.Now().Unix()),
			VerifiedBy: nil,
		})
	} else {
		// Update lastseen value.
		b, rn := c.node.GetNodeByPublicKey(p.PublicKey)

		if !b {
			return
		}

		rn.Lastseen = int(time.Now().Unix())

		c.node.UpdateNodeForGivenPublicKey(rn.PublicKey, rn)
	}

	m.Conn.Write([]byte("pong"))
}

// Callback function for sync nodes.
func syncNodesResp(m core.IncommingMessage, c *client) {
	var sn core.SyncNodesData

	err := sn.UnmarshalJson(m.Content.Data)
	if err != nil {
		c.logger.Error.Println(err)
		return
	}

	// Test if node if existed in routing table, if not, add it into routing table.
	for _, n := range sn.Nodes {
		c.node.CheckAndAddNodeToRoutingTable(n)
	}

	m.Conn.Write([]byte("pong"))
}

// Callback function for sync transactions.
func syncTransactionsResp(m core.IncommingMessage, c *client) {
	var st core.SyncTransactionsData

	err := st.UnmarshalJson(m.Content.Data)
	if err != nil {
		c.logger.Error.Println(err)
		return
	}

	for _, tr := range st.Transactions {
		// TODO: check transactions.
		c.node.CheckAndAddTransactionToPool(tr)
	}

	m.Conn.Write([]byte("pong"))
}

// Callback for send transaction.
func sendTransactionResp(m core.IncommingMessage, c *client) {
	var st core.SendTransactionData

	err := st.UnmarshalJson(m.Content.Data)
	if err != nil {
		c.logger.Error.Println(err)
		return
	}

	t := st.Transaction

	if c.node.VerifyTransaction(t) {
		// Add it into transactions pool
		c.node.CheckAndAddTransactionToPool(t)
	}
}

// Run web server.
func (c *client) runWebServer(port int) {
	http.HandleFunc("/", c.indexHandler)

	c.logger.Error.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))
}

// Maintain routing table.
func (c *client) maintainRoutingTable() {
	for {
		time.Sleep(pingPeriod * time.Second)

		_, nodes := c.node.GetNodesOfRoutingTable()

		p := core.NewPingMessage(c.node.PublicKey(), c.node.Addr())
		pjson, err := p.MarshalJson()
		if err != nil {
			continue
		}

		for _, n := range nodes {
			err = c.node.Send(n.Address, pjson, func([]byte) error { return nil })

			if err != nil {
				// Delete if it doesn't respond for a long time.
				if int(time.Now().Unix())-n.Lastseen > invalidPeriod {
					c.node.RemoveNodeByPublicKey(n.PublicKey)
				}

				continue
			}
		}
	}
}

// Broadcast known nodes.
func (c *client) broadcastKnownNodes() {
	for {
		time.Sleep(broadcastRoutingTablePeriod * time.Second)

		nodes := c.collectNodesFromRoutingTable()

		m := core.NewSyncNodesMessage(nodes)

		mjson, err := m.MarshalJson()
		if err != nil {
			return
		}

		c.node.Broadcast(mjson, func([]byte) error { return nil })
	}
}

// Broadcast transactions pool.
func (c *client) broadcastTransactionsPool() {
	for {
		time.Sleep(broadcastTransactionsPoolPeriod * time.Second)

		_, ts := c.node.GetTransactionsOfPool()

		m := core.NewSyncTransactionsMessage(ts)

		mjson, err := m.MarshalJson()
		if err != nil {
			return
		}

		c.node.Broadcast(mjson, func([]byte) error { return nil })
	}
}

// Collect nodes from routing table.
func (c *client) collectNodesFromRoutingTable() []core.RemoteNode {
	var nodes []core.RemoteNode

	_, nodes = c.node.GetNodesOfRoutingTable()

	// Add client itself to []Node.
	nodes = append(nodes, core.RemoteNode{
		PublicKey: c.node.Keypair.Public,
		Address:   c.node.Addr(),
		Lastseen:  int(time.Now().Unix()),
	})

	return nodes
}

func (c *client) indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is a website server by a Go HTTP server.")
}
