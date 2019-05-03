package core

import (
	"fmt"
	"math/rand"
	"testing"
)

// GenRandomRemoteNode ... Generate random remote node.
func GenRandomRemoteNode() RemoteNode {
	return RemoteNode{
		PublicKey:  GenRandomBytes(64),
		Address:    "localhost:3000",
		Lastseen:   rand.Intn(10000),
		VerifiedBy: nil,
	}
}

// GenRandomRemoteNodes ... Generate random remote nodes.
func GenRandomRemoteNodes(n int) []RemoteNode {
	var nodes []RemoteNode

	for i := 0; i < n; i++ {
		nodes = append(nodes, GenRandomRemoteNode())
	}

	return nodes
}

// Test RemoteNode marshal function.
func TestRemoteNodeMarshalJson(t *testing.T) {
	rn1 := GenRandomRemoteNode()

	rn1json, err := rn1.MarshalJson()
	if err != nil {
		panic(fmt.Errorf("(RemoteNode) MarshalJson() testing failed"))
	}

	var rn2 RemoteNode

	err = rn2.UnmarshalJson(rn1json)
	if err != nil {
		panic(fmt.Errorf("(*RemoteNode) UnmarshalJson() testing failed"))
	}

	if !rn1.EqualWith(rn2) {
		panic(fmt.Errorf("(RemoteNode) MarshalJson()/UnmarshalJson() testing failed"))
	}
}
