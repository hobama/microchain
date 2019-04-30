package main

import (
	"flag"
	"fmt"
)

var nodeIp = flag.String("addr", "localhost", "ip address that node runs on")
var nodePort = flag.Int("port", 3000, "port that node binds to")

func init() {
	flag.Parse()
}

func main() {
	c, err := newClient(*nodeIp, *nodePort)
	if err != nil {
		fmt.Println(err)
	}

	c.repl()
}
