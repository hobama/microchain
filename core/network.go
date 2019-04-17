package core

import (
	"net"
)

type ConnQueue chan string
type NodeQueue chan *Node
type NodesMap map[string]*Node
type Node struct {
	Conn     *net.TCPConn
	Lastseen int
}

type Network struct {
	Nodes              NodesMap     // Contacts
	ConnectionsQueue   ConnQueue    // Connections
	Address            string       // TCP-4 address
	ConnectionCallBack NodeQueue    // Connections callback
	BroadCastQueue     chan Message // Broadcast queue
	IncommingMessages  chan Message // Incomming messages
}

func (nm NodesMap) AddNode(n *Node) bool {
	k := n.Conn.RemoteAddr().String()

	if k != Self.Network.Address && nm[k] == nil {
		Info.Println("Node connected: ", k)
		nm[k] = n

		// go Handle(node)

		return true
	}

	return false
}

// Get IP address by hostname.
func GetIPAddr(hostname string) []string {
	addrs, err := net.LookupHost(hostname)
	if err != nil {
		return nil
	}

	return addrs
}
