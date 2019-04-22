package core

import (
	"net"
)

type NodesMap map[string]*Node
type Node struct {
	Conn      *net.TCPConn // TCP connection to node
	Lastseen  int          // The seconds since last time seen this node
	PublicKey string       // Public key of this node
}

type Network struct {
	Nodes              NodesMap     // Contacts
	ConnectionsQueue   chan string  // Connections
	Listener           net.Listener // Listener
	Address            string       // TCP-4 address
	ConnectionCallBack chan *Node   // Connections callback
	BroadcastQueue     chan Message // Broadcast queue
	IncommingMessages  chan Message // Incomming messages
}

type Peer struct {
	*KeyPair // Public key and private Key
	*Network // Network
}

// TODO: func NewNode() Node {}

// Add node to peer's network.
func (p *Peer) AddNode(n Node) bool {
	pub := n.PublicKey

	if pub != string(p.KeyPair.Public) && p.Network.Nodes[pub] == nil {
		p.Network.Nodes[pub] = &n
		return true
	}

	return false
}

// Send message to specific node.
func (p *Peer) Send(n Node, m *Message) error {
	mBytes, err := m.MarshalBinary()
	if err != nil {
		return err
	}

	_, err = n.Conn.Write(mBytes)
	if err != nil {
		return err
	}

	return nil
}

// Broadcast messages to nodes.
func (p *Peer) BroadcastMessage(m *Message) error {
	return nil
}
