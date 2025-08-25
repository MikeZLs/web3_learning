package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
)

func main() {
	// 连接到以太坊网络
	client, err := ethclient.Dial("http://127.0.0.1:7545")
	if err != nil {
		log.Fatal("Failed to connect:", err)
	}
	defer client.Close()

	// 1. 检查连接
	fmt.Println("=== Connection Test ===")
	checkConnection(client)

	// 2. 查询最新区块
	fmt.Println("\n=== Latest Block ===")
	queryBlockByNumber(client, nil)

	// 3. 查询指定区块的交易
	fmt.Println("\n=== Block Transactions ===")
	blockNum := big.NewInt(1000000) // 选择一个有交易的区块
	queryBlockTransactions(client, blockNum)

	// 4. 查询特定交易（需要真实的交易哈希）
	// txHash := "真实的交易哈希值"
	// queryTransaction(client, txHash)
	// queryTransactionReceipt(client, txHash)
}
