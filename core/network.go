package core

import (
	"net"
	"strconv"
)

// Represent other nodes.
type RemoteNode struct {
	PublicKey  []byte        // Public key
	Address    *net.TCPAddr  // Address
	lastseen   int           // The unix time of seeing this node last time
	verifiedBy []*RemoteNode // Nodes that verify this node
}

// Received packect.
type Packet struct {
	Content []byte // Raw bytes
	Conn    *net.TCPConn
	// From    *net.TCPAddr // Address that send this packet
}

// Represent ourselves.
type Node struct {
	Keypair         *KeyPair               // Key pair
	IP              string                 // IP address
	Port            int                    // Port
	RoutingTable    map[string]*RemoteNode // Routing table (public key, node)
	Listerner       *net.TCPListener       // TCP listener
	IncommingPacket chan Packet            // Incomming packet
}

func NewNode(ip string, port int) *Node {
	kp, _ := NewECDSAKeyPair()

	return &Node{
		Keypair:         kp,
		IP:              ip,
		Port:            port,
		RoutingTable:    make(map[string]*RemoteNode),
		Listerner:       new(net.TCPListener),
		IncommingPacket: make(chan Packet),
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
		n.IncommingPacket <- p
	}
}
