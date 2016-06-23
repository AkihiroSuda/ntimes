% NTIMES(1) General Commands Manual
% ntimes
% JUNE 2016
# NAME
ntimes \- **time(1)** with average time, flaky rate, ..

# SYNOPSIS
**ntimes** [OPTIONS] COMMAND [arg...]

# DESCRIPTION

ntimes is an improved version of **time(1)**.

You can execute an command`N times, and measure the average/max/min time taken for the execution.
You can also measure the "flaky" rate (i.e. failure rate).


# OPTIONS
**--help**
  Print usage statement

**-f**, **--format**=*""*
  Format string (in golang text/template, e.g. "{{json .}}")

**-n**, **--repeat-n-timesg**=*1*
  Number of times

# Examples

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


You can specify the report format using Go's **text/template** syntax.
Additionally to the standard functions provided by **text/template**, the **json** function is available.
Note that a **time.Duration** value is expressed in nanoseconds.


    $ ntimes --format "{{json .}}" -n 10 dd if=/dev/urandom of=/dev/null bs=512 count=1000
    1000+0 records in
    1000+0 records out
    512000 bytes (512 kB, 500 KiB) copied, 0.0448978 s, 11.4 MB/s
	...
    {"average":{"real":41427898,"user":400000,"system":37600000},"max":{"real":48461571,"user":4000000,"system":44000000},"min":{"real":37805623,"user":0,"system":32000000},"flaky":0}

Please refer to **ntimes --help** for the detailed help.

    $ ./ntimes --help
    Usage: ./ntimes [OPTIONS] COMMAND [ARG...]
      -f, --format string         format string (in golang text/template, e.g. "{{json .}}")
      -n, --repeat-n-times uint   number of times (default 1)

# AUTHOR
June 2016, writen by Akihiro Suda
https://github.com/AkihiroSuda/ntimes
