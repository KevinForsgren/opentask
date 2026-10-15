package main

import (
	"context"
	"fmt"
	"os"

	"github.com/kevinforsgren/opentask/internal/cli"
)

func main() {
	if err := cli.Execute(context.Background()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(cli.ExitCode(err))
	}
}
