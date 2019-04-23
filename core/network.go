package core

import (
	"io"
	"net"
	"time"
)

type NodesMap map[string]*Node
type Node struct {
	Conn      *net.TCPConn // Should use generic Conn, so that we could use various conn type
	Lastseen  int          // The seconds since last time seen this node
	PublicKey string       // Public key of this node
	Address   string       // TCP-4 Address
}

type Network struct {
	Nodes              NodesMap         // Contacts
	ConnectionsQueue   chan string      // Connections
	Listener           *net.TCPListener // Listener
	Address            string           // TCP-4 address
	ConnectionCallBack chan *Node       // Connections callback
	BroadcastQueue     chan Message     // Broadcast queue
	IncommingMessages  chan Message     // Incomming messages
}

type Peer struct {
	*KeyPair    // Public key and private Key
	*Network    // Network
	*Blockchain // Blockchain
}

// TODO: func NewNode() Node {}

// Use peer as node.
// Though this function is not so commonly used, it's still useful for testing.
func (p Peer) AsNode() Node {
	return Node{nil, 0, string(p.KeyPair.Public), p.Network.Address}
}

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
func (p *Peer) Send(n Node, m Message) (int, error) {
	mBytes, err := m.MarshalBinary()
	if err != nil {
		return 0, err
	}

	bufLen, err := (*n.Conn).Write(mBytes)
	if err != nil {
		return 0, err
	}

	return bufLen, nil
}

// Receive Message from connected nodes.
func (p *Peer) Recv(n Node) (Message, error) {
	buf := make([]byte, 1024)

	bufLen, err := (*p.Nodes[n.PublicKey].Conn).Read(buf)
	if err != nil {
		if err != io.EOF {
			return Message{}, err
		}
	}

	m := new(Message)

	err = m.UnmarshalBinary(buf[:bufLen])
	if err != nil {
		return Message{}, err
	}

	return *m, nil
}

// Ping node.
func (p *Peer) Ping(n Node) bool {
	// TODO: Implement ping message.
	return true
}

// Listen.
func (p *Peer) Listen(address string, errch chan error) (chan Node, error) {
	cb := make(chan Node)

	addr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return nil, err
	}

	p.Network.Listener, err = net.ListenTCP("tcp4", addr)
	if err != nil {
		return nil, err
	}

	go func(l *net.TCPListener) {
		for {
			conn, err := l.AcceptTCP()
			if err != nil {
				if err != io.EOF {
					errch <- err
				}
			}

			cb <- Node{conn, int(time.Now().Unix()), "", conn.RemoteAddr().String()}
		}
	}(p.Network.Listener)

	return cb, nil
}

// Parse tcp message.

// Broadcast messages to nodes.
func (p *Peer) BroadcastMessage(m Message, errch chan error) {
	for _, n := range p.Nodes {
		// We cannot interupt this process, if there is a error, we should
		// emit this error in error channel.
		_, err := p.Send(*n, m)
		if err != nil {
			errch <- err
		}
	}
}
