package public

import (
	"fmt"

	"github.com/popodidi/cob-token-cli/utils"

	"github.com/urfave/cli"
)

func cobBalanceAction(c *cli.Context) error {
	address, err := utils.AskForETHAddress()
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	balance, err := utils.GetCobBalanceOf(address)
	if err != nil {
		return err
	}
	fmt.Println(balance, "COBs")
	return nil
}
