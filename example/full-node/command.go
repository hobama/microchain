package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/bosoncat/microchain/core"
)

// Options
var queryNodesOpt = regexp.MustCompile(`nodes`)
var pingNodeOpt = regexp.MustCompile(`ping`)
var joinNetworkOpt = regexp.MustCompile(`join`)
var sendTransactionOpt = regexp.MustCompile(`tran`)
var genesisOpt = regexp.MustCompile(`genesis`)
var queryPendingJobsOpt = regexp.MustCompile(`pending`)
var queryTransactionsOpt = regexp.MustCompile(`transactions`)
var confirmReqOpt = regexp.MustCompile(`confirm`)

func checkQueryNodesCommand(s string) (bool, string) {
	if s != "nodes" {
		return false, fmt.Sprintf("Unknown command: %s, do you mean: nodes ?\n", s)
	}

	return true, ""
}

func checkPingNodeCommand(s string) (bool, string, []string) {
	if !strings.HasPrefix(s, "ping") {
		return false, fmt.Sprintf("Unknown command: %s, do you mean: ping ?\n", s), []string{}
	}

	// Remove `ping`
	s = strings.TrimSpace(s[4:])

	addrs := addrRegex.FindAllString(s, -1)

	if len(addrs) == 0 {
		return false, "Please check the address that you want to ping\n", []string{}
	}

	return true, "", addrs
}

func checkQueryPendingCommand(s string) (bool, string) {
	if s != "pending" {
		return false, fmt.Sprintf("Unknown command: %s, do you mean: pending ?\n", s)
	}

	return true, ""
}

func checkQueryTransactionsCommand(s string) (bool, string) {
	if s != "transactions" {
		return false, fmt.Sprintf("Unknown command: %s, do you mean: transactions ?\n", s)
	}

	return true, ""
}

func checkSendTransactionCommand(s string) (bool, string, []byte, string) {
	if !strings.HasPrefix(s, "tran") {
		return false, fmt.Sprintf("Unknown command: %s, do you mean: tran ?\n", s), nil, ""
	}

	// Remove `tran`
	s = strings.TrimSpace(s[4:])

	// TODO:
	tokens := strings.Fields(s)

	if len(tokens) != 2 {
		return false, fmt.Sprintf("Do you mean: tran id data ?\n"), nil, ""
	}

	id := core.Base58Decode(tokens[0])
	if len(id) == 0 {
		return false, fmt.Sprintf("Invalid node id: %s\n", tokens[0]), nil, ""
	}

	return true, "", id, tokens[1]
}

func checkGenesisCommand(s string) (bool, string, string) {
	if !strings.HasPrefix(s, "genesis") {
		return false, fmt.Sprintf("Unknown command: %s, do you mean: genesis ?\n", s), ""
	}

	tokens := strings.Fields(s)

	if len(tokens) != 2 {
		return false, fmt.Sprintf("Invalid genesis\n"), ""
	}

	return true, "", tokens[1]
}

func checkConfirmCommand(s string) (bool, string, []string) {
	if !strings.HasPrefix(s, "confirm") {
		return false, fmt.Sprintf("Unknown command: %s, do you mean: confirm ?\n", s), []string{}
	}

	// Remove `confirm`
	s = strings.TrimSpace(s[7:])

	// TODO:

	return true, "", []string{}
}
