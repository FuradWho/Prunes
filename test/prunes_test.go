package test

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/FuradWho/prunes/router"
	"github.com/FuradWho/prunes/uniswap"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"io/ioutil"
	"log"
	"math/big"
	"testing"
)

var client *ethclient.Client

func init() {
	var err error
	client, err = ethclient.Dial("https://eth-rinkeby.alchemyapi.io/v2/rI1IyR1NG8iHx-1RLlxPTh0FO1rCRzr3")

	// https://eth-rinkeby.alchemyapi.io/v2/rI1IyR1NG8iHx-1RLlxPTh0FO1rCRzr3
	// https://eth-mainnet.alchemyapi.io/v2/_VdQM-WC9xSSBwTDjV8SKjx2W5OjpOOi
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

func Test_Swap(t *testing.T) {

	contractAddr := common.HexToAddress("0xB4e16d0168e52d35CaCD2c6185b44281Ec28C9Dc")

	newUniswap, err := uniswap.NewUniswap(contractAddr, client)
	if err != nil {
		fmt.Println("new uniswap: ", err)
		return
	}

	nonce, err := client.NonceAt(context.Background(), common.HexToAddress("0x35d7D53cA0b9E122AE65E273cFE10F4a8f04454B"), nil)
	if err != nil {
		fmt.Println("get nonce: ", err)
	}

	privateKey, err := crypto.HexToECDSA("")
	if err != nil {
		log.Fatal(err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(1))
	if err != nil {
		fmt.Println("new auth: ", err)
		return
	}

	//newUniswapTransactor, err := uniswap.NewUniswapTransactor(contractAddr, client)
	//if err != nil {
	//	return
	//}

	session := uniswap.UniswapSession{
		Contract: newUniswap,
		CallOpts: bind.CallOpts{
			Pending: true,
		},
		TransactOpts: bind.TransactOpts{
			From:   auth.From,
			Signer: auth.Signer,
			Nonce:  big.NewInt(int64(nonce)),
		},
	}
	//name, err := session.
	//if err != nil {
	//	fmt.Println("get name: ", err)
	//	return
	//}

	swap, err := session.Swap(big.NewInt(1), big.NewInt(1), auth.From, nil)
	if err != nil {
		fmt.Println("swap err: ", err)
		return
	}
	fmt.Println(swap.Nonce())

	//fmt.Println(name)
	//
	//token0, _ := session.Token0()
	//token1, _ := session.Token1()
	//
	//approve0, err := session.Approve(token0, big.NewInt(100))
	//if err != nil {
	//	fmt.Println("approve token0: ", err)
	//	return
	//}
	//approve1, err := session.Approve(token1, big.NewInt(100))
	//if err != nil {
	//	fmt.Println("approve token1: ", err)
	//	return
	//}
	//
	//fmt.Println(approve1, approve0)

	//tx, err := session.Swap(big.NewInt(1), big.NewInt(1), auth.From, nil)
	//if err != nil {
	//	fmt.Println("tx: ", err)
	//	return
	//}
	//fmt.Println(tx.Hash())

	//transaction, err := newUniswap.Swap(&bind.TransactOpts{
	//	From:    auth.From,
	//	Signer:  auth.Signer,
	//	Nonce:   big.NewInt(int64(nonce)),
	//	Context: context.Background(),
	//}, big.NewInt(1), big.NewInt(1), auth.From, nil)
	//if err != nil {
	//	fmt.Println("transaction : ", err)
	//	return
	//}
	//fmt.Println(transaction.Hash())

	//
	//newUniswap.Token0()

	//initialize, err := newUniswapTransactor.Initialize(&bind.TransactOpts{
	//	From:    auth.From,
	//	Signer:  auth.Signer,
	//	Nonce:   big.NewInt(int64(nonce)),
	//	Context: context.Background(),
	//}, common.HexToAddress("0x956F47F50A910163D8BF957Cf5846D573E7f87CA"), common.HexToAddress("0xc7283b66Eb1EB5FB86327f08e1B5816b0720212B"))
	//if err != nil {
	//	fmt.Println("initialize: ", err)
	//	return
	//}
	//fmt.Println(initialize.Nonce())
	//fmt.Println(auth.From.Hex())
	//transaction, err := newUniswapTransactor.Transfer(&bind.TransactOpts{
	//	From:    auth.From,
	//	Signer:  auth.Signer,
	//	Nonce:   big.NewInt(int64(nonce)),
	//	Context: context.Background(),
	//}, auth.From, big.NewInt(1))
	//if err != nil {
	//	fmt.Println("transaction : ", err)
	//	return
	//}
	//fmt.Println(transaction.Hash())

}

func Test_Router(t *testing.T) {
	contractAddr := common.HexToAddress("0x7a250d5630B4cF539739dF2C5dAcb4c659F2488D")
	newRouter, err := router.NewRouter(contractAddr, client)
	if err != nil {
		fmt.Println("new router: ", err)
		return
	}

	privateKey, err := crypto.HexToECDSA("")
	if err != nil {
		log.Fatal(err)
	}

	// auth := bind.NewKeyedTransactor(privateKey)

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(1))
	if err != nil {
		fmt.Println("new auth: ", err)
		return
	}

	nonce, err := client.NonceAt(context.Background(), common.HexToAddress("0x35d7D53cA0b9E122AE65E273cFE10F4a8f04454B"), nil)
	if err != nil {
		fmt.Println("get nonce: ", err)
		return
	}

	session := router.RouterSession{
		Contract: newRouter,
		CallOpts: bind.CallOpts{
			Pending: true,
		},
		TransactOpts: bind.TransactOpts{
			From:   auth.From,
			Signer: auth.Signer,
			Nonce:  big.NewInt(int64(nonce)),
		},
	}

	fmt.Println(session.Receive())
}
func Test_RouterV2(t *testing.T) {

	contractAddr := common.HexToAddress("0x7a250d5630B4cF539739dF2C5dAcb4c659F2488D")

	newRouter, err := router.NewRouter(contractAddr, client)
	if err != nil {
		t.Error("new router : ", err)
		return
	}
	privateKey, err := crypto.HexToECDSA("")
	if err != nil {
		t.Error(err)
	}

	// auth := bind.NewKeyedTransactor(privateKey)

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(4))
	if err != nil {
		t.Error("new auth: ", err)
		return
	}

	nonce, err := client.NonceAt(context.Background(), common.HexToAddress("0x35d7D53cA0b9E122AE65E273cFE10F4a8f04454B"), nil)
	if err != nil {
		t.Error("get nonce: ", err)
		return
	}

	session := router.RouterSession{
		Contract: newRouter,
		CallOpts: bind.CallOpts{
			Pending: true,
		},
		TransactOpts: bind.TransactOpts{
			From:   auth.From,
			Signer: auth.Signer,
			Nonce:  big.NewInt(int64(nonce)),
		},
	}

	fmt.Println(session.Receive())

}
