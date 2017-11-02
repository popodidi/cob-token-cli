package private

import (
	"github.com/urfave/cli"
	"github.com/popodidi/cob-token-cli/utils"
	"github.com/shopspring/decimal"
	"math/big"
	"github.com/ethereum/go-ethereum/core/types"
	"fmt"
	"errors"
)

func sendETHAction(c *cli.Context) error {
	privateKey, err := utils.AskForPrivateKey()
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	var toAddress string
	toAddress, err = utils.AskForString("To address")
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	var ethFloat float64
	ethFloat, err = utils.AskForFloat("ETH Value")
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	ethDecimal := decimal.NewFromFloat(ethFloat)
	ethDecimal = ethDecimal.Mul(decimal.New(1, 18))
	ethAmount := big.NewInt(ethDecimal.IntPart())

	var gasPrice *big.Int
	gasPrice, err = utils.AskForGasPriceGwei()
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	if !utils.AskForConfirm("START") {
		return cli.NewExitError(errors.New("user stopped"), 1)
	}

	var tx *types.Transaction
	tx, err = utils.SendETH(privateKey, toAddress, ethAmount, big.NewInt(21000), gasPrice)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	fmt.Println("Transaction sent\nTX HASH: ", tx.Hash().Hex())
	return nil
}
