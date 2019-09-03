/*
 * Copyright (c) 2019.
 * by Steve Brush, Iridium Developers
 */

// Iridium node JSON RPC API for golang tests
package iridiumdRPC

import (
	"encoding/json"
	"net"
	"strconv"
	"testing"
)

// Colorize output...
const ok = "\033[32m[OK] : \033[0m"
const er = "\033[31m[ERROR] : \033[0m"

// iridium node address for tests
var node = Iridiumd{
	address: "127.0.0.1",
	port:    13007,
}

// test the returning version
func TestVersion(t *testing.T) {
	name, major, minor, patch := Version()
	version := strconv.Itoa(major) + "." + strconv.Itoa(minor) + "." + strconv.Itoa(patch)
	t.Logf("%sPackage %s v%s found", ok, name, version)
}

// validate node constructor and resolver
func TestValidateAddress(t *testing.T) {
	// constructor test
	node1 := Iridiumd{
		address: "127.0.0.1",
		port:    13007,
	}

	node2 := Iridiumd{
		address: "127.0.0.1",
		port:    13007,
	}

	node3 := Iridiumd{
		address: "nodes.ird.cash",
		port:    13007,
	}

	//this one shouldn't resolve
	node4 := Iridiumd{
		address: "do.not.resolve",
		port:    13007,
	}

	// Validate constructor
	if node1 != node2 {
		t.Errorf("%sIridiumd struct error : want %s:%d, got %s:%d", er, node2.address, node2.port, node1.address, node1.port)
	} else {
		t.Logf("%sIridiumd struct ok", ok)
	}

	// dns resolver ok
	addr, err := net.ResolveIPAddr("ip", node3.address)
	if err != nil {
		t.Errorf("%sResolution error : %s", er, err.Error())
	}
	t.Logf("%sResolved address %s is %s", ok, node3.address, addr.String())

	// dns resolver error
	addr, err = net.ResolveIPAddr("ip", node4.address)
	if err != nil {
		t.Logf("%sResolution error : %s, this is expected.", ok, err.Error())
	}
}

func TestIridiumd_GetHeight(t *testing.T) {
	resp, err := node.GetHeight()
	if err != nil {
		t.Errorf("%s %s", er, err)
	} else {
		t.Logf("%sGetHeight returns :\n%v", ok, resp)
	}
}

func TestIridiumd_GetInfo(t *testing.T) {
	resp, err := node.GetInfo()
	if err != nil {
		t.Errorf("%s %s", er, err)
	} else {
		t.Logf("%sGetInfo returns :\n%v", ok, resp)
	}
}

func TestIridiumd_GetTransactions(t *testing.T) {
	txArray := []string{
		"d56f9b6e3257568151de667b679c5fc3c03b02ac6ce9a28d346e6c0f6beafd5c",
		"d56f9b6e3257568151de667b679c5fc3c03b02ac6ce9a28d346e6c0f6beafd56"}

	resp, err := node.GetTransactions(txArray)
	if err != nil {
		t.Errorf("%s %s", er, err)
	} else {
		t.Logf("%sGetTransactions returns :\n%v", ok, resp)
	}
}

func TestIridiumd_GetGeneratedCoins(t *testing.T) {
	resp, err := node.GetGeneratedCoins()
	if err != nil {
		t.Errorf("%s %s", er, err)
	} else {
		t.Logf("%sGetGeneratedCoins returns :\n%v", ok, resp)
		/*t.Logf("%sGetGeneratedCoins returns :\n%v",ok,printJson(resp, true) )*/
	}
}

// POST methods

func TestIridiumd_GetBlockCount(t *testing.T) {
	// without id
	resp, err := node.GetBlockCount()
	if err != nil {
		t.Errorf("%s %s", er, err)
	} else {
		t.Logf("%sGetBlockCount without id returns :\n%v", ok, resp)
	}
	// with an id
	resp, err = node.GetBlockCount("withID")
	if err != nil {
		t.Errorf("%s %s", er, err)
	} else {
		t.Logf("%sGetBlockCount with id returns :\n%v", ok, resp)
	}
}

func TestIridiumd_GetCurrencyid(t *testing.T) {
	// without id
	resp, err := node.GetCurrencyid()
	if err != nil {
		t.Errorf("%s %s", er, err)
	} else {
		t.Logf("%sGetCurrencyid without id returns :\n%v", ok, resp)
	}
	// with an id
	resp, err = node.GetCurrencyid("withID")
	if err != nil {
		t.Errorf("%s %s", er, err)
	} else {
		t.Logf("%sGetCurrencyid with id returns :\n%v", ok, resp)
	}
}

