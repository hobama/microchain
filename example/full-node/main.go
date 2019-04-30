package main

func main() {
	c, err := newClient("localhost", 3000)
	if err != nil {
	}

	c.repl()
}
