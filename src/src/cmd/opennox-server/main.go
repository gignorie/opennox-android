package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/opennox/opennox/v1"
)

func main() {
	if err := opennox.RunServerArgs(os.Args); err != nil && err != flag.ErrHelp {
		if code, ok := err.(opennox.ErrExit); ok {
			os.Exit(int(code))
		} else {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}
