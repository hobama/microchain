package core

import (
	"bytes"
	"encoding/json"
	"net"
	"strconv"
	"sync"
)

// RemoteNode ... Represent other nodes.
type RemoteNode struct {
	PublicKey  []byte        `json:"public_key"` // Public key
	Address    string        `json:"address"`    // Address
	Lastseen   int           `json:"Lastseen"`   // The unix time of seeing this node last time
	VerifiedBy []*RemoteNode `json:"-"`          // Nodes that verify this node
}

// Addr ... Get address of remote node.
func (rn *RemoteNode) Addr() string {
	return rn.Address
}

// Packet ... Received packect.
type Packet struct {
	Content []byte       // Raw bytes
	Conn    *net.TCPConn // TCP connection
}

// IncommingMessage ...
type IncommingMessage struct {
	Content Message      // Message
	Conn    *net.TCPConn // TCP connection
}

// Node ... Represent ourselves.
type Node struct {
	KeypairLock          sync.RWMutex           // Key pair lock, we should change keypair of a node after it sends a transaction.
	Keypair              *KeyPair               // Key pair
	IP                   string                 // IP address
	Port                 int                    // Port
	RoutingTableLock     sync.RWMutex           // Routing table read write lock
	RoutingTable         map[string]*RemoteNode // Routing table (public key, node)
	TransactionsPoolLock sync.RWMutex           // Transactions pool read write lock
	TransactionsPool     TransactionSlice       // Transactions pool
	Listerner            *net.TCPListener       // TCP listener
	MessageChannel       chan IncommingMessage  // Incomming message
}

// NewNode ... Generate new node.
func NewNode(ip string, port int) (*Node, error) {
	kp, err := NewECDSAKeyPair()
	if err != nil {
		return nil, err
	}

	return &Node{
		KeypairLock:          sync.RWMutex{},
		Keypair:              kp,
		IP:                   ip,
		Port:                 port,
		RoutingTableLock:     sync.RWMutex{},
		RoutingTable:         make(map[string]*RemoteNode),
		TransactionsPoolLock: sync.RWMutex{},
		TransactionsPool:     TransactionSlice{},
		Listerner:            new(net.TCPListener),
		MessageChannel:       make(chan IncommingMessage),
	}, nil
}

// Run ... Run a simple TCP server.
func (n *Node) Run() error {
	listener, err := n.TCPListener()
	if err != nil {
		return err
	}

	n.Listerner = listener

	incommingPacket := make(chan Packet)

	go n.receivePacket(incommingPacket)
	go n.processPacket(incommingPacket)

	return nil
}

// Addr ... Get address of node.
func (n *Node) Addr() string {
	return n.IP + ":" + strconv.Itoa(n.Port)
}

// TCPListener ... Get TCP listener of node.
func (n *Node) TCPListener() (*net.TCPListener, error) {
	addr, err := n.TCPAddr()
	if err != nil {
		return nil, err
	}

	return net.ListenTCP("tcp", addr)
}

// TCPAddr ... Get TCP address of node.
func (n *Node) TCPAddr() (*net.TCPAddr, error) {
	return net.ResolveTCPAddr("tcp", n.Addr())
}

// PublicKey ... Get public key of node.
func (n *Node) PublicKey() []byte {
	n.KeypairLock.RLock()
	defer n.KeypairLock.RUnlock()

	kp := n.Keypair

	return kp.Public
}

// UpdateKeyPair ... Update key pair for node.
func (n *Node) UpdateKeyPair() error {
	n.KeypairLock.Lock()
	defer n.KeypairLock.Unlock()

	kp, err := NewECDSAKeyPair()
	if err != nil {
		return err
	}

	n.Keypair = kp

	return nil
}

// receivePacket ... Listen on binding address.
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
}

// processPacket ... Process packet.
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

// Send ... Send message to given node.
func (n *Node) Send(address string, data []byte, handleCallback func([]byte) error) error {
	addr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return err
	}

	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Write(data)
	if err != nil {
		return err
	}

	buf := make([]byte, 4096)

	buflen, err := conn.Read(buf)
	if err != nil {
		return err
	}

	err = handleCallback(buf[:buflen])
	if err != nil {
		return err
	}

	return nil
}

// Broadcast ... Broadcast data to all nodes.
func (n *Node) Broadcast(data []byte, handleCallback func([]byte) error) {
	_, nodes := n.GetNodesOfRoutingTable()

	for _, rn := range nodes {
		n.Send(rn.Addr(), data, handleCallback)
	}
}

// CheckAndAddNodeToRoutingTable ... Check and add node to routing table.
func (n *Node) CheckAndAddNodeToRoutingTable(rn RemoteNode) {
	n.RoutingTableLock.Lock()
	defer n.RoutingTableLock.Unlock()

	if n.RoutingTable[Base58Encode(rn.PublicKey)] != nil {
		return
	}

	n.RoutingTable[Base58Encode(rn.PublicKey)] = &rn
}

// UpdateNodeForGivenPublicKey ... Update node for given public key.
func (n *Node) UpdateNodeForGivenPublicKey(pk []byte, rn RemoteNode) {
	n.RoutingTableLock.Lock()
	defer n.RoutingTableLock.Unlock()

	n.RoutingTable[Base58Encode(pk)] = &rn
}

// GetNodeByPublicKey ... Get node by public key.
func (n *Node) GetNodeByPublicKey(pk []byte) RemoteNode {
	n.RoutingTableLock.RLock()
	defer n.RoutingTableLock.RUnlock()

	rn := n.RoutingTable[Base58Encode(pk)]

	return *rn
}

// GetNodesOfRoutingTable ... Get nodes of routing table.
func (n *Node) GetNodesOfRoutingTable() (int, []RemoteNode) {
	n.RoutingTableLock.RLock()
	defer n.RoutingTableLock.RUnlock()

	var nodes []RemoteNode

	for _, rn := range n.RoutingTable {
		nodes = append(nodes, *rn)
	}

	return len(nodes), nodes
}

// IsInRoutingTable .. Is in routing table.
func (n *Node) IsInRoutingTable(pk []byte) bool {
	n.RoutingTableLock.RLock()
	defer n.RoutingTableLock.RUnlock()

	if n.RoutingTable[Base58Encode(pk)] != nil {
		return true
	}

	return false
}

// RemoveNodeByPublicKey ... Remove node by public key.
func (n *Node) RemoveNodeByPublicKey(pk []byte) {
	n.RoutingTableLock.Lock()
	defer n.RoutingTableLock.Unlock()

	delete(n.RoutingTable, Base58Encode(pk))
}

// MarshalJson ... Serialize RemoteNode into Json.
func (rn RemoteNode) MarshalJson() ([]byte, error) {
	return json.Marshal(rn)
}

// UnmarshalJson ... Read RemoteNode from Json.
func (rn *RemoteNode) UnmarshalJson(data []byte) error {
	return json.Unmarshal(data, &rn)
}

// EqualWith ... Test if two remote nodes are equal.
func (rn RemoteNode) EqualWith(temp RemoteNode) bool {
	if !bytes.Equal(rn.PublicKey, temp.PublicKey) {
		return false
	}

	if rn.Address != temp.Address {
		return false
	}

	if rn.Lastseen != temp.Lastseen {
		return false
	}

	return true
}
