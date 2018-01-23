package utils

import (
	"bytes"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"
	"github.com/shopspring/decimal"
	"golang.org/x/net/context"
)

const erc20ABI string = `[{"constant":false,"inputs":[{"name":"spender","type":"address"},{"name":"value","type":"uint256"}],"name":"approve","outputs":[{"name":"","type":"bool"}],"payable":false,"type":"function"},{"constant":true,"inputs":[],"name":"totalSupply","outputs":[{"name":"","type":"uint256"}],"payable":false,"type":"function"},{"constant":false,"inputs":[{"name":"from","type":"address"},{"name":"to","type":"address"},{"name":"value","type":"uint256"}],"name":"transferFrom","outputs":[{"name":"","type":"bool"}],"payable":false,"type":"function"},{"constant":true,"inputs":[{"name":"who","type":"address"}],"name":"balanceOf","outputs":[{"name":"","type":"uint256"}],"payable":false,"type":"function"},{"constant":false,"inputs":[{"name":"to","type":"address"},{"name":"value","type":"uint256"}],"name":"transfer","outputs":[{"name":"","type":"bool"}],"payable":false,"type":"function"},{"constant":true,"inputs":[{"name":"owner","type":"address"},{"name":"spender","type":"address"}],"name":"allowance","outputs":[{"name":"","type":"uint256"}],"payable":false,"type":"function"},{"anonymous":false,"inputs":[{"indexed":true,"name":"owner","type":"address"},{"indexed":true,"name":"spender","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"Approval","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"from","type":"address"},{"indexed":true,"name":"to","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"Transfer","type":"event"}]`
const cobAddressHex string = "0xb2F7EB1f2c37645bE61d73953035360e768D81E6"

func NewClient() (*ethclient.Client, error) {
	return ethclient.Dial("https://mainnet.infura.io")
}

func StringToWei(s string) (*big.Int, error) {
	d, err := decimal.NewFromString(s)
	if err != nil {
		return nil, err
	}
	d = d.Mul(decimal.New(1, 18))
	d = d.Truncate(0)

	tenPower := big.NewInt(0)
	exp := big.NewInt(int64(d.Exponent()))
	tenPower.Exp(big.NewInt(10), exp, big.NewInt(0))

	var wei = big.NewInt(0)
	wei.Mul(d.Coefficient(), tenPower)

	return wei, nil
}

func GetEthBalanceOf(address string) (*decimal.Decimal, error) {
	client, err := NewClient()
	if err != nil {
		return nil, err
	}
	ctx, _ := context.WithTimeout(context.Background(), time.Minute*1)
	addr := common.HexToAddress(address)

	var balance *big.Int
	balance, err = client.BalanceAt(ctx, addr, nil)
	if err != nil {
		return nil, err
	}

	balanceDecimal := decimal.NewFromBigInt(balance, -18)
	return &balanceDecimal, nil
}

func GetCobBalanceOf(address string) (*decimal.Decimal, error) {
	client, err := NewClient()
	if err != nil {
		return nil, err
	}
	ctx, _ := context.WithTimeout(context.Background(), time.Minute*1)

	// Compose ERC20 method ABI
	abi, err := abi.JSON(bytes.NewReader([]byte(erc20ABI)))
	if err != nil {
		return nil, err
	}

	var data []byte
	data, err = abi.Pack("balanceOf", common.HexToAddress(address))
	if err != nil {
		return nil, err
	}

	// Compose smart contract call transaction message
	cobAddress := common.HexToAddress(cobAddressHex)
	contractMessage := ethereum.CallMsg{
		From:     common.HexToAddress(address),
		To:       &cobAddress,
		Gas:      nil,
		GasPrice: nil,
		Value:    nil,
		Data:     data,
	}

	// Call smart contract method
	var balanceBytes []byte
	balanceBytes, err = client.CallContract(ctx, contractMessage, nil)
	if err != nil {
		return nil, err
	}

	// Generate balance from result
	balance := big.NewInt(0).SetBytes(balanceBytes)
	balanceDecimal := decimal.NewFromBigInt(balance, -18)

	return &balanceDecimal, nil
}

