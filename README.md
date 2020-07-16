# ntimes: `time(1)` with average time, flaky rate, ..

[![Go Report Card](https://goreportcard.com/badge/github.com/AkihiroSuda/ntimes)](https://goreportcard.com/report/github.com/AkihiroSuda/ntimes)

`ntimes` is an improved version of the  [`time(1)`](http://linux.die.net/man/1/time) command.

You can execute an command `N` times, and measure the average/max/min time taken for the execution.
You can also measure the "flaky" rate (i.e. failure rate).

## Features
Metrics (performance):

- [X] average time
- [X] max and min
- [X] standard deviation
- [X] percentile
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

For Linux and macOS (experimental):

    curl -L https://github.com/AkihiroSuda/ntimes/releases/download/v0.1.0/ntimes-`uname -s`-`uname -m` >/usr/local/bin/ntimes && \
    chmod +x /usr/local/bin/ntimes

For Windows (experimental): https://github.com/AkihiroSuda/ntimes/releases/download/v0.1.0/ntimes-windows-x86_64.exe

Latest development version (requires Go):

    $ go get github.com/AkihiroSuda/ntimes

## Usage

Example usage:

	$ ntimes -n 10 bash -c 'sleep=$((RANDOM%5)); fail=$((RANDOM%2)); echo "id=$NTIMES_ID, sleep=$sleep, fail=$fail"; sleep $sleep; exit $fail'
    id=0, sleep=4, fail=0
    id=1, sleep=0, fail=0
    id=2, sleep=3, fail=1
    id=3, sleep=4, fail=0
    id=4, sleep=3, fail=1
    id=5, sleep=0, fail=1
    id=6, sleep=3, fail=1
    id=7, sleep=0, fail=0
    id=8, sleep=3, fail=1
    id=9, sleep=3, fail=1
    
    real average: 2.303595053s, max: 4.004607594s, min: 3.455164ms, std dev: 1.636182569s
    real 99 percentile: 4.004607594s, 95 percentile: 4.004607594s, 50 percentile: 3.002190904s
    user average: 0, max: 0, min: 0, std dev: 0
    sys  average: 0, max: 0, min: 0, std dev: 0
    flaky: 60%

You can specify the report format using Go's [`text/template`](https://golang.org/pkg/text/template/) syntax.
Additionally to the standard functions provided by `text/template`, the `json` function is available as well.
Note that a `time.Duration` value is expressed in nanoseconds.


    $ ntimes --format "{{json .}}" -n 10 dd if=/dev/urandom of=/dev/null bs=512 count=1000
    1000+0 records in
    1000+0 records out
    512000 bytes (512 kB, 500 KiB) copied, 0.0448978 s, 11.4 MB/s
	...
    {"real":{"average":44155207,"max":68222928,"min":38143407,"stddev":9421337,"percentiles":{"50":39855284,"95":68222928,"99":68222928}},"user":{"average":0,"max":0,"min":0,"stddev":0},"system":{"average":36000000,"max":36000000,"min":36000000,"stddev":0},"flaky":0}


Practical example for debugging flaky tests with Namazu (`nmz`, [osrg/namazu](https://github.com/osrg/namazu)):

    $ cd some_maven_project
    $ sudo ntimes -n 10 --storage /tmp/logs nmz inspectors -cmd "mvn test"
    ...
    Flaky: 10%
    
    $ find /tmp/logs -name result.json | xargs jq .successful
    true
    true
    true
    false
    true
    true
    true
    true
    true
    true
	

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

## Specification

- When running a command, `NTIMES_ID` is passed to the command as an environment variable.
The value of `NTIMES_ID` denotes a non-negative decimal value corresponding to the iteration count.
(0, 1, .., n-1)

- The format of statistics report is defined as `Stat` structure in [`common.go`](common.go).

- If `--storage dir` is specified, following files are created in `dir/$NTIMES_ID`:
    - `stdout`: containes the standard output
	- `stderr`: containes the standard err
	- `result.json`: JSON representation of a `Result` structure (defined in [`common.go`](common.go))

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
    $ ./hack.sh all
    $ git commit -a -s
    
