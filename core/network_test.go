package core

import (
	"errors"
	"net"
	"testing"
)

// Generate new peer.
func GeneratePeer(addr string) Peer {
	kp, _ := NewECDSAKeyPair()

	p := Peer{kp, new(Network), nil}

	p.Network.Address = addr
	p.Network.Nodes = make(map[string]*Node)

	return p
}

// Generate random message.
func GenRandomMessage() Message {
	return Message{GenRandomBytes(1)[0], GenRandomBytes(2), nil}
}

// Create simple echo server.
func SimpleEchoServer(listener *net.TCPListener, address string, handle func(net.Conn)) {
	addr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		panic(err)
	}

	listener, err = net.ListenTCP("tcp4", addr)
	if err != nil {
		panic(err)
	}

	for {
		conn, _ := (*listener).Accept()
		defer conn.Close()

		go handle(conn)
	}
}

// Test add node to network.
func TestAddNode(t *testing.T) {
	p := GeneratePeer("localhost:3000")

	kp1, _ := NewECDSAKeyPair()
	n := Node{new(net.TCPConn), 0, string(kp1.Public), p.Network.Address}

	if !p.AddNode(n) {
		panic(errors.New("(*Peer) AddNode() testing failed"))
	}

	if p.Nodes[string(kp1.Public)] == nil {
		panic(errors.New("(*Peer) AddNode() testing failed"))
	}
}

// Test sending/receiving messages.
func TestSendMessage(t *testing.T) {
	n0 := GeneratePeer("localhost:3000")
	n1 := GeneratePeer("localhost:3001")
	n2 := GeneratePeer("localhost:3002")

	// Add nodes.
	n0.AddNode(n1.AsNode())
	n0.AddNode(n2.AsNode())
	n1.AddNode(n0.AsNode())
	n1.AddNode(n2.AsNode())
	n2.AddNode(n0.AsNode())
	n2.AddNode(n1.AsNode())

	handleFunc := func(conn net.Conn) {
		buf := make([]byte, 1024)

		bufLen, _ := conn.Read(buf)
		if bufLen != 0 {
			conn.Write(buf[:bufLen])
		}
	}

	// Run nodes: client & server.
	go SimpleEchoServer(n0.Network.Listener, n0.Network.Address, handleFunc)
	go SimpleEchoServer(n1.Network.Listener, n1.Network.Address, handleFunc)
	go SimpleEchoServer(n2.Network.Listener, n2.Network.Address, handleFunc)

	// Test Send()/Recv()
	testSend := func(p Peer, n Node) {
		addr, _ := net.ResolveTCPAddr("tcp", n.Address)

		conn, err := net.DialTCP("tcp", nil, addr)
		if err != nil {
			panic(errors.New("Cannot dial " + n.Address + " via TCP."))
		}
		defer conn.Close()

		p.Nodes[n.PublicKey].Conn = conn

		randomMessage := GenRandomMessage()

		p.Send(n, randomMessage)

		recvMessage, err := p.Recv(n)
		if err != nil {
			panic(err)
		}

		if !randomMessage.EqualWith(recvMessage) {
			panic(errors.New("(*Peer) Send()/Recv() testing failed."))
		}
	}

	for i := 0; i < 5; i++ {
		go testSend(n0, n1.AsNode())
		go testSend(n1, n0.AsNode())
		go testSend(n2, n1.AsNode())
	}
}

func SimpleNodeProcess(p Peer) {
	handleFunc := func(conn net.Conn) {
		buf := make([]byte, 1024)

		bufLen, _ := conn.Read(buf)
		if bufLen != 0 {
			conn.Write(buf[:bufLen])
		}
	}

	// SimpleEchoServer
	go SimpleEchoServer(p.Network.Listener, p.Network.Address, handleFunc)

	// Connect to nodes.
	for _, n := range p.Network.Nodes {
		addr, err := net.ResolveTCPAddr("tcp", n.Address)
		if err != nil {
			panic(err)
		}

		conn, err := net.DialTCP("tcp", nil, addr)
		if err != nil {
			panic(errors.New("Cannot dial " + n.Address + " via TCP."))
		}
		defer conn.Close()

		n.Conn = conn
	}

	errch := make(chan error)
	m := GenRandomMessage()
	go p.BroadcastMessage(m, errch)

	for _, n := range p.Network.Nodes {
		recvMessage, err := p.Recv(*n)
		if err != nil {
			panic(err)
		}

		if !recvMessage.EqualWith(m) {
			panic(errors.New("(*Peer) BroadcastMessage() testing failed."))
		}
	}
}

func TestBroadcast(t *testing.T) {
	n0 := GeneratePeer("localhost:3000")
	n1 := GeneratePeer("localhost:3001")
	n2 := GeneratePeer("localhost:3000")

	// Add nodes.
	n0.AddNode(n1.AsNode())
	n0.AddNode(n2.AsNode())
	n1.AddNode(n0.AsNode())
	n1.AddNode(n2.AsNode())
	n2.AddNode(n0.AsNode())
	n2.AddNode(n1.AsNode())

	go SimpleNodeProcess(n0)
	go SimpleNodeProcess(n1)
	go SimpleNodeProcess(n2)
}
