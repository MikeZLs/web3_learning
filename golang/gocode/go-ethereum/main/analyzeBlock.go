package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
	"time"
)

func main() {
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

	analyzeBlock(client, big.NewInt(1))
}

func analyzeBlock(client *ethclient.Client, blockNumber *big.Int) {
	ctx := context.Background()

	// 获取区块（包含交易详情）
	block, err := client.BlockByNumber(ctx, blockNumber)
	if err != nil {
		log.Fatal("Failed to get block:", err)
	}

	fmt.Printf("=== Detailed Block Analysis ===\n")
	fmt.Printf("Block #%d\n", block.Number().Uint64())
	fmt.Printf("Hash: %s\n", block.Hash().Hex())

	// 时间戳转换
	timestamp := time.Unix(int64(block.Time()), 0)
	fmt.Printf("Timestamp: %s\n", timestamp.Format("2006-01-02 15:04:05 UTC"))

	// Gas 使用率
	gasUsagePercent := float64(block.GasUsed()) / float64(block.GasLimit()) * 100
	fmt.Printf("Gas Usage: %d/%d (%.2f%%)\n",
		block.GasUsed(), block.GasLimit(), gasUsagePercent)

	// 区块大小
	fmt.Printf("Block Size: %.2f KB\n", float64(block.Size())/1024)

	// 交易统计
	totalValue := big.NewInt(0)
	totalGasFees := big.NewInt(0)

	for _, tx := range block.Transactions() {
		totalValue.Add(totalValue, tx.Value())

		// 计算 gas 费用（需要获取交易收据来得到实际 gas 使用量）
		receipt, err := client.TransactionReceipt(ctx, tx.Hash())
		if err == nil {
			gasFee := new(big.Int).Mul(tx.GasPrice(), big.NewInt(int64(receipt.GasUsed)))
			totalGasFees.Add(totalGasFees, gasFee)
		}
	}

	// 转换为 ETH 单位 (Wei to ETH)
	ethUnit := big.NewFloat(1e18)
	totalValueETH, _ := new(big.Float).Quo(new(big.Float).SetInt(totalValue), ethUnit).Float64()
	totalFeesETH, _ := new(big.Float).Quo(new(big.Float).SetInt(totalGasFees), ethUnit).Float64()

	fmt.Printf("Total Transaction Value: %.6f ETH\n", totalValueETH)
	fmt.Printf("Total Gas Fees: %.6f ETH\n", totalFeesETH)
}
