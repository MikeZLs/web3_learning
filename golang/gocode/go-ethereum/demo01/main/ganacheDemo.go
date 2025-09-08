package main

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	//// 连接到本地 Ganache
	//client, err := ethclient.Dial("http://127.0.0.1:7545")
	//if err != nil {
	//	log.Fatal("Failed to connect to Ganache:", err)
	//}
	//defer client.Close()
	//
	//fmt.Println("Successfully connected to Ethereum network!")
	//
	//// 测试连接 - 获取网络 ID
	//ChainID, err := client.ChainID(context.Background())
	//if err != nil {
	//	log.Fatal("Failed to get chainID ID:", err)
	//}
	//
	//fmt.Printf("Connected to hainID ID: %s\n", ChainID.String())
	//
	//// 查看预设账户
	//listGanacheAccounts(client)
	//
	//// 测试基础功能
	//testBasicOperations(client)
}

func listGanacheAccounts(client *ethclient.Client) {
	fmt.Println("\n=== Ganache Accounts ===")

	// Ganache 默认创建的账户地址（示例）
	defaultAccounts := []string{
		"0x708BCCE29EEA98b91419613e592Dc9C4640a59e0",
		"0x0FF0879FEAd84EB3572FEdfA3d51913Fd1c324DC",
		"0x766009e0c9F10e8A7397834dC5523579C630d91f",
		"0x5cf9118CA06017598cDECD207a67d580a807F02c",
		"0x6c4636ce17bc7BF172eD99d079B88a8259B61A3a",
	}

	ctx := context.Background()

	for i, addressStr := range defaultAccounts {
		address := common.HexToAddress(addressStr)
		balance, err := client.BalanceAt(ctx, address, nil)
		if err != nil {
			fmt.Println(err)
			fmt.Printf("Account %d: %s - Error getting balance\n", i, addressStr)
			continue
		}

		// 转换为 ETH
		ethBalance := new(big.Float).Quo(new(big.Float).SetInt(balance), big.NewFloat(1e18))
		fmt.Printf("Account %d: %s - Balance: %s ETH\n", i, addressStr, ethBalance.Text('f', 2))
	}
}

func testBasicOperations(client *ethclient.Client) {
	ctx := context.Background()

	fmt.Println("\n=== Basic Operations Test ===")

	// 获取最新区块
	block, err := client.BlockByNumber(ctx, nil) // 返回一个完整的 区块对象，包含区块头、交易列表、时间戳、gas 信息等等(第二个1参数传 nil，表示获取最新区块)
	if err != nil {
		log.Printf("Failed to get block number: %v", err)
		return
	}
	fmt.Printf("Latest Block: %d\n", block)

	blockNumber, err := client.BlockNumber(ctx) // 只返回 最新区块号
	if err != nil {
		log.Printf("Failed to get block number: %v", err)
		return
	}
	fmt.Printf("Latest Block: %d\n", blockNumber)

	// 获取ChainID
	chainID, err := client.ChainID(ctx)
	if err != nil {
		log.Printf("Failed to get ChainID ID: %v", err)
		return
	}
	fmt.Printf("ChainID: %s\n", chainID)

	// 获取 Gas 价格
	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		log.Printf("Failed to get gas price: %v", err)
		return
	}
	gwei := new(big.Float).Quo(new(big.Float).SetInt(gasPrice), big.NewFloat(1e9))
	fmt.Printf("Gas Price: %s Gwei\n", gwei.Text('f', 2))
}
