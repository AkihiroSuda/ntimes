package main

import (
	"time"
)

const (
	// Version is the version
	Version = "0.0.1-dev"
)

// Result is a result of an experiment
type Result struct {
	ID         int           `json:"id"`
	Successful bool          `json:"successful"`
	Real       time.Duration `json:"real"`
	User       time.Duration `json:"user"`
	System     time.Duration `json:"system"`
	// status is currently unused.
	// its type is syscall.WaitStatus on UNIX.
	status interface{}
	// rusage is currently unused.
	// its type is *syscall.Rusage on UNIX.
	rusage interface{}
}

// TimeStat can be included in Stat
type TimeStat struct {
	Average time.Duration `json:"average"`
	Max     time.Duration `json:"max"`
	Min     time.Duration `json:"min"`
	// StdDev is a sample standard deviation,
	// i.e. $\sqrt(varianceNumerator/(n-1))$
	StdDev time.Duration `json:"stddev"`
	sum    time.Duration
	// varianceNumerator is $\sum (t-Average)^2$
	varianceNumerator float64
}

// Stat is a JSON-compatible stat report for ntimes.
// The format is not fixed yet.
type Stat struct {
	Real   *TimeStat `json:"real,omitempty"`
	User   *TimeStat `json:"user,omitempty"`
	System *TimeStat `json:"system,omitempty"`
	// Flaky is the percentage of failures (0.0-100.0).
	Flaky float64 `json:"flaky"`
	// nflakies is the number of failures
	nflakies int
}
