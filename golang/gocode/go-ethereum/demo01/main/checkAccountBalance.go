package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
)

func main() {
	client, err := ethclient.Dial("HTTP://127.0.0.1:7545")
	if err != nil {
		log.Fatal(err)
	}

	account := common.HexToAddress("0xD760Dc7A2F919bfF1FDfCDc887549bB744d989A5")

	// 查询最新区块的余额
	balance, err := client.BalanceAt(context.Background(), account, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("账户余额 (Wei): %s\n", balance.String())

	ethBalance := WeiToEth(balance)
	fmt.Printf("账户余额 (ETH): %s\n", ethBalance)

	// 待处理的余额
	pendingBalance, err := client.PendingBalanceAt(context.Background(), account)
	fmt.Println(pendingBalance)

}
