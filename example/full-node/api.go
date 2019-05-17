package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/vgxbj/microchain/core"
)

const (
	apiURL = "/api/" + apiVersion + "/"
)

// Run web server.
func (c *client) runWebServer(port int) {

	http.HandleFunc(apiURL+"nodes", c.getNodesHandler)
	http.HandleFunc(apiURL+"pendings", c.getPendingTransactionsHandler)
	http.HandleFunc(apiURL+"transactions", c.getTransactionsHandler)

	http.HandleFunc(apiURL+"confirm", c.confirmPendingTransactionHandler)
	http.HandleFunc(apiURL+"send_transaction", c.sendTransactionHandler)

	http.HandleFunc("/", c.indexHandler)

	c.logger.Error.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))
}

func (c *client) indexHandler(w http.ResponseWriter, r *http.Request) {

	t, err := template.ParseFiles("static/templates/index.tmpl")
	if err != nil {
		c.logger.Info.Println(err)
	}

	t.Execute(w, &struct{ URL string }{URL: "http://localhost:" + strconv.Itoa(c.webport) + apiURL})
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

func (c *client) confirmPendingTransactionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()

		if r.Form["pending_id"] != nil && r.Form["confirm"] != nil {
			pendingID := r.Form["pending_id"][0]
			confirm := r.Form["confirm"][0]

			pendingIDBytes := core.Base58Decode(pendingID)

			if len(pendingIDBytes) == 0 {
				http.Redirect(w, r, "/", http.StatusSeeOther)
			}

			c.checkAndProcessPendingTransaction(pendingIDBytes, confirm)
		}
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (c *client) sendTransactionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()

		if r.Form["node_id"] != nil && r.Form["data"] != nil {
			nodeID := r.Form["node_id"][0]
			data := r.Form["data"][0]

			nodeIDBytes := core.Base58Decode(nodeID)

			if len(nodeIDBytes) == 0 {
				http.Redirect(w, r, "/", http.StatusSeeOther)
			}

			if !bytes.Equal(nodeIDBytes, c.node.PublicKey()) {
				t := c.newPendingTransaction(nodeIDBytes, data)

				go c.sendPendingTransaction(t)
			} else {
				go c.broadcastGenesisTransaction(data)
			}
		}
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (c *client) checkAndProcessPendingTransaction(id []byte, confirm string) {
	if confirm == "1" {
		c.confirmPendingTransaction(id)
	}
}
