package public

import (
	"github.com/urfave/cli"
)

func All() []cli.Command {
	return []cli.Command{
		{
			Name:     "eth-balance",
			Category: "public",
			Usage:    "check ETH balance of address",
			Action:   ethBalanceAction,
		},
		{
			Name:     "cob-balance",
			Category: "public",
			Usage:    "check COB balance of address",
			Action:   cobBalanceAction,
		},
	}
}
