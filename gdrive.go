package main

import (
	"os"

	"github.com/grandeto/gdrive/cli"
	"github.com/grandeto/gdrive/constants"
	"github.com/grandeto/gdrive/loader"
	"github.com/grandeto/gdrive/util"
)

func main() {
	globalFlags := loader.LoadGlobalFlags()

	handlers := loader.LoadHandlers(globalFlags)

	cli.SetHandlers(handlers)

	if ok := cli.Handle(os.Args[1:]); !ok {
		util.ExitF("No valid arguments given, use '%s help' to see available commands", constants.Name)
	}
}
