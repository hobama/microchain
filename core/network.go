package core

import (
	"bytes"
	"encoding/json"
	"net"
	"sort"
	"strconv"
	"sync"
	"time"
)

const (
	MAX_RECV_PACKET = 4096 // Max receiving packet length
)

// RemoteNode ... Represent other nodes.
type RemoteNode struct {
	PublicKey  []byte        `json:"public_key"` // Public key
	Address    string        `json:"address"`    // Address
	Lastseen   int           `json:"lastseen"`   // The unix time of seeing this node last time
	VerifiedBy []*RemoteNode `json:"-"`          // Nodes that verify this node
}

// Addr ... Get address of remote node.
func (rn *RemoteNode) Addr() string {
	return rn.Address
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
	Keypair                 *KeyPair                // Key pair
	IP                      string                  // IP address
	Port                    int                     // Port
	RoutingTableLock        sync.RWMutex            // Routing table read write lock
	RoutingTable            map[string]*RemoteNode  // Routing table (public key, node)
	TransactionsPoolLock    sync.RWMutex            // Transactions pool read write lock
	TransactionsPool        map[string]*Transaction // Transactions pool
	PendingTransactionsLock sync.RWMutex            // Pending transactions lock
	PendingTransactions     map[string]*Transaction // Pending transactions
	PreviousTransaction     *Transaction            // Previous transaction
	ChainLock               sync.RWMutex            // Blockchain lock
	Chain                   Blockchain              // Blockchain
	Listerner               *net.TCPListener        // TCP listener
	MessageChannel          chan IncommingMessage   // Incomming message
}

