package main

import (
	"time"
)

const (
	// Version is the version
	Version = "0.1.1-dev"
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
	// Percentiles are percentile values.
	// Uses "linear interpolation between closest ranks method, C=1/2".
	// https://en.wikipedia.org/w/index.php?title=Percentile&oldid=724036224#First_Variant.2C_.7F.27.22.60UNIQ--postMath-0000002C-QINU.60.22.27.7F
	Percentiles map[string]time.Duration `json:"percentiles,omitempty"`
	sum         time.Duration
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
