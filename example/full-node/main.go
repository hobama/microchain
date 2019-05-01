package main

import (
	"flag"
	"os"

	"github.com/bosoncat/microchain/core"
)

var nodeIP = flag.String("addr", "localhost", "ip address that node runs on")
var nodePort = flag.Int("port", 3000, "port that node binds to")
var l *core.Logger

func init() {
	l = core.InitLogger(os.Stdout)
	flag.Parse()
}

func main() {
	c, err := newClient(*nodeIP, *nodePort, l)
	if err != nil {
		l.Error.Println(err)
		return
	}

	c.repl()
}
