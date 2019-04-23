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
func SimpleEchoServer(listener *net.Listener, addr string, handle func(net.Conn)) {
	*listener, _ = net.Listen("tcp", addr)

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
	n := Node{new(net.Conn), 0, string(kp1.Public), p.Network.Address}

	if !p.AddNode(n) {
		panic(errors.New("(*Peer) AddNode() testing failed"))
	}

	if p.Nodes[string(kp1.Public)] == nil {
		panic(errors.New("(*Peer) AddNode() testing failed"))
	}
}

// Test sending/receiving messages.
func TestSendMessage(t *testing.T) {
	s := GeneratePeer("localhost:3000")
	c := GeneratePeer("localhost:3001")

	// Add nodes.
	s.AddNode(c.AsNode())
	c.AddNode(s.AsNode())

	handleFunc := func(conn net.Conn) {
		buf := make([]byte, 1024)

		bufLen, _ := conn.Read(buf)
		if bufLen != 0 {
			conn.Write(buf[:bufLen])
		}
	}

	// Run nodes: client & server.
	go SimpleEchoServer(&s.Network.Listener, s.Network.Address, handleFunc)
	go SimpleEchoServer(&c.Network.Listener, c.Network.Address, handleFunc)

	// Test Send()/Recv()
	testSend := func(p Peer, n Node) {
		conn, err := net.Dial("tcp", n.Address)
		if err != nil {
			panic(errors.New("Cannot dial " + n.Address + " via TCP."))
		}
		defer conn.Close()

		p.Nodes[n.PublicKey].Conn = &conn

		randomMessage := GenRandomMessage()

		p.Send(n, &randomMessage)

		recvMessage, err := p.Recv(n)
		if err != nil {
			panic(err)
		}

		if !randomMessage.EqualWith(recvMessage) {
			panic(errors.New("(*Peer) Send()/Recv() testing failed."))
		}
	}

	for i := 0; i < 5; i++ {
		go testSend(c, s.AsNode())
		go testSend(s, c.AsNode())
	}
}
