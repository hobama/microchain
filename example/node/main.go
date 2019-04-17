package main

import (
	"github.com/bosoncat/microchain/core"
	"fmt"
)

func main() {
	h := core.GetIPAddr("localhost")
	fmt.Println(h)
}
