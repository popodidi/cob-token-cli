package commands

import (
	"github.com/popodidi/cob-token-cli/commands/private"
	"github.com/popodidi/cob-token-cli/commands/public"

	"github.com/urfave/cli"
)

func All() []cli.Command {
	cmds := make([]cli.Command, 0)
	cmds = append(cmds, public.All()...)
	cmds = append(cmds, private.All()...)
	return cmds
}
