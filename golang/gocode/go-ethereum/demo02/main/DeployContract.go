package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	client, err := ethclient.Dial("http://127.0.0.1:7545")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	privateKey, err := crypto.HexToECDSA("d239eade3f9c50148ac30354ea72f9dd34dded2444b2f6de1fa12b1464e760bf")
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA := publicKey.(*ecdsa.PublicKey)
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	// 检查余额
	balance, err := client.BalanceAt(context.Background(), fromAddress, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Balance: %s ETH\n", new(big.Float).Quo(new(big.Float).SetInt(balance), big.NewFloat(1e18)).Text('f', 4))

	// 获取网络信息
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Chain ID: %s\n", chainID.String())

	// 简单的字节码测试 - 只包含最基本的构造函数
	// 这是一个几乎什么都不做的合约字节码
	simpleCode := "0x608060405234801561001057600080fd5b50603f80601e6000396000f3fe6080604052600080fd"

	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// 创建合约创建交易
	tx := types.NewContractCreation(
		nonce,
		big.NewInt(0),
		uint64(3000000),
		gasPrice,
		common.FromHex(simpleCode),
	)

	// 签名交易
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal("Signing failed:", err)
	}

	fmt.Printf("Sending transaction: %s\n", signedTx.Hash().Hex())

	// 发送交易
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal("Send transaction failed:", err)
	}

	// 计算合约地址
	contractAddress := crypto.CreateAddress(fromAddress, nonce)
	fmt.Printf("Contract address: %s\n", contractAddress.Hex())
}
