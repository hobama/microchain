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
	// Response for Ping
	core.Ping: pingResp,

	// Response for Join.
	core.Join: joinResp,

	core.SyncNodes: syncNodesResp,
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
		rn := c.node.GetNodeByPublicKey(p.PublicKey)
		rn.Lastseen = int(time.Now().Unix())

		c.node.CheckAndAddNodeToRoutingTable(rn)
	}

	m.Conn.Write([]byte("pong"))
}

// Callback function for join request.
func joinResp(m core.IncommingMessage, c *client) {
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

// Run web server.
func (c *client) runWebServer(port int) {
	http.HandleFunc("/", c.indexHandler)

	c.logger.Error.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))
}

// Maintain routing table.
func (c *client) maintainRoutingTable() {
	for {
		// FIXME: We should be able to change this value.
		time.Sleep(5 * time.Second)

		_, nodes := c.node.GetNodesOfRoutingTable()

		p := core.NewPingMessage(c.node.Keypair.Public, c.node.IP+":"+strconv.Itoa(c.node.Port))
		pjson, err := p.MarshalJson()
		if err != nil {
			continue
		}

		for _, n := range nodes {

			err = c.node.Send(n.Address, pjson, func([]byte) error { return nil })

			if err != nil {
				// Delete if it doesn't respond for a long time.
				if int(time.Now().Unix())-n.Lastseen > 50 {
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
		// FIXME: We should be able to change this value.
		time.Sleep(7 * time.Second)

		nodes := c.collectNodesFromRoutingTable()

		p := core.NewSyncNodesMessage(nodes)

		pjson, err := p.MarshalJson()
		if err != nil {
			continue
		}

		for _, n := range c.node.RoutingTable {
			err = c.node.Send(n.Address, pjson, func(data []byte) error { return nil })
		}
	}
}

// Collect nodes from routing table.
func (c *client) collectNodesFromRoutingTable() []core.RemoteNode {
	var nodes []core.RemoteNode

	_, nodes = c.node.GetNodesOfRoutingTable()

	nodes = append(nodes, core.RemoteNode{
		PublicKey: c.node.Keypair.Public,
		Address:   c.node.IP + ":" + strconv.Itoa(c.node.Port),
		Lastseen:  int(time.Now().Unix()),
	})

	return nodes
}

func (c *client) indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is a website server by a Go HTTP server.")
}
