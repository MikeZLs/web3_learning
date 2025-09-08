package main

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
)

func main() {
	// 连接到以太坊网络
	client, err := ethclient.Dial("https://eth.w3node.com/8f9e39edde7c2aa1737461a34384908393dd1d4216de22085e9e214fe0f04c73/api")
	if err != nil {
		log.Fatal("Failed to connect:", err)
	}
	defer client.Close()

	//// 1. 检查连接
	//fmt.Println("=== Connection Test ===")
	//checkConnection(client)
	//
	//// 2. 查询最新区块
	//fmt.Println("\n=== Latest Block ===")
	//queryBlockByNumber(client, nil)
	//
	//// 3. 查询指定区块的交易
	//fmt.Println("\n=== Block Transactions ===")
	//blockNum := big.NewInt(1000000) // 选择一个有交易的区块
	//queryBlockTransactions(client, blockNum)
	//
	//// 4. 查询特定交易（需要真实的交易哈希）
	txHash := "0x8c66b82c9a5c726c079923149c08055bb692e6f3bbb96490734c5c092ae60740"
	//queryTransaction(client, txHash)
	QueryTransactionReceipt(client, txHash)
}
