# ntimes: `[time(1)](http://linux.die.net/man/1/time)` with average time, flaky rate, ..

[![Join the chat at https://gitter.im/AkihiroSuda/ntimes](https://img.shields.io/badge/GITTER-join%20chat-green.svg)](https://gitter.im/AkihiroSuda/ntimes?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)
[![Build Status](https://travis-ci.org/AkihiroSuda/ntimes.svg?branch=master)](https://travis-ci.org/AkihiroSuda/ntimes)
[![Go Report Card](https://goreportcard.com/badge/github.com/AkihiroSuda/ntimes)](https://goreportcard.com/report/github.com/AkihiroSuda/ntimes)

`ntimes` is an improved version of the  `[time(1)](http://linux.die.net/man/1/time)` command.

You can execute an command `N` times, and measure the average/max/min time taken for the execution.
You can also measure the "flaky" rate (i.e. failure rate).

## Features

- [X] average time
- [X] max and min
- [X] flaky rate
- [X] JSON output
- [ ] stdout and stderr storage
- [ ] skip "warming up" (e.g. cache) iterations
- [ ] percentile
- [ ] histogram (TBD output format)
- [ ] handle signals (e.g. ^C)
- [ ] parallel execution

## Install

    $ go get github.com/AkihiroSuda/ntimes

Currently, `ntimes` is only tested on Linux.
But it should work on macOS and on Windows as well.

## Usage

Example usage:

	$ ntimes -n 10 bash -c 'sleep=$((RANDOM%5)); fail=$((RANDOM%2)); echo "id=$NTIMES_ID, sleep=$sleep, fail=$fail"; sleep $sleep; exit $fail'
    id=0, sleep=3, fail=1
    id=1, sleep=4, fail=1
    id=2, sleep=3, fail=0
    id=3, sleep=1, fail=1
    id=4, sleep=0, fail=1
    id=5, sleep=3, fail=1
    id=6, sleep=2, fail=0
    id=7, sleep=0, fail=0
    id=8, sleep=4, fail=0
    id=9, sleep=1, fail=0
    
    average: 2.103107186s (user: 0, sys: 0)
    max: 4.00273996s (user: 0, sys: 0)
    min: 2.480336ms (user: 0, sys: 0)
    flaky: 50%


You can specify the report format using Go's [`text/template`] syntax (https://golang.org/pkg/text/template/).
Additionally to the standard functions provided by `text/template`, the `json` function is available.
Note that a `time.Duration` value is expressed in nanoseconds.


    $ ntimes --format "{{json .}}" -n 10 dd if=/dev/urandom of=/dev/null bs=512 count=1000
    1000+0 records in
    1000+0 records out
    512000 bytes (512 kB, 500 KiB) copied, 0.0448978 s, 11.4 MB/s
	...
    {"average":{"real":41427898,"user":400000,"system":37600000},"max":{"real":48461571,"user":4000000,"system":44000000},"min":{"real":37805623,"user":0,"system":32000000},"flaky":0}

Please refer to `ntimes --help` for the detailed help.

    $ ./ntimes --help
    Usage: ./ntimes [OPTIONS] COMMAND [ARG...]
      -f, --format string         format string (in golang text/template, e.g. "{{json .}}")
      -n, --repeat-n-times uint   number of times (default 1)

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
    
