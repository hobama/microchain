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
	Nodes             NodesMap         // Contacts
	ConnectionsQueue  chan string      // Connections
	Listener          *net.TCPListener // Listener
	Address           string           // TCP-4 address
	ConnectionPool    chan *Node       // Connections callback
	BroadcastQueue    chan Message     // Broadcast queue
	IncommingMessages chan Message     // Incomming messages
}

type Peer struct {
	*KeyPair    // Public key and private Key
	*Network    // Network
	*Blockchain // Blockchain
}