func SendETH(fromPrivKey string, toAddress string, amount, gasLimit, gasPrice *big.Int) (*types.Transaction, error) {
	// validate toAddress
	err := ValidateAddress(toAddress)
	if err != nil {
		return nil, err
	}

	client, err := NewClient()
	if err != nil {
		return nil, err
	}

	var fromPrivECDSAKey *ecdsa.PrivateKey
	fromPrivECDSAKey, err = crypto.HexToECDSA(fromPrivKey)
	if err != nil {
		return nil, err
	}
	fromPubECDSAKey := fromPrivECDSAKey.PublicKey
	fromAddress := crypto.PubkeyToAddress(fromPubECDSAKey)

	ctx, _ := context.WithTimeout(context.Background(), 1*time.Minute)

	var nonce uint64
	nonce, err = client.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		return nil, err
	}

	signer := types.NewEIP155Signer(params.MainnetChainConfig.ChainId)
	tx := types.NewTransaction(nonce, common.HexToAddress(toAddress), amount, gasLimit, gasPrice, common.FromHex("0x"))

	var signedTx *types.Transaction
	signedTx, _ = types.SignTx(tx, signer, fromPrivECDSAKey)

	err = client.SendTransaction(ctx, signedTx)
	if err != nil {
		return nil, err
	}

	return waitTxMined(client, signedTx.Hash())
}

func SendCOB(fromPrivKey string, toAddress string, amount, gasLimit, gasPrice *big.Int) (*types.Transaction, error) {
	// validate toAddress
	err := ValidateAddress(toAddress)
	if err != nil {
		return nil, err
	}

	client, err := NewClient()
	if err != nil {
		return nil, err
	}

	var fromPrivECDSAKey *ecdsa.PrivateKey
	fromPrivECDSAKey, err = crypto.HexToECDSA(fromPrivKey)
	if err != nil {
		return nil, err
	}
	fromPubECDSAKey := fromPrivECDSAKey.PublicKey
	fromAddress := crypto.PubkeyToAddress(fromPubECDSAKey)

	ctx, _ := context.WithTimeout(context.Background(), 1*time.Minute)

	var nonce uint64
	nonce, err = client.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		return nil, err
	}

	// Compose ERC20 method ABI
	abi, err := abi.JSON(bytes.NewReader([]byte(erc20ABI)))
	if err != nil {
		return nil, err
	}
	txData, _ := abi.Pack("transfer", common.HexToAddress(toAddress), amount)

	signer := types.NewEIP155Signer(params.MainnetChainConfig.ChainId)
	tx := types.NewTransaction(nonce, common.HexToAddress(cobAddressHex), big.NewInt(0), gasLimit, gasPrice, txData)

	var signedTx *types.Transaction
	signedTx, _ = types.SignTx(tx, signer, fromPrivECDSAKey)

	err = client.SendTransaction(ctx, signedTx)
	if err != nil {
		return nil, err
	}

	time.Sleep(time.Second * 5)
	return waitTxMined(client, signedTx.Hash())
}

func ValidateAddress(addressStr string) error {
	addr := common.HexToAddress(addressStr)
	if addr.Hex() == addressStr {
		return nil
	}

	return fmt.Errorf("invalid address. \"%s\"", addressStr)
}

func waitTxMined(client *ethclient.Client, txHash common.Hash) (*types.Transaction, error) {
	done := false
	for !done {
		ctx, _ := context.WithTimeout(context.Background(), 1*time.Minute)
		tx, pending, err := client.TransactionByHash(ctx, txHash)
		if err != nil {
			done = true
			return nil, err
		} else if !pending {
			done = true
			return tx, nil
		}
		time.Sleep(time.Second)
	}
	return nil, errors.New("unexpected error")
}
