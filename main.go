package main

import (
	"os"
	"time"

	"github.com/popodidi/cob-token-cli/commands"

	"github.com/urfave/cli"
)

func main() {
	cliApp := NewApp()
	cliApp.Run(os.Args)
}

func NewApp() *cli.App {
	app := cli.NewApp()
	app.Name = "cob-token-cli"
	app.Version = "0.1.9"
	app.Compiled = time.Now()
	app.Usage = "A COB token mangement command line tool"
	app.Commands = commands.All()
	return app
}
