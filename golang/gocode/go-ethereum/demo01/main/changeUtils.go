package main

import "math/big"

// 辅助函数：Wei 转 ETH
func WeiToEth(wei *big.Int) string {
	eth := new(big.Float).Quo(new(big.Float).SetInt(wei), big.NewFloat(1e18))
	return eth.Text('f', 6)
}

// 辅助函数：Wei 转 Gwei
func WeiToGwei(wei *big.Int) string {
	gwei := new(big.Float).Quo(new(big.Float).SetInt(wei), big.NewFloat(1e9))
	return gwei.Text('f', 2)
}
