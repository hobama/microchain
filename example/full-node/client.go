package main

import (
	"fmt"
	"net"
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

	// initialize web server
	go c.runWebServer(webPort)

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
	if c.node.RoutingTable[string(p.PublicKey)] == nil {
		c.node.RoutingTable[string(p.PublicKey)] = &core.RemoteNode{
			PublicKey:  p.PublicKey,
			Address:    m.Conn.RemoteAddr().(*net.TCPAddr),
			Lastseen:   int(time.Now().Unix()),
			VerifiedBy: nil,
		}
	} else {
		// Update lastseen value.
		c.node.RoutingTable[string(p.PublicKey)].Lastseen = int(time.Now().Unix())
	}

	m.Conn.Write([]byte("pong"))
}

// Callback function for join request.
func joinResp(m core.IncommingMessage, c *client) {
}

// Run web server
func (c *client) runWebServer(port int) {
	http.HandleFunc("/", c.indexHandler)

	c.logger.Error.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))
}

func (c *client) indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is a website server by a Go HTTP server.")
}
