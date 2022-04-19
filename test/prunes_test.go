package test

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/FuradWho/prunes/factory"
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
	"time"
)

var client *ethclient.Client

var pairAddr common.Address
var token0 common.Address
var token1 common.Address

const (
	accountAddr   = "0x1d32A826a2Bf24C4d4CD7C8b5c72faC83C8ec3B3"
	factoryV2Addr = "0x5C69bEe701ef814a2B6a3EDD4B1652CB9cc5aA6f"
	routerV2Addr  = "0x7a250d5630B4cF539739dF2C5dAcb4c659F2488D"
	rinkebyAddr   = ""
	mainnetAddr   = ""
	key           = ""
)

func init() {
	var err error
	client, err = ethclient.Dial(rinkebyAddr)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("we have a connection")
}

func Test_FactoryV2(t *testing.T) {

	// 获取合约地址
	contractAddr := common.HexToAddress(factoryV2Addr)

	newFactory, err := factory.NewFactory(contractAddr, client)
	if err != nil {
		t.Error("new factory: ", err)
		return
	}

	privateKey, err := crypto.HexToECDSA(key)
	if err != nil {
		t.Error(err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(4))
	if err != nil {
		t.Error("new auth: ", err)
		return
	}

	session := factory.FactorySession{
		Contract: newFactory,
		CallOpts: bind.CallOpts{
			Pending: true,
		},
		TransactOpts: bind.TransactOpts{
			From:   auth.From,
			Signer: auth.Signer,
		},
	}

	// Returns the address of the nth pair (0-indexed) created through the factory, or address(0) (0x0000000000000000000000000000000000000000) if not enough pairs have been created yet.
	pairAddr, err = session.AllPairs(big.NewInt(1))
	if err != nil {
		t.Error("get all pair: ", err)
		return
	}

	fmt.Println("all pairs: ", pairAddr)

	// Returns the total number of pairs created through the factory so far.
	length, err := session.AllPairsLength()
	if err != nil {
		t.Error("get all pair len: ", err)
		return
	}
	fmt.Println("get all pair len: ", length)

}

func Test_UniSwap(t *testing.T) {

	contractAddr := common.HexToAddress("0x80f07c368BCC7F8CbAC37E79Ec332c1D84e9335D")

	newUniswap, err := uniswap.NewUniswap(contractAddr, client)
	if err != nil {
		t.Error("new uniSwap: ", err)
		return
	}

	privateKey, err := crypto.HexToECDSA(key)
	if err != nil {
		t.Error(err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(4))
	if err != nil {
		t.Error("new auth: ", err)
		return
	}

	session := uniswap.UniswapSession{
		Contract: newUniswap,
		CallOpts: bind.CallOpts{
			Pending: true,
		},
		TransactOpts: bind.TransactOpts{
			From:   auth.From,
			Signer: auth.Signer,
		},
	}

	token0, err = session.Token0()
	if err != nil {
		t.Error("get Token0: ", err)
		return
	}
	fmt.Println("token0 addr: ", token0)

	token1, err = session.Token1()
	if err != nil {
		t.Error("get Token1: ", err)
		return
	}
	fmt.Println("token1 addr: ", token1)

}

func Test_RouterV2(t *testing.T) {
	contractAddr := common.HexToAddress(routerV2Addr)

	newRouter, err := router.NewRouter(contractAddr, client)
	if err != nil {
		t.Error("new router: ", err)
		return
	}

	privateKey, err := crypto.HexToECDSA(key)
	if err != nil {
		t.Error(err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(4))
	if err != nil {
		t.Error("new auth: ", err)
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
			Value:  big.NewInt(200),
		},
	}
	// Returns the canonical WETH address on the Ethereum mainnet, or the Ropsten, Rinkeby, Görli, or Kovan testnets.
	wethAddr, _ := session.WETH()
	fmt.Println("weth addr: ", wethAddr)

	/*
		token0 addr:  0xc778417E063141139Fce010982780140Aa0cD5Ab
		token1 addr:  0xF9bA5210F91D0474bd1e1DcDAeC4C58E359AaD85
	*/

	var addrs []common.Address
	addrs = append(addrs,
		common.HexToAddress("0xc778417E063141139Fce010982780140Aa0cD5Ab"),
		common.HexToAddress("0xF9bA5210F91D0474bd1e1DcDAeC4C58E359AaD85"),
	)

	out, err := session.GetAmountsIn(big.NewInt(1), addrs)
	if err != nil {
		t.Error("get amount out: ", err)
		return
	}
	fmt.Println(out)

	transaction, err := session.SwapETHForExactTokens(big.NewInt(1), addrs, common.HexToAddress(accountAddr), big.NewInt(time.Now().Unix()+100))
	if err != nil {
		t.Error("swap eth :", err)
		return
	}
	fmt.Println(transaction.Hash())
}

func Test_BalanceAt(t *testing.T) {

	account := common.HexToAddress(accountAddr)

	balance, err := client.BalanceAt(context.Background(), account, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(balance) // 25893180161173005034
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
