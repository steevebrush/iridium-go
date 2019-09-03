# Iridium go RPC libraries
Modules to be used in go to query core Iridium softwares
Two parts are available : 
 * IridiumdRPC for the node daemon
 * IridiumWalletdRPC for the payment gateway daemon.

## IridiumdRPC

GET methods :
 * GetHeight()
 * GetInfo()
 * GetTransactions(txsHashes []string)
 * GetGeneratedCoins()

POST methods :
 * GetBlockCount(id ...string)
 * GetCurrencyid(id ...string)
 * GetLastBlockheader(id ...string)
 * GetBlockHeaderByHeight(height uint32, id ...string)
 * GetBlocksList(height uint32, id ...string)
 * GetBlockDetails(hash string, id ...string)
 * GetTransactionDetails(hash string, id ...string)
 * GetTransactionsPool(id ...string)

all methods returns a map from the JSON response

The iridiumsRPC_test.go contains all the methods, tested.
you can launch tests with
```bash
# cd IridiumdRPC
# go test -v
```

if your node is not running locally,  configure your node rpc api address and port here :
```go
// iridium node address for tests
var node = Iridiumd{
	address: "127.0.0.1",
	port:    13007,
}
```

## IridiumWalletdRPC
not ready yet