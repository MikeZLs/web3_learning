package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
)

func queryTransactionReceipt(client *ethclient.Client, txHashStr string) {
	ctx := context.Background()

	txHash := common.HexToHash(txHashStr)

	// 查询交易收据
	receipt, err := client.TransactionReceipt(ctx, txHash)
	if err != nil {
		log.Fatal("Failed to get transaction receipt:", err)
	}

	fmt.Printf("=== Transaction Receipt ===\n")
	fmt.Printf("Transaction Hash: %s\n", receipt.TxHash.Hex())
	fmt.Printf("Block Number: %d\n", receipt.BlockNumber.Uint64())
	fmt.Printf("Block Hash: %s\n", receipt.BlockHash.Hex())
	fmt.Printf("Transaction Index: %d\n", receipt.TransactionIndex)

	// 交易状态
	status := "Failed"
	if receipt.Status == 1 {
		status = "Success"
	}
	fmt.Printf("Status: %s\n", status)

	// Gas 使用情况
	fmt.Printf("Gas Used: %d\n", receipt.GasUsed)

	// 累计 Gas 使用（区块中到此交易为止的总 Gas）
	fmt.Printf("Cumulative Gas Used: %d\n", receipt.CumulativeGasUsed)

	// 合约地址（如果是合约创建交易）
	if receipt.ContractAddress != (common.Address{}) {
		fmt.Printf("Contract Address: %s\n", receipt.ContractAddress.Hex())
	}

	// Logs 和事件
	fmt.Printf("Log Entries: %d\n", len(receipt.Logs))

	if len(receipt.Logs) > 0 {
		fmt.Println("Logs:")
		for i, log := range receipt.Logs {
			fmt.Printf("  Log #%d:\n", i+1)
			fmt.Printf("    Address: %s\n", log.Address.Hex())
			fmt.Printf("    Topics: %d\n", len(log.Topics))
			for j, topic := range log.Topics {
				fmt.Printf("      Topic %d: %s\n", j, topic.Hex())
			}
			fmt.Printf("    Data: %s\n", common.Bytes2Hex(log.Data))
		}
	}

	// 计算实际 Gas 费用
	tx, _, err := client.TransactionByHash(ctx, txHash)
	if err == nil {
		gasFee := new(big.Int).Mul(tx.GasPrice(), big.NewInt(int64(receipt.GasUsed)))
		fmt.Printf("Gas Fee Paid: %s ETH\n", weiToEth(gasFee))

		// Gas 效率
		gasEfficiency := float64(receipt.GasUsed) / float64(tx.Gas()) * 100
		fmt.Printf("Gas Efficiency: %.2f%% (%d used / %d limit)\n",
			gasEfficiency, receipt.GasUsed, tx.Gas())
	}
}
