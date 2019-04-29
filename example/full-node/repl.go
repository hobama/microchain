package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func (c *client) repl() {
	reader := bufio.NewReader(os.Stdin)

	for {
		c.terminal <- "> "
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "nodes":
			c.terminal <- "Nodes\n"

			for i := 0; i < 5; i++ {
				c.terminal <- fmt.Sprintf("aaa\n")
			}
		}
	}
}

func (c *client) printLoop() {
	for s := range c.terminal {
		fmt.Printf("%s", s)
	}
}
