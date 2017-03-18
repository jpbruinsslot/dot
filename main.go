package main

import (
	"os"

	"github.com/urfave/cli"
)

const (
	Name    = "dot"
	Usage   = "simple dotfiles tracking"
	Version = "0.2"
)

func main() {
	app := cli.NewApp()

	app.Name = Name
	app.Usage = Usage
	app.Author = "@erroneousboat"
	app.Version = Version
	app.Commands = CommandArray

	app.Run(os.Args)
}
