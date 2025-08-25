package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	// 连接到以太坊测试节点（使用Ganache本地测试）
	client, err := ethclient.Dial("http://127.0.0.1:7545")
	if err != nil {
		log.Fatal("Failed to connect to Ethereum client:", err)
	}
	defer client.Close()

	fmt.Println("Successfully connected to Ethereum network!")

	// 测试连接 - 获取chainID
	//client.NetworkID(context.Background())
	chainID, err := client.ChainID(context.Background()) // ChainID与NetworkID通常是等价的，但ChainID安全性更高，更常用
	if err != nil {
		log.Fatal("Failed to get chainID:", err)
	}

	fmt.Printf("Connected to chainID: %s\n", chainID.String())

	checkConnection(client)

	//////// 按区块高度查询 ////////
	// 查询最新区块
	fmt.Println("=== Latest Block ===")
	queryBlockByNumber(client, nil) // nil 表示最新区块

	// 查询指定区块
	//fmt.Println("=== Block 10000 ===")
	//blockNum := big.NewInt(10000)
	//queryBlockByNumber(client, blockNum)

	// 查询最新区块号
	latestBlockNumber, err := client.BlockNumber(context.Background())
	if err != nil {
		log.Fatal("Failed to get latest block number:", err)
	}
	fmt.Printf("Latest block number: %d\n", latestBlockNumber)

	//////// 按区块哈希查询 ////////
	//queryBlockByHash(client, getBlockHashFromLatest(client))

}

func getBlockHashFromLatest(client *ethclient.Client) string {
	ctx := context.Background()

	// 获取最新区块
	latestBlock, err := client.BlockByNumber(ctx, nil)
	if err != nil {
		log.Fatal("获取最新区块失败:", err)
	}

	fmt.Printf("最新区块哈希: %s\n", latestBlock.Hash().Hex())
	fmt.Printf("父区块哈希: %s\n", latestBlock.ParentHash().Hex())

	// 返回父区块哈希用于后续查询
	return latestBlock.ParentHash().Hex()
}

// 检查连接状态
func checkConnection(client *ethclient.Client) {
	ctx := context.Background()

	// 检查是否能获取最新区块号
	blockNumber, err := client.BlockNumber(ctx)
	if err != nil {
		log.Fatal("Connection test failed:", err)
	}

	fmt.Printf("Connection OK! Latest block: %d\n", blockNumber)

	// 检查网络同步状态
	syncProgress, err := client.SyncProgress(ctx)
	if err != nil {
		log.Fatal("Failed to get sync progress:", err)
	}

	if syncProgress != nil {
		fmt.Printf("节点同步中: %d/%d blocks\n",
			syncProgress.CurrentBlock, syncProgress.HighestBlock)
	} else {
		fmt.Println("节点已完全同步")
	}
}

// 按区块高度查询
func queryBlockByNumber(client *ethclient.Client, blockNumber *big.Int) {
	ctx := context.Background()

	// 查询指定高度的区块
	block, err := client.BlockByNumber(ctx, blockNumber)
	if err != nil {
		log.Fatal("Failed to get block:", err)
	}

	fmt.Printf("=== Block Information ===\n")
	fmt.Printf("Block Number: %d\n", block.Number().Uint64())
	fmt.Printf("Block Hash: %s\n", block.Hash().Hex())
	fmt.Printf("Parent Hash: %s\n", block.ParentHash().Hex())
	fmt.Printf("Timestamp: %d\n", block.Time())
	fmt.Printf("Difficulty: %s\n", block.Difficulty().String())
	fmt.Printf("Gas Limit: %d\n", block.GasLimit())
	fmt.Printf("Gas Used: %d\n", block.GasUsed())
	fmt.Printf("Miner: %s\n", block.Coinbase().Hex())
	fmt.Printf("Transaction Count: %d\n", len(block.Transactions()))
	fmt.Printf("Uncle Count: %d\n", len(block.Uncles()))
	fmt.Printf("Size: %d bytes\n", block.Size())
}

// 按区块哈希查询
func queryBlockByHash(client *ethclient.Client, hashStr string) {
	ctx := context.Background()

	// 将字符串转换为哈希
	blockHash := common.HexToHash(hashStr)

	// 按哈希查询区块
	block, err := client.BlockByHash(ctx, blockHash)
	if err != nil {
		log.Fatal("Failed to get block by hash:", err)
	}

	fmt.Printf("=== Block by Hash ===\n")
	fmt.Printf("Block Hash: %s\n", block.Hash().Hex())
	fmt.Printf("Block Number: %d\n", block.Number().Uint64())
	fmt.Printf("Timestamp: %d\n", block.Time())
	fmt.Printf("Transaction Count: %d\n", len(block.Transactions()))

	// 列出所有交易哈希
	fmt.Println("Transactions:")
	for i, tx := range block.Transactions() {
		fmt.Printf("  %d: %s\n", i, tx.Hash().Hex())
	}
}