func TestIridiumd_GetLastBlockheader(t *testing.T) {
	// without id
	resp, err := node.GetLastBlockheader()
	if err != nil {
		t.Errorf("%s %s", er, err)
	} else {
		t.Logf("%sGetLastBlockheader without id returns :\n%v", ok, resp)
	}
	// with an id
	resp, err = node.GetLastBlockheader("withID")
	if err != nil {
		t.Errorf("%s %s", er, err)
	} else {
		t.Logf("%sGetLastBlockheader with id returns :\n%v", ok, resp)
	}
}

func TestIridiumd_GetBlockHeaderByHash(t *testing.T) {
	// without id
	resp, err := node.GetBlockHeaderByHash("9d59c3ac5acc80eef180cc15c6cd49febfc7f7131ed38b080c8436021b9caf43")
	if err != nil {
		t.Errorf("%s %s", er, err)
	} else {
		t.Logf("%sGetBlockHeaderByHash without id returns :\n%v", ok, resp)
	}
	// with id
	resp, err = node.GetBlockHeaderByHash("9d59c3ac5acc80eef180cc15c6cd49febfc7f7131ed38b080c8436021b9caf43", "withID")
	if err != nil {
		t.Errorf("%s %s", er, err)
	} else {
		t.Logf("%sGetBlockHeaderByHash with id returns :\n%v", ok, resp)
	}
}

func TestIridiumd_GetBlockHeaderByHeight(t *testing.T) {
	// without id (this is the genesis block ;-)
	resp, err := node.GetBlockHeaderByHeight(0)
	if err != nil {
		t.Errorf("%s %s", er, err)
	} else {
		t.Logf("%sGetBlockHeaderByHeight without id returns :\n%v", ok, resp)
	}
	// with id
	resp, err = node.GetBlockHeaderByHeight(100, "withID")
	if err != nil {
		t.Errorf("%s %s", er, err)
	} else {
		t.Logf("%sGetBlockHeaderByHeight with id returns :\n%v", ok, resp)
	}
}

func TestIridiumd_GetBlocksList(t *testing.T) {
	// without id (this is the genesis block ;-)
	resp, err := node.GetBlocksList(30)
	if err != nil {
		t.Errorf("%s %s", er, err)
	} else {
		t.Logf("%sGetBlocksList without id returns :\n%v", ok, resp)
	}
	// with id
	resp, err = node.GetBlocksList(100, "withID")
	if err != nil {
		t.Errorf("%s %s", er, err)
	} else {
		t.Logf("%sGetBlocksList with id returns :\n%v", ok, resp)
	}
}

func TestIridiumd_GetBlockDetails(t *testing.T) {
	// without id
	resp, err := node.GetBlockDetails("9d59c3ac5acc80eef180cc15c6cd49febfc7f7131ed38b080c8436021b9caf43")
	if err != nil {
		t.Errorf("%s %s", er, err)
	} else {
		t.Logf("%sGetBlockDetails without id returns :\n%v", ok, resp)
	}
	// with id
	resp, err = node.GetBlockDetails("9d59c3ac5acc80eef180cc15c6cd49febfc7f7131ed38b080c8436021b9caf43", "withID")
	if err != nil {
		t.Errorf("%s %s", er, err)
	} else {
		t.Logf("%sGetBlockDetails with id returns :\n%v", ok, resp)
	}
}

func TestIridiumd_GetTransactionDetails(t *testing.T) {
	// without id
	resp, err := node.GetTransactionDetails("ba29fad80ab5eb6741bac01e5326f7c28ced3238c3d7bce1abbd97180aa20ec2")
	if err != nil {
		t.Errorf("%s %s", er, err)
	} else {
		t.Logf("%sGetTransactionDetails without id returns :\n%v", ok, resp)
	}
	// with id
	resp, err = node.GetTransactionDetails("ba29fad80ab5eb6741bac01e5326f7c28ced3238c3d7bce1abbd97180aa20ec2", "withID")
	if err != nil {
		t.Errorf("%s %s", er, err)
	} else {
		t.Logf("%sGetTransactionDetails with id returns :\n%v", ok, resp)
	}
}

func TestIridiumd_GetTransactionsPool(t *testing.T) {
	// without id
	resp, err := node.GetTransactionsPool()
	if err != nil {
		t.Errorf("%s %s", er, err)
	} else {
		t.Logf("%sGetTransactionsPool without id returns :\n%v", ok, resp)
	}
	// with an id
	resp, err = node.GetTransactionsPool("withID")
	if err != nil {
		t.Errorf("%s %s", er, err)
	} else {
		t.Logf("%sGetTransactionsPool with id returns :\n%v", ok, resp)
	}
}

// returns a map[string] as json with or without indentation (indent bool parameter), mainly for debugging
func printJson(m interface{}, indent bool) string {
	var b []byte
	var err error
	if indent {
		b, err = json.MarshalIndent(m, "", "  ")
	} else {
		b, err = json.Marshal(m)
	}
	if err != nil {
		return err.Error()
	}
	return string(b)
}
