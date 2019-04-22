package core

import (
	"bytes"
	"errors"
	"net"
	"testing"
)

// Generate new peer.
func GeneratePeer(addr string) Peer {
	kp, _ := NewECDSAKeyPair()

	p := Peer{kp, new(Network)}

	p.Network.Address = addr
	p.Network.Nodes = make(map[string]*Node)

	return p
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
	n := Node{new(net.TCPConn), 0, string(kp1.Public)}

	if !p.AddNode(n) {
		panic(errors.New("(*Peer) AddNode() testing failed"))
	}

	if p.Nodes[string(kp1.Public)] == nil {
		panic(errors.New("(*Peer) AddNode() testing failed"))
	}
}

// Test sending messages.
func TestSendMessage(t *testing.T) {
	testPattern := []string{"hello", "world"}

	s := GeneratePeer("localhost:3000")

	handleFunc := func(conn net.Conn) {
		buf := make([]byte, 1024)

		buflen, _ := conn.Read(buf)
		if buflen != 0 {
			conn.Write(buf[:buflen])
		}
	}

	go SimpleEchoServer(&s.Network.Listener, s.Network.Address, handleFunc)

	for _, str := range testPattern {
		conn, err := net.Dial("tcp", s.Network.Address)
		if err != nil {
			panic(errors.New("Cannot dial " + s.Network.Address + " via TCP"))
		}
		defer conn.Close()

		_, err = conn.Write([]byte(str))
		if err != nil {
			panic(err)
		}

		buf := make([]byte, 1024)
		buflen, _ := conn.Read(buf)

		if !bytes.Equal(buf[:buflen], []byte(str)) {
			panic(errors.New("(*Peer) Send() testing failed."))
		}
	}
}
