package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/pflag"
)

const defaultFormatTemplate = `
average: {{.Average.Real}} (user: {{.Average.User}}, sys: {{.Average.System}})
max: {{.Max.Real}} (user: {{.Max.User}}, sys: {{.Max.System}})
min: {{.Min.Real}} (user: {{.Min.User}}, sys: {{.Min.System}})
flaky: {{.Flaky}}%`

func main() {
	errh := func(err error) {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	nt, f, err := parseArgs(os.Args)
	if err != nil {
		errh(err)
	}
	if err = f.Init(); err != nil {
		errh(err)
	}
	if err = nt.Run(); err != nil {
		errh(err)
	}
	report, err := nt.Report()
	if err != nil {
		errh(err)
	}
	if err = f.Execute(report); err != nil {
		errh(err)
	}
	if _, err = f.Writer.Write([]byte("\n")); err != nil {
		errh(err)
	}
	if err = f.Writer.Flush(); err != nil {
		errh(err)
	}
}

func parseArgs(args []string) (nt *ntimes, f *formatter, err error) {
	var (
		n      uint
		format string
		debug  bool
	)
	fs := pflag.NewFlagSet(args[0], pflag.ExitOnError)
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS] COMMAND [ARG...]\n", args[0])
		fs.PrintDefaults()
	}
	fs.SetInterspersed(false)
	fs.UintVarP(&n, "repeat-n-times", "n", 1, "number of times")
	fs.StringVarP(&format, "format", "f", "", "format string (in golang text/template, e.g. \"{{json .}}\")")
	fs.BoolVarP(&debug, "debug", "", false, "do not use")
	if err = fs.MarkHidden("debug"); err != nil {
		return
	}
	if err = fs.Parse(args[1:]); err != nil {
		return
	}
	parsedArgs := fs.Args()
	if n == 0 {
		err = fmt.Errorf("n must be > 0")
		return
	}
	if len(parsedArgs) < 1 {
		err = fmt.Errorf("no command specified")
		return
	}
	command := parsedArgs[0]
	var commandArgs []string
	if len(parsedArgs) > 1 {
		commandArgs = parsedArgs[1:]
	}
	// prepare cmd os/exec.(*Cmd)
	cmd := exec.Command(command, commandArgs...)
	cmd.Env = os.Environ()
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	// prepare nt *ntimes
	nt = &ntimes{
		N:     int(n),
		Cmd:   cmd,
		debug: debug,
	}
	// prepare f *formatter
	if format == "" {
		format = defaultFormatTemplate
	}
	f = &formatter{
		Format: format,
		// time(1) uses stderr. so we us stderr as well here.
		Writer: bufio.NewWriter(os.Stderr),
	}
	return
}
