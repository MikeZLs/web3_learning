package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
)

func queryBlockTransactions(client *ethclient.Client, blockNumber *big.Int) {
	ctx := context.Background()

	// 获取区块
	block, err := client.BlockByNumber(ctx, blockNumber)
	if err != nil {
		log.Fatal("Failed to get block:", err)
	}

	fmt.Printf("=== Block #%d Transactions ===\n", block.Number().Uint64())
	fmt.Printf("Total Transactions: %d\n\n", len(block.Transactions()))

	totalValue := big.NewInt(0)
	contractCreations := 0

	for i, tx := range block.Transactions() {
		fmt.Printf("Transaction #%d\n", i+1)
		fmt.Printf("  Hash: %s\n", tx.Hash().Hex())
		fmt.Printf("  From: %s\n", getTransactionSender(client, tx))

		if tx.To() == nil {
			fmt.Printf("  To: Contract Creation\n")
			contractCreations++
		} else {
			fmt.Printf("  To: %s\n", tx.To().Hex())
		}

		fmt.Printf("  Value: %s ETH\n", weiToEth(tx.Value()))
		fmt.Printf("  Gas Price: %s Gwei\n", weiToGwei(tx.GasPrice()))

		totalValue.Add(totalValue, tx.Value())
		fmt.Println()
	}

	fmt.Printf("=== Summary ===\n")
	fmt.Printf("Total Value Transferred: %s ETH\n", weiToEth(totalValue))
	fmt.Printf("Contract Creations: %d\n", contractCreations)
	fmt.Printf("Regular Transfers: %d\n", len(block.Transactions())-contractCreations)
}
