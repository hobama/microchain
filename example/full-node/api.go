package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

const (
	apiURL = "/api/" + apiVersion + "/"
)

// Run web server.
func (c *client) runWebServer(port int) {

	http.HandleFunc(apiURL+"nodes", c.getNodesHandler)
	http.HandleFunc(apiURL+"pendings", c.getPendingTransactionsHandler)
	http.HandleFunc(apiURL+"transactions", c.getTransactionsHandler)

	http.HandleFunc("/", c.indexHandler)

	c.logger.Error.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))
}

func (c *client) indexHandler(w http.ResponseWriter, r *http.Request) {

	t, _ := template.ParseFiles("static/templates/index.tmpl")
	t.Execute(w, &struct{ URL string }{URL: "http://localhost:5001/api/v1/"})
}

func (c *client) getNodesHandler(w http.ResponseWriter, r *http.Request) {
	_, ns := c.node.GetNodesOfRoutingTable()

	nsjson, _ := json.Marshal(ns)
	fmt.Fprintf(w, string(nsjson))
}

func (c *client) getPendingTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	_, ts := c.node.GetPendingTransactions()

	tsjson, _ := json.Marshal(ts)
	fmt.Fprintf(w, string(tsjson))
}

func (c *client) getTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	_, ts := c.node.GetTransactionsOfPool()

	tsjson, _ := json.Marshal(ts)
	fmt.Fprintf(w, string(tsjson))
}
