package test

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/FuradWho/prunes/uniswap"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"io/ioutil"
	"log"
	"math/big"
	"testing"
)

var client *ethclient.Client

func init() {
	var err error
	client, err = ethclient.Dial("https://cloudflare-eth.com")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("we have a connection")
}

func Test_BalanceAt(t *testing.T) {

	account := common.HexToAddress("0xC75650fe4D14017b1e12341A97721D5ec51D5340")

	// 0x5C69bEe701ef814a2B6a3EDD4B1652CB9cc5aA6f
	// 0x35d7D53cA0b9E122AE65E273cFE10F4a8f04454B
	balance, err := client.BalanceAt(context.Background(), account, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(balance) // 25893180161173005034

	//address := common.HexToAddress("0x10ED43C718714eb63d5aA57B78B54704E256024E\n")
	//instance, err := store.NewStore(address, client)
	//if err != nil {
	//	log.Fatal(err)
	//}

	newUniswap, err := uniswap.NewUniswap(account, client)
	if err != nil {
		fmt.Println("newUniswap：", err)
		return
	}

	nonce, err := client.NonceAt(context.Background(), common.HexToAddress("0x35d7D53cA0b9E122AE65E273cFE10F4a8f04454B"), nil)
	if err != nil {
		fmt.Println("get nonce: ", err)
	}
	fmt.Println(nonce)
	symbol, err := newUniswap.Symbol(&bind.CallOpts{
		From:        common.HexToAddress("0x35d7D53cA0b9E122AE65E273cFE10F4a8f04454B"),
		Pending:     true,
		BlockNumber: big.NewInt(int64(nonce)),
		Context:     context.Background(),
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(symbol)

	//contractWithAbi, err := callContractWithAbi(client, privateKey, account)
	//if err != nil {
	//	fmt.Println("contractWithAbi :", err)
	//	return
	//}
	//
	//fmt.Println(contractWithAbi)

}

func callContractWithAbi(client *ethclient.Client, privKey *ecdsa.PrivateKey, from common.Address) (string, error) {
	nonce, err := client.NonceAt(context.Background(), from, nil)
	if err != nil {
		fmt.Println("get nonce: ", err)
		return "", err
	}

	_, err = client.SuggestGasPrice(context.Background())
	if err != nil {
		fmt.Println("gas price: ", err)
		return "", err
	}

	abiData, err := ioutil.ReadFile("E:\\projects\\prunes\\abi.txt")
	if err != nil {
		fmt.Println("read file: ", err)
		return "", err
	}

	contractABI, err := abi.JSON(bytes.NewReader(abiData))
	if err != nil {
		fmt.Println("abi json: ", err)
		return "", err
	}

	callData, err := contractABI.Pack("price0CumulativeLast")
	if err != nil {
		fmt.Println("abi pack: ", err)
		return "", err
	}

	tx := types.NewTransaction(nonce, common.HexToAddress("0x5C69bEe701ef814a2B6a3EDD4B1652CB9cc5aA6f"), big.NewInt(0), uint64(0), big.NewInt(0), callData)

	signTx, err := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(1)), privKey)
	if err != nil {
		fmt.Println("sign tx: ", err)
		return "", err
	}

	err = client.SendTransaction(context.Background(), signTx)
	if err != nil {
		fmt.Println("send tx :", err)
		return "", err
	}

	return signTx.Hash().Hex(), nil

}

func uniswap_func() {

}
