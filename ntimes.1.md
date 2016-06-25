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

**--storage**
  Path to stdout,stderr storage	

**--version**
  Print version to stdout and exit

**--warm-up**
  Skip first n iterations for stat

# Examples

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
    
You can specify the report format using Go's **text/template** syntax.
Additionally to the standard functions provided by **text/template**, the **json** function is available.
Note that a **time.Duration** value is expressed in nanoseconds.

    $ ntimes --format "{{json .}}" -n 10 dd if=/dev/urandom of=/dev/null bs=512 count=1000
    1000+0 records in
    1000+0 records out
    512000 bytes (512 kB, 500 KiB) copied, 0.0448978 s, 11.4 MB/s
	...
    {"real":{"average":51550496,"max":84482238,"min":38731345,"stddev":16170565},"user":{"average":800000,"max":4000000,"min":0,"stddev":1686548},"system":{"average":39600000,"max":52000000,"min":32000000,"stddev":7876829},"flaky":0}

# SEE ALSO
**time(1)**

# AUTHOR
June 2016, writen by Akihiro Suda
https://github.com/AkihiroSuda/ntimes
