package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/spf13/pflag"
)

const defaultFormatTemplate = `
real average: {{.Real.Average}}, max: {{.Real.Max}}, min: {{.Real.Min}}, std dev: {{.Real.StdDev}}
real 99 percentile: {{index .Real.Percentiles "99"}}, 95 percentile: {{index .Real.Percentiles "95"}}, 50 percentile: {{index .Real.Percentiles "50"}}
user average: {{.User.Average}}, max: {{.User.Max}}, min: {{.User.Min}}, std dev: {{.User.StdDev}}
sys  average: {{.System.Average}}, max: {{.System.Max}}, min: {{.System.Min}}, std dev: {{.System.StdDev}}
flaky: {{.Flaky}}%`

type parsed struct {
	args    []string
	n       uint
	format  string
	storage string
	warmup  uint
	version bool
	debug   bool
	fs      *pflag.FlagSet
}

func parseArgs(args []string, stdin io.Reader, stdout, stderr io.Writer) (*parsed, error) {
	p := &parsed{}
	p.fs = pflag.NewFlagSet(args[0], pflag.ContinueOnError)
	p.fs.SetOutput(stderr)
	p.fs.Usage = func() {
		fmt.Fprintf(stderr, "Usage: %s [OPTIONS] COMMAND [ARG...]\n", args[0])
		p.fs.PrintDefaults()
	}
	p.fs.SetInterspersed(false)
	p.fs.UintVarP(&p.n, "repeat-n-times", "n", 1, "number of times")
	p.fs.StringVarP(&p.format, "format", "f", "", "format string (in golang text/template, e.g. \"{{json .}}\")")
	p.fs.StringVar(&p.storage, "storage", "", "path to stdout, stderr storage")
	p.fs.UintVar(&p.warmup, "warm-up", 0, "skip first n iterations for stat")
	p.fs.BoolVar(&p.version, "version", false, "print version to stdout and exit")
	p.fs.BoolVarP(&p.debug, "debug", "", false, "do not use")
	if err := p.fs.MarkHidden("debug"); err != nil {
		return p, err
	}
	if err := p.fs.Parse(args[1:]); err != nil {
		return p, err
	}
	p.args = p.fs.Args()
	return p, nil
}

// xmain is main function made friendly to unit testing.
// xmain may return *Stat.
func xmain(args []string, stdin io.Reader, stdout, stderr io.Writer) (interface{}, error) {
	p, err := parseArgs(args, stdin, stdout, stderr)
	if err != nil {
		return nil, err
	}
	if p.version {
		// should we use stderr here?
		fmt.Fprintf(stdout, "%s\n", Version)
		return nil, nil
	}
	if p.n == 0 {
		return nil, fmt.Errorf("n must be > 0")
	}
	if p.warmup >= p.n {
		return nil, fmt.Errorf("warm-up must be < n")
	}
	if len(p.args) < 1 {
		return nil, fmt.Errorf("no command specified."+
			"Try '%s --help' for more information.",
			args[0])
	}
	command := p.args[0]
	var commandArgs []string
	if len(p.args) > 1 {
		commandArgs = p.args[1:]
	}
	// prepare cmd os/exec.(*Cmd)
	cmd := exec.Command(command, commandArgs...)
	cmd.Env = os.Environ()
	cmd.Stdin = stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	// prepare nt *ntimes
	nt := &ntimes{
		N:       int(p.n),
		Cmd:     cmd,
		Storage: p.storage,
		WarmUp:  int(p.warmup),
		debug:   p.debug,
	}
	// prepare f *formatter
	if p.format == "" {
		p.format = defaultFormatTemplate
	}
	f := &formatter{
		Format: p.format,
		// time(1) uses stderr. so we us stderr as well here.
		Writer: stderr,
	}
	return doit(nt, f)
}

func doit(nt *ntimes, f *formatter) (interface{}, error) {
	if err := f.Init(); err != nil {
		return nil, err
	}
	if err := nt.Run(); err != nil {
		return nil, err
	}
	stat, err := nt.Stat()
	if err != nil {
		return stat, err
	}
	if err = f.Execute(stat); err != nil {
		return stat, err
	}
	_, err = f.Writer.Write([]byte("\n"))
	return stat, err
}
