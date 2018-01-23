package private

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/popodidi/cob-token-cli/utils"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/urfave/cli"
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
	err = utils.ValidateAddress(toAddress)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	var ethValueString string
	ethValueString, err = utils.AskForString("ETH Value")
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	var ethAmount *big.Int
	ethAmount, err = utils.StringToWei(ethValueString)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

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
