package core

import (
	"net"
	"strconv"
)

// Represent other nodes.
type RemoteNode struct {
	PublicKey  []byte        // Public key
	Address    *net.TCPAddr  // Address
	Lastseen   int           // The unix time of seeing this node last time
	VerifiedBy []*RemoteNode // Nodes that verify this node
}

// Received packect.
type Packet struct {
	Content []byte       // Raw bytes
	Conn    *net.TCPConn // TCP connection
}

type IncommingMessage struct {
	Content Message      // Message
	Conn    *net.TCPConn // TCP connection
}

// Represent ourselves.
type Node struct {
	Keypair        *KeyPair               // Key pair
	IP             string                 // IP address
	Port           int                    // Port
	RoutingTable   map[string]*RemoteNode // Routing table (public key, node)
	Listerner      *net.TCPListener       // TCP listener
	MessageChannel chan IncommingMessage  // Incomming message
}

func NewNode(ip string, port int) *Node {
	kp, _ := NewECDSAKeyPair()

	return &Node{
		Keypair:        kp,
		IP:             ip,
		Port:           port,
		RoutingTable:   make(map[string]*RemoteNode),
		Listerner:      new(net.TCPListener),
		MessageChannel: make(chan IncommingMessage),
	}
}

// Run a simple TCP server.
func (n *Node) Run() {
	// TODO: Error handling
	addr, _ := net.ResolveTCPAddr("tcp", n.IP+":"+strconv.Itoa(n.Port))

	// TODO: Error handling
	listener, _ := net.ListenTCP("tcp", addr)

	n.Listerner = listener

	incommingPacket := make(chan Packet)

	go n.receivePacket(incommingPacket)
	go n.processPacket(incommingPacket)
}

// Listen on binding address.
func (n *Node) receivePacket(packetch chan Packet) {
	for {
		// TODO: Error handling
		conn, _ := n.Listerner.AcceptTCP()

		buf := make([]byte, 4096)

		// TODO: Error handling
		bufLen, _ := conn.Read(buf)

		p := Packet{
			Content: buf[:bufLen],
			Conn:    conn,
		}

		// send packet to channel
		packetch <- p
	}

	defer n.Listerner.Close()
}

// Process packet.
func (n *Node) processPacket(packetch chan Packet) {
	for p := range packetch {
		var m Message
		// TODO: Error handling
		err := m.UnmarshalJson(p.Content)
		if err != nil {
			// We just drop the malformed message
			continue
		}

		n.MessageChannel <- IncommingMessage{Content: m, Conn: p.Conn}
	}
}
