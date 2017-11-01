package public

import (
	"math/big"
	"github.com/urfave/cli"
	"fmt"
	"context"
	"github.com/popodidi/cob-token-cli/utils"
	"github.com/ethereum/go-ethereum/common"
	"time"
	"github.com/shopspring/decimal"
	"gopkg.in/AlecAivazis/survey.v1"
)

func ethBalanceAction(c *cli.Context) error {
	address := ""
	addressPrompt := &survey.Input{
		Message: "ETH address",
	}
	survey.AskOne(addressPrompt, &address, nil)
	balance, err := getEthBalanceOf(address)
	if err != nil {
		return err
	}
	fmt.Println(balance, "ETHs")
	return nil
}

func getEthBalanceOf(address string) (*decimal.Decimal, error) {
	client, err := utils.NewClient()
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
