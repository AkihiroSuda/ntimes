package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"
)

type ntimes struct {
	N       int
	Cmd     *exec.Cmd
	Storage string
	WarmUp  int
	debug   bool
	results []Result
}

func (nt *ntimes) Run() error {
	if nt.Storage != "" {
		// MkdirAll does nothing if storage exists
		if err := os.MkdirAll(nt.Storage, 0755); err != nil {
			return err
		}
	}
	// we don't support parallelizing at the moment
	for id := 0; id < nt.N; id++ {
		if err := nt.prepareStorageDir(id); err != nil {
			return err
		}
		r, err := nt.run(id)
		if err != nil {
			// err is a really critical error.
			// if cmd failed, result.status will be
			// set to some error value, but err is
			// not returned.
			return fmt.Errorf("id %d: %+v", id, err)
		}
		if err = nt.storeResult(id, r); err != nil {
			return err
		}
		if id >= nt.WarmUp {
			nt.results = append(nt.results, r)
		} else {
			if nt.debug {
				log.Printf("id %d warm-up", id)
			}
		}
	}
	return nil
}

func fileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func (nt *ntimes) storageDir(id int) string {
	if nt.Storage != "" {
		return filepath.Join(nt.Storage, strconv.Itoa(id))
	}
	return ""
}

func (nt *ntimes) prepareStorageDir(id int) error {
	dir := nt.storageDir(id)
	if dir != "" {
		exists, err := fileExists(dir)
		if err != nil {
			return err
		}
		if exists {
			return os.ErrExist
		}
		if err := os.Mkdir(dir, 0755); err != nil {
			return err
		}
	}
	return nil
}

func (nt *ntimes) prepareStdio(id int) (io.Reader, io.Writer, io.Writer, error) {
	dir := nt.storageDir(id)
	if dir != "" {
		stdout, err := os.Create(filepath.Join(dir, "stdout"))
		if err != nil {
			return nil, nil, nil, err
		}
		stderr, err := os.Create(filepath.Join(dir, "stderr"))
		if err != nil {
			return nil, nil, nil, err
		}
		return nt.Cmd.Stdin, stdout, stderr, nil
	}
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

func (nt *ntimes) run(id int) (Result, error) {
	r := Result{ID: id}
	cmd, err := nt.prepareCmd(id)
	if err != nil {
		return r, err
	}
	begin := time.Now()
	err = cmd.Run()
	end := time.Now()

	ps := cmd.ProcessState
	if ps == nil {
		if err != nil {
			return r, err
		}
		return r, fmt.Errorf("ProcessState is nil for %#v", cmd)
	}
	r.Successful = ps.Success()
	r.Real, r.User, r.System = end.Sub(begin), ps.UserTime(), ps.SystemTime()

	// maybe ps is useful. so we check err here.
	if err != nil {
		if _, ok := err.(*exec.ExitError); !ok {
			return r, err
		}
	}
	r.status, r.rusage = ps.Sys(), ps.SysUsage()
	if nt.Storage != "" {
		if err = cmd.Stdout.(io.WriteCloser).Close(); err != nil {
			return r, err
		}
		if err = cmd.Stderr.(io.WriteCloser).Close(); err != nil {
			return r, err
		}
	}
	return r, nil
}

func (nt *ntimes) storeResult(id int, result Result) error {
	if id != result.ID {
		return fmt.Errorf("id %d mismatch: %+v", id, result)
	}
	dir := nt.storageDir(id)
	if dir == "" {
		// NOP, not an error
		return nil
	}
	file, err := os.Create(filepath.Join(dir, "result.json"))
	if err != nil {
		return err
	}
	enc := json.NewEncoder(file)
	if err = enc.Encode(result); err != nil {
		return err
	}
	return file.Close()
}

func (nt *ntimes) Stat() (*Stat, error) {
	stat := &Stat{}
	sels := []selector{selReal, selUser, selSystem}
	stat.Real, stat.User, stat.System =
		&TimeStat{}, &TimeStat{}, &TimeStat{}
	updateStatPhase0(stat, nt.results, sels)
	updateStatPhase1(stat, nt.results, sels)
	return stat, nil
}
