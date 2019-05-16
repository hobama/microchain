package main

const (
	invalidPeriod                   = 50 // If one node does not respond for 50 seconds, it will be invalid.
	pingPeriod                      = 5  // Ping other nodes, every 5 seconds.
	broadcastRoutingTablePeriod     = 7  // Broadcast routing table, every 7 seconds.
	broadcastTransactionsPoolPeriod = 7  // Broadcast transactions pool, every 7 seconds.

	apiVersion = "v1" // API version
)