// NewNode ... Generate new node.
func NewNode(ip string, port int) (*Node, error) {
	kp, err := NewECDSAKeyPair()
	if err != nil {
		return nil, err
	}

	return &Node{
		Keypair:                 kp,
		IP:                      ip,
		Port:                    port,
		RoutingTableLock:        sync.RWMutex{},
		RoutingTable:            make(map[string]*RemoteNode),
		TransactionsPoolLock:    sync.RWMutex{},
		TransactionsPool:        make(map[string]*Transaction),
		PendingTransactionsLock: sync.RWMutex{},
		PendingTransactions:     make(map[string]*Transaction),
		PreviousTransaction:     nil,
		ChainLock:               sync.RWMutex{},
		Chain:                   Blockchain{},
		Listerner:               new(net.TCPListener),
		MessageChannel:          make(chan IncommingMessage),
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
	return n.Keypair.Public
}

// Sign ... Sign.
func (n *Node) Sign(data []byte) []byte {
	sig, _ := n.Keypair.Sign(data)
	return sig
}

// SignTransaction ... Sign transaction.
func (n *Node) SignTransaction(t Transaction) Transaction {
	t.Header.RequesteeSignature = n.Sign(SHA256(t.Meta))
	return t
}

// UpdateKeyPair ... Update key pair for node.
func (n *Node) UpdateKeyPair() error {
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

		buf := make([]byte, MAX_RECV_PACKET)

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

	buf := make([]byte, MAX_RECV_PACKET)

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
func (n *Node) GetNodeByPublicKey(pk []byte) (bool, RemoteNode) {
	n.RoutingTableLock.RLock()
	defer n.RoutingTableLock.RUnlock()

	rn := n.RoutingTable[Base58Encode(pk)]

	if rn == nil {
		return false, RemoteNode{}
	}

	return true, *rn
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

// CheckAndAddTransactionToPool ... Check and add transaction to pool.
func (n *Node) CheckAndAddTransactionToPool(t Transaction) {
	n.TransactionsPoolLock.Lock()
	defer n.TransactionsPoolLock.Unlock()

	if n.TransactionsPool[Base58Encode(t.ID())] != nil {
		return
	}

	n.TransactionsPool[Base58Encode(t.ID())] = &t
}

// VerifyTransaction ... Verify a given transaction.
func (n *Node) VerifyTransaction(t Transaction) bool {
	if t.IsGenesisTransaction() {
		// This is genesis transaction.
		if t.Accepted() != 1 || t.Rejected() != 0 {
			return false
		}

		return t.VerifyTransactionID() && t.VerifyRequesterSig() && t.VerifyRequesteeSig()
	}

	// This is not genesis transaction.
	// TODO: Verify credits.

	return t.VerifyTransactionID() && t.VerifyRequesterSig() && t.VerifyRequesteeSig()
}

// VerifyPendingTransaction ... Verify a pending transaction.
func (n *Node) VerifyPendingTransaction(t Transaction) bool {
	// First check requestee public key
	if !bytes.Equal(t.RequesteePK(), n.PublicKey()) {
		return false
	}

	if !t.VerifyRequesterSig() {
		return false
	}

	return true
}

// SetGenesisTransaction ... Set genesis transaction.
// NOTE: This function can only be called once.
func (n *Node) SetGenesisTransaction(t Transaction) bool {
	if n.PrevTransaction() != nil {
		return false
	}

	n.PreviousTransaction = &t

	return true
}

// UpdatePrevTransaction ... Update previous transaction.
func (n *Node) UpdatePrevTransaction(t Transaction) bool {
	if n.PrevTransaction() != nil {
		n.PreviousTransaction = &t

		return true
	}

	return false
}

// PrevTransaction ... Get previous transaction.
func (n *Node) PrevTransaction() *Transaction {
	return n.PreviousTransaction
}

// PrevTransactionOf ... Get previoud transaction of given transaction.
func (n *Node) PrevTransactionOf(tr Transaction) (bool, Transaction) {
	n.ChainLock.RLock()
	defer n.ChainLock.RUnlock()

	if t, tr := n.Chain.GetTransactionByID(tr.PreviousID()); t {
		return true, tr
	}

	return false, Transaction{}
}

// GetTransactionByIDFromChain ... Get transaction by id.
func (n *Node) GetTransactionByIDFromChain(id []byte) (bool, Transaction) {
	n.ChainLock.RLock()
	defer n.ChainLock.RUnlock()

	if t, tr := n.Chain.GetTransactionByID(id); t {
		return true, tr
	}

	return false, Transaction{}
}

// GetTransactionByIDFromPool ... Get transaction by id.
func (n *Node) GetTransactionByIDFromPool(id []byte) (bool, Transaction) {
	n.TransactionsPoolLock.RLock()
	defer n.TransactionsPoolLock.RUnlock()

	tr := n.TransactionsPool[Base58Encode(id)]

	if tr == nil {
		return false, Transaction{}
	}

	return true, *tr
}

// GetTransactionsOfPool ... Get transactions of transactions pool.
func (n *Node) GetTransactionsOfPool() (int, TransactionSlice) {
	n.TransactionsPoolLock.RLock()
	defer n.TransactionsPoolLock.RUnlock()

	var trs TransactionSlice

	// We should collect transactions in order.
	for _, tr := range n.TransactionsPool {
		trs = append(trs, *tr)
	}

	// Sort by time.
	sort.Sort(trs)

	return len(trs), trs
}

// IsInTransactionsPool ... Is in transactions pool.
func (n *Node) IsInTransactionsPool(id []byte) bool {
	n.TransactionsPoolLock.RLock()
	defer n.TransactionsPoolLock.RUnlock()

	if n.TransactionsPool[Base58Encode(id)] != nil {
		return true
	}

	return false
}

// RemoveTransactionByIDFromPool ... Remove transaction by id.
func (n *Node) RemoveTransactionByIDFromPool(id []byte) {
	n.TransactionsPoolLock.Lock()
	defer n.TransactionsPoolLock.Unlock()

	delete(n.TransactionsPool, Base58Encode(id))
}

// CheckAndAddPendingTransaction ... Check and add transaction to pending pool.
func (n *Node) CheckAndAddPendingTransaction(t Transaction) {
	n.PendingTransactionsLock.Lock()
	defer n.PendingTransactionsLock.Unlock()

	if n.PendingTransactions[Base58Encode(t.ID())] != nil {
		return
	}

	n.PendingTransactions[Base58Encode(t.ID())] = &t
}

// GetPendingTransactionByID ... Get pending transaction by id.
func (n *Node) GetPendingTransactionByID(id []byte) (bool, Transaction) {
	n.PendingTransactionsLock.RLock()
	defer n.PendingTransactionsLock.RUnlock()

	tr := n.PendingTransactions[Base58Encode(id)]

	if tr == nil {
		return false, Transaction{}
	}

	return true, *tr
}

// GetPendingTransactions ... Get pending transactions.
func (n *Node) GetPendingTransactions() (int, TransactionSlice) {
	n.PendingTransactionsLock.RLock()
	defer n.PendingTransactionsLock.RUnlock()

	var trs TransactionSlice

	// We should collect transactions in order.
	for _, tr := range n.PendingTransactions {
		trs = append(trs, *tr)
	}

	// Sort by time.
	sort.Sort(trs)

	return len(trs), trs
}

// IsInPendingTransactions ... Is in pending transactions.
func (n *Node) IsInPendingTransactions(id []byte) bool {
	n.PendingTransactionsLock.RLock()
	defer n.PendingTransactionsLock.RUnlock()

	if n.PendingTransactions[Base58Encode(id)] != nil {
		return true
	}

	return false
}

// RemovePendingTransactionByID ... Remove pending transaction by id.
func (n *Node) RemovePendingTransactionByID(id []byte) {
	n.PendingTransactionsLock.Lock()
	defer n.PendingTransactionsLock.Unlock()

	delete(n.PendingTransactions, Base58Encode(id))
}

// NewGenesisTransaction ... Generate new genesis transaction.
func (n *Node) NewGenesisTransaction(data []byte) Transaction {
	timestamp := int(time.Now().Unix())
	timestampByte := UInt64ToBytes(uint64(timestamp))
	id := SHA256(JoinBytes(n.PublicKey(), n.PublicKey(), timestampByte))
	sig := n.Sign(SHA256(data))

	h := TransactionHeader{
		TransactionID:      id,
		Timestamp:          timestamp,
		PrevTransactionID:  id,
		RequesterPublicKey: n.PublicKey(),
		RequesterSignature: sig,
		RequesteePublicKey: n.PublicKey(),
		RequesteeSignature: sig,
	}

	txo := TXOutput{
		Accepted: 1,
		Rejected: 0,
	}

	return Transaction{
		Header: h,
		Meta:   data,
		Output: txo,
	}
}
