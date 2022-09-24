package handlers

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"text/tabwriter"

	"github.com/grandeto/gdrive/cli"
	"github.com/grandeto/gdrive/constants"
	"github.com/grandeto/gdrive/util"
)

func PrintVersion(ctx cli.Context) {
	fmt.Printf("%s: %s\n", constants.Name, constants.Version)
	fmt.Printf("Golang: %s\n", runtime.Version())
	fmt.Printf("OS/Arch: %s/%s\n", runtime.GOOS, runtime.GOARCH)
}

func PrintHelp(ctx cli.Context) {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 0, 3, ' ', 0)

	fmt.Fprintf(w, "%s usage:\n\n", constants.Name)

	for _, h := range ctx.Handlers() {
		fmt.Fprintf(w, "%s %s\t%s\n", constants.Name, h.Pattern, h.Description)
	}

	w.Flush()
}

func PrintCommandHelp(ctx cli.Context) {
	args := ctx.Args()
	PrintCommandPrefixHelp(ctx, args.String("command"))
}

func PrintSubCommandHelp(ctx cli.Context) {
	args := ctx.Args()
	PrintCommandPrefixHelp(ctx, args.String("command"), args.String("subcommand"))
}

func PrintCommandPrefixHelp(ctx cli.Context, prefix ...string) {
	handler := getHandler(ctx.Handlers(), prefix)

	if handler == nil {
		util.ExitF("Command not found")
	}

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 0, 3, ' ', 0)

	fmt.Fprintf(w, "%s\n", handler.Description)
	fmt.Fprintf(w, "%s %s\n", constants.Name, handler.Pattern)
	for _, group := range handler.FlagGroups {
		fmt.Fprintf(w, "\n%s:\n", group.Name)
		for _, flag := range group.Flags {
			boolFlag, isBool := flag.(cli.BoolFlag)
			if isBool && boolFlag.OmitValue {
				fmt.Fprintf(w, "  %s\t%s\n", strings.Join(flag.GetPatterns(), ", "), flag.GetDescription())
			} else {
				fmt.Fprintf(w, "  %s <%s>\t%s\n", strings.Join(flag.GetPatterns(), ", "), flag.GetName(), flag.GetDescription())
			}
		}
	}

	w.Flush()
}

func getHandler(handlers []*cli.Handler, prefix []string) *cli.Handler {
	for _, h := range handlers {
		pattern := stripOptionals(h.SplitPattern())

		if len(prefix) > len(pattern) {
			continue
		}

		if util.Equal(prefix, pattern[:len(prefix)]) {
			return h
		}
	}

	return nil
}

// Strip optional groups (<...>) from pattern
func stripOptionals(pattern []string) []string {
	newArgs := []string{}

	for _, arg := range pattern {
		if strings.HasPrefix(arg, "[") && strings.HasSuffix(arg, "]") {
			continue
		}
		newArgs = append(newArgs, arg)
	}
	return newArgs
}
