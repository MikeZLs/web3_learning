package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
)

func analyzeReceiptsBatch(client *ethclient.Client, txHashes []string) {
	ctx := context.Background()

	var (
		successCount = 0
		failedCount  = 0
		totalGasUsed = uint64(0)
		totalGasFees = big.NewInt(0)
	)

	fmt.Printf("=== Batch Receipt Analysis ===\n")
	fmt.Printf("Analyzing %d transactions...\n\n", len(txHashes))

	for i, hashStr := range txHashes {
		txHash := common.HexToHash(hashStr)

		receipt, err := client.TransactionReceipt(ctx, txHash)
		if err != nil {
			fmt.Printf("Transaction %d: Failed to get receipt - %v\n", i+1, err)
			continue
		}

		// 获取原始交易信息
		tx, _, err := client.TransactionByHash(ctx, txHash)
		if err != nil {
			fmt.Printf("Transaction %d: Failed to get transaction - %v\n", i+1, err)
			continue
		}

		// 统计
		if receipt.Status == 1 {
			successCount++
		} else {
			failedCount++
		}

		totalGasUsed += receipt.GasUsed
		gasFee := new(big.Int).Mul(tx.GasPrice(), big.NewInt(int64(receipt.GasUsed)))
		totalGasFees.Add(totalGasFees, gasFee)

		// 显示摘要
		status := "✓"
		if receipt.Status != 1 {
			status = "✗"
		}

		fmt.Printf("%s Tx %d: Gas %d, Fee %s ETH\n",
			status, i+1, receipt.GasUsed, weiToEth(gasFee))
	}

	fmt.Printf("\n=== Summary ===\n")
	fmt.Printf("Successful: %d\n", successCount)
	fmt.Printf("Failed: %d\n", failedCount)
	fmt.Printf("Total Gas Used: %d\n", totalGasUsed)
	fmt.Printf("Total Fees Paid: %s ETH\n", weiToEth(totalGasFees))

	if len(txHashes) > 0 {
		avgGas := totalGasUsed / uint64(len(txHashes))
		fmt.Printf("Average Gas per Transaction: %d\n", avgGas)
	}
}
