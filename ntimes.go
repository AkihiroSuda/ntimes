package main

import (
	"fmt"
	"io"
	"log"
	"os/exec"
	"time"
)

type result struct {
	id         int
	successful bool

	real   time.Duration
	user   time.Duration
	system time.Duration

	// status is currently unused.
	// its type is syscall.WaitStatus on UNIX.
	status interface{}
	// rusage is currently unused.
	// its type is *syscall.Rusage on UNIX.
	rusage interface{}
}

type ntimes struct {
	N       int
	Cmd     *exec.Cmd
	debug   bool
	results []result
}

func (nt *ntimes) Run() error {
	// we don't support parallelizing at the moment
	for id := 0; id < nt.N; id++ {
		r, err := nt.run(id)
		if err != nil {
			// err is a really critical error.
			// if cmd failed, result.status will be
			// set to some error value, but err is
			// not returned.
			return fmt.Errorf("error in id %d: %+v", id, err)
		}
		if nt.debug {
			log.Printf("id %d result: %#v (rusage: %#v)",
				id, r, r.rusage)
		}
		nt.results = append(nt.results, r)
	}
	return nil
}

func (nt *ntimes) prepareStdio(id int) (io.Reader, io.Writer, io.Writer, error) {
	return nt.Cmd.Stdin, nt.Cmd.Stdout, nt.Cmd.Stderr, nil
}

func (nt *ntimes) prepareCmd(id int) (*exec.Cmd, error) {
	if nt.Cmd.ExtraFiles != nil {
		err := fmt.Errorf("ExtraFiles is not supported: %v", nt.Cmd)
		return nil, err
	}
	stdin, stdout, stderr, err := nt.prepareStdio(id)
	if err != nil {
		return nil, err
	}
	env := append(nt.Cmd.Env, fmt.Sprintf("NTIMES_ID=%d", id))
	cmd := &exec.Cmd{
		Path:   nt.Cmd.Path,
		Args:   nt.Cmd.Args,
		Env:    env,
		Dir:    nt.Cmd.Dir,
		Stdin:  stdin,
		Stdout: stdout,
		Stderr: stderr,
	}
	return cmd, nil
}

func (nt *ntimes) run(id int) (result, error) {
	r := result{id: id}
	cmd, err := nt.prepareCmd(id)
	if err != nil {
		return r, nil
	}
	begin := time.Now()
	err = cmd.Run()
	end := time.Now()

	ps := cmd.ProcessState
	if ps == nil {
		return r, fmt.Errorf("ProcessState is nil for %#v", cmd)
	}
	r.successful = ps.Success()
	r.real, r.user, r.system = end.Sub(begin), ps.UserTime(), ps.SystemTime()

	if err != nil {
		if _, ok := err.(*exec.ExitError); !ok {
			return r, err
		}
	}
	r.status, r.rusage = ps.Sys(), ps.SysUsage()
	return r, nil
}

func updateMin(report *Report, result *result) {
	if report.Min.Real == 0 ||
		result.real < report.Min.Real {
		report.Min.Real = result.real
	}
	if report.Min.User == 0 ||
		result.user < report.Min.User {
		report.Min.User = result.user
	}
	if report.Min.System == 0 ||
		result.system < report.Min.System {
		report.Min.System = result.system
	}
}

func updateMax(report *Report, result *result) {
	if result.real > report.Max.Real {
		report.Max.Real = result.real
	}
	if result.user > report.Max.User {
		report.Max.User = result.user
	}
	if result.system > report.Max.System {
		report.Max.System = result.system
	}
}

func (nt *ntimes) Report() (*Report, error) {
	report := &Report{}
	sReal, sUser, sSystem := int64(0), int64(0), int64(0)
	flakies := 0
	n := len(nt.results)
	if n != nt.N && nt.debug {
		log.Printf("WARNING: len(results)=%d, but N=%d",
			n, nt.N)
	}
	// we do not need to parallelize this loop.
	// fixme: ugly code clone for {real,user,system}
	for id, result := range nt.results {
		if id != result.id && nt.debug {
			log.Printf("WARNING: id=%d, but result.id=%d",
				id, result.id)
		}
		sReal += int64(result.real)
		sUser += int64(result.user)
		sSystem += int64(result.system)

		updateMin(report, &result)
		updateMax(report, &result)

		if !result.successful {
			flakies++
		}
	}
	report.Average.Real = time.Duration(sReal / int64(n))
	report.Average.User = time.Duration(sUser / int64(n))
	report.Average.System = time.Duration(sSystem / int64(n))
	report.Flaky = 100 * float64(flakies) / float64(n)
	return report, nil
}
