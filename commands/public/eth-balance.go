package public

import (
	"github.com/urfave/cli"
	"fmt"
	"github.com/popodidi/cob-token-cli/utils"
)

func ethBalanceAction(c *cli.Context) error {
	address, err := utils.AskForETHAddress()
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	balance, err := utils.GetEthBalanceOf(address)
	if err != nil {
		return err
	}
	fmt.Println(balance, "ETHs")
	return nil
}
