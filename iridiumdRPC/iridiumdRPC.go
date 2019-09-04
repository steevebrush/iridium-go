/*
 * Copyright (c) 2019.
 * by Steve Brush, Iridium Developers
 */

// Iridium node JSON RPC API for golang

package iridiumdRPC

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

// Version, returns version major, minor and patch
// todo : get git revision
func Version() (name string, major int, minor int, patch int) {
	return "iridiumdRPC", 0, 0, 1
}

// node json/rpc api address, port and minimum version needed
type Iridiumd struct {
	Address string
	Port    int
}

// Perform server request
func doRequest(req *http.Request) (*http.Response, error) {
	// use custom client : default timeout is "no timeout", this mean unlimited...
	netClient := &http.Client{
		Timeout: time.Second * 30,
	}
	resp, err := netClient.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Check server response
func handleServerResponse(resp *http.Response) (interface{}, error) {
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Server response : " + resp.Status)
	}
	defer resp.Body.Close()
	// handle responses, errors
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if len(body) == 0 {
		return nil, errors.New("body is empty")
	}
	var mapBody interface{}
	if err = json.Unmarshal(body, &mapBody); err != nil {
		return nil, err
	}
	return mapBody, nil
}

func (node *Iridiumd) makeGetRequest(method string, params map[string]interface{}) (interface{}, error) {
	// json parameters to send
	var jsonPayload []byte
	if params != nil {
		var err error
		params["jsonrpc"] = "2.0"
		jsonPayload, err = json.Marshal(params)
		if err != nil {
			return nil, err
		}
	}

	// construct request
	req, err := http.NewRequest("GET", "http://"+node.Address+":"+strconv.Itoa(node.Port)+"/"+method, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, err
	}

	// perform request
	resp, err := doRequest(req)
	if err != nil {
		return nil, err
	}

	// Handle server response
	body, err := handleServerResponse(resp)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (node *Iridiumd) makePostRequest(method string, params map[string]interface{}) (interface{}, error) {
	// json parameters
	payload := make(map[string]interface{})
	payload["jsonrpc"] = "2.0"
	payload["method"] = method
	payload["params"] = params

	// check if an id exists and use it
	if idValue, exist := params["id"]; exist {
		payload["id"] = idValue
		delete(params, "id")
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	// construct request
	req, err := http.NewRequest("POST", "http://"+node.Address+":"+strconv.Itoa(node.Port)+"/json_rpc", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, err
	}

	// perform request
	resp, err := doRequest(req)
	if err != nil {
		return nil, err
	}

	// Handle server response
	body, err := handleServerResponse(resp)
	if err != nil {
		return nil, err
	}

	// check if req(id) and resp(id match)
	if payload["id"] != body.(map[string]interface{})["id"] {
		return body, errors.New("Warning : ids doesn't match ")
	}

	return body, nil
}

// node GET methods

/*
/getheight, returns current node height, current network height and status
output example : map[height:357588 network_height:357588 status:OK]
*/
func (node *Iridiumd) GetHeight() (map[string]interface{}, error) {
	resp, err := node.makeGetRequest("getheight", nil)
	if err != nil {
		return nil, err
	}
	return resp.(map[string]interface{}), nil
}

/*
/getinfo, returns global node informations
output example : map[alt_blocks_count:0 difficulty:2e+09 grey_peerlist_size:1480 height:357588 incoming_connections_count:14 last_known_block_index:357587 outgoing_connections_count:8 status:OK synced:true tx_count:532754 tx_pool_size:26 version:5.0.0 (0f3cc89) Release white_peerlist_size:108]
*/
func (node *Iridiumd) GetInfo() (map[string]interface{}, error) {
	resp, err := node.makeGetRequest("getinfo", nil)
	if err != nil {
		return nil, err
	}
	return resp.(map[string]interface{}), nil
}

/*
/gettransactions, returns found or not tx and node status
input params : txs_hashes []string
output example: map[missed_tx:[] status:OK txs_as_hex:[]
*/
func (node *Iridiumd) GetTransactions(txsHashes []string) (map[string]interface{}, error) {
	payload := make(map[string]interface{})
	payload["txs_hashes"] = txsHashes
	resp, err := node.makeGetRequest("gettransactions", payload)
	if err != nil {
		return nil, err
	}
	return resp.(map[string]interface{}), nil
}

/*
/get_generated_coins, returns status and current circulating coins * denomination unit : 100000000
output example :  map[alreadyGeneratedCoins:1860645432419705 status:OK]
this mean : 1860645432419705/100000000 = 18 606 454,3241971 IRD
*/
func (node *Iridiumd) GetGeneratedCoins() (map[string]interface{}, error) {
	resp, err := node.makeGetRequest("get_generated_coins", nil)
	if err != nil {
		return nil, err
	}
	return resp.(map[string]interface{}), nil
}

// POST methods

/*
getblockcount, returns current height (including current mined block),
id is optional but when specified, request id and response id are compared
output :  map[jsonrpc:2.0 result:map[count:357655 status:OK]]
or with id : map[id:withID jsonrpc:2.0 result:map[count:357655 status:OK]]
*/
func (node *Iridiumd) GetBlockCount(id ...string) (map[string]interface{}, error) {
	payload := make(map[string]interface{})
	if len(id) != 0 {
		payload["id"] = id[0]
	}
	resp, err := node.makePostRequest("getblockcount", payload)
	if err != nil {
		return nil, err
	}
	return resp.(map[string]interface{}), nil
}

/*
getcurrencyid, returns genesis block hash
id is optional but when specified, request id and response id are compared
output : map[jsonrpc:2.0 result:map[currency_id_blob:9d59c3ac5acc80eef180cc15c6cd49febfc7f7131ed38b080c8436021b9caf43]]
or with id : map[id:withID jsonrpc:2.0 result:map[currency_id_blob:9d59c3ac5acc80eef180cc15c6cd49febfc7f7131ed38b080c8436021b9caf43]]
*/
func (node *Iridiumd) GetCurrencyid(id ...string) (map[string]interface{}, error) {
	payload := make(map[string]interface{})
	if len(id) != 0 {
		payload["id"] = id[0]
	}
	resp, err := node.makePostRequest("getcurrencyid", payload)
	if err != nil {
		return nil, err
	}
	return resp.(map[string]interface{}), nil
}

/*
getlastblockheader, last mined block header
id is optional but when specified, request id and response id are compared
output :  map[jsonrpc:2.0 result:map[block_header:map[depth:0 difficulty:1.75e+09 hash:982addbdf1dc687e886baffcbdf6c29fabab8269787980991ee1c34242433cdd height:357765 major_version:5 minor_version:0 nonce:31207 orphan_status:false prev_hash:3b7ca9fdf86ad62a7411827299b4028e5d2701cfa3ccae1a3390faebecfa5b42 reward:2.437467306e+09 timestamp:1.567538093e+09] status:OK]]
or with id : map[id:withID jsonrpc:2.0 result:map[block_header:map[depth:0 difficulty:1.75e+09 hash:982addbdf1dc687e886baffcbdf6c29fabab8269787980991ee1c34242433cdd height:357765 major_version:5 minor_version:0 nonce:31207 orphan_status:false prev_hash:3b7ca9fdf86ad62a7411827299b4028e5d2701cfa3ccae1a3390faebecfa5b42 reward:2.437467306e+09 timestamp:1.567538093e+09] status:OK]]
*/
func (node *Iridiumd) GetLastBlockheader(id ...string) (map[string]interface{}, error) {
	payload := make(map[string]interface{})
	if len(id) != 0 {
		payload["id"] = id[0]
	}
	resp, err := node.makePostRequest("getlastblockheader", payload)
	if err != nil {
		return nil, err
	}
	return resp.(map[string]interface{}), nil
}

/*
getblockheaderbyhash, returns the block header by hash
id is optional but when specified, request id and response id are compared
input : hash string
output : map[jsonrpc:2.0 result:map[block_header:map[depth:357767 difficulty:1 hash:9d59c3ac5acc80eef180cc15c6cd49febfc7f7131ed38b080c8436021b9caf43 height:0 major_version:1 minor_version:0 nonce:70 orphan_status:false prev_hash:0000000000000000000000000000000000000000000000000000000000000000 reward:9.536743164e+09 timestamp:0] status:OK]]
or with id :  map[id:withID jsonrpc:2.0 result:map[block_header:map[depth:357767 difficulty:1 hash:9d59c3ac5acc80eef180cc15c6cd49febfc7f7131ed38b080c8436021b9caf43 height:0 major_version:1 minor_version:0 nonce:70 orphan_status:false prev_hash:0000000000000000000000000000000000000000000000000000000000000000 reward:9.536743164e+09 timestamp:0] status:OK]]
*/
func (node *Iridiumd) GetBlockHeaderByHash(hash string, id ...string) (map[string]interface{}, error) {
	payload := make(map[string]interface{})
	payload["hash"] = hash
	if len(id) != 0 {
		payload["id"] = id[0]
	}
	resp, err := node.makePostRequest("getblockheaderbyhash", payload)
	if err != nil {
		return nil, err
	}
	return resp.(map[string]interface{}), nil
}

/*
getblockheaderbyheight, returns the block header at desired height
id is optional but when specified, request id and response id are compared
input : height uint32
output : map[jsonrpc:2.0 result:map[block_header:map[depth:357672 difficulty:1 hash:9d59c3ac5acc80eef180cc15c6cd49febfc7f7131ed38b080c8436021b9caf43 height:0 major_version:1 minor_version:0 nonce:70 orphan_status:false prev_hash:0000000000000000000000000000000000000000000000000000000000000000 reward:9.536743164e+09 timestamp:0] status:OK]]
or with id : map[id:withID jsonrpc:2.0 result:map[block_header:map[depth:357572 difficulty:10857 hash:823f7bf3e6ccf9818c7b58aebebde7bf79f25b5118dc02233ac38333340ed894 height:100 major_version:1 minor_version:0 nonce:4.29672923e+08 orphan_status:false prev_hash:3f726a8f697c1cc03f54bf0f1d609ef677b7b8597f24a67aa733bd8a810f023c reward:9.533105872e+09 timestamp:1.504560271e+09] status:OK]]
*/
func (node *Iridiumd) GetBlockHeaderByHeight(height uint32, id ...string) (map[string]interface{}, error) {
	payload := make(map[string]interface{})
	payload["height"] = height
	if len(id) != 0 {
		payload["id"] = id[0]
	}
	resp, err := node.makePostRequest("getblockheaderbyheight", payload)
	if err != nil {
		return nil, err
	}
	return resp.(map[string]interface{}), nil
}

/*
f_blocks_list_json, returns 30 blocks headers from desired height to height - 30
id is optional but when specified, request id and response id are compared
input : height uint32
output : map[jsonrpc:2.0 result:map[blocks:[map[cumul_size:410 difficulty:9997 hash:f8efb98beee5930a403b16a12bb212f88a06c84dadb330090eb0e22528b3c90f height:30 reward:9.53565183e+09 timestamp:1.504551188e+09 tx_count:1], etc...
or with id : map[id:withID jsonrpc:2.0 result:map[blocks:[map[cumul_size:408 difficulty:10857 hash:823f7bf3e6ccf9818c7b58aebebde7bf79f25b5118dc02233ac38333340ed894 height:100 reward:9.533105872e+09 timestamp:1.504560271e+09 tx_count:1], etc...
*/
func (node *Iridiumd) GetBlocksList(height uint32, id ...string) (map[string]interface{}, error) {
	payload := make(map[string]interface{})
	payload["height"] = height
	if len(id) != 0 {
		payload["id"] = id[0]
	}
	resp, err := node.makePostRequest("f_blocks_list_json", payload)
	if err != nil {
		return nil, err
	}
	return resp.(map[string]interface{}), nil
}

/*
f_block_json, returns block detail at  desired hash
id is optional but when specified, request id and response id are compared
input : hash string
output :  map[jsonrpc:2.0 result:map[block:map[alreadyGeneratedCoins:9536743164 alreadyGeneratedTransactions:1 baseReward:9.536743164e+09 blockSize:118 depth:357785 difficulty:1 effectiveSizeMedian:20000 hash:9d59c3ac5acc80eef180cc15c6cd49febfc7f7131ed38b080c8436021b9caf43 height:0 major_version:1 minor_version:0 nonce:70 orphan_status:false penalty:0 prev_hash:0000000000000000000000000000000000000000000000000000000000000000 reward:9.536743164e+09 sizeMedian:0 timestamp:0 totalFeeAmount:0 transactions:[map[amount_out:9.536743164e+09 fee:0 hash:f3fe271b4edceebf60a29d535a8dec957809baf4c69549a09ae113eb88a5f1ad size:78]] transactionsCumulativeSize:78] status:OK]]
or with id : same with [id:withID...
*/
func (node *Iridiumd) GetBlockDetails(hash string, id ...string) (map[string]interface{}, error) {
	payload := make(map[string]interface{})
	payload["hash"] = hash
	if len(id) != 0 {
		payload["id"] = id[0]
	}
	resp, err := node.makePostRequest("f_block_json", payload)
	if err != nil {
		return nil, err
	}
	return resp.(map[string]interface{}), nil
}

/*
f_transaction_json, returns block detail at  desired hash
id is optional but when specified, request id and response id are compared
input : hash string
output :  map[jsonrpc:2.0 result:map[block:map[cumul_size:8430 difficulty:4.75677472e+08 hash:f0c002e703d8dc19f6a5ca844805015213f223f3a52a258d3fa22a69e8177213 height:357782 reward:3.85792735e+08 timestamp:1.567540598e+09 tx_count:7] status:OK tx:map[extra:01e3a1b6040ba307633583398b8f59f37d71f643e8c06ad761a546d0feeeccccad021100000000e53a8eb3000000000000000000 unlock_time:357802 version:1 vin:[map[type:ff value:map[height:357782]]] vout:[map[amount:241 target:map[data:map[key:5f7a83a576e401ad0c1a10b0b9050be3e02640604be5ddca986693636bf84a58] type:02]] map[amount:9000 target:map[data:map[key:4b40a1f97a411881321022de6a8f532f30c4c404c39c345e26250e162935af79] type:02]] map[amount:600000 target:map[data:map[key:c7c6b07057750fff3ba1520ffbd71aeeb080f796315764ac2021f7f28d064b42] type:02]] map[amount:7e+06 target:map[data:map[key:638e1cc733f3566091bfbf696c680264a53dec7e75f82b5d09bd96b7e0e788d0] type:02]] map[amount:3e+07 target:map[data:map[key:fc0b8cb0e7bf0a7fd8bee0765add136a03b72ddca9c5bd83b755ac3900cc4822] type:02]] map[amount:4e+08 target:map[data:map[key:e9caba0c8f62238ebe81997faa69212c26a21d26199f2c2ba3aa125e3ec33c2c] type:02]] map[amount:2e+09 target:map[data:map[key:d1592617829f57545c0929e4cb286c019da0dd1d1910b7531a25a396466b0bda] type:02]]]] txDetails:map[amount_out:2.437609241e+09 fee:0 hash:ba29fad80ab5eb6741bac01e5326f7c28ced3238c3d7bce1abbd97180aa20ec2 mixin:0 paymentId: size:319]]]
or with id : same with [id:withID...
*/
func (node *Iridiumd) GetTransactionDetails(hash string, id ...string) (map[string]interface{}, error) {
	payload := make(map[string]interface{})
	payload["hash"] = hash
	if len(id) != 0 {
		payload["id"] = id[0]
	}
	resp, err := node.makePostRequest("f_transaction_json", payload)
	if err != nil {
		return nil, err
	}
	return resp.(map[string]interface{}), nil
}

/*
f_on_transactions_pool_json, get mem pool status
id is optional but when specified, request id and response id are compared
output :  7c0902baf3e50f7f3a202bdafe5091a1eb6befb9eedd89b6cbfa2b051450a73a
or with id : map[id:withID jsonrpc:2.0 ...
*/
func (node *Iridiumd) GetTransactionsPool(id ...string) (map[string]interface{}, error) {
	payload := make(map[string]interface{})
	if len(id) != 0 {
		payload["id"] = id[0]
	}
	resp, err := node.makePostRequest("f_on_transactions_pool_json", payload)
	if err != nil {
		return nil, err
	}
	return resp.(map[string]interface{}), nil
}
