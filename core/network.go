package core

import (
	"net"
)

type RemoteNode struct {
	PublicKey  []byte        // Public key
	Address    *net.TCPAddr  // Address
	lastseen   int           // The unix time of seeing this node last time
	verifiedBy []*RemoteNode // Nodes that verify this node
}

type Node struct {
	PublicKey     []byte                 // Public key of this node
	IP            string                 // TCP-4 Address
	Port          int                    // Port
	AddressToNode map[string]*RemoteNode // Remote nodes
}
