package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// ETh 转账
func main() {
	client, err := ethclient.Dial("HTTP://127.0.0.1:7545")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// 自动移除 0x 前缀
	privateKeyHex := "0xd239eade3f9c50148ac30354ea72f9dd34dded2444b2f6de1fa12b1464e760bf"
	privateKeyHex = strings.TrimPrefix(privateKeyHex, "0x")

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	value := big.NewInt(1e17) // in wei (1 eth)
	gasLimit := uint64(21000) // in units
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	toAddress := common.HexToAddress("0x9e0A541dA6CD5F45B93f6B8d1527E7ee9a92108a")
	var data []byte
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)

	chainID, err := client.ChainID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("fromAddress:%s 已发送到 toAddress:%s %s个ETH \n", fromAddress, toAddress, weiToEth(value, 18))

	fmt.Printf("tx sent: %s", signedTx.Hash().Hex())
}

func weiToEth(wei *big.Int, decimals int) string {
	return new(big.Float).
		SetInt(wei).
		Quo(new(big.Float).SetInt(wei), new(big.Float).SetInt(big.NewInt(1e18))).
		Text('f', decimals)
}
