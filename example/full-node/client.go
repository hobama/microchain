package main

type client struct {
	terminal chan string
}

func newClient() *client {
	c := &client{
		terminal: make(chan string)}

	go c.printLoop()

	return c
}

func main() {
	c := newClient()

	c.repl()
}
