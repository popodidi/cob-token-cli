package private

import (
	"github.com/urfave/cli"
	"github.com/shopspring/decimal"
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

	var cobFloat float64
	cobFloat, err = utils.AskForFloat("COB Value")
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	cobDecimal := decimal.NewFromFloat(cobFloat)
	cobDecimal = cobDecimal.Mul(decimal.New(1, 18))
	cobAmount := big.NewInt(cobDecimal.IntPart())

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
