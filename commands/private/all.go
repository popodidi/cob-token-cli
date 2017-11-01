package private

import (
	"github.com/urfave/cli"
)

func All() []cli.Command {
	return []cli.Command{
		{
			Name:     "send-eth",
			Category: "private",
			Usage:    "send ETHs",
			Action:   sendETHAction,
		},
		{
			Name:     "send-cob",
			Category: "private",
			Usage:    "send COBs",
			Action:   sendCOBAction,
		},
		{
			Name:     "allocate-cob",
			Category: "private",
			Usage:    "allocate COBs to multiple addresses",
			Action:   allocateCOBAction,
		},
	}
}
