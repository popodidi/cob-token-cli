package public

import (
	"github.com/urfave/cli"
	"fmt"
	"github.com/popodidi/cob-token-cli/utils"
	"gopkg.in/AlecAivazis/survey.v1"
)

func cobBalanceAction(c *cli.Context) error {
	address := ""
	addressPrompt := &survey.Input{
		Message: "ETH address",
	}
	survey.AskOne(addressPrompt, &address, nil)

	balance, err := utils.GetCobBalanceOf(address)
	if err != nil {
		return err
	}
	fmt.Println(balance, "COBs")
	return nil
}
