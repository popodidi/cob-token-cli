package private

import (
	"github.com/urfave/cli"
	"math/big"
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/popodidi/cob-token-cli/utils"
	"errors"
)

func sendCOBAction(c *cli.Context) error {
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

	var cobValueString string
	cobValueString, err = utils.AskForString("COB Value")
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	var cobAmount *big.Int
	cobAmount, err = utils.StringToWei(cobValueString)
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
	tx, err = utils.SendCOB(privateKey, toAddress, cobAmount, big.NewInt(510000), gasPrice)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	fmt.Println("Transaction sent\nTX HASH: ", tx.Hash().Hex())
	return nil
}
