package core

import ()

type NodesMap map[string]*Node
type Node struct {
	/*
		In this project, I would like to use HTTP protocol, rather than TCP.
		Conn      *net.TCPConn  // Should use generic Conn, so that we could use various conn type
	*/
	Lastseen  int    // The seconds since last time seen this node
	PublicKey string // Public key of this node
	Address   string // TCP-4 Address
}

type Network struct {
	Nodes            NodesMap    // Contacts
	ConnectionsQueue chan string // Connections
	/*
		In this project, I would like to use HTTP protocol, rather than TCP.
		Listener          *net.TCPListener // Listener
	*/
	Address           string       // Address
	ConnectionPool    chan *Node   // Connections callback
	BroadcastQueue    chan Message // Broadcast queue
	IncommingMessages chan Message // Incomming messages
}

type Peer struct {
	*KeyPair    // Public key and private Key
	*Network    // Network
	*Blockchain // Blockchain
}
