package main

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"
)

func main() {
	if _, err := xmain(os.Args,
		os.Stdin, os.Stdout, os.Stderr); err != nil {
		if err != pflag.ErrHelp {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
		}
		os.Exit(1)
	}
}
