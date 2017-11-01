package private

import (
	"github.com/urfave/cli"
	"gopkg.in/AlecAivazis/survey.v1"
	"github.com/shopspring/decimal"
	"math/big"
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/popodidi/cob-token-cli/utils"
)

func sendCOBAction(c *cli.Context) error {
	var qs = []*survey.Question{
		{
			Name:     "from-private-key",
			Prompt:   &survey.Password{Message: "From private key"},
			Validate: survey.Required,
		},
		{
			Name:     "to-address",
			Prompt:   &survey.Input{Message: "To address",},
			Validate: survey.Required,
		},
		{
			Name:     "cob-value",
			Prompt:   &survey.Input{Message: "COB Value"},
			Validate: survey.Required,
		},
		{
			Name:     "gas-price",
			Prompt:   &survey.Input{Message: "Gas Price (Gwei)"},
			Validate: survey.Required,
		},
	}

	// the answers will be written to this struct
	answers := struct {
		FromPrivKey string  `survey:"from-private-key"` // survey will match the question and field names
		ToAddress   string  `survey:"to-address"`       // or you can tag fields to match a specific name
		COBValue    float64 `survey:"cob-value"`        // if the types don't match exactly, survey will try to convert for you
		GasPrice    int64   `survey:"gas-price"`
	}{}

	// perform the questions
	err := survey.Ask(qs, &answers)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	cobValue := decimal.NewFromFloat(answers.COBValue)
	cobValue = cobValue.Mul(decimal.New(1, 18))

	cobAmount := big.NewInt(cobValue.IntPart())

	gasPrice := big.NewInt(1)
	gasPrice.Mul(big.NewInt(answers.GasPrice), big.NewInt(1000000000))

	var tx *types.Transaction
	tx, err = utils.SendCOB(answers.FromPrivKey, answers.ToAddress, cobAmount, big.NewInt(500000), gasPrice)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	fmt.Println("Transaction sent\nTX HASH: ", tx.Hash().Hex())
	return nil
}
