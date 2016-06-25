# ntimes: `time(1)` with average time, flaky rate, ..

[![Join the chat at https://gitter.im/AkihiroSuda/ntimes](https://img.shields.io/badge/GITTER-join%20chat-green.svg)](https://gitter.im/AkihiroSuda/ntimes?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)
[![Build Status](https://travis-ci.org/AkihiroSuda/ntimes.svg?branch=master)](https://travis-ci.org/AkihiroSuda/ntimes)
[![Go Report Card](https://goreportcard.com/badge/github.com/AkihiroSuda/ntimes)](https://goreportcard.com/report/github.com/AkihiroSuda/ntimes)

`ntimes` is an improved version of the  [`time(1)`](http://linux.die.net/man/1/time) command.

You can execute an command `N` times, and measure the average/max/min time taken for the execution.
You can also measure the "flaky" rate (i.e. failure rate).

## Features
Metrics (performance):

- [X] average time
- [X] max and min
- [X] standard deviation
- [ ] percentile
- [ ] histogram (TBD output format)

Metrics (flakiness):

- [X] flaky rate

Others:

- [X] JSON output
- [X] stdout and stderr storage
- [X] skip "warming up" (e.g. cache) iterations for stat
- [ ] handle signals (e.g. ^C)
- [ ] parallel execution

## Install

    $ go get github.com/AkihiroSuda/ntimes

Currently, `ntimes` is only tested on Linux.
But it should work on macOS and on Windows as well.

## Usage

Example usage:

	$ ntimes -n 10 bash -c 'sleep=$((RANDOM%5)); fail=$((RANDOM%2)); echo "id=$NTIMES_ID, sleep=$sleep, fail=$fail"; sleep $sleep; exit $fail'
    id=0, sleep=1, fail=1
    id=1, sleep=1, fail=1
    id=2, sleep=3, fail=1
    id=3, sleep=4, fail=0
    id=4, sleep=3, fail=0
    id=5, sleep=3, fail=1
    id=6, sleep=0, fail=0
    id=7, sleep=4, fail=0
    id=8, sleep=1, fail=1
    id=9, sleep=3, fail=1
    
    real average: 2.304027218s, max: 4.004967847s, min: 3.962992ms, std dev: 1.418004182s
    user average: 0, max: 0, min: 0, std dev: 0
    sys  average: 0, max: 0, min: 0, std dev: 0
    flaky: 60%
    
You can specify the report format using Go's [`text/template`](https://golang.org/pkg/text/template/) syntax.
Additionally to the standard functions provided by `text/template`, the `json` function is available.
Note that a `time.Duration` value is expressed in nanoseconds.


    $ ntimes --format "{{json .}}" -n 10 dd if=/dev/urandom of=/dev/null bs=512 count=1000
    1000+0 records in
    1000+0 records out
    512000 bytes (512 kB, 500 KiB) copied, 0.0448978 s, 11.4 MB/s
	...
    {"real":{"average":51550496,"max":84482238,"min":38731345,"stddev":16170565},"user":{"average":800000,"max":4000000,"min":0,"stddev":1686548},"system":{"average":39600000,"max":52000000,"min":32000000,"stddev":7876829},"flaky":0}

Please refer to `ntimes --help` for the detailed help.

    $ ./ntimes --help
    Usage: ./ntimes [OPTIONS] COMMAND [ARG...]
      -f, --format string         format string (in golang text/template, e.g. "{{json .}}")
      -n, --repeat-n-times uint   number of times (default 1)
      --storage string            path to stdout,stderr storage
      --version                   print version to stdout and exit
      --warm-up uint              skip first n iterations for stat

## Motivation

Originally, `ntimes` was designed so that it can be combined with [osrg/namazu](https://github.com/osrg/namazu).
Namazu can be used for controlling non-deternimism and increasing reproducibility of flaky test failures.
`ntimes` can be used for measuring the reproducibility of flaky test failures.

## Related Project

Some projects similar to `ntimes`:

- [jmcabo/avgtime](https://github.com/jmcabo/avgtime)
- [ryanmjacobs/avgtime](https://github.com/ryanmjacobs/avgtime)
- [kevinstreit/avgtime](https://github.com/kevinstreit/avgtime)


## How to Contribute
Please feel free to send your pull requests on github!

    $ git clone https://github.com/AkihiroSuda/ntimes
    $ cd ntimes
    $ git checkout -b your-branch
    $ vim foo.go
    $ ./hack.sh
    $ git commit -a -s
    
