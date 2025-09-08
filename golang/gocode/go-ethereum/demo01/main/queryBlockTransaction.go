package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
)

// 按交易哈希查询交易数据
func main() {
	client, err := ethclient.Dial("https://eth.w3node.com/8f9e39edde7c2aa1737461a34384908393dd1d4216de22085e9e214fe0f04c73/api")
	if err != nil {
		log.Fatal("Failed to connect to Ethereum client:", err)
	}
	defer client.Close()
	queryTransaction(client, "0x8c66b82c9a5c726c079923149c08055bb692e6f3bbb96490734c5c092ae60740")
}

func queryTransaction(client *ethclient.Client, txHashStr string) {
	ctx := context.Background()

	// 转换交易哈希
	txHash := common.HexToHash(txHashStr)

	// 查询交易
	tx, isPending, err := client.TransactionByHash(ctx, txHash)
	if err != nil {
		log.Fatal("Failed to get transaction:", err)
	}

	fmt.Printf("=== Transaction Information ===\n")
	fmt.Printf("Hash: %s\n", tx.Hash().Hex())
	fmt.Printf("Status: %s\n", map[bool]string{true: "Pending", false: "Confirmed"}[isPending])
	fmt.Printf("Block Number: %s\n", func() string {
		if isPending {
			return "Pending"
		}
		// 获取交易所在区块
		receipt, err := client.TransactionReceipt(ctx, txHash)
		if err != nil {
			return "Unknown"
		}
		return fmt.Sprintf("%d", receipt.BlockNumber.Uint64())
	}())

	fmt.Printf("From: %s\n", getTransactionSender(client, tx))
	fmt.Printf("To: %s\n", func() string {
		if tx.To() == nil {
			return "Contract Creation"
		}
		return tx.To().Hex()
	}())

	// 转换 Wei 到 ETH
	ethValue := WeiToEth(tx.Value())
	fmt.Printf("Value: %s ETH\n", ethValue)

	fmt.Printf("Gas Limit: %d\n", tx.Gas())
	fmt.Printf("Gas Price: %s Gwei\n", WeiToGwei(tx.GasPrice()))
	fmt.Printf("Nonce: %d\n", tx.Nonce())

	// 数据字段
	if len(tx.Data()) > 0 {
		fmt.Printf("Input Data: %s\n", common.Bytes2Hex(tx.Data()))
		fmt.Printf("Data Size: %d bytes\n", len(tx.Data()))
	} else {
		fmt.Println("Input Data: None (Simple Transfer)")
	}
}

// 辅助函数：获取交易发送者
func getTransactionSender(client *ethclient.Client, tx *types.Transaction) string {
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		return "Unknown"
	}

	signer := types.NewEIP155Signer(chainID)
	sender, err := types.Sender(signer, tx)
	if err != nil {
		return "Unknown"
	}

	return sender.Hex()
}
